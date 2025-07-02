package ini_test

import (
	"testing"

	"github.com/ncpa0cpl/ini"
)

type TestConfig struct {
	K    string  `ini:"k" json:"k,omitempty"`
	K1   int     `ini:"k1" json:"k1,omitempty"`
	K2   float64 `ini:"k2"`
	K3   int64   `ini:"k3"`
	User User    `ini:"user"`
}

type TestConfig2 struct {
	K    string `ini:"k"`
	User *User  `ini:"user"`
}

type User struct {
	Name string `ini:"name"`
	Age  int    `ini:"age"`
}

type OnlySection struct {
	SomeSection struct {
		Value string `ini:"value"`
	} `init:"some_section"`
}

type BasicKv struct {
	K string `ini:"k"`
}

func TestIniUnmarshal(t *testing.T) {
	expect := expect(t)

	docStr := `
k=v
k1=2
k2=2.2
k3=3

[user]
name=tom
age=-23
`

	cfg := TestConfig{}

	expect(ini.Unmarshal(docStr, &cfg)).NoErr()

	expect(cfg.K).ToBe("v")
	expect(cfg.K1).ToBe(int(2))
	expect(cfg.K2).ToBe(float64(2.2))
	expect(cfg.K3).ToBe(int64(3))
	expect(cfg.User.Name).ToBe("tom")
	expect(cfg.User.Age).ToBe(int(-23))
}

func TestIniUnmarshal2(t *testing.T) {
	expect := expect(t)

	docStr := `
k=123
isTrue=true
isFalse=false

[user]
name=Barbara
age=54
canDrink=true
isDead=false
`

	type User2 struct {
		Name     string `ini:"name"`
		Age      int    `ini:"age"`
		CanDrink bool   `ini:"canDrink"`
		IsDead   bool   `ini:"isDead"`
	}

	type TestConfig4 struct {
		K       string `ini:"k"`
		User    *User2 `ini:"user"`
		IsTrue  bool   `ini:"isTrue"`
		IsFalse bool   `ini:"isFalse"`
	}

	cfg := TestConfig4{}

	expect(ini.Unmarshal(docStr, &cfg)).NoErr()

	expect(cfg.K).ToBe("123")
	expect(cfg.User.Name).ToBe("Barbara")
	expect(cfg.User.Age).ToBe(54)
	expect(cfg.IsTrue).ToBe(true)
	expect(cfg.IsFalse).ToBe(false)
	expect(cfg.User.CanDrink).ToBe(true)
	expect(cfg.User.IsDead).ToBe(false)
}

func TestIniUnmarshal3(t *testing.T) {
	expect := expect(t)

	doc := `
k=Foo Bar Baz I can have whitespaces

[user]
name=Bara bara
age=54
`

	cfg := TestConfig2{}

	expect(ini.Unmarshal(doc, &cfg)).NoErr()

	expect(cfg.K).ToBe("Foo Bar Baz I can have whitespaces")
	expect(cfg.User.Name).ToBe("Bara bara")
	expect(cfg.User.Age).ToBe(54)
}

func TestIniUnmarshal4(t *testing.T) {
	expect := expect(t)

	doc := `[MapStrToStr]
foo=bar
baz=quux
corge=gorge
`

	type WithMapSection struct {
		MapStrToStr map[string]string
	}

	cfg := WithMapSection{}

	expect(ini.Unmarshal(doc, &cfg)).NoErr()

	expect(cfg.MapStrToStr["foo"]).ToBe("bar")
	expect(cfg.MapStrToStr["baz"]).ToBe("quux")
	expect(cfg.MapStrToStr["corge"]).ToBe("gorge")
}

func TestIniUnmarshal5(t *testing.T) {
	expect := expect(t)

	doc := `[MapStrToInt]
foo=1
baz=-312
corge=6969
`

	type WithMapSection struct {
		MapStrToInt map[string]int
	}

	cfg := WithMapSection{}

	expect(ini.Unmarshal(doc, &cfg)).NoErr()

	expect(cfg.MapStrToInt["foo"]).ToBe(1)
	expect(cfg.MapStrToInt["baz"]).ToBe(-312)
	expect(cfg.MapStrToInt["corge"]).ToBe(6969)
}

func TestIniUnmarshal6(t *testing.T) {
	expect := expect(t)

	doc := `[MapStrToInterface]
foo=1
bar=0.0001
baz=hello
corge=true
`

	type WithMapSection struct {
		MapStrToInterface map[string]any
	}

	cfg := WithMapSection{}

	expect(ini.Unmarshal(doc, &cfg)).NoErr()

	expect(cfg.MapStrToInterface["foo"]).ToBe("1")
	expect(cfg.MapStrToInterface["bar"]).ToBe("0.0001")
	expect(cfg.MapStrToInterface["baz"]).ToBe("hello")
	expect(cfg.MapStrToInterface["corge"]).ToBe("true")
}

func TestIniMarshal(t *testing.T) {
	expect := expect(t)
	cfg := &TestConfig{
		K:  " foobar ",
		K1: 1234,
		K2: 420.69,
		K3: -9999,
		User: User{
			Age:  100,
			Name: "Brian",
		},
	}

	doc, err := ini.Marshal(cfg)
	expect(err).NoErr()

	expectedResult := `k=foobar
k1=1234
k2=420.69
k3=-9999

[user]
name=Brian
age=100
`

	expect(doc).ToBe(expectedResult)
}

func TestIniMarshal2(t *testing.T) {
	expect := expect(t)
	cfg := TestConfig2{
		K: "1234%",
		User: &User{
			Age:  12,
			Name: "Tom",
		},
	}

	doc, err := ini.Marshal(cfg)
	expect(err).NoErr()

	expectedResult := `k=1234%

[user]
name=Tom
age=12
`

	expect(doc).ToBe(expectedResult)
}

func TestIniMarshal3(t *testing.T) {
	expect := expect(t)
	cfg := TestConfig2{
		K:    "1234%",
		User: nil,
	}

	doc, err := ini.Marshal(cfg)
	expect(err).NoErr()

	expectedResult := `k=1234%
`

	expect(doc).ToBe(expectedResult)
}

func TestIniMarshal4(t *testing.T) {
	expect := expect(t)
	strct := OnlySection{
		SomeSection: struct {
			Value string `ini:"value"`
		}{
			Value: "1234%",
		},
	}

	doc, err := ini.Marshal(&strct)
	expect(err).NoErr()

	expectedResult := `
[SomeSection]
value=1234%
`

	expect(doc).ToBe(expectedResult)
}

func TestIniMarshal5(t *testing.T) {
	type WithMapSection struct {
		TopLevelValue string

		MapSection map[string]string
	}

	expect := expect(t)
	strct := WithMapSection{
		TopLevelValue: "hello",
		MapSection: map[string]string{
			"foo": "bar",
			"1":   "world",
		},
	}

	doc, err := ini.Marshal(&strct)
	expect(err).NoErr()

	expectedResult := `TopLevelValue=hello

[MapSection]
1=world
foo=bar
`

	expect(doc).ToBe(expectedResult)
}

func TestIniMarshal6(t *testing.T) {
	type WithMapSection struct {
		TopLevelValue string

		MapSection map[string]int64
	}

	expect := expect(t)
	strct := WithMapSection{
		TopLevelValue: "hello",
		MapSection: map[string]int64{
			"bar": 420,
			"foo": 69,
		},
	}

	doc, err := ini.Marshal(&strct)
	expect(err).NoErr()

	expectedResult := `TopLevelValue=hello

[MapSection]
bar=420
foo=69
`

	expect(doc).ToBe(expectedResult)
}

func TestIniMarshal7(t *testing.T) {
	type WithMapSection struct {
		TopLevelValue string

		MapSection map[string]any
	}

	expect := expect(t)
	strct := WithMapSection{
		TopLevelValue: "hello",
		MapSection: map[string]any{
			"foo":   "bar",
			"baz":   420,
			"quux":  true,
			"corge": false,
			"gorge": struct{ Field string }{"fieldvalue"},
			"array": []string{"foobar"},
		},
	}

	doc, err := ini.Marshal(&strct)
	expect(err).NoErr()

	expectedResult := `TopLevelValue=hello

[MapSection]
baz=420
corge=false
foo=bar
quux=true
`

	expect(doc).ToBe(expectedResult)
}

type CustomMarshalIni struct {
	f   int64
	br  bool
	bz  string
	Qux string  `ini:"qux"`
	Sec BasicKv `ini:"sec"`
}

func (m *CustomMarshalIni) UnmarshalINI(doc ini.DocOrSection) error {
	m.f, _ = doc.GetInt("foo")
	m.br, _ = doc.GetBool("bar")
	m.bz = doc.Get("baz")
	return nil
}

func (m *CustomMarshalIni) MarshalINI() (ini.DocOrSection, error) {
	doc := ini.NewDoc()
	doc.SetInt("foo", m.f)
	doc.SetBool("bar", m.br)
	doc.Set("baz", m.bz)
	return doc, nil
}

func TestCustomUnmarshalDocStruct(t *testing.T) {
	expect := expect(t)

	docStr := `foo=-420
bar=true
baz=hello world
qux=value

[sec]
k=v`

	myini := CustomMarshalIni{}
	expect(ini.Unmarshal(docStr, &myini)).NoErr()

	expect(myini.f).ToBe(int64(-420))
	expect(myini.br).ToBe(true)
	expect(myini.bz).ToBe("hello world")
	expect(myini.Qux).ToBe("")   // CustomMarshalIni.UnmarshalINI does not parse the Qux field
	expect(myini.Sec.K).ToBe("") // CustomMarshalIni.UnmarshalINI does not parse the Section
}

type CustomMarshalSection struct {
	f   int64
	br  bool
	bz  string
	Qux string `ini:"qux"`
}

func (m *CustomMarshalSection) UnmarshalINI(doc ini.DocOrSection) error {
	m.f, _ = doc.GetInt("foo")
	m.br, _ = doc.GetBool("bar")
	m.bz = doc.Get("baz")
	return nil
}

func (m CustomMarshalSection) MarshalINI() (ini.DocOrSection, error) {
	section := ini.NewSection()
	section.SetInt("foo", m.f)
	section.SetBool("bar", m.br)
	section.Set("baz", m.bz)
	return section, nil
}

func TestCustomUnmarshalSectionStruct(t *testing.T) {
	expect := expect(t)

	type Ini struct {
		Top           string
		CustomSection *CustomMarshalSection
	}

	docStr := `Top=abc

[CustomSection]
foo=1024
bar=true
baz = Lorem ipsum dolor sit amet
`

	myini := Ini{}
	expect(ini.Unmarshal(docStr, &myini)).NoErr()

	expect(myini.Top).ToBe("abc")
	expect(myini.CustomSection.f).ToBe(int64(1024))
	expect(myini.CustomSection.br).ToBe(true)
	expect(myini.CustomSection.bz).ToBe("Lorem ipsum dolor sit amet")
}

func TestCustomMarshalDocStruct(t *testing.T) {
	expect := expect(t)

	myini := CustomMarshalIni{
		f:   512,
		br:  true,
		bz:  "|string|",
		Qux: "(string)",
		Sec: BasicKv{
			K: "v",
		},
	}

	docStr, err := ini.Marshal(&myini)
	expect(err).NoErr()

	expectedResult := `foo=512
bar=true
baz=|string|
`

	expect(docStr).ToBe(expectedResult)
}

func TestCustomMarshalSectionStruct(t *testing.T) {
	expect := expect(t)

	type Ini struct {
		Top           string
		CustomSection CustomMarshalSection
	}

	myini := Ini{
		Top: "TOP",
		CustomSection: CustomMarshalSection{
			f:   -7,
			br:  false,
			bz:  "bzbzbz",
			Qux: "string",
		},
	}

	docStr, err := ini.Marshal(myini)
	expect(err).NoErr()

	expectedResult := `Top=TOP

[CustomSection]
foo=-7
bar=false
baz=bzbzbz
`

	expect(docStr).ToBe(expectedResult)
}

type CustomMarshalSection2 struct {
	Value string
}

func (m *CustomMarshalSection2) MarshalINI() (ini.DocOrSection, error) {
	doc := ini.NewDoc()
	doc.SetInt("version", 1)
	doc.Set("value1", m.Value)
	return doc, nil
}

func TestCustomMarshalSectionStruct2(t *testing.T) {
	expect := expect(t)

	type Ini struct {
		Top           string
		CustomSection *CustomMarshalSection2
		BasicKV       BasicKv
	}

	myini := Ini{
		Top: "TOP2",
		CustomSection: &CustomMarshalSection2{
			Value: "foobar",
		},
		BasicKV: BasicKv{
			K: "somevalue",
		},
	}

	docStr, err := ini.Marshal(myini)
	expect(err).NoErr()

	expectedResult := `Top=TOP2

[CustomSection]
version=1
value1=foobar

[BasicKV]
k=somevalue
`

	expect(docStr).ToBe(expectedResult)
}

func TestMarshalSubsections(t *testing.T) {
	expect := expect(t)

	type C struct {
		Key string
	}

	type B struct {
		Key string
		C   C
	}

	type A struct {
		Key string
		B   B
	}

	type IniFile struct {
		A A
	}

	iniFile := IniFile{
		A: A{
			Key: "value a",
			B: B{
				Key: "value b",
				C: C{
					Key: "value c",
				},
			},
		},
	}

	docStr, err := ini.Marshal(iniFile)
	expect(err).NoErr()

	expectedResult := `
[A]
Key=value a

[A.B]
Key=value b

[A.B.C]
Key=value c
`

	expect(docStr).ToBe(expectedResult)
}

func TestMarshalSubsectionsWithMaps(t *testing.T) {
	expect := expect(t)

	type B struct {
		Key string
		C   map[string]string
	}

	type A struct {
		Key string
		B   B
	}

	type IniFile struct {
		A A
	}

	iniFile := IniFile{
		A: A{
			Key: "value a",
			B: B{
				Key: "value b",
				C: map[string]string{
					"mapElem1": "mapVal1",
					"mapElem2": "mapVal2",
				},
			},
		},
	}

	docStr, err := ini.Marshal(iniFile)
	expect(err).NoErr()

	expectedResult := `
[A]
Key=value a

[A.B]
Key=value b

[A.B.C]
mapElem1=mapVal1
mapElem2=mapVal2
`

	expect(docStr).ToBe(expectedResult)
}

func TestUnmarshalSubsections(t *testing.T) {
	expect := expect(t)

	type C struct {
		Key string
	}

	type B struct {
		Key string
		C   C
	}

	type A struct {
		Key string
		B   B
	}

	type IniFile struct {
		A A
	}

	docStr := `
[A]
Key=value a

[A.B]
Key=value b

[A.B.C]
Key=value c
`

	iniFile := IniFile{}
	expect(ini.Unmarshal(docStr, &iniFile)).NoErr()

	expect(iniFile.A.Key).ToBe("value a")
	expect(iniFile.A.B.Key).ToBe("value b")
	expect(iniFile.A.B.C.Key).ToBe("value c")
}

func TestUnmarshalSubsectionsWithMaps(t *testing.T) {
	expect := expect(t)

	type B struct {
		Key  string
		Key2 string
		C    map[string]string
	}

	type A struct {
		Key string
		B   B
	}

	type IniFile struct {
		A A
	}

	docStr := `
[A]
Key=value a

[A.B]
Key=value b

[A.B.C]
Key=value c
anotherKey=1

[A.B]
Key2=second value b
`

	iniFile := IniFile{}
	expect(ini.Unmarshal(docStr, &iniFile)).NoErr()

	expect(iniFile.A.Key).ToBe("value a")
	expect(iniFile.A.B.Key).ToBe("value b")
	expect(iniFile.A.B.Key2).ToBe("second value b")
	expect(iniFile.A.B.C["Key"]).ToBe("value c")
	expect(iniFile.A.B.C["anotherKey"]).ToBe("1")
}

type MarshalableWithSubsections struct {
	top     string
	subK    string
	subK2   int
	subsubK string
}

func (self MarshalableWithSubsections) MarshalINI() (ini.DocOrSection, error) {
	sec := ini.NewSection()
	sec.Set("TOP", self.top)

	sub := sec.Section("SA")
	sub.Set("K", self.subK)
	sub.SetInt("K2", int64(self.subK2))

	subsub := sub.Section("B")
	subsub.Set("K", self.subsubK)

	return sec, nil
}

func TestCustomSubsectionMarshaling(t *testing.T) {
	expect := expect(t)

	ini1 := MarshalableWithSubsections{
		top:     "hello",
		subK:    "world",
		subK2:   23,
		subsubK: "reeee",
	}

	docStr1, err := ini.Marshal(ini1)
	expect(err).NoErr()

	expect(docStr1).ToBe(`TOP=hello

[SA]
K=world
K2=23

[SA.B]
K=reeee
`)

	type IniFile struct {
		FooBar MarshalableWithSubsections
	}

	ini2 := IniFile{
		FooBar: MarshalableWithSubsections{
			top:     "hello",
			subK:    "world",
			subK2:   23,
			subsubK: "reeee",
		},
	}

	docStr2, err := ini.Marshal(ini2)
	expect(err).NoErr()

	expect(docStr2).ToBe(`
[FooBar]
TOP=hello

[FooBar.SA]
K=world
K2=23

[FooBar.SA.B]
K=reeee
`)

	type FooBarWrapper struct {
		FooBar MarshalableWithSubsections
	}

	type IniFile2 struct {
		Wrapper FooBarWrapper
	}

	ini3 := IniFile2{
		Wrapper: FooBarWrapper{MarshalableWithSubsections{
			top:     "hello",
			subK:    "world",
			subK2:   23,
			subsubK: "reeee",
		}},
	}

	docStr3, err := ini.Marshal(ini3)
	expect(err).NoErr()

	expect(docStr3).ToBe(`
[Wrapper.FooBar]
TOP=hello

[Wrapper.FooBar.SA]
K=world
K2=23

[Wrapper.FooBar.SA.B]
K=reeee
`)
}
