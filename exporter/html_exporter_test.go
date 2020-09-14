package exporter

import (
	"testing"

	"github.com/pedroppinheiro/fwf/configuration"
)

var exporter = HTMLExporter{
	htmlTemplate: "<template>{{.}}</template>",
	htmlMarker: configuration.Marker{
		ObtainInitialMarker: func(field configuration.Field) string {
			return "<!--"
		},
		ObtainEndMarker: func(field configuration.Field) string {
			return "-->"
		},
	},
}

func TestHTMLExporter_MarkFieldsOnString(t *testing.T) {
	type args struct {
		fields []configuration.Field
		s      string
	}
	tests := []struct {
		name     string
		exporter Exporter
		args     args
		want     string
	}{
		{
			"Should mark fields correctly",
			exporter,
			args{
				[]configuration.Field{{Name: "", Initial: 7, End: 11}},
				"hello world",
			},
			"hello <!--world-->",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.exporter.MarkFieldsOnString(tt.args.fields, tt.args.s); got != tt.want {
				t.Errorf("HTMLExporter.MarkFieldsOnString() = %v, want %v", got, tt.want)
			}
		})
	}
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
