package configuration

import (
	"regexp"
	"testing"
)

func TestRecords_IsMatch(t *testing.T) {
	yaml := `
      records:
        - name: "record A"
          regex: "(?m)^A.*$"
          fields:
            - name: "field 1"
              initial: 1
              end: 1
            - name: "field 2"
              initial: 2
              end: 2`

	configuration, err := ReadConfiguration([]byte(yaml))
	if err != nil {
		panic(err)
	}

	type args struct {
		s string
	}
	tests := []struct {
		name    string
		records Record
		args    args
		want    bool
	}{
		{
			"correctly matches a string",
			Record{Regex: Regex{regex: regexp.MustCompile("test")}},
			args{s: "test"},
			true,
		},
		{
			"correctly matches a string",
			Record{Regex: Regex{regex: regexp.MustCompile("^A.*$")}},
			args{s: "Athequickbrownfoxjumpsoverthelazydog"},
			true,
		},
		{
			"correctly matches a string after parsing a yaml into a configuration", //integration test
			configuration.Records[0],
			args{s: "Athequickbrownfoxjumpsoverthelazydog"},
			true,
		},
		{
			"string is not matched",
			Record{Regex: Regex{regex: regexp.MustCompile("test")}},
			args{s: "findme"},
			false,
		},
		{
			"string is matched when no regex was given",
			Record{},
			args{s: "findme"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.records.IsMatch(tt.args.s); got != tt.want {
				t.Errorf("Records.IsMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
