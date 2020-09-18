package exporter

import (
	"testing"

	"github.com/pedroppinheiro/fwf/yamlconfig"
)

var exporter = HTMLExporter{
	htmlTemplate: "<template>{{.}}</template>",
}

func TestHTMLExporter_ExportVisualization(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name     string
		exporter Exporter
		args     args
		want     string
	}{
		{
			"Should export template correctly",
			exporter,
			args{
				"hello <!--world-->",
			},
			"<template>hello <!--world--></template>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.exporter.ExportVisualization(tt.args.s); got != tt.want {
				t.Errorf("HTMLExporter.ExportVisualization() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTMLExporter_MarkRecordsOnString(t *testing.T) {
	var differentRecords = []yamlconfig.Record{
		{
			Name:  "record A",
			Regex: yamlconfig.MustCreateRegex("^A.*$"),
			Fields: []yamlconfig.Field{
				{
					Name:    "field 1",
					Initial: 1,
					End:     1,
				},
			},
		},
		{
			Name:  "record B",
			Regex: yamlconfig.MustCreateRegex("^B.*$"),
			Fields: []yamlconfig.Field{
				{
					Name:    "",
					Initial: 1,
					End:     1,
				},
			},
		},
	}

	type args struct {
		records []yamlconfig.Record
		s       string
	}
	tests := []struct {
		name     string
		exporter HTMLExporter
		args     args
		want     string
	}{
		{
			"Should correctly mark the fields of the first record",
			GetHTMLExporter(),
			args{differentRecords, "Athequickbrownfoxjumpsoverthelazydog"},
			"<span><div class='tooltip'>A<span class='tooltiptext'>field 1</span></div>thequickbrownfoxjumpsoverthelazydog</span>",
		},
		{
			"Should correctly mark the fields of the second record",
			GetHTMLExporter(),
			args{differentRecords, "Bthequickbrownfoxjumpsoverthelazydog"},
			"<span><div class='tooltip'>B<span class='tooltiptext'></span></div>thequickbrownfoxjumpsoverthelazydog</span>",
		},
		{
			"Should not mark due to not match any record",
			GetHTMLExporter(),
			args{differentRecords, "Cthequickbrownfoxjumpsoverthelazydog"},
			"<span>Cthequickbrownfoxjumpsoverthelazydog</span>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.exporter.MarkRecordsOnString(tt.args.records, tt.args.s); got != tt.want {
				t.Errorf("HTMLExporter.MarkRecordsOnString() = %v, want %v", got, tt.want)
			}
		})
	}
}
