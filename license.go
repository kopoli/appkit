package appkit

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func licenseFormat(pkg, name, text string, format string, a ...interface{}) string {
	s := format
	pkg = strings.ReplaceAll(pkg, "%", "%%")
	name = strings.ReplaceAll(name, "%", "%%")
	text = strings.ReplaceAll(text, "%", "%%")
	s = strings.ReplaceAll(s, "%Lp", pkg)
	s = strings.ReplaceAll(s, "%Ln", name)
	s = strings.ReplaceAll(s, "%Lt", text)

	return fmt.Sprintf(s, a...)
}

// LicenseString returns a string of the licenses printed out in the given
// format.
//
// The licenses is expected to be of the following kind of type:
//
// map[string]struct {
//   Name string
//   Text string
// }
//
// The map key is the package name, Name is the license name, and Text is the
// license text.
//
// The format string is similar to fmt's but with the following additional
// verbs:
// %Lp    Package name
// %Ln    License name
// %Lt    License text
//
// The licenses are printed in the alphabetical order of the package name.
//
// As an implementation detail, https://github.com/kopoli/licrep generates a
// compatible map.
func LicenseFormatString(licenses interface{}, format string, a ...interface{}) (string, error) {
	invalidTypeErr := fmt.Errorf("invalid license map argument")

	licmap := reflect.ValueOf(licenses)
	if licmap.Kind() != reflect.Map {
		return "", invalidTypeErr
	}

	keys := licmap.MapKeys()
	names := []string{}

	if len(keys) == 0 {
		return "", nil
	}

	if keys[0].Kind() != reflect.String {
		return "", invalidTypeErr
	}

	for i := range keys {
		names = append(names, keys[i].String())
	}
	sort.Strings(names)

	sb := strings.Builder{}

	for _, pkg := range names {
		v := licmap.MapIndex(reflect.ValueOf(pkg))
		if v.Kind() != reflect.Struct {
			return "", invalidTypeErr
		}

		name := v.FieldByName("Name")
		if name.Kind() != reflect.String {
			return "", invalidTypeErr
		}

		text := v.FieldByName("Text")
		if text.Kind() != reflect.String {
			return "", invalidTypeErr
		}

		s := licenseFormat(pkg, name.String(), text.String(), format, a...)
		_, err := sb.WriteString(s)
		if err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}

func LicenseString(licenses interface{}) (string, error) {
	return LicenseFormatString(licenses, "* %Lp: %Ln\n\n%Lt\n\n")
}
