package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const EventCollection = "events"

type Event struct {
	ID                  string               `bson:"_id,omitempty" json:"_id" query:"_id" form:"_id"`
	Name                string               `bson:"name" json:"name" query:"name" form:"name"`
	Description         string               `bson:"description" json:"description" query:"description" form:"description"`
	Template            string               `bson:"template" json:"template" query:"template" form:"template"`
	StartTime           string               `bson:"startTime" json:"startTime" query:"startTime" form:"startTime"`
	EndTime             string               `bson:"endTime" json:"endTime" query:"endTime" form:"endTime"`
	Date                time.Time            `bson:"date" json:"date" query:"date" form:"date"`
	TeamID              string               `bson:"teamId" json:"teamId" query:"teamId" form:"teamId"`
	PositionAssignments []PositionAssignment `bson:"positionAssignments" json:"positionAssignments" query:"positionAssignments" form:"positionAssignments"`
}

type PositionAssignment struct {
	PositionName string `bson:"positionName" json:"positionName" query:"positionName" form:"positionName"`
	Member       string `bson:"member" json:"member" query:"member" form:"member"`
}

func GetAllEvents(db *mongo.Database) ([]Event, error) {
	collection := db.Collection(EventCollection)
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var events []Event
	for cursor.Next(context.TODO()) {
		var event Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func InsertEvent(db *mongo.Database, event *Event) (*mongo.InsertOneResult, error) {
	collection := db.Collection(EventCollection)
	res, err := collection.InsertOne(context.TODO(), event)
	return res, err
}

func UpdateEvent(db *mongo.Database, event *Event) (*mongo.UpdateResult, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"_id": event.ID}
	update := bson.M{
		"$set": bson.M{
			"name":                event.Name,
			"description":         event.Description,
			"template":            event.Template,
			"startTime":           event.StartTime,
			"endTime":             event.EndTime,
			"positionAssignments": event.PositionAssignments,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func DeleteEvent(db *mongo.Database, eventID string) (*mongo.DeleteResult, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"_id": eventID}
	res, err := collection.DeleteOne(context.TODO(), filter)
	return res, err
}

func GetEventsByService(db *mongo.Database, templateID string) ([]Event, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"template": templateID}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var events []Event
	for cursor.Next(context.TODO()) {
		var event Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func GetEventByID(db *mongo.Database, eventID string) (*Event, error) {
	collection := db.Collection(EventCollection)
	var event Event
	err := collection.FindOne(context.TODO(), bson.M{"_id": eventID}).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func GetEventsByMember(db *mongo.Database, memberID string) ([]Event, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"positionAssignments.member": memberID}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var events []Event
	for cursor.Next(context.TODO()) {
		var event Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func DeleteEventsByService(db *mongo.Database, templateID string) (*mongo.DeleteResult, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"template": templateID}
	res, err := collection.DeleteMany(context.TODO(), filter)
	return res, err
}

func AddPositionAssignment(db *mongo.Database, eventID string, assignment PositionAssignment) (*mongo.UpdateResult, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"_id": eventID}
	update := bson.M{
		"$addToSet": bson.M{
			"positionAssignments": assignment,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func RemovePositionAssignment(db *mongo.Database, eventID string, positionName string) (*mongo.UpdateResult, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"_id": eventID}
	update := bson.M{
		"$pull": bson.M{
			"positionAssignments": bson.M{"positionName": positionName},
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}
