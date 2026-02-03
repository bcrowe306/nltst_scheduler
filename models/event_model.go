package models

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const EventCollection = "events"

//   Event:
//     id: objectId
//     service: objectId
//     startTime: datetime
//     endTime: datetime
//     positionAssignments:
//       - positionName: string
//         member: objectId

type Event struct {
	ID                  bson.ObjectID `bson:"_id,omitempty"`
	Service             bson.ObjectID
	StartTime           string
	EndTime             string
	PositionAssignments []PositionAssignment
}

type PositionAssignment struct {
	PositionName string
	Member       bson.ObjectID
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
			"service":             event.Service,
			"startTime":           event.StartTime,
			"endTime":             event.EndTime,
			"positionAssignments": event.PositionAssignments,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func DeleteEvent(db *mongo.Database, eventID bson.ObjectID) (*mongo.DeleteResult, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"_id": eventID}
	res, err := collection.DeleteOne(context.TODO(), filter)
	return res, err
}

func GetEventsByService(db *mongo.Database, serviceID bson.ObjectID) ([]Event, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"service": serviceID}
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

func GetEventByID(db *mongo.Database, eventID bson.ObjectID) (*Event, error) {
	collection := db.Collection(EventCollection)
	var event Event
	err := collection.FindOne(context.TODO(), bson.M{"_id": eventID}).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func GetEventsByMember(db *mongo.Database, memberID bson.ObjectID) ([]Event, error) {
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

func DeleteEventsByService(db *mongo.Database, serviceID bson.ObjectID) (*mongo.DeleteResult, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"service": serviceID}
	res, err := collection.DeleteMany(context.TODO(), filter)
	return res, err
}

func AddPositionAssignment(db *mongo.Database, eventID bson.ObjectID, assignment PositionAssignment) (*mongo.UpdateResult, error) {
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

func RemovePositionAssignment(db *mongo.Database, eventID bson.ObjectID, positionName string) (*mongo.UpdateResult, error) {
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
