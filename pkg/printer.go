package pkg

import (
	"io"
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Printer interface {
	PrintDocument(doc interface{}) string
	PrintDocuments(docs interface{}) string
}

type Options struct {
	Style        table.Style
	EnableStdout bool
}

type printer struct {
	options Options
}

func NewPrinter(options Options) Printer {
	return &printer{
		options: options,
	}
}

func (p *printer) PrintDocument(doc interface{}) string {
	t := p.newTableWriter()

	hr := generateHeaderRow(doc)
	dr := generateDataRow(doc)

	t.AppendHeader(hr)
	t.AppendRow(dr)
	return t.Render()
}

func (p *printer) PrintDocuments(docs interface{}) string {
	t := p.newTableWriter()

	val := reflect.ValueOf(docs)
	if val.Kind() != reflect.Slice || val.Len() == 0 {
		return ""
	}

	doc := val.Index(0).Interface()

	hr := generateHeaderRow(doc)
	drs := generateDataRows(docs)

	t.AppendHeader(hr)
	for _, dr := range drs {
		t.AppendRow(dr)
	}
	return t.Render()
}

func (p *printer) newTableWriter() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(p.getWriter())
	t.SetStyle(p.options.Style)
	return t
}

func (p *printer) getWriter() io.Writer {
	if p.options.EnableStdout {
		return os.Stdout
	}

	return nil
}

func generateHeaderRow(doc interface{}) table.Row {
	hr := table.Row{}
	val := reflect.ValueOf(doc)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for idx := 0; idx < val.NumField(); idx++ {
		hr = append(hr, val.Type().Field(idx).Name)
	}

	return hr
}

func generateDataRows(docs interface{}) []table.Row {
	drs := []table.Row{}

	val := reflect.ValueOf(docs)
	if val.Kind() != reflect.Slice || val.Len() == 0 {
		return drs
	}

	for idx := 0; idx < val.Len(); idx++ {
		doc := val.Index(idx).Interface()
		dr := generateDataRow(doc)
		drs = append(drs, dr)
	}

	return drs
}

func generateDataRow(doc interface{}) table.Row {
	val := reflect.ValueOf(doc)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	dr := table.Row{}
	for idx := 0; idx < val.NumField(); idx++ {
		dr = append(dr, val.Field(idx).Interface())
	}
	return dr
}
