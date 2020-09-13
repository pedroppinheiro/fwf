package exporter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"text/template"

	"../configuration"
)

var (
	htmlTemplate = `
		<!DOCTYPE html>
		<html>
			<style>
				.tooltip {
					position: relative;
					display: inline-block;
					/*border: 2px dotted black*/
					box-shadow: 0 0 5px rgb(0,100,0,1);
				}

				.tooltip .tooltiptext {
					visibility: hidden;
					background-color: #555;
					color: #fff;
					text-align: center;
					border-radius: 6px;
					position: fixed;
					z-index: 1;
					bottom: 0;
					left: 0;
					right: 0;
					opacity: 0;
					transition: opacity 0.3s;    
					font-size: 1.5em;
				}
				
				.tooltip:hover {
					box-shadow: 0 0 5px rgb(100,0,0,1);
					/*box-shadow: 0 0 5px rgba(0,0,0,0.5);*/
				}

				.tooltip:hover .tooltiptext {
					visibility: visible;
					opacity: 1;
				}

				pre {
					font-family: 'Courier New', Courier, monospace;
					height: 100%;
					padding: 4px;
					counter-reset: line;
					overflow-x: hidden;
					overflow-y: hidden;
				}
			</style>
			<body><pre>{{.}}</pre>
			</body>
		</html>`

	htmlMarker = configuration.Marker{
		ObtainInitialMarker: func(field configuration.Field) string {
			return "<div class='tooltip'>"
		},
		ObtainEndMarker: func(field configuration.Field) string {
			return fmt.Sprintf("<span class='tooltiptext'>%v</span></div>", field.Name)
		},
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
	t, err := template.New("customTemplate").Parse(exporter.htmlTemplate)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, s)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())
	return buf.String()
}

// SaveToFile saves a given string to a given path
func (exporter HTMLExporter) SaveToFile(s string, path string) (generatedFilePath string, err error) {
	generatedFilePath = path + exporter.defaultFileName
	err = ioutil.WriteFile(generatedFilePath, []byte(s), 0777)
	return
}