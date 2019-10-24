package html

import (
	"bytes"
	"fmt"
	"io"

	"pg_flame/pkg/plan"
)

type Flame struct {
	Name     string  `json:"name"`
	Value    float64 `json:"value"`
	Time     float64 `json:"time"`
	Detail   string  `json:"detail"`
	Color    string  `json:"color"`
	InitPlan bool    `json:"init_plan"`
	Children []Flame `json:"children"`
}

const detailSpan = "<span>%s</span>"

const colorPlan = "#00C05A"
const colorInit = "#C0C0C0"

func Generate(w io.Writer, p plan.Plan) error {
	err, f := buildFlame(p)
	if err != nil {
		return err
	}

	err = templateHTML.Execute(w, f)
	if err != nil {
		return err
	}

	return nil
}

func buildFlame(p plan.Plan) (error, Flame) {
	planningFlame := Flame{
		Name:   "Query Planning",
		Value:  p.PlanningTime,
		Time:   p.PlanningTime,
		Detail: fmt.Sprintf(detailSpan, "Time to generate the query plan"),
		Color:  colorPlan,
	}

	err, executionFlame := convertPlanNode(p.ExecutionTree, "")
	if err != nil {
		return err, Flame{}
	}

	return nil, Flame{
		Name:     "Total",
		Value:    planningFlame.Value + executionFlame.Value,
		Time:     planningFlame.Time + executionFlame.Time,
		Detail:   fmt.Sprintf(detailSpan, "Includes planning and execution time"),
		Children: []Flame{planningFlame, executionFlame},
	}
}

func convertPlanNode(n plan.Node, color string) (error, Flame) {
	initPlan := n.ParentRelationship == "InitPlan"
	value := n.TotalTime

	if initPlan {
		color = colorInit
	}

	var childFlames []Flame
	for _, childNode := range n.Children {

		// Pass the color forward for grey InitPlan trees
		err, f := convertPlanNode(childNode, color)
		if err != nil {
			return err, Flame{}
		}

		// Add to the total value if the child is an InitPlan node
		if f.InitPlan {
			value += f.Value
		}

		childFlames = append(childFlames, f)
	}

	err, d := detail(n)
	if err != nil {
		return err, Flame{}
	}

	return nil, Flame{
		Name:     name(n),
		Value:    value,
		Time:     n.TotalTime,
		Detail:   d,
		Color:    color,
		InitPlan: initPlan,
		Children: childFlames,
	}
}

func name(n plan.Node) string {
	switch {
	case n.Table != "" && n.Index != "":
		return fmt.Sprintf("%s using %s on %s", n.Method, n.Index, n.Table)
	case n.Table != "":
		return fmt.Sprintf("%s on %s", n.Method, n.Table)
	default:
		return n.Method
	}
}

func detail(n plan.Node) (error, string) {
	var b bytes.Buffer

	err := templateTable.Execute(&b, n)
	if err != nil {
		return err, ""
	}

	return nil, b.String()
}
