package storage

import "go.mongodb.org/mongo-driver/mongo"

type Storage struct {
	Client *mongo.Client
}

func New(c *mongo.Client) *Storage {
	return &Storage{Client: c}
}
