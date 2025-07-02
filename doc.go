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
	root    *IniDoc
	name    string
	lines   []iniLine
	comment string
}

type IniDoc struct {
	lines    []iniLine
	sections []*IniSection
}

func NewDoc() *IniDoc {
	return &IniDoc{
		lines:    []iniLine{},
		sections: []*IniSection{},
	}
}

func NewSection() *IniSection {
	return &IniSection{
		lines: []iniLine{},
	}
}

func (d *IniDoc) createSectionIfNotExist(sectionName string) {
	for _, sect := range d.sections {
		if sect.name == sectionName {
			return
		}
	}

	section := IniSection{
		root:  d,
		name:  sectionName,
		lines: []iniLine{},
	}

	d.sections = append(d.sections, &section)
}

func (d *IniDoc) putSection(section *IniSection) {
	defer func() {
		if section.root != nil && section.root != d {
			// copy over any subsections
			for _, subSection := range section.root.sections {
				added := false
				subSection.name = fmt.Sprintf("%s.%s", section.name, subSection.name)
				for idx, dsection := range d.sections {
					if dsection.name == subSection.name {
						d.sections[idx] = subSection
						added = true
						break
					}
				}
				if !added {
					d.sections = append(d.sections, subSection)
				}
			}
		}
	}()

	for idx, dsection := range d.sections {
		if dsection.name == section.name {
			d.sections[idx] = section
			return
		}
	}
	d.sections = append(d.sections, section)
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
	for _, dsection := range d.sections {
		if dsection.name == sectionName {
			return dsection
		}
	}

	section := IniSection{
		root:  d,
		name:  sectionName,
		lines: []iniLine{},
	}

	d.sections = append(d.sections, &section)

	if strings.Contains(sectionName, ".") {
		segments := strings.Split(sectionName, ".")
		nextSection := ""
		for _, seg := range segments {
			nextSection += seg
			if nextSection != "" {
				d.createSectionIfNotExist(nextSection)
			}
			nextSection += "."
		}
	}

	return &section
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

func (d *IniDoc) SectionNames(includeSubsections ...bool) []string {
	names := make([]string, 0, len(d.sections))
	if len(includeSubsections) > 0 && includeSubsections[0] {
		for idx := range d.sections {
			names = append(names, d.sections[idx].name)
		}
	} else {
		for _, section := range d.sections {
			if !strings.Contains(section.name, ".") {
				names = append(names, section.name)
			}
		}
	}
	return names
}

func (d *IniSection) putSubSection(section *IniSection) {
	section.name = fmt.Sprintf("%s.%s", d.name, section.name)
	// if section.root != nil {
	// 	for _, s := range section.root.sections {
	// 		if s.name != section.name {
	// 			s.name = fmt.Sprintf("%s.%s", section.name, s.name)
	// 		}
	// 	}
	// }
	d.root.putSection(section)
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

func (d *IniSection) Section(sectionName string) *IniSection {
	if d.root == nil {
		d.root = &IniDoc{}
	}

	if d.name != "" {
		return d.root.Section(fmt.Sprintf("%s.%s", d.name, sectionName))
	}

	return d.root.Section(sectionName)
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

func (d *IniSection) SubsectionNames(includeSubsections ...bool) []string {
	allSectionNames := d.root.SectionNames(true)

	result := make([]string, 0, len(allSectionNames))
	if len(includeSubsections) > 0 && includeSubsections[0] {
		for _, sectName := range allSectionNames {
			if strings.HasPrefix(sectName, d.name+".") {
				result = append(result, sectName[len(d.name)+1:])
			}
		}
	} else {
		for _, sectName := range allSectionNames {
			if strings.HasPrefix(sectName, d.name+".") {
				if !strings.Contains(sectName[len(d.name)+1:], ".") {
					result = append(result, sectName[len(d.name)+1:])
				}
			}
		}
	}

	return result
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
	if len(s.lines) == 0 {
		return ""
	}

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
		secStr := section.ToString()
		if secStr != "" {
			v += "\n"
			v += secStr
		}
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
