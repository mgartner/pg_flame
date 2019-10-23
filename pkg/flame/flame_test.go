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

	t.Run("returns filter details", func(t *testing.T) {
		assert.Equal(t, "Filter: (id = 123)", detail(plan.Node{Filter: "(id = 123)"}))
		assert.Equal(t, "Join Filter: (id = 123)", detail(plan.Node{JoinFilter: "(id = 123)"}))
		assert.Equal(t, "Hash Cond: (id = 123)", detail(plan.Node{HashCond: "(id = 123)"}))
		assert.Equal(t, "Index Cond: (id = 123)", detail(plan.Node{IndexCond: "(id = 123)"}))
		assert.Equal(t, "Recheck Cond: (id = 123)", detail(plan.Node{RecheckCond: "(id = 123)"}))
	})

	t.Run("returns buffer details", func(t *testing.T) {
		n := plan.Node{
			BuffersHit:  8,
			BuffersRead: 5,
		}

		assert.Equal(t, "Buffers Shared Hit: 8, Buffers Shared Read: 5", detail(n))
	})

	t.Run("returns hash details", func(t *testing.T) {
		n := plan.Node{
			MemoryUsage: 12,
			HashBuckets: 1024,
			HashBatches: 1,
		}

		assert.Equal(t, "Buckets: 1024, Batches: 1, Memory Usage: 12kB", detail(n))
	})

	t.Run("returns all information if available", func(t *testing.T) {
		n := plan.Node{
			Filter:      "(id = 123)",
			BuffersHit:  8,
			BuffersRead: 5,
			MemoryUsage: 12,
			HashBuckets: 1024,
			HashBatches: 1,
		}

		expected := "Filter: (id = 123) | Buffers Shared Hit: 8, Buffers Shared Read: 5 | Buckets: 1024, Batches: 1, Memory Usage: 12kB"
		assert.Equal(t, expected, detail(n))
	})

}
