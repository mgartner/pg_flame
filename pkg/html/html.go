package html

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"strings"

	"pg_flame/pkg/plan"
)

type Flame struct {
	Name     string  `json:"name"`
	Value    float64 `json:"value"`
	Time     float64 `json:"time"`
	Detail   string  `json:"detail"`
	Color    string  `json:"color"`
	InitPlan bool
	Children []Flame `json:"children"`
}

const colorPlan = "#00C05A"
const colorInit = "#C0C0C0"

func Generate(w io.Writer, p plan.Plan) error {
	f := buildFlame(p)

	t, err := template.New("pg_flame").Parse(templateHTML)
	if err != nil {
		return err
	}

	flameJSON, err := json.Marshal(f)
	if err != nil {
		return err
	}

	data := struct {
		Data template.JS
	}{
		Data: template.JS(flameJSON),
	}

	err = t.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func buildFlame(p plan.Plan) Flame {
	planningFlame := Flame{
		Name:   "Query Planning",
		Value:  p.PlanningTime,
		Time:   p.PlanningTime,
		Detail: "Time to generate the query plan",
		Color:  colorPlan,
	}

	executionFlame := convertPlanNode(p.ExecutionTree, "")

	return Flame{
		Name:     "Total",
		Value:    planningFlame.Value + executionFlame.Value,
		Time:     planningFlame.Time + executionFlame.Time,
		Detail:   "This node includes planning and execution time",
		Children: []Flame{planningFlame, executionFlame},
	}
}

func convertPlanNode(n plan.Node, color string) Flame {
	initPlan := n.ParentRelationship == "InitPlan"
	value := n.TotalTime

	if initPlan {
		color = colorInit
	}

	var childFlames []Flame
	for _, childNode := range n.Children {

		// Pass the color forward for grey InitPlan trees
		f := convertPlanNode(childNode, color)

		// Add to the total value if the child is an InitPlan node
		if f.InitPlan {
			value += f.Value
		}

		childFlames = append(childFlames, f)
	}

	return Flame{
		Name:     name(n),
		Value:    value,
		Time:     n.TotalTime,
		Detail:   detail(n),
		Color:    color,
		InitPlan: initPlan,
		Children: childFlames,
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

	fmt.Fprintf(&b, rowTemplate, "Parent Relationship", n.ParentRelationship)
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
