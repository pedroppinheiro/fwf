package yamlconfig

import (
	"fmt"
	"log"
	"sort"
)

// Field holds the data of the a field on a record.
type Field struct {
	Name    string
	Initial int
	End     int
}

// Marker needs to be implemented in order to get the initial and end marker. These markers are placed before and after a string (field)
type Marker interface {
	ObtainInitialMarker(field Field) string

	ObtainEndMarker(field Field) string
}

// isValid returns true if the field is valid, false otherwise. A valid field is one where
// all of its positions (initial and end) are positive and initial cannot be 0
func (field Field) isValid() bool {
	if field.Initial >= 1 && field.Initial <= field.End {
		return true
	}
	return false
}

// sortFieldsByInitialPositionAsc sorts a slice of fields in ascending order based on its initial position value
func sortFieldsByInitialPositionAsc(fields []Field) {
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Initial < fields[j].Initial
	})
}

// existsConflict returns true if there is a conflict between fields of a given slice of Field, false otherwise
func existsConflictOnFields(fields []Field) (bool, error) {
	sortFieldsByInitialPositionAsc(fields)
	for i := 0; i < len(fields); i++ {

		if i != 0 {
			existsConflict, err := existsConflict(fields[i-1], fields[i])
			if err != nil || existsConflict {
				return existsConflict, err
			}
		}
	}
	return false, nil
}

// existsConflict returns true if there is a conflict between two given fields. A conflict is when
// an initial and end position exists in the same range of the positions of another field
func existsConflict(field1 Field, field2 Field) (bool, error) {
	if !field1.isValid() || !field2.isValid() {
		return false, fmt.Errorf("existsConflict(): error - cannot check conflicts on invalid fields: field1(%v) field2(%v)", field1, field2)
	}

	if field1.Initial > field2.Initial {
		field1, field2 = field2, field1
	}

	isField1InitialOutOfField2Range := field1.Initial < field2.Initial && field1.Initial < field2.End
	isField1EndOutOfField2Range := field1.End < field2.Initial && field1.End < field2.End

	return !isField1InitialOutOfField2Range || !isField1EndOutOfField2Range, nil
}

// getStringBeforeField returns the string that exists before a given field.
// For instance, given the string "thequickbrownfox", and a field with initial 4 and end 8
// the resulting string will be "the"
func getStringBeforeField(s string, field Field) string {
	if !field.isValid() {
		panic("Error - the given field is invalid")
	}

	if s == "" || field.Initial == 1 {
		return ""
	}

	runes := []rune(s)
	return string(runes[0 : field.Initial-1])
}

// getStringBeforeField returns the string of a given field.
// For instance, given the string "thequickbrownfox", and a field with initial 4 and end 8
// the resulting string will be "quick"
func getStringOfField(s string, field Field) string {
	if !field.isValid() {
		panic("Error - the given field is invalid")
	}

	if s == "" {
		return ""
	}

	var end int
	if field.End > len(s) {
		log.Printf("Warning getStringOfField() - Field \"%v\" with end %v is higher then the length of the given string \"%v\"", field.Name, field.End, s)
		end = len(s)
	} else {
		end = field.End
	}

	runes := []rune(s)
	return string(runes[field.Initial-1 : end])
}

// getStringBeforeField returns the string after a given field.
// For instance, given the string "thequickbrownfox", and a field with initial 4 and end 8
// the resulting string will be "brownfox"
func getStringAfterField(s string, field Field) string {
	if !field.isValid() {
		panic("Error - the given field is invalid")
	}

	if field.End > len(s) {
		log.Printf("Warning getStringAfterField() - Field \"%v\" with end %v is higher then the length of the given string \"%v\"", field.Name, field.End, s)
		return ""
	}

	if s == "" || field.End == len(s) {
		return ""
	}

	runes := []rune(s)
	return string(runes[field.End:])
}

// ApplyMarkerToFieldsOnString returns a string that is the result of applying a field marker to the fields on a string
// For instance, given a marker "<" and ">", and given the string "thequickbrownfox" with a field with initial 4 and end 8,
// the resulting string will be "the<quick>brownfox"
func ApplyMarkerToFieldsOnString(marker Marker, fields []Field, s string) string {
	sortFieldsByInitialPositionAsc(fields)

	var (
		finalString          string
		tempString           string
		lastFieldEndPosition int
	)

	for i, field := range fields {
		tempString = getStringBeforeField(s, field)
		tempString += marker.ObtainInitialMarker(field)
		tempString += getStringOfField(s, field)
		tempString += marker.ObtainEndMarker(field)

		if i != 0 {
			lastFieldEndPosition = fields[i-1].End
			tempString = tempString[lastFieldEndPosition:]
		}

		finalString += tempString
	}

	lastField := fields[len(fields)-1]
	finalString += getStringAfterField(s, lastField)
	return finalString
}
