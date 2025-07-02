package ini_test

import (
	"os"
	"testing"

	"github.com/ncpa0cpl/ini"
	"github.com/stretchr/testify/assert"
)

func TestIni0(t *testing.T) {
	assert := assert.New(t)

	doc := `
[section]
k =v
`
	section := ini.Parse(doc).Section("section")

	assert.Equal("v", section.Get("k"))
}

func TestIni2(t *testing.T) {
	assert := assert.New(t)

	docStr := `
[section]
k=4.20
k1 = 1
`
	doc := ini.Parse(docStr)
	section := doc.Section("section")

	k, err := section.GetFloat("k")
	assertNoError(err)
	assert.Equal(4.20, k)

	k1, err := section.GetInt("k1")
	assertNoError(err)
	assert.Equal(int64(1), k1)

	k1u, err := section.GetUint("k1")
	assertNoError(err)
	assert.Equal(uint64(1), k1u)
}

func TestIni3(t *testing.T) {
	assert := assert.New(t)

	docStr := `
a =b
c=d
[section]
kKkK0 =abc def kgoeirj
`
	doc := ini.Parse(docStr)

	assert.Equal("abc def kgoeirj", doc.Section("section").Get("kKkK0"))

	assert.Equal("b", doc.Get("a"))

	assert.Equal("d", doc.Get("c"))
}

func TestIni4(t *testing.T) {
	assert := assert.New(t)

	docStr := `
a =b
c= d

a1 = 2.1
`
	doc := ini.Parse(docStr)
	assert.Equal("b", doc.Get("a"))
	assert.Equal("d", doc.Get("c"))

	a1, err := doc.GetFloat("a1")
	assertNoError(err)
	assert.Equal(2.1, a1)
}

func TestIni5(t *testing.T) {
	assert := assert.New(t)

	docStr := `
a =b
[s1]
k=v
k1 = v12

[s2]
k2=v2
k2= v22

[s3]
k =v
a= b
`
	doc := ini.Parse(docStr)

	assert.Equal([]string{"a"}, doc.Keys())
	assert.Equal([]string{"s1", "s2", "s3"}, doc.SectionNames())
	assert.Equal([]string{"k", "k1"}, doc.Section("s1").Keys())
	assert.Equal([]string{"k2"}, doc.Section("s2").Keys())
	assert.Equal([]string{"k", "a"}, doc.Section("s3").Keys())

	s2k2 := doc.Section("s2").Get("k2")
	assert.Equal("v22", s2k2)
}

func TestIniDelete(t *testing.T) {
	assert := assert.New(t)

	docStr := `
k=v
a=b
c=d
[section]

`
	ini := ini.Parse(docStr)

	ini.Del("a")
	assert.Equal(`k=v
c=d

[section]
`, ini.ToString())

	ini.Del("c")
	assert.Equal(`k=v

[section]
`, ini.ToString())

	ini.Del("k")
	assert.Equal(`
[section]
`, ini.ToString())
}

func TestIniSet(t *testing.T) {
	assert := assert.New(t)

	docStr := `
k =v
[section]
a=b
c=d
`
	doc := ini.Parse(docStr)

	section := doc.Section("section")
	section.SetInt("a", 11)
	section.SetFloat("c", 12.3)

	expectedResult := `k=v

[section]
a=11
c=12.3
`

	assert.Equal(expectedResult, doc.ToString())

	v, err := doc.Section("section").GetInt("a")
	assertNoError(err)
	assert.Equal(int64(11), v)

	v1, err := doc.Section("section").GetFloat("c")
	assertNoError(err)
	assert.Equal(12.3, v1)

	v2 := doc.Get("k")
	assert.Equal("v", v2)
}

func TestIniSaveAndLoad(t *testing.T) {
	assert := assert.New(t)

	docStr := `
; 123
c11=d12312312
# 434

[section]
k=v
; dsfads
;123
#3452345


[section1]
k1=v1

[section3]
k3=v3
`
	ini.Parse(docStr).Save("./save.ini")

	doc, err := ini.Load("./save.ini")
	assertNoError(err)

	expectedResult := `; 123
c11=d12312312
# 434

[section]
k=v
; dsfads
; 123
# 3452345

[section1]
k1=v1

[section3]
k3=v3
`

	assert.Equal(expectedResult, doc.ToString())
}

func TestIniFile(t *testing.T) {
	file := "./test.ini"
	doc, err := ini.Load(file)
	assertNoError(err)

	a := assert.New(t)
	a.Equal("'23'34?::'<>,.'", doc.Get("a"))
	a.Equal("d", doc.Get("c"))
	a.Equal("fdasf", doc.Section("s1").Get("k1"))

	s1k, err := doc.Section("s1").GetInt("k")
	assertNoError(err)
	a.Equal(int64(67676), s1k)

	s2 := doc.Section("s2")
	a.Equal("3", s2.Get("k"))
	a.Equal("945", s2.Get("k2"))
	a.Equal("-435", s2.Get("k3"))
	a.Equal("0.0.0.0", s2.Get("k4"))
	a.Equal("127.0.0.1", s2.Get("k5"))
	a.Equal("levene@github.com", s2.Get("k6"))
	a.Equal("~/.path.txt", s2.Get("k7"))
	a.Equal("./34/34/uh.txt", s2.Get("k8"))
	a.Equal("234@!@#$%^&*()324", s2.Get("k9"))
	a.Equal("'23'34?::'<>,.'", s2.Get("k10"))
}

func TestIniSave2(t *testing.T) {
	assert := assert.New(t)

	filename := "./save.ini"
	doc := ini.NewDoc()
	doc.SetInt("a1", 1)
	doc.Save(filename)

	bts, _ := os.ReadFile(filename)

	assert.Equal("a1=1\n", string(bts))
}

func TestIniSave3(t *testing.T) {
	assert := assert.New(t)

	filename := "./save.ini"
	doc := ini.NewDoc()
	doc.SetInt("a1", 985123)
	doc.Section("s1").Set("a2", "v2")
	doc.Save(filename)

	bts, _ := os.ReadFile(filename)

	assert.Equal("a1=985123\n\n[s1]\na2=v2\n", string(bts))
}

func TestIni6(t *testing.T) {
	assert := assert.New(t)

	docStr := `k=v`
	doc := ini.Parse(docStr)

	assert.Equal("v", doc.Get("k"))
}

func TestIni7(t *testing.T) {
	assert := assert.New(t)

	docStr := `k=v ;this is comment`
	doc := ini.Parse(docStr)

	assert.Equal("v", doc.Get("k"))
	assert.Equal("this is comment", doc.GetComment("k"))
}

func TestMultilineValues(t *testing.T) {
	assert := assert.New(t)

	doc := ini.NewDoc()

	doc.Set("foo", "foo")
	doc.Set("multiline", `Lorem ipsum
dolor sit amet,
consectetur adipiscing elit.`)
	doc.Set("bar", "bar")

	docStr := doc.ToString()

	doc2 := ini.Parse(docStr)

	assert.Equal("foo", doc2.Get("foo"))
	assert.Equal("Lorem ipsum\ndolor sit amet,\nconsectetur adipiscing elit.", doc2.Get("multiline"))
	assert.Equal("bar", doc2.Get("bar"))
}

func TestMultilineValues2(t *testing.T) {
	assert := assert.New(t)

	doc := ini.NewDoc()

	doc.Set("foo", "foo")
	doc.Set("multiline", `Line One
Line Two \\N|
Line Three`)
	doc.Set("bar", "bar")

	docStr := doc.ToString()

	doc2 := ini.Parse(docStr)

	assert.Equal("foo", doc2.Get("foo"))
	assert.Equal("Line One\nLine Two \\N|\nLine Three", doc2.Get("multiline"))
	assert.Equal("bar", doc2.Get("bar"))
}

func TestDocSetCharEscaping(t *testing.T) {
	assert := assert.New(t)

	doc := ini.NewDoc()

	doc.Set("key", "value;not a comment # also not a comment = foobar")
	doc.SetFieldComment("key", "this is a comment")

	docStr := doc.ToString()

	expectedResult := "key=value\\;not a comment \\# also not a comment = foobar ;this is a comment\n"

	assert.Equal(expectedResult, docStr)

	doc2 := ini.Parse(docStr)

	assert.Equal("value;not a comment # also not a comment = foobar", doc2.Get("key"))
}
