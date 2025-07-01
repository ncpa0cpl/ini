package ini

import (
	"fmt"
	"reflect"
	"strings"
)

type Marshalable interface {
	MarshalINI() (DocOrSection, error)
}

type Unmarshalable interface {
	UnmarshalINI(DocOrSection) error
}

func Unmarshal(data string, v interface{}) error {
	if v == nil {
		return fmt.Errorf("given struct is nil")
	}

	doc := Parse(data)

	vUnmarshalable, ok := v.(Unmarshalable)
	if ok {
		err := vUnmarshalable.UnmarshalINI(doc)
		return err
	}

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
		fkind := f.Type.Kind()

		switch fkind {
		case reflect.Bool:
			strvalue := doc.Get(fieldInfo.Alias)
			if strvalue == "true" {
				vElem.FieldByName(fieldInfo.Name).SetBool(true)
			}
		case reflect.String:
			value := doc.Get(fieldInfo.Alias)
			vElem.FieldByName(fieldInfo.Name).SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value, err := doc.GetInt(fieldInfo.Alias)
			if err != nil {
				return err
			}
			vElem.FieldByName(fieldInfo.Name).SetInt(value)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value, err := doc.GetUint(fieldInfo.Alias)
			if err != nil {
				return err
			}
			vElem.FieldByName(fieldInfo.Name).SetUint(value)
		case reflect.Float32, reflect.Float64:
			value, err := doc.GetFloat(fieldInfo.Alias)
			if err != nil {
				return err
			}
			vElem.FieldByName(fieldInfo.Name).SetFloat(value)
		case reflect.Struct, reflect.Ptr:
			fieldVal := vElem.FieldByName(fieldInfo.Name)

			sectElem := fieldVal
			sectType := f.Type
			sectKind := sectType.Kind()

			if fieldVal.IsZero() && sectKind == reflect.Ptr {
				fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
			}

			docSection := doc.Section(fieldInfo.Alias)

			vUnmarshalable, ok := fieldVal.Interface().(Unmarshalable)
			if ok {
				err := vUnmarshalable.UnmarshalINI(docSection)
				if err != nil {
					return err
				}
				continue
			}

			if sectKind == reflect.Ptr {
				sectType = sectType.Elem()
				sectKind = sectType.Kind()
				sectElem = fieldVal.Elem()
			}

			if sectKind != reflect.Struct {
				continue
			}

			sectFields := reflect.VisibleFields(sectType)
			for _, sectF := range sectFields {
				secFieldInfo := parseField("ini", sectF)
				switch sectF.Type.Kind() {
				case reflect.Bool:
					strvalue := docSection.Get(secFieldInfo.Alias)
					if strvalue == "true" {
						sectElem.FieldByName(secFieldInfo.Name).SetBool(true)
					}
				case reflect.String:
					value := docSection.Get(secFieldInfo.Alias)
					sectElem.FieldByName(secFieldInfo.Name).SetString(value)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					value, err := docSection.GetInt(secFieldInfo.Alias)
					if err != nil {
						return err
					}
					sectElem.FieldByName(secFieldInfo.Name).SetInt(value)
				case reflect.Float32, reflect.Float64:
					value, err := docSection.GetFloat(secFieldInfo.Alias)
					if err != nil {
						return err
					}
					sectElem.FieldByName(secFieldInfo.Name).SetFloat(value)
				}
			}
		case reflect.Map:
			fieldVal := vElem.FieldByName(fieldInfo.Name)

			keyKind := fieldVal.Type().Key().Kind()
			if keyKind != reflect.String {
				continue
			}

			docSection := doc.Section(fieldInfo.Alias)
			docKeys := docSection.Keys()

			if fieldVal.IsZero() {
				fieldVal.Set(reflect.MakeMapWithSize(
					reflect.MapOf(
						reflect.TypeOf(""),
						f.Type.Elem(),
					),
					len(docKeys),
				))
			}

			mapElemType := f.Type.Elem()
			switch mapElemType.Kind() {
			case reflect.String:
				for _, key := range docKeys {
					value := docSection.Get(key)
					fieldVal.SetMapIndex(
						reflect.ValueOf(key),
						reflect.ValueOf(value),
					)
				}
			case reflect.Bool:
				for _, key := range docKeys {
					value, err := docSection.GetBool(key)
					if err != nil {
						return err
					}
					fieldVal.SetMapIndex(
						reflect.ValueOf(key),
						reflect.ValueOf(value),
					)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				for _, key := range docKeys {
					value, err := docSection.GetInt(key)
					if err != nil {
						return err
					}
					fieldVal.SetMapIndex(
						reflect.ValueOf(key),
						reflect.ValueOf(value).Convert(mapElemType),
					)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				for _, key := range docKeys {
					value, err := docSection.GetUint(key)
					if err != nil {
						return err
					}
					fieldVal.SetMapIndex(
						reflect.ValueOf(key),
						reflect.ValueOf(value).Convert(mapElemType),
					)
				}
			case reflect.Float32, reflect.Float64:
				for _, key := range docKeys {
					value, err := docSection.GetFloat(key)
					if err != nil {
						return err
					}
					fieldVal.SetMapIndex(
						reflect.ValueOf(key),
						reflect.ValueOf(value).Convert(mapElemType),
					)
				}
			case reflect.Interface:
				for _, key := range docKeys {
					var value any = docSection.Get(key)
					fieldVal.SetMapIndex(
						reflect.ValueOf(key),
						reflect.ValueOf(value),
					)
				}
			}
		}
	}

	return nil
}

func Marshal(v any) (string, error) {
	if v == nil {
		return "", fmt.Errorf("given struct is nil")
	}

	vUnmarshalable, ok := v.(Marshalable)
	if ok {
		doc, err := vUnmarshalable.MarshalINI()
		return doc.ToString(), err
	}

	doc := NewDoc()

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
		case reflect.Bool:
			value := vElem.FieldByName(fieldInfo.Name).Bool()
			doc.SetBool(fieldInfo.Alias, value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value := vElem.FieldByName(fieldInfo.Name).Int()
			doc.SetInt(fieldInfo.Alias, value)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value := vElem.FieldByName(fieldInfo.Name).Uint()
			doc.SetUint(fieldInfo.Alias, value)
		case reflect.Float32, reflect.Float64:
			value := vElem.FieldByName(fieldInfo.Name).Float()
			doc.SetFloat(fieldInfo.Alias, value)
		case reflect.Struct, reflect.Ptr:
			fieldVal := vElem.FieldByName(fieldInfo.Name)
			sectElem := fieldVal
			sectType := f.Type
			sectKind := sectType.Kind()

			if fieldVal.IsZero() && sectKind == reflect.Ptr {
				continue
			}

			vUnmarshalable, ok := fieldVal.Interface().(Marshalable)
			if ok {
				secOrDoc, err := vUnmarshalable.MarshalINI()
				if err != nil {
					return "", err
				}

				section, ok := secOrDoc.(*IniSection)
				if ok {
					if section.name == "" {
						section.name = fieldInfo.Alias
					}
					doc.putSection(section)
					continue
				}

				asDoc, ok := secOrDoc.(*IniDoc)
				if ok {
					section := docToSection(asDoc)
					section.name = fieldInfo.Alias
					doc.putSection(section)
					continue
				}
			}

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
				case reflect.Bool:
					value := sectElem.FieldByName(secFieldInfo.Name).Bool()
					docSection.SetBool(secFieldInfo.Alias, value)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					value := sectElem.FieldByName(secFieldInfo.Name).Int()
					docSection.SetInt(secFieldInfo.Alias, value)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					value := sectElem.FieldByName(secFieldInfo.Name).Uint()
					docSection.SetUint(secFieldInfo.Alias, value)
				case reflect.Float32, reflect.Float64:
					value := sectElem.FieldByName(secFieldInfo.Name).Float()
					docSection.SetFloat(secFieldInfo.Alias, value)
				}
			}
		case reflect.Map:
			fieldVal := vElem.FieldByName(fieldInfo.Name)

			if fieldVal.IsZero() {
				continue
			}

			keyKind := fieldVal.Type().Key().Kind()
			if keyKind != reflect.String {
				continue
			}

			docSection := doc.Section(fieldInfo.Alias)

			for _, key := range fieldVal.MapKeys() {
				value := fieldVal.MapIndex(key)
				valueKind := value.Kind()

				if valueKind == reflect.Interface {
					value = value.Elem()
					valueKind = value.Kind()
				}

				keyV := key.String()
				switch valueKind {
				case reflect.String:
					value := value.String()
					docSection.Set(keyV, value)
				case reflect.Bool:
					value := value.Bool()
					docSection.SetBool(keyV, value)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					value := value.Int()
					docSection.SetInt(keyV, value)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					value := value.Uint()
					docSection.SetUint(keyV, value)
				case reflect.Float32, reflect.Float64:
					value := value.Float()
					docSection.SetFloat(keyV, value)
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
