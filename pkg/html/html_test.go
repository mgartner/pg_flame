package html

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"pg_flame/pkg/plan"
)

func TestNew(t *testing.T) {

	t.Run("writes an HTML flamegraph based on a Flame", func(t *testing.T) {
		p := plan.Plan{
			ExecutionTree: plan.Node{
				Table:           "bears",
				ActualTotalTime: 0.022,
			},
		}

		b := new(bytes.Buffer)

		err := Generate(b, p)

		assert.NoError(t, err)

		assert.Contains(t, b.String(), p.ExecutionTree.Table)
	})

}

func Test_buildFlame(t *testing.T) {

	t.Run("creates a new Flame from a Plan", func(t *testing.T) {
		p := plan.Plan{
			PlanningTime: 0.01,
			ExecutionTree: plan.Node{
				Method:          "Limit",
				ActualTotalTime: 0.123,
				Children: []plan.Node{
					{
						Method:          "Seq Scan",
						Table:           "bears",
						ActualTotalTime: 0.022,
					},
				},
			},
		}

		f, err := buildFlame(p)

		assert.NoError(t, err)

		assert.Equal(t, "Total", f.Name)
		assert.Equal(t, 0.133, f.Value)
		assert.Equal(t, 0.133, f.Time)
		assert.Equal(t, "<span>Includes planning and execution time</span>", f.Detail)

		assert.Equal(t, "Query Planning", f.Children[0].Name)
		assert.Equal(t, colorPlan, f.Children[0].Color)
		assert.Equal(t, 0.01, f.Children[0].Value)
		assert.Equal(t, 0.01, f.Children[0].Time)
		assert.Equal(t, "<span>Time to generate the query plan</span>", f.Children[0].Detail)

		assert.Equal(t, "Limit", f.Children[1].Name)
		assert.Equal(t, 0.123, f.Children[1].Value)
		assert.Equal(t, 0.123, f.Children[1].Time)

		assert.Equal(t, "Seq Scan on bears", f.Children[1].Children[0].Name)
		assert.Equal(t, 0.022, f.Children[1].Children[0].Value)
		assert.Equal(t, 0.022, f.Children[1].Children[0].Time)
	})

	t.Run("handles InitPlan nodes", func(t *testing.T) {
		p := plan.Plan{
			ExecutionTree: plan.Node{
				Method:          "Seq Scan",
				ActualTotalTime: 0.12,
				Children: []plan.Node{
					{
						Method:             "Seq Scan",
						Table:              "bears",
						ParentRelationship: "InitPlan",
						ActualTotalTime:    0.2,
						Children: []plan.Node{
							{
								Method:          "Seq Scan",
								ActualTotalTime: 0.12,
							},
						},
					},
				},
			},
		}

		f, err := buildFlame(p)

		assert.NoError(t, err)

		assert.Equal(t, "Total", f.Name)
		assert.Equal(t, 0.32, f.Value)
		assert.Equal(t, 0.12, f.Time)

		assert.Equal(t, "Seq Scan", f.Children[1].Name)
		assert.Equal(t, 0.32, f.Children[1].Value)
		assert.Equal(t, 0.12, f.Children[1].Time)
		assert.Equal(t, "", f.Children[1].Color)
		assert.False(t, f.Children[1].InitPlan)

		assert.Equal(t, "Seq Scan on bears", f.Children[1].Children[0].Name)
		assert.Equal(t, 0.2, f.Children[1].Children[0].Value)
		assert.Equal(t, 0.2, f.Children[1].Children[0].Time)
		assert.Equal(t, colorInit, f.Children[1].Children[0].Color)
		assert.True(t, f.Children[1].Children[0].InitPlan)

		assert.Equal(t, colorInit, f.Children[1].Children[0].Children[0].Color)
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

	t.Run("returns the method, index, and table if table exists", func(t *testing.T) {
		n := plan.Node{
			Method: "Index Scan",
			Table:  "bears",
			Index:  "bears_pkey",
		}

		assert.Equal(t, "Index Scan using bears_pkey on bears", name(n))
	})

	t.Run("returns the method if there is no table", func(t *testing.T) {
		n := plan.Node{Method: "Seq Scan"}

		assert.Equal(t, "Seq Scan", name(n))
	})

}

func Test_detail(t *testing.T) {

	t.Run("returns a table of details", func(t *testing.T) {
		n := plan.Node{
			Filter:      "(id = 123)",
			MemoryUsage: 12,
		}

		d, err := detail(n)

		assert.NoError(t, err)

		assert.Contains(t, d, n.Filter)
		assert.Contains(t, d, "12 kB")
	})

}
