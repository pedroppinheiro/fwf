package exporter

import (
	"github.com/pedroppinheiro/fwf/yamlconfig"
)

//Exporter defines the interface for all exporters
type Exporter interface {
	// ExportVisualization will take a given string and may add specific content to aid in the visualizing of the end result
	ExportVisualization(string) string

	// SaveToFile saves a given string to a given path
	SaveToFile(s string, path string) (generatedFilePath string, err error)

	MarkRecordsOnString(records []yamlconfig.Record, s string) string

	yamlconfig.Marker
}
