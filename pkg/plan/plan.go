package plan

import (
	"encoding/json"
	"errors"
	"io"
)

type Plan struct {
	Root Node `json:"Plan"`
	// TODO PlanningTime float64 `json:"Planning Time"`
}

type Node struct {
	Method string `json:"Node Type"`
	Table  string `json:"Relation Name"`
	// TODO StartupTime float64     `json:"Actual Startup Time"`
	TotalTime float64 `json:"Actual Total Time"`
	Children  []Node  `json:"Plans"`
	// TODO conditionals and more information
}

var ErrEmptyPlanJSON = errors.New("empty plan JSON")
var ErrInvalidPlanJSON = errors.New("invalid plan JSON")

func New(r io.Reader) (error, Plan) {
	var plans []Plan

	err := json.NewDecoder(r).Decode(&plans)
	var e *json.UnmarshalTypeError
	if errors.As(err, &e) {
		return ErrInvalidPlanJSON, Plan{}
	} else if err != nil {
		return err, Plan{}
	}

	if len(plans) < 1 {
		return ErrEmptyPlanJSON, Plan{}
	}

	return nil, plans[0]
}
