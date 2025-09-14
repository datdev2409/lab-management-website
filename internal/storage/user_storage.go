package storage

import (
	"context"
	"errors"
	"log"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrUserNotFound = errors.New("user not found")

func (m *MongoStorage) CreateUser(ctx context.Context, user *models.User) (string, error) {
	coll := m.db.Collection("users")
	return MongoInsert(ctx, coll, user)
}

func (m *MongoStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	coll := m.db.Collection("users")
	var user models.User
	err := coll.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	log.Println("GetUserByUsername error:", err)
	if err == mongo.ErrNoDocuments {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
