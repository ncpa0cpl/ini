package ini_test

import (
	"testing"

	"github.com/ncpa0cpl/ini"
)

func TestDoc(t *testing.T) {
	expect := expect(t)

	doc := ini.NewDoc()
	doc.AddComment("this is a comment")
	doc.Set("str", "hello world")
	doc.SetBool("bl", true)
	doc.SetBool("bl2", false)
	doc.AddComment("numbers:")
	doc.SetFloat("flt", 1.2345)
	doc.SetInt("int", 420)
	doc.SetFieldComment("bl", "this is a boolean")

	stringifiedDoc := doc.ToString()

	expectedResult := `; this is a comment
str=hello world
bl=true ;this is a boolean
bl2=false
; numbers:
flt=1.2345
int=420
`
	expect(stringifiedDoc).ToBe(expectedResult)
}

func TestDocSection(t *testing.T) {
	expect := expect(t)

	doc := ini.NewDoc()
	doc.Set("top", "lorem ipsum")

	section1 := doc.Section("MySection")
	section1.SetFloat("number", 4.2069)
	section1.Set("foobar", "bazquux")
	section1.SetBool("is", true)

	section2 := doc.Section("AnotherSection")
	section2.SetBool("is", false)
	section2.Set("bazquux", "foobar")
	section2.SetInt("number", -999)

	stringifiedDoc := doc.ToString()

	expectedResult := `top=lorem ipsum

[MySection]
number=4.2069
foobar=bazquux
is=true

[AnotherSection]
is=false
bazquux=foobar
number=-999
`
	expect(stringifiedDoc).ToBe(expectedResult)
}

func TestSections(t *testing.T) {
	expect := expect(t)

	docStr := `
[sect]
foo=1
bar=abc

[sect.sub1]
foo=2
baz=true

[sect.sub1.nestedSub]
abc=def

[sect.sub2]
ghi=jkl`

	doc := ini.Parse(docStr)

	expect(doc.Section("sect").Keys()).ToBe([]string{"foo", "bar"})
	expect(doc.Section("sect.sub1").Keys()).ToBe([]string{"foo", "baz"})
	expect(doc.Section("sect.sub1.nestedSub").Keys()).ToBe([]string{"abc"})
	expect(doc.Section("sect.sub2").Keys()).ToBe([]string{"ghi"})
	expect(doc.Section("sect").Section("sub1").Keys()).ToBe([]string{"foo", "baz"})
	expect(doc.Section("sect").Section("sub1").Section("nestedSub").Keys()).ToBe([]string{"abc"})
	expect(doc.Section("sect").Section("sub1.nestedSub").Keys()).ToBe([]string{"abc"})
	expect(doc.Section("sect").Section("sub2").Keys()).ToBe([]string{"ghi"})

	expect(doc.Section("sect").Get("foo")).ToBe("1")
	expect(doc.Section("sect").Get("bar")).ToBe("abc")
	expect(doc.Section("sect").Section("sub1").Get("foo")).ToBe("2")
	expect(doc.Section("sect").Section("sub1").Get("baz")).ToBe("true")
	expect(doc.Section("sect").Section("sub1").Section("nestedSub").Get("abc")).ToBe("def")
	expect(doc.Section("sect").Section("sub1.nestedSub").Get("abc")).ToBe("def")
	expect(doc.Section("sect").Section("sub2").Get("ghi")).ToBe("jkl")

	expect(doc.SectionNames()).ToBe([]string{"sect"})
	expect(doc.SectionNames(true)).ToBe([]string{"sect", "sect.sub1", "sect.sub1.nestedSub", "sect.sub2"})

	expect(doc.Section("sect").SubsectionNames()).ToBe([]string{"sub1", "sub2"})
	expect(doc.Section("sect").SubsectionNames(true)).ToBe([]string{"sub1", "sub1.nestedSub", "sub2"})
	expect(doc.Section("sect.sub1").SubsectionNames()).ToBe([]string{"nestedSub"})
}

func TestSections2(t *testing.T) {
	expect := expect(t)

	doc := ini.NewDoc()

	sect := doc.Section("sect")
	sect.SetInt("foo", 1)
	sect.Set("bar", "abc")

	sectSub1 := sect.Section("sub1")
	sectSub1.SetInt("foo", 2)
	sectSub1.SetBool("baz", true)

	nestedSub := sectSub1.Section("nestedSub")
	nestedSub.Set("abc", "def")

	sub2 := sect.Section("sub2")
	sub2.Set("ghi", "jkl")

	expectedResult := `
[sect]
foo=1
bar=abc

[sect.sub1]
foo=2
baz=true

[sect.sub1.nestedSub]
abc=def

[sect.sub2]
ghi=jkl
`

	expect(doc.ToString()).ToBe(expectedResult)
}
