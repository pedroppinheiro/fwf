package exporter

import "github.com/pedroppinheiro/fwf/configuration"

//Exporter defines the interface for all exporters
type Exporter interface {
	// MarkFieldsOnString will mark all the fields on a given string
	MarkFieldsOnString([]configuration.Field, string) string

	// ExportVisualization will take a given string and may add specific content to aid in the visualizing of the end result
	ExportVisualization(string) string

	// SaveToFile saves a given string to a given path
	SaveToFile(s string, path string) (generatedFilePath string, err error)
}
