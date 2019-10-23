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
	Children []Flame `json:"children"`
}

func New(p plan.Plan) Flame {
	// TODO add planning time frame and total
	// TODO handle CTE InitPlan

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
	condWords := make([]string, 0, 5)
	if n.Filter != "" {
		condWords = append(condWords, fmt.Sprintf("Filter: %s", n.Filter))
	}
	if n.JoinFilter != "" {
		condWords = append(condWords, fmt.Sprintf("Join Filter: %s", n.JoinFilter))
	}
	if n.HashCond != "" {
		condWords = append(condWords, fmt.Sprintf("Hash Cond: %s", n.HashCond))
	}
	if n.IndexCond != "" {
		condWords = append(condWords, fmt.Sprintf("Index Cond: %s", n.IndexCond))
	}
	if n.RecheckCond != "" {
		condWords = append(condWords, fmt.Sprintf("Recheck Cond: %s", n.RecheckCond))
	}
	cond := strings.Join(condWords, ", ")

	bufferWords := make([]string, 0, 2)
	if n.BuffersHit != 0 {
		bufferWords = append(bufferWords, fmt.Sprintf("Buffers Shared Hit: %v", n.BuffersHit))
	}
	if n.BuffersRead != 0 {
		bufferWords = append(bufferWords, fmt.Sprintf("Buffers Shared Read: %v", n.BuffersRead))
	}
	buffer := strings.Join(bufferWords, ", ")

	hashWords := make([]string, 0, 3)
	if n.HashBuckets != 0 {
		hashWords = append(hashWords, fmt.Sprintf("Buckets: %v", n.HashBuckets))
	}
	if n.HashBatches != 0 {
		hashWords = append(hashWords, fmt.Sprintf("Batches: %v", n.HashBatches))
	}
	if n.MemoryUsage != 0 {
		hashWords = append(hashWords, fmt.Sprintf("Memory Usage: %vkB", n.MemoryUsage))
	}
	hash := strings.Join(hashWords, ", ")

	sections := make([]string, 0, 3)
	if cond != "" {
		sections = append(sections, cond)
	}
	if buffer != "" {
		sections = append(sections, buffer)
	}
	if hash != "" {
		sections = append(sections, hash)
	}

	return strings.Join(sections, " | ")
}
