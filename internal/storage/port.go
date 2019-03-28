package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/pkg/errors"
	"github.com/port-domain/internal/models"
)

const (
	database   = "admin"
	collection = "port-details"
)

func (s *Storage) InsertOrUpdatePort(ctx context.Context, port *models.Port) error {
	t := time.Now()
	port.CreatedAt = t
	port.UpdateddAt = t

	coll := s.Client.Database(database).Collection(collection)

	var p *models.Port
	err := coll.FindOne(ctx, bson.D{{Key: "_id", Value: port.ID}}).Decode(&p)
	switch errors.Cause(err) {
	case mongo.ErrNoDocuments:
		if _, err := coll.InsertOne(ctx, &port); err != nil {
			return errors.WithStack(err)
		}
	case nil:
		if _, err := coll.ReplaceOne(ctx, bson.D{{Key: "_id", Value: port.ID}}, &port); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// todo return cursor
func (s *Storage) GetPorts(ctx context.Context, limit int) ([]*models.Port, error) {
	cur, err := s.Client.Database(database).Collection(collection).Find(ctx, bson.D{}, options.Find().SetBatchSize(10).SetLimit(int64(limit)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer cur.Close(ctx)
	result := make([]*models.Port, 0)
	for cur.Next(ctx) {
		var p *models.Port
		if err := cur.Decode(&p); err != nil {
			return nil, errors.WithStack(err)
		}
		result = append(result, p)
	}
	return result, nil
}
