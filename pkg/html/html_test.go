package html

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"pg_flame/pkg/flame"
)

func TestNew(t *testing.T) {

	t.Run("writes an HTML flamegraph based on a Flame", func(t *testing.T) {
		f := flame.Flame{
			Name:  "Seq Scan on bears",
			Value: 0.022,
		}

		b := new(bytes.Buffer)

		err := Generate(b, f)

		assert.NoError(t, err)
		assert.Contains(t, b.String(), f.Name)
	})

}
