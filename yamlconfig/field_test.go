package yamlconfig

import (
	"reflect"
	"testing"
)

func TestField_isValid(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		want  bool
	}{
		{
			name:  "Field is invalid because is using 0 as initial",
			field: Field{Initial: 0, End: 1},
			want:  false,
		},
		{
			name:  "Field is invalid because is using 0 as final",
			field: Field{Initial: 1, End: 0},
			want:  false,
		},
		{
			name:  "Field should be invalid because both positions are zero",
			field: Field{},
			want:  false,
		},
		{
			name:  "Field should be invalid because both positions are zero",
			field: Field{Initial: 0, End: 0},
			want:  false,
		},
		{
			name:  "Field should be valid",
			field: Field{Initial: 1, End: 1},
			want:  true,
		},
		{
			name:  "Field should be valid",
			field: Field{Initial: 1, End: 2},
			want:  true,
		},
		{
			name:  "Field should be invalid because end is negative",
			field: Field{Initial: 1, End: -1},
			want:  false,
		},
		{
			name:  "Field should be invalid because initial is negative",
			field: Field{Initial: -1, End: 1},
			want:  false,
		},
		{
			name:  "Field should be invalid because both initial and end are negative",
			field: Field{Initial: -2, End: -1},
			want:  false,
		},
		{
			name:  "Field should be invalid because initial is higher then end",
			field: Field{Initial: 3, End: 2},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.isValid(); got != tt.want {
				t.Errorf("Field.isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortFieldsByInitialPositionAsc(t *testing.T) {
	type args struct {
		fields []Field
	}
	tests := []struct {
		name string
		args args
		want []Field
	}{
		{
			name: "Empty slice should continue be empty",
			args: args{[]Field{}},
			want: []Field{},
		},
		{
			name: "Slice with one Field should remain the same",
			args: args{[]Field{
				{"", 5, 5},
			}},
			want: []Field{{"", 5, 5}},
		},
		{
			name: "Slice sorted by Initial desc should be sorted by Initial asc",
			args: args{[]Field{
				{"", 5, 5},
				{"", 1, 1},
			}},
			want: []Field{{"", 1, 1}, {"", 5, 5}},
		},
		{
			name: "Slice sorted by Initial asc should remain the same",
			args: args{[]Field{
				{"", 1, 1},
				{"", 5, 5},
			}},
			want: []Field{{"", 1, 1}, {"", 5, 5}},
		},
		{
			name: "Unsorted slice should be sorted",
			args: args{[]Field{
				{"", 5, 5},
				{"", 1, 1},
				{"", 10, 10},
			}},
			want: []Field{{"", 1, 1}, {"", 5, 5}, {"", 10, 10}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortFieldsByInitialPositionAsc(tt.args.fields)

			if !reflect.DeepEqual(tt.args.fields, tt.want) {
				t.Errorf("sortFieldsByInitialPositionAsc() = %v, want %v", tt.args.fields, tt.want)
			}
		})
	}
}

func Test_existsConflict(t *testing.T) {
	type args struct {
		field1 Field
		field2 Field
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Should not detect any conflict",
			args: args{field1: Field{"", 1, 2}, field2: Field{"", 3, 4}},
			want: false,
		},
		{
			name: "Should not detect any conflict",
			args: args{field1: Field{"", 3, 4}, field2: Field{"", 1, 2}},
			want: false,
		},
		{
			name: "Should detect conflict - field2's initial is the same as field1's end",
			args: args{field1: Field{"", 1, 2}, field2: Field{"", 2, 3}},
			want: true,
		},
		{
			name: "Should detect conflict - field2's positions are inside field1's",
			args: args{field1: Field{"", 1, 5}, field2: Field{"", 2, 3}},
			want: true,
		},
		{
			name: "Should detect conflict - field1's positions are inside field2's",
			args: args{field1: Field{"", 2, 3}, field2: Field{"", 1, 5}},
			want: true,
		},
		{
			name:    "Should give error due to invalid field",
			args:    args{field1: Field{"", 0, 1}, field2: Field{"", 2, 5}},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := existsConflict(tt.args.field1, tt.args.field2)

			if (err != nil) != tt.wantErr {
				t.Errorf("existsConflict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("existsConflict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_existsConflictOnFields(t *testing.T) {
	var fieldsWithConflicts = []Field{
		{"", 1, 1},
		{"", 2, 3},
		{"", 4, 5},
		{"", 5, 6},
	}

	var unsortedfieldsWithConflicts = []Field{
		{"", 4, 5},
		{"", 2, 3},
		{"", 5, 6},
		{"", 1, 1},
	}

	type args struct {
		fields []Field
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Should not detect any conflict on empty slice",
			args: args{fieldsWithConflicts[0:0]},
			want: false,
		},
		{
			name: "Should not detect any conflict one just one element",
			args: args{fieldsWithConflicts[0:1]},
			want: false,
		},
		{
			name: "Should not detect any conflict",
			args: args{fieldsWithConflicts[0:2]},
			want: false,
		},
		{
			name: "Should not detect any conflict",
			args: args{fieldsWithConflicts[0:3]},
			want: false,
		},
		{
			name: "Should detect conflict",
			args: args{fieldsWithConflicts[0:4]},
			want: true,
		},
		{
			name: "Should detect conflict",
			args: args{unsortedfieldsWithConflicts[0:4]},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := existsConflictOnFields(tt.args.fields)

			if (err != nil) != tt.wantErr {
				t.Errorf("existsConflictOnFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("existsConflictOnFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStringBeforeField(t *testing.T) {

	type args struct {
		s     string
		field Field
	}
	tests := []struct {
		name        string
		args        args
		want        string
		expectPanic bool
	}{
		{
			name: "String before field should be empty string",
			args: args{s: "", field: Field{"", 1, 2}},
			want: "",
		},
		{
			name: "String before field should be empty string",
			args: args{s: "thequickbrownfoxjumpsoverthelazydog", field: Field{"", 1, 35}},
			want: "",
		},
		{
			name:        "Should panic due to invalid field",
			args:        args{s: "", field: Field{"", 0, 2}},
			expectPanic: true,
		},
		{
			name: "Should get the string before the field",
			args: args{s: "thequickbrownfoxjumpsoverthelazydog", field: Field{"", 2, 2}},
			want: "t",
		},
		{
			name: "Should get the string before the field",
			args: args{s: "thequickbrownfoxjumpsoverthelazydog", field: Field{"", 4, 4}},
			want: "the",
		},
		{
			name: "Should get the string before the field",
			args: args{s: "thequickbrownfoxjumpsoverthelazydog", field: Field{"", 4, 100}},
			want: "the",
		},
		{
			name: "Should correctly get the string with accents before the field",
			args: args{s: "ÇÇÇÇÇuickbrownfoxjumpsoverthelazydog", field: Field{"", 2, 5}},
			want: "Ç",
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

			if got := getStringBeforeField(tt.args.s, tt.args.field); got != tt.want {
				t.Errorf("getStringBeforeField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStringOfField(t *testing.T) {
	type args struct {
		s     string
		field Field
	}
	tests := []struct {
		name        string
		args        args
		want        string
		expectPanic bool
	}{
		{
			name: "String of field should be empty string",
			args: args{s: "", field: Field{"", 1, 2}},
			want: "",
		},
		{
			name:        "Should panic due to invalid field",
			args:        args{s: "", field: Field{"", 0, 2}},
			expectPanic: true,
		},
		{
			name: "Should get the string of the field",
			args: args{s: "thequickbrownfoxjumpsoverthelazydog", field: Field{"", 1, 1}},
			want: "t",
		},
		{
			name: "Should get the string of the field",
			args: args{s: "thequickbrownfoxjumpsoverthelazydog", field: Field{"", 2, 2}},
			want: "h",
		},
		{
			name: "Should get the string of the field",
			args: args{s: "thequickbrownfoxjumpsoverthelazydog", field: Field{"", 4, 8}},
			want: "quick",
		},
		{
			name: "Should get the string of the field",
			args: args{s: "thequickbrownfoxjumpsoverthelazydog", field: Field{"", 4, 100}},
			want: "quickbrownfoxjumpsoverthelazydog",
		},
		{
			name: "Should get the string of the field",
			args: args{s: "thequickbrownfoxjumpsoverthelazydog", field: Field{"", 1, 35}},
			want: "thequickbrownfoxjumpsoverthelazydog",
		},
		{
			name: "Should correctly get the string with accents of field",
			args: args{s: "ÇÇÇÇÇuickbrownfoxjumpsoverthelazydog", field: Field{"", 2, 5}},
			want: "ÇÇÇÇ",
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

			if got := getStringOfField(tt.args.s, tt.args.field); got != tt.want {
				t.Errorf("getStringOfField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStringAfterField(t *testing.T) {
	type args struct {
		s     string
		field Field
	}
	tests := []struct {
		name        string
		args        args
		want        string
		expectPanic bool
	}{
		{
			name: "String after field should be empty string",
			args: args{s: "", field: Field{"", 1, 2}},
			want: "",
		},
		{
			name: "String after field should be empty string",
			args: args{s: "thequickbrownfoxjumpsoverthelazydog", field: Field{"", 1, 35}},
			want: "",
		},
		{
			name:        "Should panic due to invalid field",
			args:        args{s: "", field: Field{"", 0, 2}},
			expectPanic: true,
		},
		{
			name: "Should get the string after the field",
			args: args{s: "thequickbrownfoxjumpsoverthelazydog", field: Field{"", 1, 1}},
			want: "hequickbrownfoxjumpsoverthelazydog",
		},
		{
			name: "Should get the string after the field",
			args: args{s: "thequickbrownfoxjumpsoverthelazydog", field: Field{"", 2, 2}},
			want: "equickbrownfoxjumpsoverthelazydog",
		},
		{
			name: "Should get the string after the field",
			args: args{s: "thequickbrownfoxjumpsoverthelazydog", field: Field{"", 4, 8}},
			want: "brownfoxjumpsoverthelazydog",
		},
		{
			name: "Should get the string after the field",
			args: args{s: "thequickbrownfoxjumpsoverthelazydog", field: Field{"", 4, 100}},
			want: "",
		},
		{
			name: "Should correctly get the string with accents after field",
			args: args{s: "ÇÇÇÇÇuickbrownfoxjumpsoverthelazydog", field: Field{"", 2, 5}},
			want: "uickbrownfoxjumpsoverthelazydog",
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

			if got := getStringAfterField(tt.args.s, tt.args.field); got != tt.want {
				t.Errorf("getStringAfterField() = %v, want %v", got, tt.want)
			}
		})
	}
}

type marker struct{}

func (m marker) ObtainInitialMarker(field Field) string {
	return "<"
}
func (m marker) ObtainEndMarker(field Field) string {
	return ">"
}

func Test_ApplyMarkerToFieldsOnString(t *testing.T) {

	customMarker := marker{}
	type args struct {
		marker Marker
		fields []Field
		s      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should mark the string correctly with the fields",
			args: args{
				marker: customMarker,
				fields: []Field{
					{"", 4, 8},
					{"", 17, 21},
				},
				s: "thequickbrownfoxjumpsoverthelazydog",
			},
			want: "the<quick>brownfox<jumps>overthelazydog",
		},
		{
			name: "Should mark the string correctly with the fields",
			args: args{
				marker: customMarker,
				fields: []Field{
					{"", 9, 13},
					{"", 4, 8},
					{"", 17, 21},
					{"", 14, 16},
				},
				s: "thequickbrownfoxjumpsoverthelazydog",
			},
			want: "the<quick><brown><fox><jumps>overthelazydog",
		},
		{
			name: "Should mark the string correctly with the fields",
			args: args{
				marker: customMarker,
				fields: []Field{
					{"", 1, 35},
				},
				s: "thequickbrownfoxjumpsoverthelazydog",
			},
			want: "<thequickbrownfoxjumpsoverthelazydog>",
		},
		{
			name: "Should mark the string correctly with the fields",
			args: args{
				marker: customMarker,
				fields: []Field{
					{"", 1, 1},
					{"", 35, 35},
				},
				s: "thequickbrownfoxjumpsoverthelazydog",
			},
			want: "<t>hequickbrownfoxjumpsoverthelazydo<g>",
		},
		{
			name: "Should mark the string correctly with the fields",
			args: args{
				marker: customMarker,
				fields: []Field{
					{"", 34, 100},
				},
				s: "thequickbrownfoxjumpsoverthelazydog",
			},
			want: "thequickbrownfoxjumpsoverthelazyd<og>",
		},
		{
			name: "Should mark the string correctly with accents",
			args: args{
				marker: customMarker,
				fields: []Field{
					{"", 2, 5},
				},
				s: "ÇÇÇÇÇuickbrownfoxjumpsoverthelazydog",
			},
			want: "Ç<ÇÇÇÇ>uickbrownfoxjumpsoverthelazydog",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApplyMarkerToFieldsOnString(tt.args.marker, tt.args.fields, tt.args.s); got != tt.want {
				t.Errorf("ApplyMarkerToFieldsOnString() = %v, want %v", got, tt.want)
			}
		})
	}
}
