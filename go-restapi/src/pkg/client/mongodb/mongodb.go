package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewClient(ctx context.Context, host, port, username, password, database, authDb string) (db *mongo.Database, err error) {
	// Connect
	// Ping

	return nil, nil
}
