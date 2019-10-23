package flame

import (
	"fmt"
	"strings"

	"pg_flame/pkg/plan"
)

type Flame struct {
	Name     string  `json:"name"`
	Value    float64 `json:"value"`
	Detail   string  `json:"detail"`
	Color    string  `json:"color"`
	Children []Flame `json:"children"`
}

func New(p plan.Plan) Flame {
	// TODO handle CTE InitPlan
	planningFlame := Flame{
		Name:   "Query Planning",
		Value:  p.PlanningTime,
		Detail: "Time to generate the query plan",
		Color:  "#00C05A",
	}

	executionFlame := convert(p.ExecutionTree)

	return Flame{
		Name:     "Total",
		Value:    planningFlame.Value + executionFlame.Value,
		Detail:   "This node includes planning and execution time",
		Children: []Flame{planningFlame, executionFlame},
	}
}

func convert(n plan.Node) Flame {
	var subFlames []Flame

	for _, subOp := range n.Children {
		subFlames = append(subFlames, convert(subOp))
	}

	return Flame{
		Name:     name(n),
		Value:    n.TotalTime,
		Detail:   detail(n),
		Children: subFlames,
	}
}

func name(n plan.Node) string {
	if n.Table != "" && n.Index != "" {
		return fmt.Sprintf("%s using %s on %s", n.Method, n.Index, n.Table)
	}

	if n.Table != "" {
		return fmt.Sprintf("%s on %s", n.Method, n.Table)
	}

	return n.Method
}

func detail(n plan.Node) string {
	var b strings.Builder
	b.WriteString(`<table class="table table-striped table-bordered"><tbody>`)

	rowTemplate := "<tr><th>%s</th><td>%v</td></tr>"

	fmt.Fprintf(&b, rowTemplate, "Filter", n.Filter)
	fmt.Fprintf(&b, rowTemplate, "Join Filter", n.JoinFilter)
	fmt.Fprintf(&b, rowTemplate, "Hash Cond", n.HashCond)
	fmt.Fprintf(&b, rowTemplate, "Index Cond", n.IndexCond)
	fmt.Fprintf(&b, rowTemplate, "Recheck Cond", n.RecheckCond)
	fmt.Fprintf(&b, rowTemplate, "Buffers Shared Hit", n.BuffersHit)
	fmt.Fprintf(&b, rowTemplate, "Buffers Shared Read", n.BuffersRead)
	fmt.Fprintf(&b, rowTemplate, "Hash Buckets", n.HashBuckets)
	fmt.Fprintf(&b, rowTemplate, "Hash Batches", n.HashBatches)
	fmt.Fprintf(&b, rowTemplate, "Memory Usage", fmt.Sprintf("%vkB", n.HashBatches))

	b.WriteString(`</tbody></table>`)

	return b.String()
}
