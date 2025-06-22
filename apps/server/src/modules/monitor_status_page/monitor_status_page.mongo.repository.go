package monitor_status_page

import (
	"context"
	"peekaping/src/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoModel struct {
	ID           primitive.ObjectID `bson:"_id"`
	StatusPageID string             `bson:"status_page_id"`
	MonitorID    string             `bson:"monitor_id"`
	Order        int                `bson:"order"`
	Active       bool               `bson:"active"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

type mongoUpdateModel struct {
	StatusPageID *string    `bson:"status_page_id,omitempty"`
	MonitorID    *string    `bson:"monitor_id,omitempty"`
	Order        *int       `bson:"order,omitempty"`
	Active       *bool      `bson:"active,omitempty"`
	UpdatedAt    *time.Time `bson:"updated_at,omitempty"`
}

func toDomainModel(mm *mongoModel) *Model {
	return &Model{
		ID:           mm.ID.Hex(),
		StatusPageID: mm.StatusPageID,
		MonitorID:    mm.MonitorID,
		Order:        mm.Order,
		Active:       mm.Active,
		CreatedAt:    mm.CreatedAt,
		UpdatedAt:    mm.UpdatedAt,
	}
}

type MongoRepositoryImpl struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func NewMongoRepository(client *mongo.Client, cfg *config.Config) Repository {
	db := client.Database(cfg.DBName)
	collection := db.Collection("monitor_status_page")
	return &MongoRepositoryImpl{client, db, collection}
}

func (r *MongoRepositoryImpl) Create(ctx context.Context, entity *CreateUpdateDto) (*Model, error) {
	mm := &mongoModel{
		ID:           primitive.NewObjectID(),
		StatusPageID: entity.StatusPageID,
		MonitorID:    entity.MonitorID,
		Order:        entity.Order,
		Active:       entity.Active,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	_, err := r.collection.InsertOne(ctx, mm)
	if err != nil {
		return nil, err
	}

	return toDomainModel(mm), nil
}

func (r *MongoRepositoryImpl) FindByID(ctx context.Context, id string) (*Model, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var mm mongoModel
	err = r.collection.FindOne(ctx, filter).Decode(&mm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return toDomainModel(&mm), nil
}

func (r *MongoRepositoryImpl) FindAll(ctx context.Context, page int, limit int, q string) ([]*Model, error) {
	var entities []*Model

	// Calculate the number of documents to skip
	skip := int64((page) * limit)
	limit64 := int64(limit)

	// Define options for pagination
	options := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit64,
		Sort:  bson.D{{Key: "created_at", Value: -1}},
	}

	filter := bson.M{}
	if q != "" {
		filter["$or"] = bson.A{
			bson.M{"status_page_id": bson.M{"$regex": q, "$options": "i"}},
			bson.M{"monitor_id": bson.M{"$regex": q, "$options": "i"}},
		}
	}

	cursor, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var mm mongoModel
		if err := cursor.Decode(&mm); err != nil {
			return nil, err
		}
		entities = append(entities, toDomainModel(&mm))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *MongoRepositoryImpl) UpdateFull(ctx context.Context, id string, entity *CreateUpdateDto) (*Model, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	mm := &mongoModel{
		ID:           objectID,
		StatusPageID: entity.StatusPageID,
		MonitorID:    entity.MonitorID,
		Order:        entity.Order,
		Active:       entity.Active,
		UpdatedAt:    time.Now().UTC(),
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": mm}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return toDomainModel(mm), nil
}

func (r *MongoRepositoryImpl) UpdatePartial(ctx context.Context, id string, entity *PartialUpdateDto) (*Model, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	update := &mongoUpdateModel{
		StatusPageID: entity.StatusPageID,
		MonitorID:    entity.MonitorID,
		Order:        entity.Order,
		Active:       entity.Active,
		UpdatedAt:    &now,
	}

	filter := bson.M{"_id": objectID}
	updateDoc := bson.M{"$set": update}

	_, err = r.collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return nil, err
	}

	// Get the updated document
	var mm mongoModel
	err = r.collection.FindOne(ctx, filter).Decode(&mm)
	if err != nil {
		return nil, err
	}

	return toDomainModel(&mm), nil
}

func (r *MongoRepositoryImpl) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}

// Relationship management methods
func (r *MongoRepositoryImpl) AddMonitorToStatusPage(ctx context.Context, statusPageID, monitorID string, order int, active bool) (*Model, error) {
	// Check if the relationship already exists
	existing, err := r.FindByStatusPageAndMonitor(ctx, statusPageID, monitorID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		// Update existing relationship
		return r.UpdatePartial(ctx, existing.ID, &PartialUpdateDto{
			Order:  &order,
			Active: &active,
		})
	}

	// Create new relationship
	mm := &mongoModel{
		ID:           primitive.NewObjectID(),
		StatusPageID: statusPageID,
		MonitorID:    monitorID,
		Order:        order,
		Active:       active,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	_, err = r.collection.InsertOne(ctx, mm)
	if err != nil {
		return nil, err
	}

	return toDomainModel(mm), nil
}

func (r *MongoRepositoryImpl) RemoveMonitorFromStatusPage(ctx context.Context, statusPageID, monitorID string) error {
	filter := bson.M{
		"status_page_id": statusPageID,
		"monitor_id":     monitorID,
	}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *MongoRepositoryImpl) GetMonitorsForStatusPage(ctx context.Context, statusPageID string) ([]*Model, error) {
	filter := bson.M{"status_page_id": statusPageID}
	options := &options.FindOptions{
		Sort: bson.D{{Key: "order", Value: 1}},
	}

	cursor, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entities []*Model
	for cursor.Next(ctx) {
		var mm mongoModel
		if err := cursor.Decode(&mm); err != nil {
			return nil, err
		}
		entities = append(entities, toDomainModel(&mm))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *MongoRepositoryImpl) GetStatusPagesForMonitor(ctx context.Context, monitorID string) ([]*Model, error) {
	filter := bson.M{"monitor_id": monitorID}
	options := &options.FindOptions{
		Sort: bson.D{{Key: "order", Value: 1}},
	}

	cursor, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entities []*Model
	for cursor.Next(ctx) {
		var mm mongoModel
		if err := cursor.Decode(&mm); err != nil {
			return nil, err
		}
		entities = append(entities, toDomainModel(&mm))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *MongoRepositoryImpl) FindByStatusPageAndMonitor(ctx context.Context, statusPageID, monitorID string) (*Model, error) {
	filter := bson.M{
		"status_page_id": statusPageID,
		"monitor_id":     monitorID,
	}
	var mm mongoModel
	err := r.collection.FindOne(ctx, filter).Decode(&mm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return toDomainModel(&mm), nil
}

func (r *MongoRepositoryImpl) UpdateMonitorOrder(ctx context.Context, statusPageID, monitorID string, order int) (*Model, error) {
	filter := bson.M{
		"status_page_id": statusPageID,
		"monitor_id":     monitorID,
	}
	update := bson.M{
		"$set": bson.M{
			"order":      order,
			"updated_at": time.Now().UTC(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Get the updated document
	var mm mongoModel
	err = r.collection.FindOne(ctx, filter).Decode(&mm)
	if err != nil {
		return nil, err
	}

	return toDomainModel(&mm), nil
}

func (r *MongoRepositoryImpl) UpdateMonitorActiveStatus(ctx context.Context, statusPageID, monitorID string, active bool) (*Model, error) {
	filter := bson.M{
		"status_page_id": statusPageID,
		"monitor_id":     monitorID,
	}
	update := bson.M{
		"$set": bson.M{
			"active":     active,
			"updated_at": time.Now().UTC(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Get the updated document
	var mm mongoModel
	err = r.collection.FindOne(ctx, filter).Decode(&mm)
	if err != nil {
		return nil, err
	}

	return toDomainModel(&mm), nil
}

func (r *MongoRepositoryImpl) DeleteAllMonitorsForStatusPage(ctx context.Context, statusPageID string) error {
	filter := bson.M{"status_page_id": statusPageID}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}
