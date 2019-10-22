package flame

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"pg_flame/pkg/plan"
)

func TestNew(t *testing.T) {

	t.Run("creates a new Flame from a Plan", func(t *testing.T) {
		p := plan.Plan{
			Root: plan.Node{
				Method:    "Limit",
				TotalTime: 0.123,
				Children: []plan.Node{
					{
						Method:    "Seq Scan",
						Table:     "bears",
						TotalTime: 0.022,
					},
				},
			},
		}

		f := New(p)

		assert.Equal(t, "Limit", f.Name)
		assert.Equal(t, 0.123, f.Value)

		assert.Equal(t, "Seq Scan on bears", f.Children[0].Name)
		assert.Equal(t, 0.022, f.Children[0].Value)
	})

}

func Test_name(t *testing.T) {

	t.Run("returns the method and table if table exists", func(t *testing.T) {
		n := plan.Node{
			Method: "Seq Scan",
			Table:  "bears",
		}

		assert.Equal(t, "Seq Scan on bears", name(n))
	})

	t.Run("returns the method if there is no table", func(t *testing.T) {
		n := plan.Node{Method: "Seq Scan"}

		assert.Equal(t, "Seq Scan", name(n))
	})

}
