package configuration

import (
	"regexp"

	"gopkg.in/yaml.v2"
)

// Configuration holds the data of the yaml configuration
type Configuration struct {
	Records []Records
}

// Records holds the data of the Records
type Records struct {
	Name   string
	Regex  Regex
	Fields []Fields
}

// MatchString reports whether the string s contains any match of the regular expression pattern.
func (records Records) MatchString(s string) bool {
	return records.Regex.regex.MatchString(s)
}

// Regex stores regex
type Regex struct {
	regexString string
	regex       *regexp.Regexp
}

// UnmarshalYAML interface is implemented to give a custom behaviour when marshalling the yaml to the "Regex" field.
// See https://godoc.org/gopkg.in/yaml.v2#Unmarshaler for more details
func (r *Regex) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	reg, err := regexp.Compile(s)
	if err != nil {
		r = &Regex{}
		return err
	}

	r.regexString = s
	r.regex = reg
	return nil
}

// Fields holds the data of the Fields
type Fields struct {
	Name    string
	Initial int
	End     int
}

// ReadConfiguration reads a YAML and returns the equivalent Configuration struct
func ReadConfiguration(yamlConfiguration []byte) (Configuration, error) {
	configuration := Configuration{}
	err := yaml.UnmarshalStrict(yamlConfiguration, &configuration)
	if err != nil {
		return Configuration{}, err
	}
	return configuration, err
}
