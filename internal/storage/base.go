package storage

import (
	"context"
	"errors"
	"math"
	"sort"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (m *MongoStorage) getCollection(collectionName string) *mongo.Collection {
	return m.db.Collection(collectionName)
}

func MongoInsert[T interface{}](ctx context.Context, col *mongo.Collection, entity T) (string, error) {
	result, err := col.InsertOne(ctx, entity)
	if err != nil {
		return "", err
	}
	switch id := result.InsertedID.(type) {
	case string:
		return id, nil
	case bson.ObjectID:
		return id.Hex(), nil
	default:
		return "", errors.New("unsupported ID type returned from MongoDB")
	}
}

func MongoGetById[T interface{}](ctx context.Context, col *mongo.Collection, id string) (*T, error) {
	var result T
	err := col.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func MongoGetByIds[T interface{}](ctx context.Context, col *mongo.Collection, ids []string) ([]*T, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []*T

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

type IDGetter interface {
	GetID() string
}

func MongoGetByIdsOrdered[T IDGetter](ctx context.Context, col *mongo.Collection, ids []string) ([]*T, error) {
	results, err := MongoGetByIds[T](ctx, col, ids)
	if err != nil {
		return nil, err
	}

	idOrderMap := make(map[string]int)
	for i, id := range ids {
		idOrderMap[id] = i
	}

	// Sort results based on the order of IDs
	sort.Slice(results, func(i, j int) bool {
		return idOrderMap[(*results[i]).GetID()] < idOrderMap[(*results[j]).GetID()]
	})

	return results, nil
}

func MongoList[T interface{}](ctx context.Context, col *mongo.Collection, filter interface{}, opts models.GenericQueryOptions) ([]*T, *models.PaginationResponse, error) {
	findOpts := BuildMongoSortAndPaginationOptions(opts)
	cursor, err := col.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	count, err := col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	pagniation := &models.PaginationResponse{
		Total:     int(count),
		Page:      opts.Page,
		PageSize:  opts.PageSize,
		TotalPage: int(math.Ceil(float64(count) / float64(opts.PageSize))),
	}

	var results []*T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, pagniation, err
	}
	return results, pagniation, nil
}

func MongoUpdateById[T interface{}](ctx context.Context, col *mongo.Collection, id string, update interface{}) error {
	updateDoc, ok := update.(bson.M)
	if !ok {
		return errors.New("update must be a bson.M type")
	}

	setDoc, ok := updateDoc["$set"]
	if !ok {
		return errors.New("$set operator is required in update")
	}

	setDocMap, ok := setDoc.(bson.M)
	if !ok {
		return errors.New("$set must be a bson.M type")
	}

	setDocMap["updated_at"] = GetCurrentTime()

	_, err := col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": setDocMap})
	return err
}

func MongoDeleteById[T interface{}](ctx context.Context, col *mongo.Collection, id string) error {
	_, err := col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
