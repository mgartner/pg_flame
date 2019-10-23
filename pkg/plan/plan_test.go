package plan

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	t.Run("decodes EXPLAIN ANALYZE plan JSON", func(t *testing.T) {
		input := strings.NewReader(planJSON)

		_, p := New(input)

		assert.Equal(t, "Nested Loop", p.ExecutionTree.Method)
		assert.Equal(t, "", p.ExecutionTree.Table)
		assert.Equal(t, 0.049, p.ExecutionTree.TotalTime)

		child := p.ExecutionTree.Children[0]

		assert.Equal(t, "Hash Join", child.Method)
		assert.Equal(t, "users", child.Table)
		assert.Equal(t, "users_pkey", child.Index)
		assert.Equal(t, "Outer", child.ParentRelationship)
		assert.Equal(t, "((title)::text ~ '.*sql.*'::text)", child.Filter)
		assert.Equal(t, "(id = 123)", child.JoinFilter)
		assert.Equal(t, "((p.user_id = c.user_id) AND (p.id = c.post_id))", child.HashCond)
		assert.Equal(t, "(id = p.user_id)", child.IndexCond)
		assert.Equal(t, "(p.user_id = 123)", child.RecheckCond)
		assert.Equal(t, 5, child.BuffersHit)
		assert.Equal(t, 1, child.BuffersRead)
		assert.Equal(t, 8, child.MemoryUsage)
		assert.Equal(t, 1024, child.HashBuckets)
		assert.Equal(t, 1, child.HashBatches)
		assert.Equal(t, 0.049, child.TotalTime)
	})

	t.Run("returns an error with empty plan JSON", func(t *testing.T) {
		input := strings.NewReader("[]")

		err, _ := New(input)

		assert.Error(t, err)
		assert.Equal(t, ErrEmptyPlanJSON, err)
	})

	t.Run("returns an error with invalid plan JSON", func(t *testing.T) {
		input := strings.NewReader("{}")

		err, _ := New(input)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPlanJSON, err)
	})

	t.Run("returns an error with invalid JSON syntax", func(t *testing.T) {
		input := strings.NewReader("[}")

		err, _ := New(input)

		assert.Error(t, err)
	})

}

const planJSON = `
[
  {
    "Plan": {
      "Node Type": "Nested Loop",
      "Parallel Aware": false,
      "Join Type": "Inner",
      "Startup Cost": 265.38,
      "Total Cost": 288.42,
      "Plan Rows": 1,
      "Plan Width": 539,
      "Actual Startup Time": 0.049,
      "Actual Total Time": 0.049,
      "Actual Rows": 0,
      "Actual Loops": 1,
      "Inner Unique": true,
      "Shared Hit Blocks": 5,
      "Shared Read Blocks": 1,
      "Shared Dirtied Blocks": 0,
      "Shared Written Blocks": 0,
      "Local Hit Blocks": 0,
      "Local Read Blocks": 0,
      "Local Dirtied Blocks": 0,
      "Local Written Blocks": 0,
      "Temp Read Blocks": 0,
      "Temp Written Blocks": 0,
      "Plans": [
        {
          "Node Type": "Hash Join",
          "Relation Name": "users",
          "Index Name": "users_pkey",
          "Parent Relationship": "Outer",
          "Parallel Aware": false,
          "Join Type": "Inner",
          "Startup Cost": 13.50,
          "Total Cost": 35.06,
          "Plan Rows": 1,
          "Plan Width": 543,
          "Actual Startup Time": 0.049,
          "Actual Total Time": 0.049,
          "Actual Rows": 0,
          "Actual Loops": 1,
          "Inner Unique": false,
          "Filter": "((title)::text ~ '.*sql.*'::text)",
          "Hash Cond": "((p.user_id = c.user_id) AND (p.id = c.post_id))",
          "Index Cond": "(id = p.user_id)",
          "Join Filter": "(id = 123)",
          "Recheck Cond": "(p.user_id = 123)",
          "Hash Buckets": 1024,
          "Hash Batches": 1,
          "Peak Memory Usage": 8,
          "Shared Hit Blocks": 5,
          "Shared Read Blocks": 1,
          "Shared Dirtied Blocks": 0,
          "Shared Written Blocks": 0,
          "Local Hit Blocks": 0,
          "Local Read Blocks": 0,
          "Local Dirtied Blocks": 0,
          "Local Written Blocks": 0,
          "Temp Read Blocks": 0,
          "Temp Written Blocks": 0,
          "Plans": [
          ]
        }
      ]
    },
    "Planning Time": 2.523,
    "Triggers": [
    ],
    "Execution Time": 0.221
  }
]

`
