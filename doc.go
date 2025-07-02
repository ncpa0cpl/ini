package ini

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	lineTypeKv = iota
	lineTypeComment
	lineTypeHashComment
)

type FieldValue struct {
	Key   string
	Value string
}

type iniLine struct {
	lineType int
	key      string
	value    string
	comment  string
}

type IniSection struct {
	name    string
	lines   []iniLine
	comment string
}

type IniDoc struct {
	lines    []iniLine
	sections []IniSection
}

func NewDoc() *IniDoc {
	return &IniDoc{
		lines:    []iniLine{},
		sections: []IniSection{},
	}
}

func NewSection() *IniSection {
	return &IniSection{
		lines: []iniLine{},
	}
}

func (d *IniDoc) putSection(section *IniSection) {
	for idx := range d.sections {
		if d.sections[idx].name == section.name {
			d.sections[idx] = *section
			return
		}
	}

	d.sections = append(d.sections, *section)
}

func (d *IniDoc) getField(key string) *iniLine {
	for idx := range d.lines {
		if d.lines[idx].lineType == lineTypeKv && d.lines[idx].key == key {
			return &d.lines[idx]
		}
	}
	return nil
}

func (d *IniDoc) addField(key, value string) {
	d.lines = append(d.lines, iniLine{
		lineType: lineTypeKv,
		key:      key,
		value:    value,
	})
}

func (d *IniDoc) AddComment(value string) {
	d.lines = append(d.lines, iniLine{
		lineType: lineTypeComment,
		value:    value,
	})
}

func (d *IniDoc) AddHashComment(value string) {
	d.lines = append(d.lines, iniLine{
		lineType: lineTypeHashComment,
		value:    value,
	})
}

func (d *IniDoc) Del(key string) {
	for idx := range d.lines {
		if d.lines[idx].lineType == lineTypeKv && d.lines[idx].key == key {
			d.lines = slices.Delete(d.lines, idx, idx+1)
			return
		}
	}
}

func (d *IniDoc) SetFieldComment(fieldKey string, value string) {
	f := d.getField(fieldKey)
	if f != nil {
		f.comment = value
	}
}

func (d *IniDoc) Set(key, value string) {
	value = strings.Trim(value, " ")
	if isKeyValid(key) {
		f := d.getField(key)
		if f == nil {
			d.addField(key, value)
		} else {
			f.value = value
		}
	}
}

func (d *IniDoc) SetInt(key string, value int64) {
	strVal := strconv.FormatInt(value, 10)
	d.Set(key, strVal)
}

func (d *IniDoc) SetUint(key string, value uint64) {
	strVal := strconv.FormatUint(value, 10)
	d.Set(key, strVal)
}

func (d *IniDoc) SetFloat(key string, value float64) {
	strVal := strconv.FormatFloat(value, 'f', -1, 64)
	d.Set(key, strVal)
}

func (d *IniDoc) SetBool(key string, value bool) {
	strVal := strconv.FormatBool(value)
	d.Set(key, strVal)
}

func (d *IniDoc) GetComment(key string) string {
	f := d.getField(key)
	if f == nil {
		return ""
	} else {
		return f.comment
	}
}

func (d *IniDoc) Get(key string) string {
	f := d.getField(key)
	if f == nil {
		return ""
	} else {
		return f.value
	}
}

func (d *IniDoc) GetInt(key string) (int64, error) {
	v := d.Get(key)
	if v == "" {
		return 0, nil
	}
	return strconv.ParseInt(v, 10, 64)
}

func (d *IniDoc) GetUint(key string) (uint64, error) {
	v := d.Get(key)
	if v == "" {
		return 0, nil
	}
	return strconv.ParseUint(v, 10, 64)
}

func (d *IniDoc) GetFloat(key string) (float64, error) {
	v := d.Get(key)
	if v == "" {
		return 0, nil
	}
	return strconv.ParseFloat(v, 64)
}

func (d *IniDoc) GetBool(key string) (bool, error) {
	v := d.Get(key)
	if v == "" {
		return false, nil
	}
	return strconv.ParseBool(v)
}

func (d *IniDoc) Section(sectionName string) *IniSection {
	for idx := range d.sections {
		if d.sections[idx].name == sectionName {
			return &d.sections[idx]
		}
	}

	section := IniSection{
		name:  sectionName,
		lines: []iniLine{},
	}

	d.sections = append(d.sections, section)
	return &d.sections[len(d.sections)-1]
}

func (d *IniDoc) Keys() []string {
	keys := make([]string, len(d.lines))
	for idx := range d.lines {
		keys[idx] = d.lines[idx].key
	}
	return keys
}

func (d *IniDoc) Values() []FieldValue {
	keys := make([]FieldValue, len(d.lines))
	for idx := range d.lines {
		keys[idx] = FieldValue{d.lines[idx].key, d.lines[idx].value}
	}
	return keys
}

func (d *IniDoc) SectionNames() []string {
	names := make([]string, len(d.sections))
	for idx := range d.sections {
		names[idx] = d.sections[idx].name
	}
	return names
}

func (d *IniSection) getField(key string) *iniLine {
	for idx := range d.lines {
		if d.lines[idx].lineType == lineTypeKv && d.lines[idx].key == key {
			return &d.lines[idx]
		}
	}
	return nil
}

func (d *IniSection) addField(key, value string) {
	d.lines = append(d.lines, iniLine{
		lineType: lineTypeKv,
		key:      key,
		value:    value,
	})
}

func (d *IniSection) AddComment(value string) {
	d.lines = append(d.lines, iniLine{
		lineType: lineTypeComment,
		value:    value,
	})
}

func (d *IniSection) AddHashComment(value string) {
	d.lines = append(d.lines, iniLine{
		lineType: lineTypeHashComment,
		value:    value,
	})
}

func (d *IniSection) Del(key string) {
	for idx := range d.lines {
		if d.lines[idx].lineType == lineTypeKv && d.lines[idx].key == key {
			d.lines = slices.Delete(d.lines, idx, idx+1)
			return
		}
	}
}

func (d *IniSection) SetFieldComment(fieldKey string, value string) {
	f := d.getField(fieldKey)
	if f != nil {
		f.comment = value
	}
}

func (d *IniSection) Set(key, value string) {
	value = strings.Trim(value, " ")
	if isKeyValid(key) {
		f := d.getField(key)
		if f == nil {
			d.addField(key, value)
		} else {
			f.value = value
		}
	}
}

func (d *IniSection) SetInt(key string, value int64) {
	strVal := strconv.FormatInt(value, 10)
	d.Set(key, strVal)
}

func (d *IniSection) SetUint(key string, value uint64) {
	strVal := strconv.FormatUint(value, 10)
	d.Set(key, strVal)
}

func (d *IniSection) SetFloat(key string, value float64) {
	strVal := strconv.FormatFloat(value, 'f', -1, 64)
	d.Set(key, strVal)
}

func (d *IniSection) SetBool(key string, value bool) {
	strVal := strconv.FormatBool(value)
	d.Set(key, strVal)
}

func (d *IniSection) GetComment(key string) string {
	f := d.getField(key)
	if f == nil {
		return ""
	} else {
		return f.comment
	}
}

func (d *IniSection) Get(key string) string {
	f := d.getField(key)
	if f == nil {
		return ""
	} else {
		return f.value
	}
}

func (d *IniSection) GetInt(key string) (int64, error) {
	v := d.Get(key)
	if v == "" {
		return 0, nil
	}
	return strconv.ParseInt(v, 10, 64)
}

func (d *IniSection) GetUint(key string) (uint64, error) {
	v := d.Get(key)
	if v == "" {
		return 0, nil
	}
	return strconv.ParseUint(v, 10, 64)
}

func (d *IniSection) GetFloat(key string) (float64, error) {
	v := d.Get(key)
	if v == "" {
		return 0, nil
	}
	return strconv.ParseFloat(v, 64)
}

func (d *IniSection) GetBool(key string) (bool, error) {
	v := d.Get(key)
	if v == "" {
		return false, nil
	}
	return strconv.ParseBool(v)
}

func (d *IniSection) Keys() []string {
	keys := make([]string, len(d.lines))
	for idx := range d.lines {
		keys[idx] = d.lines[idx].key
	}
	return keys
}

func (d *IniSection) Values() []FieldValue {
	keys := make([]FieldValue, len(d.lines))
	for idx := range d.lines {
		keys[idx] = FieldValue{d.lines[idx].key, d.lines[idx].value}
	}
	return keys
}

func (d *IniSection) SetName(name string) {
	d.name = name
}

// serialization

func escapeIniValue(value string) string {
	escapedV := make([]rune, 0, len(value)+8)

	for _, char := range value {
		switch char {
		case ';', '#':
			escapedV = append(escapedV, '\\', char)
		case '\n':
			escapedV = append(escapedV, '\\', 'N')
		default:
			escapedV = append(escapedV, char)
		}
	}

	return string(escapedV)
}

func (f *iniLine) ToString() string {
	switch f.lineType {
	case lineTypeKv:
		v := fmt.Sprintf("%s=%s", f.key, escapeIniValue(f.value))
		if f.comment != "" {
			v += fmt.Sprintf(" ;%s", f.comment)
		}
		return v + "\n"
	case lineTypeComment:
		var v string
		lines := strings.Split(f.value, "\n")
		for _, line := range lines {
			v = fmt.Sprintf("; %s\n", line)
		}
		return v
	case lineTypeHashComment:
		var v string
		lines := strings.Split(f.value, "\n")
		for _, line := range lines {
			v = fmt.Sprintf("# %s\n", line)
		}
		return v
	}

	panic("invalid line type: " + strconv.FormatInt(int64(f.lineType), 10))
}

func (s *IniSection) ToString() string {
	var v string = ""

	if s.comment != "" {
		v += fmt.Sprintf("; %s\n", strings.ReplaceAll(s.comment, "\n", "\\n"))
	}

	v += fmt.Sprintf("[%s]\n", s.name)

	for _, line := range s.lines {
		v += line.ToString()
	}

	return v
}

func (d *IniDoc) ToString() string {
	var v string = ""

	for _, line := range d.lines {
		v += line.ToString()
	}

	for _, section := range d.sections {
		v += "\n"
		v += section.ToString()
	}

	return v
}

func docToSection(doc *IniDoc) *IniSection {
	sec := IniSection{
		lines: make([]iniLine, len(doc.lines)),
	}
	copy(sec.lines, doc.lines)
	return &sec
}
