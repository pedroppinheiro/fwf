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
	var err error
	if err = unmarshal(&matchedString); err != nil {
		return err
	}

	*regex, err = CreateRegex(matchedString)
	return err
}

// MustCreateRegex creates a compiled regex based on a given string, but panics if anything goes wrong
func MustCreateRegex(s string) (regex Regex) {
	regex, err := CreateRegex(s)
	if err != nil {
		panic(err)
	}

	return regex
}

// CreateRegex creates a compiled regex based on a given string
func CreateRegex(s string) (Regex, error) {
	compiledRegex, err := regexp.Compile(s)
	if err != nil {
		return Regex{}, err
	}

	return Regex{s, compiledRegex}, nil
}

// FindFirstRecordThatMatchesString returns the first record, in a given slice of records, in which its
// regex matches the given line. If a record is found it returs the found record and true.
// if it does not find it returns an empty Record and false
func FindFirstRecordThatMatchesString(records []Record, line string) (Record, bool) {
	for _, record := range records {
		if record.IsMatch(line) {
			return record, true
		}
	}

	return Record{}, false
}
