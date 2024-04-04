package lib

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
)

type fieldInfo struct {
	Name        string
	Type        string
	IsComposite bool
}

func GenerateTableName(table interface{}) string {
	tableType := reflect.TypeOf(table)
	tableValue := reflect.ValueOf(table)
	if tableType.Kind() == reflect.Ptr {
		tableType = tableType.Elem()
		tableValue = reflect.New(tableType)
	}
	if tableType.Kind() == reflect.Slice {
		tableType = tableType.Elem()
		tableValue = reflect.New(tableType)
	}
	tableName := tableType.Name()
	tableName = strcase.ToSnake(tableName)
	if method := tableValue.MethodByName("TableName"); method.IsValid() {
		rs := method.Call([]reflect.Value{})
		tableName = rs[0].String()
	}

	return tableName
}

func getTableInfo(table interface{}) ([]fieldInfo, error) {
	tableType := reflect.TypeOf(table)
	tableValue := reflect.ValueOf(table)

	if tableType.Kind() == reflect.Ptr {
		tableType = tableType.Elem()
		tableValue = tableValue.Elem()
	}

	if tableType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("getTableInfo: table is expected to be struct, got %v", tableType.Kind())
	}

	var fields []fieldInfo
	for i := 0; i < tableType.NumField(); i++ {
		fieldType := tableType.Field(i)
		fieldValue := tableValue.Field(i)
		if !fieldType.IsExported() {
			fmt.Println("non-exported fields: ", fieldType.Name)
			continue
		}

		field := getFieldInfo(fieldType, fieldValue)

		fields = append(fields, field)
	}

	return fields, nil
}

func createTable(ctx context.Context, conn *sql.Conn, tableName string, fields []fieldInfo) error {
	query := createTableQuery(tableName, fields)
	_, err := ExecuteWriteQuery(ctx, query, conn)
	return err
}

func addMissingColumns(ctx context.Context, conn *sql.Conn, tableName string, fields []fieldInfo) error {
	columns, err := getExistingColumns(ctx, conn, tableName)
	if err != nil {
		return err
	}

	missingColumns := getMissingColumns(fields, columns)
	if len(missingColumns) > 0 {
		alterQuery := generateAddColumnQuery(tableName, missingColumns)
		_, err = ExecuteWriteQuery(ctx, alterQuery, conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func getFieldInfo(fieldType reflect.StructField, fieldValue reflect.Value) fieldInfo {
	fieldName := getFieldName(fieldType)
	columnConstraint, autoincr, isComposite := getFieldConstraint(fieldType)
	if columnConstraint != "" {
		columnConstraint = " " + columnConstraint
	}
	sqlType := getSQLType(fieldValue.Type(), autoincr)
	return fieldInfo{
		Name:        fieldName,
		Type:        sqlType + columnConstraint,
		IsComposite: isComposite,
	}
}

func getFieldName(fieldType reflect.StructField) string {
	fieldName := fieldType.Name
	if dbTag := fieldType.Tag.Get("db"); dbTag != "" {
		colName := strings.Split(dbTag, ",")[0]
		if colName != "" {
			fieldName = colName
		}
	}

	return strcase.ToSnake(fieldName)
}

func getFieldConstraint(fieldType reflect.StructField) (fc string, autoincr bool, isComposite bool) {
	constraints := []string{}
	if dbTag := fieldType.Tag.Get("db"); dbTag != "" {
		tagParts := strings.Split(dbTag, ",")
		if len(tagParts) > 1 {
			for _, part := range strings.Fields(tagParts[1]) {
				switch strings.ToUpper(part) {
				case "PK":
					constraints = append(constraints, "PRIMARY KEY")
				case "UQ":
					constraints = append(constraints, "UNIQUE")
				case "UQS":
					isComposite = true
				case "AUTOINCR":
					autoincr = true
				}
			}
		}
	}

	return strings.Join(constraints, " "), autoincr, isComposite
}

func getUniqueColumnGroups(t reflect.Type) [][]string {
	groups := map[int][]string{}
	groupIndex := 0
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if dbTag := field.Tag.Get("db"); dbTag != "" {
			tagParts := strings.Split(dbTag, ",")
			for _, part := range tagParts[1:] {
				if strings.ToUpper(part) == "UQS" {
					groups[groupIndex] = append(groups[groupIndex], getFieldName(field))
					groupIndex++
				}
			}
		}
	}

	result := [][]string{}
	for _, group := range groups {
		result = append(result, group)
	}

	return result
}

func getSQLType(fieldType reflect.Type, autoincr bool) string {
	if autoincr {
		switch fieldType.Kind() {
		case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint64:
			return "SERIAL"
		}
	}

	switch fieldType.Kind() {
	case reflect.Int, reflect.Int32:
		return "INTEGER"
	case reflect.Int64, reflect.Uint64:
		return "BIGINT"
	case reflect.Float32, reflect.Float64:
		return "FLOAT"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.String:
		return "VARCHAR(255)"
	case reflect.Struct:
		if fieldType == reflect.TypeOf(time.Time{}) {
			return "TIMESTAMP WITH TIME ZONE"
		}
	}

	return ""
}

func tableExists(ctx context.Context, conn *sql.Conn, tableName string) (bool, error) {
	tableQuery := "" +
		"SELECT EXISTS (" +
		"    SELECT FROM " +
		"        information_schema.tables " +
		"    WHERE " +
		"        table_schema LIKE 'public' AND " +
		"        table_name = $1" +
		");"

	var exists bool
	err := conn.QueryRowContext(ctx, tableQuery, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if table exists: %v", err)
	}

	return exists, nil
}

func createTableQuery(tableName string, fields []fieldInfo) string {
	var columnDefs []string
	var compositeKeyGroup []string
	for _, field := range fields {
		columnDefs = append(columnDefs, fmt.Sprintf("%s %s", field.Name, field.Type))
		if field.IsComposite {
			compositeKeyGroup = append(compositeKeyGroup, field.Name)
		}
	}

	columnSQL := strings.Join(columnDefs, ", ")
	if len(compositeKeyGroup) > 0 {
		compositeKeySQL := fmt.Sprintf("UNIQUE(%s)", strings.Join(compositeKeyGroup, ", "))
		columnSQL += ", " + compositeKeySQL
	}

	return fmt.Sprintf("CREATE TABLE \"%s\" (%s);", tableName, columnSQL)
}

func generateAddColumnQuery(tableName string, missingColumns []string) string {
	alterQuery := fmt.Sprintf("ALTER TABLE \"%s\" ", tableName)
	var addColumns []string
	for _, col := range missingColumns {
		addColumns = append(addColumns, fmt.Sprintf("ADD COLUMN %s", col))
	}

	alterQuery += strings.Join(addColumns, ", ")
	return alterQuery
}

func getExistingColumns(ctx context.Context, conn *sql.Conn, tableName string) ([]string, error) {
	var columns []string

	rows, err := conn.QueryContext(ctx, fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name='%s'", tableName))
	if err != nil {
		return nil, fmt.Errorf("error getting columns for table %s: %v", tableName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var column string
		err = rows.Scan(&column)
		if err != nil {
			return nil, fmt.Errorf("error scanning column for table %s: %v", tableName, err)
		}
		columns = append(columns, column)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting columns for table %s: %v", tableName, err)
	}

	return columns, nil
}

func getMissingColumns(fields []fieldInfo, columns []string) []string {
	var missingColumns []string

	for _, f := range fields {
		if !contains(columns, f.Name) {
			missingColumns = append(missingColumns, fmt.Sprintf("%s %s", f.Name, f.Type))
		}
	}

	return missingColumns
}

func getUniqueConstraints(ctx context.Context, conn *sql.Conn, tableName string) ([][]string, error) {
	query := `
	SELECT kcu.column_name
	FROM information_schema.table_constraints tc
	JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name
	WHERE tc.table_name = $1 AND tc.constraint_type = 'UNIQUE'
	ORDER BY kcu.ordinal_position;
	`

	rows, err := conn.QueryContext(ctx, query, tableName)
	if err != nil {
		return nil, fmt.Errorf("error getting unique constraints for table %s: %v", tableName, err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string
		err = rows.Scan(&column)
		if err != nil {
			return nil, fmt.Errorf("error scanning unique constraint for table %s: %v", tableName, err)
		}
		columns = append(columns, column)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting unique constraints for table %s: %v", tableName, err)
	}

	var result [][]string
	if len(columns) > 0 {
		result = append(result, columns)
	}

	return result, nil
}

func generateDropConstraintStatement(tableName string, uqConstraints [][]string) string {
	sql := fmt.Sprintf("ALTER TABLE \"%s\" ", tableName)

	var dropConstraints []string
	for i := range uqConstraints {
		dropConstraints = append(dropConstraints,
			fmt.Sprintf("DROP CONSTRAINT IF EXISTS %s_uq_%d", tableName,
				i))
	}

	sql += strings.Join(dropConstraints, ", ")

	return sql
}

func generateAddConstraintStatement(tableName string,
	uqGroups [][]string) string {

	sql := fmt.Sprintf("ALTER TABLE \"%s\" ", tableName)

	var addConstraints []string
	for i, group := range uqGroups {
		addConstraints = append(addConstraints,
			fmt.Sprintf("ADD CONSTRAINT %s_uq_%d UNIQUE(%s)",
				tableName,
				i,
				strings.Join(group,
					", ")))
	}

	sql += strings.Join(addConstraints,
		", ")

	return sql
}

func contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func SyncTable(ctx context.Context, conn *sql.Conn, table any) error {
	tableName := GenerateTableName(table)
	fields, err := getTableInfo(table)
	if err != nil {
		return err
	}

	if exist, err := tableExists(ctx, conn, tableName); err != nil {
		return err
	} else if !exist {
		if err = createTable(ctx, conn, tableName, fields); err != nil {
			return err
		}
	} else {
		if err = addMissingColumns(ctx, conn, tableName, fields); err != nil {
			return err
		}
		//if err = updateUniqueCompositeConstraints(ctx, p.conn, tableName); err != nil {
		//	return err
		//}
	}
	return nil
}
