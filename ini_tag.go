package ini

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Unmarshal(data []byte, v interface{}) error {
	if data == nil {
		return nil
	}

	mp := New().Load(data).Marshal2Map()

	bindTag("ini", v, mp)
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

// Bind binds the content of data into the struct s
func bindTag(tagName string, s interface{}, data map[string]interface{}) interface{} {
	if s == nil {
		return nil
	}

	t := reflect.TypeOf(s)
	tk := t.Kind()

	if tk != reflect.Ptr {
		return nil
	}

	t = t.Elem()
	tk = t.Kind()

	if tk != reflect.Struct {
		return nil
	}

	v := reflect.ValueOf(s).Elem()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fi := parseField(tagName, f)

		fv := v.FieldByName(fi.Name)
		fk := fv.Kind()
		ft := fv.Type()

		if v, ok := data[fi.Alias]; ok {
			vt := reflect.TypeOf(v)
			vk := vt.Kind()

			if !vt.AssignableTo(ft) {

				change_fun := func(k reflect.Kind, v interface{}) (val reflect.Value, err error) {

					err = errors.New("unsupper kind")

					var str_v string
					var ok bool
					if str_v, ok = v.(string); !ok {
						return
					}

					// str_v = v.(string)

					if k == reflect.Int {
						rv, er := strconv.Atoi(str_v)
						if er == nil {
							val = reflect.ValueOf(rv)
							err = nil
						}
					} else if k == reflect.Int8 {
						rv, er := strconv.ParseInt(str_v, 10, 8)
						if er == nil {
							val = reflect.ValueOf(rv)
							err = nil
						}
					} else if k == reflect.Int16 {
						rv, er := strconv.ParseInt(str_v, 10, 16)
						if er == nil {
							val = reflect.ValueOf(rv)
							err = nil
						}
					} else if k == reflect.Int32 {
						rv, er := strconv.ParseInt(str_v, 10, 32)
						if er == nil {
							val = reflect.ValueOf(rv)
							err = nil
						}
					} else if k == reflect.Int64 {
						rv, er := strconv.ParseInt(str_v, 10, 64)
						if er == nil {
							val = reflect.ValueOf(rv)
							err = nil
						}
					} else if k == reflect.Float32 {
						rv, er := strconv.ParseFloat(str_v, 32)
						if er == nil {
							val = reflect.ValueOf(float32(rv))
							err = nil
						}
					} else if k == reflect.Float64 {
						rv, er := strconv.ParseFloat(str_v, 64)
						if er == nil {
							val = reflect.ValueOf(rv)
							err = nil
						}
					}

					return val, err
				}

				val, err := change_fun(fk, v)

				if err == nil {
					fv.Set(val)
					continue
				}
			} else {
				fv.Set(reflect.ValueOf(v))
				continue
			}

			if fk == reflect.Struct && vk == reflect.Map {

				if fv.CanInterface() {
					bindTag(tagName, fv.Addr().Interface(), v.(map[string]interface{}))
					continue
				}
			}
		}
	}

	return s
}

func bindSlice(tagName string, s reflect.Value, data []interface{}) {
	// sk := s.Kind()
	et := s.Type().Elem()
	ek := et.Kind()

	ret := reflect.MakeSlice(et, s.Len(), s.Cap())
	vet := reflect.TypeOf(data).Elem()
	if vet.AssignableTo(et) {
		for i := 0; i < s.Len(); i++ {
			ret.Index(i).Set(reflect.ValueOf(data[i]))
		}
	} else if ek == reflect.Struct {
		for i := 0; i < s.Len(); i++ {
			// v := Bind(tagName, ret.Index(i).Addr().Interface(), data[i].(map[string]interface{}))
			v := bindTag(tagName, ret.Index(i).Addr().Interface(), data[i].(map[string]interface{}))
			ret.Index(i).Set(reflect.ValueOf(v))
		}
	}
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
