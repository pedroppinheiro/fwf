package main

import (
	"bufio"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/pedroppinheiro/fwf/exporter"
	"github.com/pedroppinheiro/fwf/yamlconfig"
)

var (
	yamlLocation         string
	fileLocation         string
	fileExportedLocation string
)

func init() {
	flag.StringVar(&yamlLocation, "yaml", "", "the full path for the yaml configuration")
	flag.StringVar(&fileLocation, "file", "", "the full path for the file to generate the visualization")
	flag.StringVar(&fileExportedLocation, "o", "./", "the path to where the exported file should be created")
	flag.Parse()
}

func main() {
	if yamlLocation == "" {
		panic("Please provide a valid yaml location with the flag \"-yaml\", use \"fwf -h\" or \"fwf --help\" for help")
	}
	if fileLocation == "" {
		panic("Please provide a valid file location location with the flag \"-file\", use \"fwf -h\" or \"fwf --help\" for help")
	}

	configuration := readConfigurationFromYAML(yamlLocation)
	file := getFile(fileLocation)
	defer file.Close()
	reader := bufio.NewReader(file)

	fileExporter := getCurrentExporter()
	exportedContent := ""
	for {
		line, err := reader.ReadString('\n')

		record, isRecordFound := yamlconfig.FindFirstRecordThatMatches(configuration.Records, line)
		if isRecordFound {
			exportedContent += yamlconfig.ApplyMarkerToFieldsOnString(fileExporter, record.Fields, line)
		} else {
			exportedContent += line
		}

		if err == io.EOF {
			break
		}
	}
	finalExportedContent := fileExporter.ExportVisualization(exportedContent)
	generatedFilePath, err := fileExporter.SaveToFile(finalExportedContent, fileExportedLocation)

	if err == nil {
		log.Printf("File created successfully on %v\n", generatedFilePath)
	} else {
		panic(err)
	}
}

func readConfigurationFromYAML(yamlLocation string) yamlconfig.Configuration {
	yamlContent := readFileContent(yamlLocation)
	configuration, err := yamlconfig.ReadConfiguration(yamlContent)
	if err != nil {
		panic(err)
	}
	return configuration
}

func readFileContent(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}
	return content
}

func getFile(fileLocation string) *os.File {
	file, err := os.Open(fileLocation)
	if err != nil {
		panic(err)
	}
	return file
}

func getCurrentExporter() exporter.Exporter {
	return exporter.GetHTMLExporter()
}
