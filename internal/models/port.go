package models

type Port struct {
	id          string        `json:"name"`
	coordinates []float32     `json:"coordinates"`
	city        string        `json:"city"`
	province    string        `json:"province"`
	country     string        `json:"country"`
	alias       []interface{} `json:"alias"`
	regions     []interface{} `json:"regions"`
	timezone    string        `json:"timezone"`
	unlocs      []string      `json:"unlocs"`
	code        string        `json:"code"`
}
