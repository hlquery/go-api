package hlquery

import (
	"math"
)

var defaultRankWeights = map[string]float64{
	"popularity_log":  1.15,
	"hit_log":         0.95,
	"popularity_sqrt": 0.25,
	"hit_log_sqrt":    0.15,
}

func ComputeRankSignal(popularity, hitLog float64, weights map[string]float64) float64 {
	w := make(map[string]float64)
	for k, v := range defaultRankWeights {
		w[k] = v
	}

	for k, v := range weights {
		if _, ok := w[k]; ok {
			w[k] = v
		}
	}

	return math.Log(popularity+1)*w["popularity_log"] + math.Log(hitLog+1)*w["hit_log"] + math.Sqrt(popularity)*w["popularity_sqrt"] + math.Sqrt(hitLog)*w["hit_log_sqrt"]
}

func AttachRankSort(params map[string]interface{}, field, direction string) {
	if params == nil {
		return
	}

	dir := direction
	if dir != "asc" && dir != "desc" {
		dir = "desc"
	}

	sortInstruction := field + ":" + dir
	current, ok := params["sort_by"].(string)
	if ok && current != "" {
		params["sort_by"] = current + "," + sortInstruction
	} else {
		params["sort_by"] = sortInstruction
	}
}
