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
	f, err := buildFlame(p)
	if err != nil {
		return err
	}

	err = templateHTML.Execute(w, f)
	if err != nil {
		return err
	}

	return nil
}

func buildFlame(p plan.Plan) (Flame, error) {
	planningFlame := Flame{
		Name:   "Query Planning",
		Value:  p.PlanningTime,
		Time:   p.PlanningTime,
		Detail: fmt.Sprintf(detailSpan, "Time to generate the query plan"),
		Color:  colorPlan,
	}

	executionFlame, err := convertPlanNode(p.ExecutionTree, "")
	if err != nil {
		return Flame{}, err
	}

	return Flame{
		Name:     "Total",
		Value:    planningFlame.Value + executionFlame.Value,
		Time:     planningFlame.Time + executionFlame.Time,
		Detail:   fmt.Sprintf(detailSpan, "Includes planning and execution time"),
		Children: []Flame{planningFlame, executionFlame},
	}, nil
}

func convertPlanNode(n plan.Node, color string) (Flame, error) {
	initPlan := n.ParentRelationship == "InitPlan"
	value := n.ActualTotalTime

	if initPlan {
		color = colorInit
	}

	var childFlames []Flame
	for _, childNode := range n.Children {

		// Pass the color forward for grey InitPlan trees
		f, err := convertPlanNode(childNode, color)
		if err != nil {
			return Flame{}, err
		}

		// Add to the total value if the child is an InitPlan node
		if f.InitPlan {
			value += f.Value
		}

		childFlames = append(childFlames, f)
	}

	d, err := detail(n)
	if err != nil {
		return Flame{}, err
	}

	return Flame{
		Name:     name(n),
		Value:    value,
		Time:     n.ActualTotalTime,
		Detail:   d,
		Color:    color,
		InitPlan: initPlan,
		Children: childFlames,
	}, nil
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

func detail(n plan.Node) (string, error) {
	var b bytes.Buffer

	err := templateTable.Execute(&b, n)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
