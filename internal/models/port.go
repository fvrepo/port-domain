package models

import (
	"time"
)

type Port struct {
	ID          string        `json:"name" bson:"_id"`
	Coordinates []float32     `json:"coordinates" bson:"coordinates"`
	City        string        `json:"city" bson:"city"`
	Province    string        `json:"province" bson:"province"`
	Country     string        `json:"country" bson:"country"`
	Alias       []interface{} `json:"alias" bson:"alias"`
	Regions     []interface{} `json:"regions" bson:"regions"`
	Timezone    string        `json:"timezone" bson:"timezone"`
	Unlocs      []string      `json:"unlocs" bson:"unlocs"`
	Code        string        `json:"code" bson:"code"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdateddAt  time.Time     `json:"updated_at" bson:"updated_at"`
}
