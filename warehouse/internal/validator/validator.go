package validator

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

type validationError struct {
	Field      string `json:"field"`
	Validation string `json:"validation"`
}

func Validate(s any) error {
	var err []validationError

	r_value := reflect.ValueOf(s)
	r_type := reflect.TypeOf(s)

	if r_value.Kind() == reflect.Ptr {
		r_value = r_value.Elem()
		r_type = r_type.Elem()
	}

	for i := 0; i < r_value.NumField(); i++ {
		field := r_type.Field(i)
		value := r_value.Field(i).Interface()
		tag := field.Tag.Get("validator")

		if tag == "" {
			continue
		}

		for rule := range strings.SplitSeq(tag, ",") {
			switch rule {
			case "required":
				if isEmpty(value) {
					err = append(err, validationError{
						Field: field.Name,
						Validation: "is required",
					})
				}
			}
		}
	}

	if len(err) == 0 {
		return nil
	}

	jsonBytes, _ := json.Marshal(err)
	return errors.New(string(jsonBytes))
}

func isEmpty(v any) bool {
	switch val := v.(type) {
	case string:
		return strings.TrimSpace(val) == ""
	case int:
		return val == 0
	default:
		return reflect.ValueOf(v).IsZero()
	}
}
