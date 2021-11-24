package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	fixturePlugin = "../fixture/test-plugin.so"
)

func TestLoader(t *testing.T) {
	tests := []struct {
		Desc      string
		SymName   string
		ExpectErr error
	}{
		{
			Desc:      "Plugin symbol is not exist",
			SymName:   "../fixture/not-exist.so",
			ExpectErr: ErrSymbolNotExist,
		},
		{
			Desc:      "Load plugin",
			SymName:   fixturePlugin,
			ExpectErr: nil,
		},
	}

	for _, test := range tests {
		l := loader{}

		callSpec, err := l.loadSymbol(test.SymName)

		if test.ExpectErr != nil {
			assert.Nil(t, callSpec)
			assert.Equal(t, test.ExpectErr, err)

			continue
		}

		assert.NotNil(t, callSpec)

		spec := Spec{
			Name:      "dummy",
			CallSpecs: callSpec,
		}

		spec.Invoke("add", "1", "10")
	}
}

// TODO: Need to check parameters?
