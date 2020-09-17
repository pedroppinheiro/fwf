package yamlconfig

import "regexp"

// Record holds the data of the Records
type Record struct {
	Name   string
	Regex  Regex
	Fields []Field
}

// IsMatch reports whether the string s contains any match of the regular expression pattern
func (record Record) IsMatch(s string) bool {
	if record.Regex.regex == nil && record.Regex.regexString == "" {
		return true
	}
	return record.Regex.regex.MatchString(s)
}

// Regex holds regexp information
type Regex struct {
	regexString string
	regex       *regexp.Regexp
}

// UnmarshalYAML interface is implemented to give a custom behaviour when marshalling the yaml to the "Regex" field.
// See https://godoc.org/gopkg.in/yaml.v2#Unmarshaler for more details
func (regex *Regex) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var matchedString string
	if err := unmarshal(&matchedString); err != nil {
		return err
	}

	compiledRegex, err := regexp.Compile(matchedString)
	if err != nil {
		*regex = Regex{}
		return err
	}

	*regex = Regex{matchedString, compiledRegex}
	return nil
}

// FindFirstRecordThatMatches returns the first record, in a given slice of records, in which its
// regex matches the given line. If a record is found it returs the found record and true.
// if it does not find it returns an empty Record and false
func FindFirstRecordThatMatches(records []Record, line string) (Record, bool) {
	for _, record := range records {
		if record.IsMatch(line) {
			return record, true
		}
	}

	return Record{}, false
}
