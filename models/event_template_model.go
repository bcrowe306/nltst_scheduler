package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const EventTemplateCollection = "event_templates"

type Position struct {
	Name        string `bson:"name" json:"name" query:"name" form:"name"`
	Description string `bson:"description" json:"description" query:"description" form:"description"`
}

type EventTemplate struct {
	ID          string     `bson:"_id,omitempty" json:"_id" query:"_id" form:"_id"`
	Name        string     `bson:"name" json:"name" query:"name" form:"name"`
	Description string     `bson:"description" json:"description" query:"description" form:"description"`
	StartTime   string     `bson:"startTime,omitempty" json:"startTime" query:"startTime" form:"startTime"`
	EndTime     string     `bson:"endTime,omitempty" json:"endTime" query:"endTime" form:"endTime"`
	Positions   []Position `bson:"positions,omitempty" json:"positions" query:"positions" form:"positions"`
	TeamID      string     `bson:"teamId,omitempty" json:"teamId" query:"teamId" form:"teamId"`
	CreatedAt   time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `bson:"updatedAt" json:"updatedAt"`
}

type EventTemplateView struct {
	ID          string     `bson:"_id,omitempty" json:"_id" query:"_id" form:"_id"`
	Name        string     `bson:"name" json:"name" query:"name" form:"name"`
	Description string     `bson:"description" json:"description" query:"description" form:"description"`
	StartTime   string     `bson:"startTime,omitempty" json:"startTime" query:"startTime" form:"startTime"`
	EndTime     string     `bson:"endTime,omitempty" json:"endTime" query:"endTime" form:"endTime"`
	Positions   []Position `bson:"positions,omitempty" json:"positions" query:"positions" form:"positions"`
	TeamID      string     `bson:"teamId,omitempty" json:"teamId" query:"teamId" form:"teamId"`
	Team        Team       `bson:"team,omitempty" json:"team" query:"team" form:"team"`
	CreatedAt   time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `bson:"updatedAt" json:"updatedAt"`
}

func GetAllEventTemplates(db *mongo.Database) ([]EventTemplateView, error) {
	collection := db.Collection(EventTemplateCollection)
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var eventTemplates []EventTemplateView

	// build aggregation pipeline to lookup team details
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "teams"},
			{Key: "localField", Value: "teamId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "team"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$team"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}

	cursor, err = collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var eventTemplate EventTemplateView
		if err := cursor.Decode(&eventTemplate); err != nil {
			return nil, err
		}
		eventTemplates = append(eventTemplates, eventTemplate)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return eventTemplates, nil
}

func InsertEventTemplate(db *mongo.Database, eventTemplate *EventTemplate) (*mongo.InsertOneResult, error) {
	eventTemplate.ID = uuid.NewString()
	eventTemplate.CreatedAt = time.Now()
	eventTemplate.UpdatedAt = time.Now()
	eventTemplate.Positions = make([]Position, 0)
	collection := db.Collection(EventTemplateCollection)
	res, err := collection.InsertOne(context.TODO(), eventTemplate)
	return res, err
}

func UpdateEventTemplate(db *mongo.Database, eventTemplate *EventTemplate) (*mongo.UpdateResult, error) {
	collection := db.Collection(EventTemplateCollection)
	filter := bson.M{"_id": eventTemplate.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        eventTemplate.Name,
			"description": eventTemplate.Description,
			"updatedAt":   time.Now(),
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func DeleteEventTemplate(db *mongo.Database, eventTemplateID string) (*mongo.DeleteResult, error) {
	collection := db.Collection(EventTemplateCollection)
	filter := bson.M{"_id": eventTemplateID}
	res, err := collection.DeleteOne(context.TODO(), filter)
	return res, err
}

func GetEventTemplateByID(db *mongo.Database, eventTemplateID string) (*EventTemplate, error) {
	collection := db.Collection(EventTemplateCollection)
	var eventTemplate EventTemplate
	err := collection.FindOne(context.TODO(), bson.M{"_id": eventTemplateID}).Decode(&eventTemplate)
	if err != nil {
		return nil, err
	}
	return &eventTemplate, nil
}

func GetEventTemplateByName(db *mongo.Database, name string) (*EventTemplate, error) {
	collection := db.Collection(EventTemplateCollection)
	var eventTemplate EventTemplate
	err := collection.FindOne(context.TODO(), bson.M{"name": name}).Decode(&eventTemplate)
	if err != nil {
		return nil, err
	}
	return &eventTemplate, nil
}

func GetEventTemplateByIDString(db *mongo.Database, idStr string) (*EventTemplate, error) {
	return GetEventTemplateByID(db, idStr)
}

func DeleteEventTemplateByIDString(db *mongo.Database, idStr string) (*mongo.DeleteResult, error) {

	return DeleteEventTemplate(db, idStr)
}

func AddPositionToEventTemplate(db *mongo.Database, eventTemplateID string, position Position) (*mongo.UpdateResult, error) {
	collection := db.Collection(EventTemplateCollection)
	filter := bson.M{"_id": eventTemplateID}
	// Use $addToSet to avoid duplicate positions. Make sure positions is array using $ifNull if null.

	update := bson.M{
		"$addToSet": bson.M{
			"positions": position,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func RemovePositionFromEventTemplate(db *mongo.Database, eventTemplateID string, positionName string) (*mongo.UpdateResult, error) {
	collection := db.Collection(EventTemplateCollection)
	filter := bson.M{"_id": eventTemplateID}
	update := bson.M{
		"$pull": bson.M{
			"positions": bson.M{"name": positionName},
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func UpdatePositionInEventTemplate(db *mongo.Database, eventTemplateID string, position Position) (*mongo.UpdateResult, error) {
	collection := db.Collection(EventTemplateCollection)
	filter := bson.M{"_id": eventTemplateID, "positions.name": position.Name}
	update := bson.M{
		"$set": bson.M{
			"positions.$.description": position.Description,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}
