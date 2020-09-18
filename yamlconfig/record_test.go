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

func Test_FindFirstRecordThatMatchesString(t *testing.T) {

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
			got, got1 := FindFirstRecordThatMatchesString(tt.args.records, tt.args.line)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindFirstRecordThatMatchesString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindFirstRecordThatMatchesString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMustCreateRegex(t *testing.T) {

	compiledRegex := regexp.MustCompile("test")

	type args struct {
		s string
	}
	tests := []struct {
		name        string
		args        args
		wantRegex   Regex
		expectPanic bool
	}{
		{
			"Must create regex",
			args{"test"},
			Regex{"test", compiledRegex},
			false,
		},
		{
			"Must panic due to invalid regex",
			args{"\\"},
			Regex{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Error - was expecting function to panic")
						return
					}
				}()
			}

			if gotRegex := MustCreateRegex(tt.args.s); !reflect.DeepEqual(gotRegex, tt.wantRegex) {
				t.Errorf("MustCreateRegex() = %v, want %v", gotRegex, tt.wantRegex)
			}
		})
	}
}

func TestCreateRegex(t *testing.T) {

	compiledRegex := regexp.MustCompile("test")

	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Regex
		wantErr bool
	}{
		{
			"Create regex successfully",
			args{"test"},
			Regex{"test", compiledRegex},
			false,
		},
		{
			"Creates empty regex when given an invalid regex",
			args{"\\"},
			Regex{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateRegex(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRegex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateRegex() = %v, want %v", got, tt.want)
			}
		})
	}
}
