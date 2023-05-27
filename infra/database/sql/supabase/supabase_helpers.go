package supabase

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/models"

	"github.com/iancoleman/strcase"
	"github.com/nedpals/supabase-go"
)

type Doc struct {
	ID string
}

type keyValue struct {
	key   string
	value string
}

func InitializeSupabase(ctx context.Context) *supabase.Client {
	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	return supabase.CreateClient(supabaseUrl, supabaseKey)
}

func generateFilters(filter interface{}) []keyValue {
	kvs := []keyValue{}
	val := reflect.ValueOf(filter)
	for idx := 0; idx < val.NumMethod(); idx++ {
		field := val.Type().Field(idx)
		if val.Field(idx).IsZero() {
			continue
		}

		key := field.Tag.Get("db")
		if key == "" {
			key = strcase.ToSnake(field.Name)
		}

		kv := keyValue{
			key:   key,
			value: toString(val.Field(idx).Interface()),
		}

		kvs = append(kvs, kv)
	}

	return kvs
}

func toDBFieldName(fieldName string) string {
	return strcase.ToSnake(fieldName)
}

func fromDBFieldName(fieldName string) string {
	return strcase.ToLowerCamel(fieldName)
}

func toString(val interface{}) string {
	var value string
	switch v := val.(type) {
	case string:
		value = fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case time.Time:
		value = fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
	case []string:
		value = fmt.Sprintf("('%s')", strings.Join(v, "', '"))
	case []any:
		value = handleSliceAny(v)
	default:
		value = fmt.Sprintf("%v", v)
	}

	return value
}

func handleSliceAny(v []any) string {
	var value string
	var vals []string
	typ := reflect.String.String()
	for _, elem := range v {
		if str, ok := elem.(string); ok {
			vals = append(vals, str)
		} else {
			typ = reflect.Interface.String()
			vals = append(vals, fmt.Sprintf("%v", elem))
		}
	}

	if typ == reflect.String.String() {
		value = fmt.Sprintf("('%s')", strings.Join(vals, "', '"))
	} else {
		value = fmt.Sprintf("(%s)", strings.Join(vals, ", "))
	}
	return value
}

func checkIDNonEmpty(id string) error {
	if id == "" {
		return models.StatusError{
			Status:  http.StatusBadRequest,
			Message: "must provide document id",
		}
	}
	return nil
}

func checkIdOrFilterNonEmpty(id string, filter interface{}) error {
	if id == "" && filter == nil {
		return models.StatusError{
			Status:  http.StatusBadRequest,
			Message: "must provide id and/or filter",
		}
	}
	return nil
}
