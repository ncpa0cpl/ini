package ini

import (
	"fmt"
	"reflect"
	"strings"
)

func Unmarshal(data []byte, v interface{}) error {
	if data == nil {
		return nil
	}
	if v == nil {
		return fmt.Errorf("given struct is nil")
	}

	doc := New().Load(data)

	vType := reflect.TypeOf(v)
	vKind := vType.Kind()
	vElem := reflect.ValueOf(v)
	if vKind == reflect.Ptr {
		vType = vType.Elem()
		vKind = vType.Kind()
		vElem = vElem.Elem()
	}

	if vKind != reflect.Struct {
		return fmt.Errorf("given value is not a struct")
	}

	vFields := reflect.VisibleFields(vType)
	for _, f := range vFields {
		fieldInfo := parseField("ini", f)
		switch f.Type.Kind() {
		case reflect.String:
			value := doc.Get(fieldInfo.Alias)
			vElem.FieldByName(fieldInfo.Name).SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value := doc.GetInt64(fieldInfo.Alias)
			vElem.FieldByName(fieldInfo.Name).SetInt(value)
		case reflect.Float32, reflect.Float64:
			value := doc.GetFloat64(fieldInfo.Alias)
			vElem.FieldByName(fieldInfo.Name).SetFloat(value)
		case reflect.Struct, reflect.Ptr:
			fieldVal := vElem.FieldByName(fieldInfo.Name)
			sectElem := fieldVal
			sectType := f.Type
			sectKind := sectType.Kind()
			if sectKind == reflect.Ptr {
				sectType = sectType.Elem()
				sectKind = sectType.Kind()
				sectElem = fieldVal.Elem()
			}

			if sectKind != reflect.Struct {
				continue
			}

			docSection := doc.Section(fieldInfo.Alias)

			sectFields := reflect.VisibleFields(sectType)
			for _, sectF := range sectFields {
				secFieldInfo := parseField("ini", sectF)
				switch sectF.Type.Kind() {
				case reflect.String:
					value := docSection.Get(secFieldInfo.Alias)
					sectElem.FieldByName(secFieldInfo.Name).SetString(value)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					value := docSection.GetInt64(secFieldInfo.Alias)
					sectElem.FieldByName(secFieldInfo.Name).SetInt(value)
				case reflect.Float32, reflect.Float64:
					value := docSection.GetFloat64(secFieldInfo.Alias)
					sectElem.FieldByName(secFieldInfo.Name).SetFloat(value)
				}
			}
		}
	}

	return nil
}

func Marshal(v interface{}) (string, error) {
	if v == nil {
		return "", fmt.Errorf("given struct is nil")
	}

	doc := New()

	vType := reflect.TypeOf(v)
	vKind := vType.Kind()
	vElem := reflect.ValueOf(v)
	if vKind == reflect.Ptr {
		vType = vType.Elem()
		vKind = vType.Kind()
		vElem = vElem.Elem()
	}

	if vKind != reflect.Struct {
		return "", fmt.Errorf("given value is not a struct")
	}

	vFields := reflect.VisibleFields(vType)
	for _, f := range vFields {
		fieldInfo := parseField("ini", f)
		switch f.Type.Kind() {
		case reflect.String:
			value := vElem.FieldByName(fieldInfo.Name).String()
			doc.Set(fieldInfo.Alias, value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value := vElem.FieldByName(fieldInfo.Name).Int()
			doc.Set(fieldInfo.Alias, value)
		case reflect.Float32, reflect.Float64:
			value := vElem.FieldByName(fieldInfo.Name).Float()
			doc.Set(fieldInfo.Alias, value)
		case reflect.Struct, reflect.Ptr:
			fieldVal := vElem.FieldByName(fieldInfo.Name)
			sectElem := fieldVal
			sectType := f.Type
			sectKind := sectType.Kind()
			if sectKind == reflect.Ptr {
				sectType = sectType.Elem()
				sectKind = sectType.Kind()
				sectElem = fieldVal.Elem()
			}

			if sectKind != reflect.Struct {
				continue
			}

			docSection := doc.Section(fieldInfo.Alias)

			sectFields := reflect.VisibleFields(sectType)
			for _, sectF := range sectFields {
				secFieldInfo := parseField("ini", sectF)
				switch sectF.Type.Kind() {
				case reflect.String:
					value := sectElem.FieldByName(secFieldInfo.Name).String()
					docSection.Set(secFieldInfo.Alias, value)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					value := sectElem.FieldByName(secFieldInfo.Name).Int()
					docSection.Set(secFieldInfo.Alias, value)
				case reflect.Float32, reflect.Float64:
					value := sectElem.FieldByName(secFieldInfo.Name).Float()
					docSection.Set(secFieldInfo.Alias, value)
				}
			}
		}
	}

	return doc.ToString(), nil
}

type fieldInfo struct {
	Alias string

	Name string
}

// ParseField parses [FieldInfo] for the given struct field [f] from struct tag with name [tagName]
func parseField(tagName string, f reflect.StructField) *fieldInfo {
	var parts []string
	alias := f.Name

	tag, tagOk := f.Tag.Lookup(tagName)
	if tagOk {
		partsTemp := strings.Split(tag, ",")
		parts = make([]string, 0, len(partsTemp))
		for i := 0; i < len(partsTemp); i++ {
			part := strings.TrimSpace(partsTemp[i])
			if len(part) != 0 {
				parts = append(parts, part)
			}
		}
	}

	if len(parts) != 0 {
		alias = parts[0]
		// TODO parse other tags
	}

	return &fieldInfo{
		Alias: alias,
		Name:  f.Name,
	}
}
