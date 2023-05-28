package pkg

import (
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
)

func Print(doc interface{}) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	hr := table.Row{}
	drs := []table.Row{}
	switch reflect.TypeOf(doc).Kind() {
	case reflect.Slice:
		hr = generateHeaderRow(reflect.ValueOf(doc).Index(0).Interface())
		drs = generateDataRows(doc)
	default:
		hr = generateHeaderRow(doc)
		dr := generateDataRow(doc)
		drs = []table.Row{dr}
	}

	t.AppendHeader(hr)
	for idx := range drs {
		t.AppendRow(drs[idx])
	}
	t.SetStyle(table.StyleLight)
	//t.SetStyle(table.StyleColoredBright)
	t.Render()
}

func generateHeaderRow(doc interface{}) table.Row {
	hr := table.Row{}
	val := reflect.ValueOf(doc)
	if reflect.TypeOf(doc).Kind() == reflect.Pointer {
		val = val.Elem()
	}
	for idx := 0; idx < val.NumField(); idx++ {
		hr = append(hr, val.Type().Field(idx).Name)
	}

	return hr
}

func generateDataRows(doc interface{}) []table.Row {
	typ := reflect.TypeOf(doc)
	val := reflect.ValueOf(doc)
	if typ.Kind() != reflect.Slice || val.Len() == 0 {
		return []table.Row{}
	}

	drs := make([]table.Row, 0, val.Len())

	for idx := 0; idx < val.Len(); idx++ {
		dr := generateDataRow(val.Index(idx).Interface())
		drs = append(drs, dr)
	}

	return drs
}

func generateDataRow(doc interface{}) table.Row {
	if reflect.TypeOf(doc).Kind() == reflect.Slice {
		return table.Row{}
	}

	val := reflect.ValueOf(doc)
	if reflect.TypeOf(doc).Kind() == reflect.Pointer {
		val = val.Elem()
	}
	dr := table.Row{}
	for idx := 0; idx < val.NumField(); idx++ {
		dr = append(dr, val.Field(idx).Interface())
	}
	return dr
}
