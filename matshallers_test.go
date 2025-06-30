package ini_test

import (
	"testing"

	"github.com/ncpa0cpl/ini"
	"github.com/stretchr/testify/assert"
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

func TestIniUnmarshal(t *testing.T) {
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

	assertNoError(ini.Unmarshal(docStr, &cfg))

	a := assert.New(t)
	a.Equal("v", cfg.K, "'k' was not parsed correctly")
	a.Equal(int(2), cfg.K1, "'k1' was not parsed correctly")
	a.Equal(float64(2.2), cfg.K2, "'k2' was not parsed correctly")
	a.Equal(int64(3), cfg.K3, "'k3' was not parsed correctly")
	a.Equal("tom", cfg.User.Name, "'user.name' was not parsed correctly")
	a.Equal(int(-23), cfg.User.Age, "'user.name' was not parsed correctly")
}

func TestIniUnmarshal2(t *testing.T) {
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

	assertNoError(ini.Unmarshal(docStr, &cfg))

	a := assert.New(t)
	a.Equal("123", cfg.K, "'k' was not parsed correctly")
	a.Equal("Barbara", cfg.User.Name, "'user.name' was not parsed correctly")
	a.Equal(54, cfg.User.Age, "'user.name' was not parsed correctly")
	a.Equal(true, cfg.IsTrue, "'user.isTrue' was not parsed correctly")
	a.Equal(false, cfg.IsFalse, "'user.isFalse' was not parsed correctly")
	a.Equal(true, cfg.User.CanDrink, "'user.canDrink' was not parsed correctly")
	a.Equal(false, cfg.User.IsDead, "'user.isDead' was not parsed correctly")
}

func TestIniUnmarshal3(t *testing.T) {
	doc := `
k=Foo Bar Baz I can have whitespaces

[user]
name=Bara bara
age=54
`

	cfg := TestConfig2{}

	assertNoError(ini.Unmarshal(doc, &cfg))

	a := assert.New(t)
	a.Equal("Foo Bar Baz I can have whitespaces", cfg.K, "'k' was not parsed correctly")
	a.Equal("Bara bara", cfg.User.Name, "'user.name' was not parsed correctly")
	a.Equal(54, cfg.User.Age, "'user.name' was not parsed correctly")
}

func TestIniUnmarshal4(t *testing.T) {
	doc := `[MapStrToStr]
foo=bar
baz=quux
corge=gorge
`

	type WithMapSection struct {
		MapStrToStr map[string]string
	}

	cfg := WithMapSection{}

	assertNoError(ini.Unmarshal(doc, &cfg))

	a := assert.New(t)
	a.Equal("bar", cfg.MapStrToStr["foo"])
	a.Equal("quux", cfg.MapStrToStr["baz"])
	a.Equal("gorge", cfg.MapStrToStr["corge"])
}

func TestIniUnmarshal5(t *testing.T) {
	doc := `[MapStrToInt]
foo=1
baz=-312
corge=6969
`

	type WithMapSection struct {
		MapStrToInt map[string]int
	}

	cfg := WithMapSection{}

	assertNoError(ini.Unmarshal(doc, &cfg))

	a := assert.New(t)
	a.Equal(1, cfg.MapStrToInt["foo"])
	a.Equal(-312, cfg.MapStrToInt["baz"])
	a.Equal(6969, cfg.MapStrToInt["corge"])
}

func TestIniUnmarshal6(t *testing.T) {
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

	assertNoError(ini.Unmarshal(doc, &cfg))

	a := assert.New(t)
	a.Equal("1", cfg.MapStrToInterface["foo"])
	a.Equal("0.0001", cfg.MapStrToInterface["bar"])
	a.Equal("hello", cfg.MapStrToInterface["baz"])
	a.Equal("true", cfg.MapStrToInterface["corge"])
}

func TestIniMarshal(t *testing.T) {
	a := assert.New(t)
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
	a.NoError(err, "marshal operation failed")

	expectedResult := `k=foobar
k1=1234
k2=420.69
k3=-9999

[user]
name=Brian
age=100
`

	a.Equal(expectedResult, doc, "TestConfig was not marshaled correctly")
}

func TestIniMarshal2(t *testing.T) {
	a := assert.New(t)
	cfg := TestConfig2{
		K: "1234%",
		User: &User{
			Age:  12,
			Name: "Tom",
		},
	}

	doc, err := ini.Marshal(cfg)
	a.NoError(err, "marshal operation failed")

	expectedResult := `k=1234%

[user]
name=Tom
age=12
`

	a.Equal(expectedResult, doc, "TestConfig was not marshaled correctly")
}

func TestIniMarshal3(t *testing.T) {
	a := assert.New(t)
	cfg := TestConfig2{
		K:    "1234%",
		User: nil,
	}

	doc, err := ini.Marshal(cfg)
	a.NoError(err, "marshal operation failed")

	expectedResult := `k=1234%
`

	a.Equal(expectedResult, doc, "TestConfig was not marshaled correctly")
}

func TestIniMarshal4(t *testing.T) {
	a := assert.New(t)
	strct := OnlySection{
		SomeSection: struct {
			Value string `ini:"value"`
		}{
			Value: "1234%",
		},
	}

	doc, err := ini.Marshal(&strct)
	a.NoError(err, "marshal operation failed")

	expectedResult := `
[SomeSection]
value=1234%
`

	a.Equal(expectedResult, doc, "TestConfig was not marshaled correctly")
}

func TestIniMarshal5(t *testing.T) {
	type WithMapSection struct {
		TopLevelValue string

		MapSection map[string]string
	}

	a := assert.New(t)
	strct := WithMapSection{
		TopLevelValue: "hello",
		MapSection: map[string]string{
			"foo": "bar",
			"1":   "world",
		},
	}

	doc, err := ini.Marshal(&strct)
	a.NoError(err, "marshal operation failed")

	expectedResult := `TopLevelValue=hello

[MapSection]
foo=bar
1=world
`

	a.Equal(expectedResult, doc, "TestConfig was not marshaled correctly")
}

func TestIniMarshal6(t *testing.T) {
	type WithMapSection struct {
		TopLevelValue string

		MapSection map[string]int64
	}

	a := assert.New(t)
	strct := WithMapSection{
		TopLevelValue: "hello",
		MapSection: map[string]int64{
			"bar": 420,
			"foo": 69,
		},
	}

	doc, err := ini.Marshal(&strct)
	a.NoError(err, "marshal operation failed")

	expectedResult := `TopLevelValue=hello

[MapSection]
bar=420
foo=69
`

	a.Equal(expectedResult, doc, "TestConfig was not marshaled correctly")
}

func TestIniMarshal7(t *testing.T) {
	type WithMapSection struct {
		TopLevelValue string

		MapSection map[string]any
	}

	a := assert.New(t)
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
	a.NoError(err, "marshal operation failed")

	expectedResult := `TopLevelValue=hello

[MapSection]
foo=bar
baz=420
quux=true
corge=false
`

	a.Equal(expectedResult, doc, "TestConfig was not marshaled correctly")
}
