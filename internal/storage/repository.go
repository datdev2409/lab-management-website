package storage

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository[T any] interface {
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
	WithTimeout(duration time.Duration) (context.Context, context.CancelFunc)
	Insert(ctx context.Context, entity T) error
	List(ctx context.Context, filter interface{}, opts models.GenericQueryOptions) ([]*T, *models.PaginationResponse, error)
	GetById(ctx context.Context, id string) (*T, error)
	GetByIds(ctx context.Context, ids []string) ([]*T, error)
	UpdateById(ctx context.Context, id string, update interface{}) error
	DeleteById(ctx context.Context, id string) error
}

func (r *BaseRepository[T]) Insert(ctx context.Context, entity T) error {
	_, err := r.col.InsertOne(ctx, entity)
	return err
}

func (r *BaseRepository[T]) GetById(ctx context.Context, id string) (*T, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result T
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *BaseRepository[T]) GetByIds(ctx context.Context, ids []string) ([]*T, error) {
	objectIds := make([]bson.ObjectID, 0, len(ids))
	for _, id := range ids {
		objectId, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIds = append(objectIds, objectId)
	}
	filter := bson.M{"_id": bson.M{"$in": objectIds}}
	cursor, err := r.col.Find(ctx, filter)
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

func (r *BaseRepository[T]) List(ctx context.Context, filter interface{}, opts models.GenericQueryOptions) ([]*T, *models.PaginationResponse, error) {
	findOpts := BuildMongoSortAndPaginationOptions(opts)
	cursor, err := r.col.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	count, err := r.col.CountDocuments(ctx, filter)
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

func (r *BaseRepository[T]) UpdateById(ctx context.Context, id string, update interface{}) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

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

	_, err = r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": setDocMap})
	return err
}

func (r *BaseRepository[T]) DeleteById(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

type BaseRepository[T any] struct {
	db  *mongo.Database
	col *mongo.Collection
}

func NewBaseRepository[T any](db *mongo.Database, collectionName string) *BaseRepository[T] {
	return &BaseRepository[T]{
		db:  db,
		col: db.Collection(collectionName),
	}
}

func (r *BaseRepository[T]) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessionContext context.Context) (interface{}, error) {
		return nil, fn(sessionContext)
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

func (r *BaseRepository[T]) WithTimeout(duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), duration)
}
