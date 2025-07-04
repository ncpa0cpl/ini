package ini_test

import (
	"os"
	"testing"

	"github.com/ncpa0cpl/ini"
)

func TestIni0(t *testing.T) {
	expect := expect(t)

	doc := `
[section]
k =v
`
	section := ini.Parse(doc).Section("section")

	expect(section.Get("k")).ToBe("v")
}

func TestIni2(t *testing.T) {
	expect := expect(t)

	docStr := `
[section]
k=4.20
k1 = 1
`
	doc := ini.Parse(docStr)
	section := doc.Section("section")

	k, err := section.GetFloat("k")
	expect(err).NoErr()
	expect(k).ToBe(4.20)

	k1, err := section.GetInt("k1")
	expect(err).NoErr()
	expect(k1).ToBe(int64(1))

	k1u, err := section.GetUint("k1")
	expect(err).NoErr()
	expect(k1u).ToBe(uint64(1))
}

func TestIni3(t *testing.T) {
	expect := expect(t)

	docStr := `
a =b
c=d
[section]
kKkK0 =abc def kgoeirj
`
	doc := ini.Parse(docStr)

	expect(doc.Section("section").Get("kKkK0")).ToBe("abc def kgoeirj")

	expect(doc.Get("a")).ToBe("b")

	expect(doc.Get("c")).ToBe("d")
}

func TestIni4(t *testing.T) {
	expect := expect(t)

	docStr := `
a =b
c= d

a1 = 2.1
`
	doc := ini.Parse(docStr)
	expect(doc.Get("a")).ToBe("b")
	expect(doc.Get("c")).ToBe("d")

	a1, err := doc.GetFloat("a1")
	expect(err).NoErr()
	expect(a1).ToBe(2.1)
}

func TestIni5(t *testing.T) {
	expect := expect(t)

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

	expect(doc.Keys()).ToBe([]string{"a"})
	expect(doc.SectionNames()).ToBe([]string{"s1", "s2", "s3"})
	expect(doc.Section("s1").Keys()).ToBe([]string{"k", "k1"})
	expect(doc.Section("s2").Keys()).ToBe([]string{"k2"})
	expect(doc.Section("s3").Keys()).ToBe([]string{"k", "a"})

	s2k2 := doc.Section("s2").Get("k2")
	expect(s2k2).ToBe("v22")
}

func TestIniDelete(t *testing.T) {
	expect := expect(t)

	docStr := `
k=v
a=b
c=d
[section]
secA=1
`
	iniDoc := ini.Parse(docStr)

	iniDoc.Del("a")
	expect(iniDoc.ToString()).ToBe(`
k=v
c=d

[section]
secA=1
`)

	iniDoc.Section("section").Del("secA")
	expect(iniDoc.ToString()).ToBe(`
k=v
c=d
`)

	iniDoc.Del("c")
	expect(iniDoc.ToString()).ToBe(`
k=v
`)

	iniDoc.Del("k")
	expect(iniDoc.ToString()).ToBe("\n")
}

func TestIniSet(t *testing.T) {
	expect := expect(t)

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

	expectedResult := `
k=v

[section]
a=11
c=12.3
`

	expect(doc.ToString()).ToBe(expectedResult)

	v, err := doc.Section("section").GetInt("a")
	expect(err).NoErr()
	expect(v).ToBe(int64(11))

	v1, err := doc.Section("section").GetFloat("c")
	expect(err).NoErr()
	expect(v1).ToBe(12.3)

	v2 := doc.Get("k")
	expect(v2).ToBe("v")
}

func TestIniSaveAndLoad(t *testing.T) {
	expect := expect(t)

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
	expect(err).NoErr()

	expectedResult := `
; 123
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

	expect(doc.ToString()).ToBe(expectedResult)

	expectedResult2 := `; 123
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

	doc.StripWhiteLines()
	expect(doc.ToString()).ToBe(expectedResult2)
}

func TestIniFile(t *testing.T) {
	expect := expect(t)

	file := "./test.ini"
	doc, err := ini.Load(file)
	expect(err).NoErr()

	expect(doc.Get("a")).ToBe("'23'34?::'<>,.'")
	expect(doc.Get("c")).ToBe("d")
	expect(doc.Section("s1").Get("k1")).ToBe("fdasf")

	s1k, err := doc.Section("s1").GetInt("k")
	expect(err).NoErr()
	expect(s1k).ToBe(int64(67676))

	s2 := doc.Section("s2")
	expect(s2.Get("k")).ToBe("3")
	expect(s2.Get("k2")).ToBe("945")
	expect(s2.Get("k3")).ToBe("-435")
	expect(s2.Get("k4")).ToBe("0.0.0.0")
	expect(s2.Get("k5")).ToBe("127.0.0.1")
	expect(s2.Get("k6")).ToBe("levene@github.com")
	expect(s2.Get("k7")).ToBe("~/.path.txt")
	expect(s2.Get("k8")).ToBe("./34/34/uh.txt")
	expect(s2.Get("k9")).ToBe("234@!@#$%^&*()324")
	expect(s2.Get("k10")).ToBe("'23'34?::'<>,.'")
}

func TestIniSave2(t *testing.T) {
	expect := expect(t)

	filename := "./save.ini"
	doc := ini.NewDoc()
	doc.SetInt("a1", 1)
	doc.Save(filename)

	bts, _ := os.ReadFile(filename)

	expect(string(bts)).ToBe("a1=1\n")
}

func TestIniSave3(t *testing.T) {
	expect := expect(t)

	filename := "./save.ini"
	doc := ini.NewDoc()
	doc.SetInt("a1", 985123)
	doc.Section("s1").Set("a2", "v2")
	doc.Save(filename)

	bts, _ := os.ReadFile(filename)

	expect(string(bts)).ToBe("a1=985123\n\n[s1]\na2=v2\n")
}

func TestIni6(t *testing.T) {
	expect := expect(t)

	docStr := `k=v`
	doc := ini.Parse(docStr)

	expect(doc.Get("k")).ToBe("v")
}

func TestIni7(t *testing.T) {
	expect := expect(t)

	docStr := `k=v ;this is comment`
	doc := ini.Parse(docStr)

	expect(doc.Get("k")).ToBe("v")
	expect(doc.GetComment("k")).ToBe("this is comment")
}

func TestMultilineValues(t *testing.T) {
	expect := expect(t)

	doc := ini.NewDoc()

	doc.Set("foo", "foo")
	doc.Set("multiline", `Lorem ipsum
dolor sit amet,
consectetur adipiscing elit.`)
	doc.Set("bar", "bar")

	docStr := doc.ToString()

	doc2 := ini.Parse(docStr)

	expect(doc2.Get("foo")).ToBe("foo")
	expect(doc2.Get("multiline")).ToBe("Lorem ipsum\ndolor sit amet,\nconsectetur adipiscing elit.")
	expect(doc2.Get("bar")).ToBe("bar")
}

func TestMultilineValues2(t *testing.T) {
	expect := expect(t)

	doc := ini.NewDoc()

	doc.Set("foo", "foo")
	doc.Set("multiline", `Line One
Line Two \\N|
Line Three`)
	doc.Set("bar", "bar")

	docStr := doc.ToString()

	doc2 := ini.Parse(docStr)

	expect(doc2.Get("foo")).ToBe("foo")
	expect(doc2.Get("multiline")).ToBe("Line One\nLine Two \\N|\nLine Three")
	expect(doc2.Get("bar")).ToBe("bar")
}

func TestDocSetCharEscaping(t *testing.T) {
	expect := expect(t)

	doc := ini.NewDoc()

	doc.Set("key", "value;not a comment # also not a comment = foobar")
	doc.SetFieldComment("key", "this is a comment")

	docStr := doc.ToString()

	expectedResult := "key=value\\;not a comment \\# also not a comment = foobar ; this is a comment\n"

	expect(docStr).ToBe(expectedResult)

	doc2 := ini.Parse(docStr)

	expect(doc2.Get("key")).ToBe("value;not a comment # also not a comment = foobar")
}

func TestSectionComments(t *testing.T) {
	expect := expect(t)

	docStr := `; hello world
; this is a section comment
[FooBar]
k=v

[FooBaz]
k2=v2

# this section also has a comment
# Lorem ipsum dolor sit amet,
[FooQux]
k3=v3

; this is a subsection of the FooQux
[FooQux.A]
k4=v4
`

	doc := ini.Parse(docStr)

	expect(doc.SectionNames(true)).ToContain("FooBar", "FooBaz", "FooQux", "FooQux.A")
	expect(doc.Section("FooBar").GetSectionComment()).ToBe("hello world\nthis is a section comment")
	expect(doc.Section("FooBaz").GetSectionComment()).ToBe("")
	expect(doc.Section("FooQux").GetSectionComment()).ToBe("this section also has a comment\nLorem ipsum dolor sit amet,")
	expect(doc.Section("FooQux.A").GetSectionComment()).ToBe("this is a subsection of the FooQux")

	expect(doc.ToString()).ToBe(`; hello world
; this is a section comment
[FooBar]
k=v

[FooBaz]
k2=v2

; this section also has a comment
; Lorem ipsum dolor sit amet,
[FooQux]
k3=v3

; this is a subsection of the FooQux
[FooQux.A]
k4=v4
`)
}

func TestSectionComments2(t *testing.T) {
	expect := expect(t)

	docStr := `foo=bar
; this is a standalone comment

; hello world
; this is a section comment
[FooBar]
k=v

; tralalala

[FooBaz]
k2=v2
`

	doc := ini.Parse(docStr)

	expect(doc.SectionNames(true)).ToContain("FooBar", "FooBaz")
	expect(doc.Section("FooBar").GetSectionComment()).ToBe("hello world\nthis is a section comment")
	expect(doc.Section("FooBaz").GetSectionComment()).ToBe("")

	expect(doc.ToString()).ToBe(`foo=bar
; this is a standalone comment

; hello world
; this is a section comment
[FooBar]
k=v

; tralalala

[FooBaz]
k2=v2
`)
}
