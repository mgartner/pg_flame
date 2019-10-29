package plan

import (
	"encoding/json"
	"errors"
	"io"
)

type Plan struct {
	PlanningTime  float64 `json:"Planning Time"`
	ExecutionTree Node    `json:"Plan"`
}

type Node struct {
	Method             string   `json:"Node Type"`
	Table              string   `json:"Relation Name"`
	Alias              string   `json:"Alias"`
	Index              string   `json:"Index Name"`
	ParentRelationship string   `json:"Parent Relationship"`
	PlanCost           float64  `json:"Total Cost"`
	PlanRows           int      `json:"Plan Rows"`
	PlanWidth          int      `json:"Plan Width"`
	ActualTotalTime    float64  `json:"Actual Total Time"`
	ActualRows         int      `json:"Actual Rows"`
	ActualLoops        int      `json:"Actual Loops"`
	Filter             string   `json:"Filter"`
	JoinFilter         string   `json:"Join Filter"`
	HashCond           string   `json:"Hash Cond"`
	IndexCond          string   `json:"Index Cond"`
	RecheckCond        string   `json:"Recheck Cond"`
	BuffersHit         int      `json:"Shared Hit Blocks"`
	BuffersRead        int      `json:"Shared Read Blocks"`
	MemoryUsage        int      `json:"Peak Memory Usage"`
	HashBuckets        int      `json:"Hash Buckets"`
	HashBatches        int      `json:"Hash Batches"`
	SortKey            []string `json:"Sort Key"`
	SortMethod         string   `json:"Sort Method"`
	SortSpaceUsed      int      `json:"Sort Space Used"`
	SortSpaceType      string   `json:"Sort Space Type"`
	Children           []Node   `json:"Plans"`
}

var ErrEmptyPlanJSON = errors.New("empty plan JSON")
var ErrInvalidPlanJSON = errors.New("invalid plan JSON")

func New(r io.Reader) (Plan, error) {
	var plans []Plan

	err := json.NewDecoder(r).Decode(&plans)
	if err != nil {
		var e *json.UnmarshalTypeError
		if errors.As(err, &e) {
			err = ErrInvalidPlanJSON
		}
		return Plan{}, err
	}

	if len(plans) < 1 {
		return Plan{}, ErrEmptyPlanJSON
	}

	return plans[0], nil
}
