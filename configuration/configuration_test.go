package configuration

import (
	"reflect"
	"regexp"
	"testing"
)

func Test_ReadConfiguration(t *testing.T) {
	//Test 1 - full yaml configuration
	var yamlString1 string = `
        records:
         - name: "record A"
           regex: "^A.*$"
           fields:
            - name: "field 1"
              initial: 1
              end: 2
            - name: "field 2"
              initial: 3
              end: 4
         - name: "record B"
           regex: "^B.*$"
           fields:
            - name: "field 3"
              initial: 5
              end: 6
            - name: "field 4"
              initial: 7
              end: 8`
	var expectedConfiguration1 Configuration = Configuration{
		Records: []Record{
			{
				Name:  "record A",
				Regex: Regex{"^A.*$", regexp.MustCompile("^A.*$")},
				Fields: []Field{
					{
						Name:    "field 1",
						Initial: 1,
						End:     2,
					},
					{
						Name:    "field 2",
						Initial: 3,
						End:     4,
					},
				},
			},
			{
				Name:  "record B",
				Regex: Regex{"^B.*$", regexp.MustCompile("^B.*$")},
				Fields: []Field{
					{
						Name:    "field 3",
						Initial: 5,
						End:     6,
					},
					{
						Name:    "field 4",
						Initial: 7,
						End:     8,
					},
				},
			},
		},
	}

	//Test 2 - few fields are set
	var yamlString2 string = `
        records:
         - name: "record A"
           fields:
            - name: "field 1"
            - name: "field 2"`

	//Test 3 -  fields are incorrectly named
	var yamlString3 string = `
        recordsWithChangedName:
         - nameIsChanged: "record A"
           fieldsIsChanged:
            - nameIsChanged: "field 1"
			- nameIsChanged: "field 2"`

	//Test 4 - unparsable regex
	var yamlString4 string = `
        records:
         - name: "record A"
           regex: "\\"
           fields:
            - name: "field 1"
            - name: "field 2"`

	type args struct {
		yamlConfiguration []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Configuration
		wantErr bool
	}{
		{
			"construct yaml with every configuration's field",
			args{yamlConfiguration: []byte(yamlString1)},
			expectedConfiguration1,
			false,
		},
		{
			"construct yaml with only some of the configuration's fields",
			args{yamlConfiguration: []byte(yamlString2)},
			Configuration{},
			true,
		},
		{
			"get error when trying to construct unproper yaml",
			args{yamlConfiguration: []byte("unparsed yaml")},
			Configuration{},
			true,
		},
		{
			"get empty configuration given an empty string",
			args{yamlConfiguration: []byte("")},
			Configuration{},
			false,
		},
		{
			"get error due to yaml not having the same configuration's field names",
			args{yamlConfiguration: []byte(yamlString3)},
			Configuration{},
			true,
		},
		{
			"get error due to yaml having an unparsable regex",
			args{yamlConfiguration: []byte(yamlString4)},
			Configuration{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfiguration(tt.args.yamlConfiguration)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfiguration_isValid(t *testing.T) {
	var validConfiguration Configuration = Configuration{
		Records: []Record{
			{
				Name:  "record A",
				Regex: Regex{"^A.*$", regexp.MustCompile("^A.*$")},
				Fields: []Field{
					{
						Name:    "field 1",
						Initial: 1,
						End:     2,
					},
					{
						Name:    "field 2",
						Initial: 3,
						End:     4,
					},
				},
			},
			{
				Name:  "record B",
				Regex: Regex{"^B.*$", regexp.MustCompile("^B.*$")},
				Fields: []Field{
					{
						Name:    "field 3",
						Initial: 5,
						End:     6,
					},
					{
						Name:    "field 4",
						Initial: 7,
						End:     8,
					},
				},
			},
		},
	}

	var configurationWithInvalidField Configuration = Configuration{
		Records: []Record{
			{
				Name:  "record A",
				Regex: Regex{"^A.*$", regexp.MustCompile("^A.*$")},
				Fields: []Field{
					{
						Name:    "field 1",
						Initial: 1,
						End:     2,
					},
					{
						Name:    "field 2",
						Initial: 3,
						End:     4,
					},
				},
			},
			{
				Name:  "record B",
				Regex: Regex{"^B.*$", regexp.MustCompile("^B.*$")},
				Fields: []Field{
					{
						Name:    "field 3",
						Initial: 5,
						End:     6,
					},
					{
						Name:    "field 4",
						Initial: 0,
						End:     0,
					},
				},
			},
		},
	}

	var configurationWithConflictOnFields Configuration = Configuration{
		Records: []Record{
			{
				Name:  "record A",
				Regex: Regex{"^A.*$", regexp.MustCompile("^A.*$")},
				Fields: []Field{
					{
						Name:    "field 1",
						Initial: 1,
						End:     2,
					},
					{
						Name:    "field 2",
						Initial: 3,
						End:     4,
					},
				},
			},
			{
				Name:  "record B",
				Regex: Regex{"^B.*$", regexp.MustCompile("^B.*$")},
				Fields: []Field{
					{
						Name:    "field 3",
						Initial: 5,
						End:     6,
					},
					{
						Name:    "field 4",
						Initial: 6,
						End:     7,
					},
				},
			},
		},
	}

	tests := []struct {
		name          string
		configuration Configuration
		want          bool
		wantErr       bool
	}{
		{
			name:          "Should be a valid configuration",
			configuration: validConfiguration,
			want:          true,
			wantErr:       false,
		},
		{
			name:          "Should be an invalid configuration due to invalid field",
			configuration: configurationWithInvalidField,
			want:          false,
			wantErr:       true,
		},
		{
			name:          "Should be an invalid configuration due to conflicts on fields",
			configuration: configurationWithConflictOnFields,
			want:          false,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := tt.configuration.isValid(); got != tt.want || ((err != nil) != tt.wantErr) {
				t.Errorf("Configuration.isValid() = %v, want %v; err = %v, wantedError = %v", got, tt.want, err, tt.wantErr)
			}
		})
	}
}
