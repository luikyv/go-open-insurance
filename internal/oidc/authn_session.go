package oidc

import (
	"context"
	"errors"

	"github.com/luikyv/go-oidc/pkg/goidc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthnSessionManager struct {
	Collection *mongo.Collection
}

func NewAuthnSessionManager(database *mongo.Database) AuthnSessionManager {
	return AuthnSessionManager{
		Collection: database.Collection("authentication_sessions"),
	}
}

func (manager AuthnSessionManager) Save(
	ctx context.Context,
	session *goidc.AuthnSession,
) error {
	shouldUpsert := true
	filter := bson.D{{Key: "_id", Value: session.ID}}
	if _, err := manager.Collection.ReplaceOne(
		ctx,
		filter,
		session,
		&options.ReplaceOptions{Upsert: &shouldUpsert},
	); err != nil {
		return err
	}

	return nil
}

func (manager AuthnSessionManager) SessionByCallbackID(
	ctx context.Context,
	callbackID string,
) (
	*goidc.AuthnSession,
	error,
) {
	return manager.getWithFilter(
		ctx,
		bson.D{{Key: "callback_id", Value: callbackID}},
	)
}

func (manager AuthnSessionManager) SessionByAuthCode(
	ctx context.Context,
	authorizationCode string,
) (
	*goidc.AuthnSession,
	error,
) {
	return manager.getWithFilter(
		ctx,
		bson.D{{Key: "auth_code", Value: authorizationCode}},
	)
}

func (manager AuthnSessionManager) SessionByPushedAuthReqID(
	ctx context.Context,
	id string,
) (
	*goidc.AuthnSession,
	error,
) {
	return manager.getWithFilter(ctx, bson.D{{Key: "pushed_auth_req_id", Value: id}})
}

func (manager AuthnSessionManager) SessionByCIBAAuthID(
	ctx context.Context,
	id string,
) (
	*goidc.AuthnSession,
	error,
) {
	return nil, errors.ErrUnsupported
}

func (manager AuthnSessionManager) Delete(
	ctx context.Context,
	id string,
) error {
	filter := bson.D{{Key: "_id", Value: id}}
	if _, err := manager.Collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (manager AuthnSessionManager) getWithFilter(
	ctx context.Context,
	filter any,
) (
	*goidc.AuthnSession,
	error,
) {

	result := manager.Collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var authnSession goidc.AuthnSession
	if err := result.Decode(&authnSession); err != nil {
		return nil, err
	}

	return &authnSession, nil
}
