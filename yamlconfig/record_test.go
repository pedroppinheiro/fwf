package yamlconfig

import (
	"reflect"
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

func Test_FindFirstRecordThatMatches(t *testing.T) {

	var differentRecords = []Record{
		{
			Name:  "record A",
			Regex: Regex{"^A.*$", regexp.MustCompile("^A.*$")},
			Fields: []Field{
				{
					Name:    "field 1",
					Initial: 1,
					End:     2,
				},
			},
		},
		{
			Name:  "record B",
			Regex: Regex{"^B.*$", regexp.MustCompile("^B.*$")},
			Fields: []Field{
				{
					Name:    "field 2",
					Initial: 1,
					End:     2,
				},
			},
		},
	}

	var recordsWithSameRegex = []Record{
		{
			Name:  "record A",
			Regex: Regex{"^A.*$", regexp.MustCompile("^A.*$")},
			Fields: []Field{
				{
					Name:    "field 1",
					Initial: 1,
					End:     2,
				},
			},
		},
		{
			Name:  "record B",
			Regex: Regex{"^A.*$", regexp.MustCompile("^A.*$")},
			Fields: []Field{
				{
					Name:    "field 2",
					Initial: 1,
					End:     2,
				},
			},
		},
	}

	type args struct {
		records []Record
		line    string
	}
	tests := []struct {
		name  string
		args  args
		want  Record
		want1 bool
	}{
		{
			name:  "Find first record correctly",
			args:  args{differentRecords, "Athequickbrownfoxjumpsoverthelazydog"},
			want:  differentRecords[0],
			want1: true,
		},
		{
			name:  "Find first record correctly even if it's the last in the array",
			args:  args{differentRecords, "Bthequickbrownfoxjumpsoverthelazydog"},
			want:  differentRecords[1],
			want1: true,
		},
		{
			name:  "Find first record correctly even if all records match line",
			args:  args{recordsWithSameRegex, "Athequickbrownfoxjumpsoverthelazydog"},
			want:  recordsWithSameRegex[0],
			want1: true,
		},
		{
			name:  "Do not find any records",
			args:  args{differentRecords, "Xthequickbrownfoxjumpsoverthelazydog"},
			want:  Record{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := FindFirstRecordThatMatches(tt.args.records, tt.args.line)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindFirstRecordThatMatches() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindFirstRecordThatMatches() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
