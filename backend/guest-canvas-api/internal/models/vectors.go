package models

type VectorData struct {
	Key         string      `json:"key"`
	VectorID    string      `json:"vector_id,omitempty"`
	Points      [][]float64 `json:"points"`
	StrokeWidth int         `json:"stroke_width"`
	StrokeColor string      `json:"stroke_color"`
}
