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
	doc := `
k=v
k1=2
k2=2.2
k3=3

[user]
name=tom
age=-23
`

	cfg := TestConfig{}

	ini.Unmarshal([]byte(doc), &cfg)

	a := assert.New(t)
	a.Equal("v", cfg.K, "'k' was not parsed correctly")
	a.Equal(int(2), cfg.K1, "'k1' was not parsed correctly")
	a.Equal(float64(2.2), cfg.K2, "'k2' was not parsed correctly")
	a.Equal(int64(3), cfg.K3, "'k3' was not parsed correctly")
	a.Equal("tom", cfg.User.Name, "'user.name' was not parsed correctly")
	a.Equal(int(-23), cfg.User.Age, "'user.name' was not parsed correctly")
}

func TestIniUnmarshal2(t *testing.T) {
	doc := `
k=123

[user]
name=Barbara
age=54
`

	cfg := TestConfig2{}

	ini.Unmarshal([]byte(doc), &cfg)

	a := assert.New(t)
	a.Equal("123", cfg.K, "'k' was not parsed correctly")
	a.Equal("Barbara", cfg.User.Name, "'user.name' was not parsed correctly")
	a.Equal(54, cfg.User.Age, "'user.name' was not parsed correctly")
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

	excpectedResult := `k=foobar
k1=1234
k2=420.69
k3=-9999

[user]
name=Brian
age=100
`

	a.Equal(excpectedResult, doc, "TestConfig was not marshaled correctly")
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

	excpectedResult := `k=1234%

[user]
name=Tom
age=12
`

	a.Equal(excpectedResult, doc, "TestConfig was not marshaled correctly")
}

func TestIniMarshal3(t *testing.T) {
	a := assert.New(t)
	cfg := TestConfig2{
		K:    "1234%",
		User: nil,
	}

	doc, err := ini.Marshal(cfg)
	a.NoError(err, "marshal operation failed")

	excpectedResult := `k=1234%
`

	a.Equal(excpectedResult, doc, "TestConfig was not marshaled correctly")
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

	excpectedResult := `[SomeSection]
value=1234%
`

	a.Equal(excpectedResult, doc, "TestConfig was not marshaled correctly")
}
