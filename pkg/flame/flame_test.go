package flame

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"pg_flame/pkg/plan"
)

func TestNew(t *testing.T) {

	t.Run("creates a new Flame from a Plan", func(t *testing.T) {
		p := plan.Plan{
			PlanningTime: 0.01,
			ExecutionTree: plan.Node{
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

		assert.Equal(t, "Total", f.Name)
		assert.Equal(t, 0.133, f.Value)

		assert.Equal(t, "Query Planning", f.Children[0].Name)
		assert.Equal(t, 0.01, f.Children[0].Value)

		assert.Equal(t, "Limit", f.Children[1].Name)
		assert.Equal(t, 0.123, f.Children[1].Value)

		assert.Equal(t, "Seq Scan on bears", f.Children[1].Children[0].Name)
		assert.Equal(t, 0.022, f.Children[1].Children[0].Value)
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
			BuffersHit:  8,
			BuffersRead: 5,
			MemoryUsage: 12,
			HashBuckets: 1024,
			HashBatches: 1,
		}

		expected := strings.Join([]string{
			"<table class=\"table table-striped table-bordered\"><tbody>",
			"<tr><th>Filter</th><td>(id = 123)</td></tr>",
			"<tr><th>Join Filter</th><td></td></tr>",
			"<tr><th>Hash Cond</th><td></td></tr>",
			"<tr><th>Index Cond</th><td></td></tr>",
			"<tr><th>Recheck Cond</th><td></td></tr>",
			"<tr><th>Buffers Shared Hit</th><td>8</td></tr>",
			"<tr><th>Buffers Shared Read</th><td>5</td></tr>",
			"<tr><th>Hash Buckets</th><td>1024</td></tr>",
			"<tr><th>Hash Batches</th><td>1</td></tr>",
			"<tr><th>Memory Usage</th><td>1kB</td></tr>",
			"</tbody></table>",
		}, "")

		assert.Equal(t, expected, detail(n))
	})

}
