package flame

import (
	"fmt"

	"pg_flame/pkg/plan"
)

type Flame struct {
	Name     string  `json:"name"`
	Value    float64 `json:"value"`
	Children []Flame `json:"children"`
}

func New(p plan.Plan) Flame {
	// TODO add planning time frame and total

	return convert(p.Root)
}

func convert(n plan.Node) Flame {
	var subFlames []Flame

	for _, subOp := range n.Children {
		subFlames = append(subFlames, convert(subOp))
	}

	return Flame{
		Name:     name(n),
		Value:    n.TotalTime,
		Children: subFlames,
	}
}

func name(n plan.Node) string {
	if n.Table != "" {
		return fmt.Sprintf("%s on %s", n.Method, n.Table)
	}

	return n.Method
}
