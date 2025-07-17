package status_page

import (
	"context"
	"errors"
	"peekaping/src/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoModel struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty"`
	Slug                 string             `bson:"slug"`
	Title                string             `bson:"title"`
	Description          string             `bson:"description"`
	Icon                 string             `bson:"icon"`
	Theme                string             `bson:"theme"`
	Published            bool               `bson:"published"`
	SearchEngineIndex    bool               `bson:"search_engine_index"`
	Password             string             `bson:"password,omitempty"`
	FooterText           string             `bson:"footer_text"`
	GoogleAnalyticsTagID string             `bson:"google_analytics_tag_id"`
	AutoRefreshInterval  int                `bson:"auto_refresh_interval"`
	Domain               string             `bson:"domain,omitempty"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func toDomainModel(m *mongoModel) *Model {
	return &Model{
		ID:                  m.ID.Hex(),
		Slug:                m.Slug,
		Title:               m.Title,
		Description:         m.Description,
		Icon:                m.Icon,
		Theme:               m.Theme,
		Published:           m.Published,
		FooterText:          m.FooterText,
		AutoRefreshInterval: m.AutoRefreshInterval,
		Domain:              m.Domain,

		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

type MongoRepository struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func NewMongoRepository(client *mongo.Client, cfg *config.Config) Repository {
	db := client.Database(cfg.DBName)
	collection := db.Collection("status_pages")

	// Create indexes
	go func() {
		_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
		if err != nil {
			// Handle error appropriately, e.g., log it
		}
	}()

	return &MongoRepository{
		client:     client,
		db:         db,
		collection: collection,
	}
}

func (r *MongoRepository) Create(ctx context.Context, statusPage *Model) (*Model, error) {
	mm := &mongoModel{
		ID:                  primitive.NewObjectID(),
		Slug:                statusPage.Slug,
		Title:               statusPage.Title,
		Description:         statusPage.Description,
		Icon:                statusPage.Icon,
		Theme:               statusPage.Theme,
		Published:           statusPage.Published,
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
		FooterText:          statusPage.FooterText,
		AutoRefreshInterval: statusPage.AutoRefreshInterval,
		Domain:              statusPage.Domain,
	}

	_, err := r.collection.InsertOne(ctx, mm)
	if err != nil {
		return nil, err
	}

	return toDomainModel(mm), nil
}

func (r *MongoRepository) FindByID(ctx context.Context, id string) (*Model, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var mm mongoModel
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&mm)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return toDomainModel(&mm), nil
}

func (r *MongoRepository) FindBySlug(ctx context.Context, slug string) (*Model, error) {
	var mm mongoModel
	err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&mm)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return toDomainModel(&mm), nil
}

func (r *MongoRepository) FindByDomain(ctx context.Context, domain string) (*Model, error) {
	var mm mongoModel
	err := r.collection.FindOne(ctx, bson.M{"domain": domain}).Decode(&mm)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return toDomainModel(&mm), nil
}

func (r *MongoRepository) FindAll(ctx context.Context, page int, limit int, q string) ([]*Model, error) {
	skip := int64(page * limit)
	limit64 := int64(limit)

	opts := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit64,
		Sort:  bson.D{{Key: "created_at", Value: -1}},
	}

	filter := bson.M{}
	if q != "" {
		filter["title"] = bson.M{"$regex": q, "$options": "i"}
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoStatusPages []*mongoModel
	if err = cursor.All(ctx, &mongoStatusPages); err != nil {
		return nil, err
	}

	var statusPages []*Model
	for _, msp := range mongoStatusPages {
		statusPages = append(statusPages, toDomainModel(msp))
	}

	return statusPages, nil
}

func (r *MongoRepository) Update(ctx context.Context, id string, statusPage *UpdateModel) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updatePayload := bson.M{}

	if statusPage.Slug != nil {
		updatePayload["slug"] = *statusPage.Slug
	}
	if statusPage.Title != nil {
		updatePayload["title"] = *statusPage.Title
	}
	if statusPage.Description != nil {
		updatePayload["description"] = *statusPage.Description
	}
	if statusPage.Icon != nil {
		updatePayload["icon"] = *statusPage.Icon
	}
	if statusPage.Theme != nil {
		updatePayload["theme"] = *statusPage.Theme
	}
	if statusPage.Published != nil {
		updatePayload["published"] = *statusPage.Published
	}
	if statusPage.FooterText != nil {
		updatePayload["footer_text"] = *statusPage.FooterText
	}
	if statusPage.AutoRefreshInterval != nil {
		updatePayload["auto_refresh_interval"] = *statusPage.AutoRefreshInterval
	}
	if statusPage.Domain != nil {
		updatePayload["domain"] = *statusPage.Domain
	}

	if len(updatePayload) == 0 {
		return nil // nothing to update
	}

	updatePayload["updated_at"] = time.Now().UTC()

	update := bson.M{
		"$set": updatePayload,
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *MongoRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
