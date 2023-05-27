package server

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/models"

	"github.com/iancoleman/strcase"
)

type fieldInfo struct {
	Name string
	Type string
}

func getTableName(table interface{}) string {
	tableType := reflect.TypeOf(table)
	tableValue := reflect.ValueOf(table)
	if tableType.Kind() == reflect.Ptr {
		tableType = tableType.Elem()
		tableValue = tableValue.Elem()
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
		// Skip any field that is not exported (starts with a lowercase letter)
		if fieldType.PkgPath != "" {
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
	_, err := executeWriteQuery(ctx, query, conn)
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
		_, err = executeWriteQuery(ctx, alterQuery, conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func getFieldInfo(fieldType reflect.StructField, fieldValue reflect.Value) fieldInfo {
	fieldName := getFieldName(fieldType)
	sqlType := getSQLType(fieldValue.Type())
	columnConstraint := getFieldConstraint(fieldName)
	return fieldInfo{
		Name: fieldName,
		Type: sqlType + columnConstraint,
	}
}

func getFieldName(fieldType reflect.StructField) string {
	fieldName := fieldType.Name
	if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
		fieldName = strings.Split(jsonTag, ",")[0]
	}

	return strcase.ToSnake(fieldName)
}

func getFieldConstraint(fieldName string) string {
	if strings.ToLower(fieldName) == "id" {
		return " PRIMARY KEY"
	}
	return ""
}

func getSQLType(fieldType reflect.Type) string {
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
	var fieldsStr []string
	for _, field := range fields {
		fieldsStr = append(fieldsStr, fmt.Sprintf("%s %s", field.Name, field.Type))
	}
	return fmt.Sprintf("CREATE TABLE \"%s\" (%s);", tableName, strings.Join(fieldsStr, ", "))
}

func generateAddColumnQuery(tableName string, missingColumns []string) string {
	alterQuery := fmt.Sprintf("ALTER TABLE %s ", tableName)
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

func contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func getModels() []interface{} {
	mds := make([]interface{}, 0)
	mds = append(mds, models.User{})
	return mds
}
