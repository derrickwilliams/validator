package main

import (
	"fmt"
	"reflect"
	"strings"
)

func Register(name string, validator ValidatorFilter) {
	registerValidator(name, validator)
}

func walkValidateArray(context *ValidatorContext, normalizedValue *NormalizedValue, errors *Errors) {
	valueType := reflect.ValueOf(normalizedValue.Value)
	for i := 0; i < valueType.Len(); i++ {
		walkValidate(context, valueType.Index(i).Interface(), errors)
	}
}

func walkValidateMap(context *ValidatorContext, normalizedValue *NormalizedValue, errors *Errors) {
	valueType := reflect.ValueOf(normalizedValue.Value)
	for _, key := range valueType.MapKeys() {
		walkValidate(context, valueType.MapIndex(key).Interface(), errors)
	}
}

func walkValidateStruct(context *ValidatorContext, normalizedValue *NormalizedValue, errors *Errors) {
	for _, field := range getFields(normalizedValue.Value, "validate") {
		normalizedFieldValue, err := normalizeValue(field.Value, false)

		if err != nil {
			errors.Add(NewValidatorError(field.Name, "MISSING_TAG_NAME", strings.Replace(err.Error(), "{field}", field.Name, 1)))
			break
		}

		context.SetField(field)
		context.SetValue(normalizedFieldValue)

		for _, tag := range field.Tags {
			if validate, err := getValidator(tag.Name); err == nil {
				if err = validate(context, tag.Options); err != nil {
					errors.Add(NewValidatorError(field.Name, tag.Name, strings.Replace(err.Error(), "{field}", field.Name, 1)))
				}
				if context.StopValidate {
					break
				}
			} else {
				errors.Add(NewValidatorError(field.Name, tag.Name, fmt.Sprintf("Validator '%s' used on field '%s' does not exist.", tag.Name, field.Name)))
				break
			}
		}

		walkValidate(context, context.Value, errors)
	}
}

func walkValidate(context *ValidatorContext, value interface{}, errors *Errors) {
	var normalizedValue *NormalizedValue

	if typedValue, ok := value.(*NormalizedValue); ok {
		normalizedValue = typedValue
	} else {
		normalizedValue, _ = normalizeValue(value, false)
	}

	switch normalizedValue.OriginalKind {
	case reflect.Array, reflect.Slice:
		walkValidateArray(context, normalizedValue, errors)
	case reflect.Map:
		walkValidateMap(context, normalizedValue, errors)
	case reflect.Struct:
		context.SetParent(normalizedValue.Value)
		walkValidateStruct(context, normalizedValue, errors)
	}
}

func Validate(value interface{}) *Errors {
	errors := NewErrors()
	context := NewValidatorContext()

	walkValidate(context, value, errors)

	return errors
}
