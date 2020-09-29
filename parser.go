package rac

import (
	"reflect"
	"strings"
)

type object map[string]string

type decodeState struct {
	data   []byte
	fields []object
}

func (s *decodeState) init(data []byte) {

	s.data = data
	s.fields = []object{}

}

func (s *decodeState) parse() {

	parts := strings.Split(string(s.data), "\n")
	//parts = strings.Split(parts, "\n")

	part := ""
	for _, str := range parts {

		part += str + "\n"

		if len(str) == 0 {

			field := s.parseObject(part)

			if len(part) == 0 {
				continue
			}

			s.fields = append(s.fields, field)

			part = ""
		}

	}

}

func (s *decodeState) parseObject(in string) object {

	result := make(object)

	for _, line := range strings.Split(in, "\n") {
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			result[strings.Trim(parts[0], " ")] = strings.Trim(parts[1], " ")
		}
	}

	return result
}

func (s *decodeState) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	s.parse()

	switch rv.Kind() {

	case reflect.Slice, reflect.Array:

	case reflect.Struct:

	}

	return nil
}

func Unmarshal(data []byte, v interface{}) error {

	var d decodeState
	//err := checkValid(data, &d.scan)
	//if err != nil {
	//	return err
	//}

	d.init(data)
	return d.unmarshal(v)
}

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "rac: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "rac: Unmarshal(nil " + e.Type.String() + ")"
}
