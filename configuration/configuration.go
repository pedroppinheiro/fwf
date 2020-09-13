package configuration

import (
	"gopkg.in/yaml.v2"
)

// Configuration is the representation of the records described on a YAML file
type Configuration struct {
	Records []Record
}

// isValid returns true if the given configuration is valid, and else otherwise.
// A valid configuration is a configuration in which all fields, of each record, are valid and
// there is no conflict between them.
func (configuration Configuration) isValid() bool {
	for _, record := range configuration.Records {
		existsConflict, err := existsConflictOnFields(record.Fields)
		if err != nil || existsConflict {
			return false
		}
	}
	return true
}

// ReadConfiguration reads a YAML content and returns the equivalent Configuration struct
func ReadConfiguration(yamlConfiguration []byte) (Configuration, error) {
	configuration := Configuration{}
	err := yaml.UnmarshalStrict(yamlConfiguration, &configuration)
	if err != nil {
		return Configuration{}, err
	}
	return configuration, err
}
