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
	Method             string  `json:"Node Type"`
	Table              string  `json:"Relation Name"`
	Index              string  `json:"Index Name"`
	ParentRelationship string  `json:"Parent Relationship"`
	Filter             string  `json:"Filter"`
	JoinFilter         string  `json:"Join Filter"`
	HashCond           string  `json:"Hash Cond"`
	IndexCond          string  `json:"Index Cond"`
	RecheckCond        string  `json:"Recheck Cond"`
	BuffersHit         int     `json:"Shared Hit Blocks"`
	BuffersRead        int     `json:"Shared Read Blocks"`
	MemoryUsage        int     `json:"Peak Memory Usage"`
	HashBuckets        int     `json:"Hash Buckets"`
	HashBatches        int     `json:"Hash Batches"`
	TotalTime          float64 `json:"Actual Total Time"`
	Children           []Node  `json:"Plans"`
}

var ErrEmptyPlanJSON = errors.New("empty plan JSON")
var ErrInvalidPlanJSON = errors.New("invalid plan JSON")

func New(r io.Reader) (error, Plan) {
	var plans []Plan

	err := json.NewDecoder(r).Decode(&plans)
	if err != nil {
		var e *json.UnmarshalTypeError
		if errors.As(err, &e) {
			err = ErrInvalidPlanJSON
		}
		return err, Plan{}
	}

	if len(plans) < 1 {
		return ErrEmptyPlanJSON, Plan{}
	}

	return nil, plans[0]
}
