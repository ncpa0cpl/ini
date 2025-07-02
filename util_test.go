package ini_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Expect struct {
	t      *testing.T
	assert *assert.Assertions
	value  any
}

func expect(t *testing.T) func(any) *Expect {
	assert := assert.New(t)
	return func(value any) *Expect {
		return &Expect{t, assert, value}
	}
}

func (e *Expect) ToBe(equalTo any) {
	pass := e.assert.Equal(equalTo, e.value)
	if !pass {
		e.t.FailNow()
	}
}

func (e *Expect) ToContain(elems ...any) {
	for _, el := range elems {
		pass := e.assert.Contains(e.value, el)
		if !pass {
			e.t.FailNow()
		}
	}
}

func (e *Expect) NoErr() {
	pass := e.assert.Equal(nil, e.value)
	if !pass {
		e.t.FailNow()
	}
}
