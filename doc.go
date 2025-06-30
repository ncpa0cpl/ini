package ini

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type FieldValue struct {
	Key   string
	Value string
}

type iniField struct {
	order   int
	key     string
	value   string
	comment string
}

type iniComment struct {
	order       int
	value       string
	commentType string
}

type IniSection struct {
	name         string
	fields       []iniField
	comments     []iniComment
	comment      string
	sectionOrder int
	nextOrder    int
}

type IniDoc struct {
	fields           []iniField
	comments         []iniComment
	sections         []IniSection
	nextFieldOrder   int
	nextSectionOrder int
}

func NewDoc() *IniDoc {
	return &IniDoc{
		fields:   []iniField{},
		comments: []iniComment{},
		sections: []IniSection{},
	}
}

func (d *IniDoc) getField(key string) *iniField {
	for idx := range d.fields {
		if d.fields[idx].key == key {
			return &d.fields[idx]
		}
	}
	return nil
}

func (d *IniDoc) addField(key, value string) {
	d.fields = append(d.fields, iniField{
		key:   key,
		value: value,
		order: d.nextFieldOrder,
	})
	d.nextFieldOrder++
}

func (d *IniDoc) AddComment(value string) {
	d.comments = append(d.comments, iniComment{
		order:       d.nextFieldOrder,
		value:       value,
		commentType: ";",
	})
	d.nextFieldOrder++
}

func (d *IniDoc) AddHashComment(value string) {
	d.comments = append(d.comments, iniComment{
		order:       d.nextFieldOrder,
		value:       value,
		commentType: "#",
	})
	d.nextFieldOrder++
}

func (d *IniDoc) Del(key string) {
	for idx := range d.fields {
		if d.fields[idx].key == key {
			d.fields = slices.Delete(d.fields, idx, idx+1)
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
		name:         sectionName,
		fields:       []iniField{},
		comments:     []iniComment{},
		sectionOrder: d.nextSectionOrder,
	}
	d.nextSectionOrder++

	d.sections = append(d.sections, section)
	return &d.sections[len(d.sections)-1]
}

func (d *IniDoc) Keys() []string {
	keys := make([]string, len(d.fields))
	for idx := range d.fields {
		keys[idx] = d.fields[idx].key
	}
	return keys
}

func (d *IniDoc) Values() []FieldValue {
	keys := make([]FieldValue, len(d.fields))
	for idx := range d.fields {
		keys[idx] = FieldValue{d.fields[idx].key, d.fields[idx].value}
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

func (d *IniSection) getField(key string) *iniField {
	for idx := range d.fields {
		if d.fields[idx].key == key {
			return &d.fields[idx]
		}
	}
	return nil
}

func (d *IniSection) addField(key, value string) {
	d.fields = append(d.fields, iniField{
		key:   key,
		value: value,
		order: d.nextOrder,
	})
	d.nextOrder++
}

func (d *IniSection) AddComment(value string) {
	d.comments = append(d.comments, iniComment{
		order:       d.nextOrder,
		value:       value,
		commentType: ";",
	})
	d.nextOrder++
}

func (d *IniSection) AddHashComment(value string) {
	d.comments = append(d.comments, iniComment{
		order:       d.nextOrder,
		value:       value,
		commentType: "#",
	})
	d.nextOrder++
}

func (d *IniSection) Del(key string) {
	for idx := range d.fields {
		if d.fields[idx].key == key {
			d.fields = slices.Delete(d.fields, idx, idx+1)
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
	keys := make([]string, len(d.fields))
	for idx := range d.fields {
		keys[idx] = d.fields[idx].key
	}
	return keys
}

func (d *IniSection) Values() []FieldValue {
	keys := make([]FieldValue, len(d.fields))
	for idx := range d.fields {
		keys[idx] = FieldValue{d.fields[idx].key, d.fields[idx].value}
	}
	return keys
}

// serialization

func (f *iniField) ToString() string {
	v := fmt.Sprintf("%s=%s", f.key, f.value)
	if f.comment != "" {
		v += fmt.Sprintf(" ;%s", f.comment)
	}
	v = strings.ReplaceAll(v, "\n", "\\n")
	return v + "\n"
}

func (c *iniComment) ToString() string {
	var v string

	lines := strings.Split(c.value, "\n")
	if c.commentType == ";" {
		for _, line := range lines {
			v = fmt.Sprintf("; %s\n", line)
		}
	} else {
		for _, line := range lines {
			v = fmt.Sprintf("# %s\n", line)
		}
	}

	return v
}

func (s *IniSection) getNextLine(prevOrder int) (line string, orderNo int, done bool) {
	var lowestField *iniField
	var lowestComment *iniComment

	for idx := range s.fields {
		if s.fields[idx].order > prevOrder {
			lowestField = &s.fields[idx]
			break
		}
	}

	for idx := range s.comments {
		if s.comments[idx].order > prevOrder {
			lowestComment = &s.comments[idx]
			break
		}
	}

	if lowestField == nil && lowestComment == nil {
		return "", 0, true
	}

	if lowestComment == nil {
		return lowestField.ToString(), lowestField.order, false
	}

	if lowestField == nil {
		return lowestComment.ToString(), lowestComment.order, false
	}

	if lowestField.order < lowestComment.order {
		return lowestField.ToString(), lowestField.order, false
	} else {
		return lowestComment.ToString(), lowestComment.order, false
	}
}

func (s *IniSection) ToString() string {
	var v string = ""

	if s.comment != "" {
		v += fmt.Sprintf("; %s\n", strings.ReplaceAll(s.comment, "\n", "\\n"))
	}

	v += fmt.Sprintf("[%s]\n", s.name)

	slices.SortFunc(s.fields, func(a, b iniField) int {
		return a.order - b.order
	})
	slices.SortFunc(s.comments, func(a, b iniComment) int {
		return a.order - b.order
	})

	var nextline string
	var prevOrder int = -1
	var done bool
	for !done {
		nextline, prevOrder, done = s.getNextLine(prevOrder)
		if nextline != "" {
			v += nextline
		}
	}

	return v
}

func (d *IniDoc) getNextLine(prevOrder int) (line string, orderNo int, done bool) {
	var lowestField *iniField
	var lowestComment *iniComment

	for idx := range d.fields {
		if d.fields[idx].order > prevOrder {
			lowestField = &d.fields[idx]
			break
		}
	}

	for idx := range d.comments {
		if d.comments[idx].order > prevOrder {
			lowestComment = &d.comments[idx]
			break
		}
	}

	if lowestField == nil && lowestComment == nil {
		return "", 0, true
	}

	if lowestComment == nil {
		return lowestField.ToString(), lowestField.order, false
	}

	if lowestField == nil {
		return lowestComment.ToString(), lowestComment.order, false
	}

	if lowestField.order < lowestComment.order {
		return lowestField.ToString(), lowestField.order, false
	} else {
		return lowestComment.ToString(), lowestComment.order, false
	}
}

func (d *IniDoc) ToString() string {
	var v string = ""

	slices.SortFunc(d.fields, func(a, b iniField) int {
		return a.order - b.order
	})
	slices.SortFunc(d.comments, func(a, b iniComment) int {
		return a.order - b.order
	})

	var nextline string
	var prevOrder int = -1
	var done bool
	for !done {
		nextline, prevOrder, done = d.getNextLine(prevOrder)
		if nextline != "" {
			v += nextline
		}
	}

	for _, section := range d.sections {
		v += "\n"
		v += section.ToString()
	}

	return v
}
