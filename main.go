package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"

	"./configuration"
	"./exporter"
)

func main() {
	var (
		yamlLocation         = "./test_yaml.yaml"
		fileLocation         = "./test_conteudo.txt"
		exporter             = getCurrentExporter()
		fileExportedLocation = "./"
	)

	configuration := readConfigurationFromYAML(yamlLocation)
	file := getFile(fileLocation)
	defer file.Close()
	reader := bufio.NewReader(file)

	exportedContent := ""
	for {
		line, err := reader.ReadString('\n')

		for _, record := range configuration.Records {
			if record.IsMatch(line) {
				exportedContent += exporter.MarkFieldsOnString(record.Fields, line)
			} else {
				exportedContent += line
			}
		}

		if err == io.EOF {
			break
		}
	}
	finalExportedContent := exporter.ExportVisualization(exportedContent)
	exporter.SaveToFile(finalExportedContent, fileExportedLocation)
}

func readConfigurationFromYAML(yamlLocation string) configuration.Configuration {
	yamlContent := readFileContent(yamlLocation)
	configuration, err := configuration.ReadConfiguration(yamlContent)
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

func saveToFile(s string, path string) {

}
