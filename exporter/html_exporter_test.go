package exporter

import (
	"testing"

	"github.com/pedroppinheiro/fwf/yamlconfig"
)

var exporter = HTMLExporter{
	htmlTemplate: "<template>{{.}}</template>",
	htmlMarker: yamlconfig.Marker{
		ObtainInitialMarker: func(field yamlconfig.Field) string {
			return "<!--"
		},
		ObtainEndMarker: func(field yamlconfig.Field) string {
			return "-->"
		},
	},
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
