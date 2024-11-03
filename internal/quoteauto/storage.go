package quoteauto

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	quoteLeadCollection *mongo.Collection
	quoteCollection     *mongo.Collection
}

func NewStorage(db *mongo.Database) Storage {
	return Storage{
		quoteLeadCollection: db.Collection("auto_quote_leads"),
		quoteCollection:     db.Collection("auto_quotes"),
	}
}

func (st Storage) saveLead(
	ctx context.Context,
	lead Lead,
) error {
	shouldUpsert := true
	filter := bson.D{{Key: "_id", Value: lead.ID}}
	if _, err := st.quoteLeadCollection.ReplaceOne(
		ctx,
		filter,
		lead,
		&options.ReplaceOptions{Upsert: &shouldUpsert},
	); err != nil {
		return err
	}

	return nil
}

func (st Storage) fetchLeadByConsentID(
	ctx context.Context,
	id string,
) (
	Lead,
	error,
) {
	filter := bson.D{{Key: "consent_id", Value: id}}

	result := st.quoteLeadCollection.FindOne(ctx, filter)
	if result.Err() != nil {
		return Lead{}, result.Err()
	}

	var lead Lead
	if err := result.Decode(&lead); err != nil {
		return Lead{}, err
	}

	return lead, nil
}

func (st Storage) saveQuote(
	ctx context.Context,
	quote Quote,
) error {
	shouldUpsert := true
	filter := bson.D{{Key: "_id", Value: quote.ID}}
	if _, err := st.quoteCollection.ReplaceOne(
		ctx,
		filter,
		quote,
		&options.ReplaceOptions{Upsert: &shouldUpsert},
	); err != nil {
		return err
	}

	return nil
}

func (st Storage) fetchQuoteByConsentID(
	ctx context.Context,
	id string,
) (
	Quote,
	error,
) {
	filter := bson.D{{Key: "consent_id", Value: id}}

	result := st.quoteCollection.FindOne(ctx, filter)
	if result.Err() != nil {
		return Quote{}, result.Err()
	}

	var quote Quote
	if err := result.Decode(&quote); err != nil {
		return Quote{}, err
	}

	return quote, nil
}
