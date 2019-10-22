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

		assert.Equal(t, "Limit", p.Root.Method)
		assert.Equal(t, "", p.Root.Table)
		assert.Equal(t, 0.022, p.Root.TotalTime)

		assert.Equal(t, "Seq Scan", p.Root.Children[0].Method)
		assert.Equal(t, "bears", p.Root.Children[0].Table)
		assert.Equal(t, 0.018, p.Root.Children[0].TotalTime)
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
      "Node Type": "Limit",
      "Parallel Aware": false,
      "Startup Cost": 0.00,
      "Total Cost": 0.11,
      "Plan Rows": 1,
      "Plan Width": 32,
      "Actual Startup Time": 0.022,
      "Actual Total Time": 0.022,
      "Actual Rows": 1,
      "Actual Loops": 1,
      "Shared Hit Blocks": 1,
      "Shared Read Blocks": 0,
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
          "Node Type": "Seq Scan",
          "Parent Relationship": "Outer",
          "Parallel Aware": false,
          "Relation Name": "bears",
          "Alias": "bears",
          "Startup Cost": 0.00,
          "Total Cost": 11.00,
          "Plan Rows": 100,
          "Plan Width": 32,
          "Actual Startup Time": 0.018,
          "Actual Total Time": 0.018,
          "Actual Rows": 1,
          "Actual Loops": 1,
          "Shared Hit Blocks": 1,
          "Shared Read Blocks": 0,
          "Shared Dirtied Blocks": 0,
          "Shared Written Blocks": 0,
          "Local Hit Blocks": 0,
          "Local Read Blocks": 0,
          "Local Dirtied Blocks": 0,
          "Local Written Blocks": 0,
          "Temp Read Blocks": 0,
          "Temp Written Blocks": 0
        }
      ]
    },
    "Planning Time": 1.756,
    "Triggers": [
    ],
    "Execution Time": 0.059
  }
]
`
