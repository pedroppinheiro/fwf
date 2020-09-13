package exporter

import (
	"fmt"
	"io/ioutil"

	"../configuration"
)

var (
	htmlTemplate = `
	<html>
		<head></head>
		<body>
			<pre>%v</pre>
		</body>
	</html>`

	htmlMarker = configuration.Marker{
		InitialMarker: "<span class='tooltip'>",
		EndMarker:     "</span>",
	}
)

// HTMLExporter is an implementation of the Exporter interface,
// in which is reponsible to mark and export a string to its html visualization
type HTMLExporter struct {
	htmlTemplate    string
	htmlMarker      configuration.Marker
	defaultFileName string
}

// GetHTMLExporter returns the initialized HTMLExporter with its custom template and marker
func GetHTMLExporter() Exporter {
	return HTMLExporter{htmlTemplate, htmlMarker, "index.html"}
}

// MarkFieldsOnString will mark all the fields on a given string using a custom html marker
func (exporter HTMLExporter) MarkFieldsOnString(fields []configuration.Field, s string) string {
	return configuration.ApplyMarkerToFieldsOnString(exporter.htmlMarker, fields, s)
}

// ExportVisualization will take a given string and will use it on a HTML template to make it better to visualize the end result on a browser
func (exporter HTMLExporter) ExportVisualization(s string) string {
	return fmt.Sprintf(exporter.htmlTemplate, s)
}

// SaveToFile saves a given string to a given path
func (exporter HTMLExporter) SaveToFile(s string, path string) (generatedFilePath string, err error) {
	generatedFilePath = path + exporter.defaultFileName
	err = ioutil.WriteFile(generatedFilePath, []byte(s), 0777)
	return
}
