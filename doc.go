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
	lineTypeWhiteLine
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
		lines:    make([]iniLine, 0, 16),
		sections: make([]*IniSection, 0, 16),
	}
}

func NewSection() *IniSection {
	return &IniSection{
		lines: make([]iniLine, 0, 16),
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

func (d *IniDoc) lastLine() *iniLine {
	if len(d.lines) == 0 {
		return nil
	}
	return &d.lines[len(d.lines)-1]
}

func (d *IniDoc) addParsedSection(name string) *IniSection {
	comment := ""

	if len(d.sections) > 0 {
		lastSection := d.sections[len(d.sections)-1]
		lastLine := lastSection.lastLine()
		if lastLine != nil && isCommentLine(lastLine) {
			lastSection.lines = lastSection.lines[:len(lastSection.lines)-1]
			comment = lastLine.value
		}
	} else {
		lastLine := d.lastLine()
		if lastLine != nil && isCommentLine(lastLine) {
			d.lines = d.lines[:len(d.lines)-1]
			comment = lastLine.value
		}
	}

	s := d.Section(name)
	s.comment = comment

	return s
}

// Adds an empty line
func (d *IniDoc) AddWhiteLine() {
	d.lines = append(d.lines, iniLine{
		lineType: lineTypeWhiteLine,
	})
}

// Adds a comment line
func (d *IniDoc) AddComment(value string) {
	lastLine := d.lastLine()
	if lastLine != nil && lastLine.lineType == lineTypeComment {
		lastLine.value += "\n" + value
		return
	}

	d.lines = append(d.lines, iniLine{
		lineType: lineTypeComment,
		value:    value,
	})
}

// Adds a comment line, but with a `#` instead of the default `;`
func (d *IniDoc) AddHashComment(value string) {
	lastLine := d.lastLine()
	if lastLine != nil && lastLine.lineType == lineTypeHashComment {
		lastLine.value += "\n" + value
		return
	}

	d.lines = append(d.lines, iniLine{
		lineType: lineTypeHashComment,
		value:    value,
	})
}

// Remove the key-value pair from the document root
func (d *IniDoc) Del(key string) {
	for idx := range d.lines {
		if d.lines[idx].lineType == lineTypeKv && d.lines[idx].key == key {
			d.lines = slices.Delete(d.lines, idx, idx+1)
			return
		}
	}
}

// Adds a comment after a key-value pair, comment will be on the same line as the property (e.x. `key=value ; comment`)
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

// Returns the current comment that's associated with the given property key
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

// Retrieves the given section, if that section does not exist it will be added
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

// Returns a list of all keys of the top level key-value pairs
func (d *IniDoc) Keys() []string {
	keys := make([]string, 0, len(d.lines))
	for idx := range d.lines {
		if d.lines[idx].lineType == lineTypeKv {
			keys = append(keys, d.lines[idx].key)
		}
	}
	return keys
}

// Returns a list of all top level key-value pairs
func (d *IniDoc) Values() []FieldValue {
	keys := make([]FieldValue, 0, len(d.lines))
	for idx := range d.lines {
		if d.lines[idx].lineType == lineTypeKv {
			keys = append(keys, FieldValue{d.lines[idx].key, d.lines[idx].value})
		}
	}
	return keys
}

// Returns a list of all sections within the document. Subsections are not listed, can be passed a `true` value
// to list those as well.
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

// Removes all unnecessary empty lines from the deocument.
func (d *IniDoc) StripWhiteLines() {
	newLines := make([]iniLine, 0, len(d.lines))
	for idx := range d.lines {
		if d.lines[idx].lineType != lineTypeWhiteLine {
			newLines = append(newLines, d.lines[idx])
		}
	}
	d.lines = newLines

	for _, section := range d.sections {
		section.StripWhiteLines()
	}
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

func (d *IniSection) lastLine() *iniLine {
	if len(d.lines) == 0 {
		return nil
	}
	return &d.lines[len(d.lines)-1]
}

// Adds an empty line
func (d *IniSection) AddWhiteLine() {
	d.lines = append(d.lines, iniLine{
		lineType: lineTypeWhiteLine,
	})
}

// Adds a comment line
func (d *IniSection) AddComment(value string) {
	lastLine := d.lastLine()
	if lastLine != nil && lastLine.lineType == lineTypeComment {
		lastLine.value += "\n" + value
		return
	}

	d.lines = append(d.lines, iniLine{
		lineType: lineTypeComment,
		value:    value,
	})
}

// Adds a comment line, but with a `#` instead of the default `;`
func (d *IniSection) AddHashComment(value string) {
	lastLine := d.lastLine()
	if lastLine != nil && lastLine.lineType == lineTypeHashComment {
		lastLine.value += "\n" + value
		return
	}

	d.lines = append(d.lines, iniLine{
		lineType: lineTypeHashComment,
		value:    value,
	})
}

// Remove the key-value pair from this section
func (d *IniSection) Del(key string) {
	for idx := range d.lines {
		if d.lines[idx].lineType == lineTypeKv && d.lines[idx].key == key {
			d.lines = slices.Delete(d.lines, idx, idx+1)
			return
		}
	}
}

// Adds a comment after a key-value pair, comment will be on the same line as the property (e.x. `key=value ; comment`)
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

// Returns the comment associated with this Section. This is the comment right above the
// section name.
func (d *IniSection) GetSectionComment() string {
	return d.comment
}

// Replace the comment associated with this Section. This is the comment right above the
// section name.
func (d *IniSection) SetSectionComment(comment string) {
	d.comment = comment
}

// Returns the current comment that's associated with the given property key
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

// Retrieves the given sub-section, if that sub-section does not exist it will be added
func (d *IniSection) Section(sectionName string) *IniSection {
	if d.root == nil {
		d.root = &IniDoc{}
	}

	if d.name != "" {
		return d.root.Section(fmt.Sprintf("%s.%s", d.name, sectionName))
	}

	return d.root.Section(sectionName)
}

// Returns a list of all keys of this section key-value pairs
func (d *IniSection) Keys() []string {
	keys := make([]string, 0, len(d.lines))
	for idx := range d.lines {
		if d.lines[idx].lineType == lineTypeKv {
			keys = append(keys, d.lines[idx].key)
		}
	}
	return keys
}

// Returns a list of all key-value pairs within this section
func (d *IniSection) Values() []FieldValue {
	keys := make([]FieldValue, 0, len(d.lines))
	for idx := range d.lines {
		if d.lines[idx].lineType == lineTypeKv {
			keys = append(keys, FieldValue{d.lines[idx].key, d.lines[idx].value})
		}
	}
	return keys
}

// Returns a list of all sub-sections. Only direct subsections are listed by default,
// `true` argument can be passed to list all subsections to any level deep.
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

// Removes all unnecessary empty lines from the deocument.
func (d *IniSection) StripWhiteLines() {
	newLines := make([]iniLine, 0, len(d.lines))
	for idx := range d.lines {
		if d.lines[idx].lineType != lineTypeWhiteLine {
			newLines = append(newLines, d.lines[idx])
		}
	}
	d.lines = newLines
}

// Change the name of this section
func (d *IniSection) SetName(name string) {
	d.name = name
}

func (d *IniSection) addParsedSection(name string) *IniSection {
	if d.root == nil {
		d.root = &IniDoc{}
	}

	if d.name != "" {
		return d.root.addParsedSection(fmt.Sprintf("%s.%s", d.name, name))
	}

	return d.root.addParsedSection(name)

	// lastLine := d.root.lastLine()
	// if d.root == from && lastLine != nil && (lastLine.lineType == lineTypeComment || lastLine.lineType == lineTypeHashComment) {
	// 	d.root.lines = d.root.lines[:len(d.root.lines)-1]
	// 	d.comment = lastLine.value
	// 	return
	// }

	// lastSection := d.root.lastSection()
	// if lastSection != nil && lastSection == from {
	// 	lastLine := lastSection.lastLine()
	// 	if lastLine != nil && (lastLine.lineType == lineTypeComment || lastLine.lineType == lineTypeHashComment) {
	// 		lastSection.lines = lastSection.lines[:len(lastSection.lines)-1]
	// 		d.comment = lastLine.value
	// 	}
	// }
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
	var v string = ""
	switch f.lineType {
	case lineTypeKv:
		v = fmt.Sprintf("%s=%s", f.key, escapeIniValue(f.value))
		if f.comment != "" {
			v += fmt.Sprintf(" ; %s", f.comment)
		}
		return v + "\n"
	case lineTypeComment:
		lines := strings.Split(f.value, "\n")
		for _, line := range lines {
			v += fmt.Sprintf("; %s\n", line)
		}
		return v
	case lineTypeHashComment:
		lines := strings.Split(f.value, "\n")
		for _, line := range lines {
			v += fmt.Sprintf("# %s\n", line)
		}
		return v
	case lineTypeWhiteLine:
		v += "\n"
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
		lines := strings.Split(s.comment, "\n")
		for _, line := range lines {
			v += fmt.Sprintf("; %s\n", line)
		}
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
			if len(v) >= 2 && v[len(v)-2:] != "\n\n" {
				v += "\n"
			}
			v += secStr
		}
	}

	return v
}

func docToSection(doc *IniDoc) *IniSection {
	sec := IniSection{
		root:  doc,
		lines: doc.lines,
	}
	doc.lines = []iniLine{}
	return &sec
}

func isCommentLine(l *iniLine) bool {
	return l.lineType == lineTypeComment || l.lineType == lineTypeHashComment
}
