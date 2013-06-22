package gomill

import (
	"net/url"
	"reflect"
	"strconv"
)

// Converts all fields in a struct to a url.Values string map, as long as the
// field has a json tag set.
func structToMap(i interface{}) (values url.Values) {
	values = url.Values{}
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		// If we can't use this field - as set by the tag - skip it
		tagName, tagOpts := parseTag(typ.Field(i).Tag.Get("json"))
		if tagName == "" || tagName == "-" {
			continue
		}
		// Convert each type into a string for the url.Values string map
		var v string
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			v = f.String()
		}
		if tagOpts.Contains("omitempty") && v == "" {
			continue
		}
		values.Set(tagName, v)
	}
	return
}
