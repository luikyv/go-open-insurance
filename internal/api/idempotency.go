package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var errIdempotencyNotFound = errors.New("idempotency not found")

type IdempotencyService struct {
	storage IdempotencyStorage
}

func NewIdempotencyService(storage IdempotencyStorage) IdempotencyService {
	return IdempotencyService{
		storage: storage,
	}
}

func (s IdempotencyService) FetchIdempotencyResponse(
	ctx context.Context,
	id string,
	req any,
) (
	resp string,
	err error,
) {
	record, err := s.storage.record(ctx, id)
	if err != nil {
		Logger(ctx).Debug(
			"error fetching the idempotency record",
			slog.String("error", err.Error()),
		)
		return "", err
	}

	payload, err := json.Marshal(req)
	if err != nil {
		Logger(ctx).Error(
			"could not marshal the request",
			slog.String("error", err.Error()),
			slog.Any("request", req),
		)
		return "", ErrInternal
	}

	if string(payload) != record.Request {
		Logger(ctx).Debug("requested payload doesn't match the previous one sent for idempotency")
		return "", errors.New("requested payload doesn't match the previous one sent for idempotency")
	}

	return record.Response, nil
}

func (s IdempotencyService) CreateIdempotency(
	ctx context.Context,
	id string,
	req any,
	resp any,
) error {
	reqPayload, err := json.Marshal(req)
	if err != nil {
		Logger(ctx).Error(
			"could not marshal the request",
			slog.String("error", err.Error()),
			slog.Any("request", resp),
		)
		return err
	}
	respPayload, err := json.Marshal(resp)
	if err != nil {
		Logger(ctx).Error(
			"could not marshal the response",
			slog.String("error", err.Error()),
			slog.Any("request", resp),
		)
		return err
	}

	Logger(ctx).Info("requested payload doesn't match the previous one sent for idempotency")
	if err := s.storage.save(ctx, idempotencyRecord{
		ID:       id,
		Request:  string(reqPayload),
		Response: string(respPayload),
	}); err != nil {
		Logger(ctx).Error(
			"could not save the idempotency record",
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}

type idempotencyRecord struct {
	ID       string `bson:"_id"`
	Request  string `bson:"request"`
	Response string `bson:"response"`
}

type IdempotencyStorage struct {
	collection *mongo.Collection
}

func NewIdempotencyStorage(db *mongo.Database) IdempotencyStorage {
	return IdempotencyStorage{
		collection: db.Collection("idempotency"),
	}
}

func (s IdempotencyStorage) save(ctx context.Context, record idempotencyRecord) error {
	shouldUpsert := true
	filter := bson.D{{Key: "_id", Value: record.ID}}
	if _, err := s.collection.ReplaceOne(
		ctx,
		filter,
		record,
		&options.ReplaceOptions{Upsert: &shouldUpsert},
	); err != nil {
		return err
	}

	return nil
}

func (s IdempotencyStorage) record(
	ctx context.Context,
	id string,
) (
	idempotencyRecord,
	error,
) {
	filter := bson.D{{Key: "_id", Value: id}}

	result := s.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return idempotencyRecord{}, errIdempotencyNotFound
		}
		return idempotencyRecord{}, result.Err()
	}

	var record idempotencyRecord
	if err := result.Decode(&record); err != nil {
		return idempotencyRecord{}, err
	}

	return record, nil
}
