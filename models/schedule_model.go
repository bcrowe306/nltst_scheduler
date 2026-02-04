package models

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const ScheduleCollection = "schedules"

type Schedule struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Name        string
	Description string
	Events      []bson.ObjectID
}

func GetAllSchedules(db *mongo.Database) ([]Schedule, error) {
	collection := db.Collection(ScheduleCollection)
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var schedules []Schedule
	for cursor.Next(context.TODO()) {
		var schedule Schedule
		if err := cursor.Decode(&schedule); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return schedules, nil
}

func InsertSchedule(db *mongo.Database, schedule *Schedule) (*mongo.InsertOneResult, error) {
	collection := db.Collection(ScheduleCollection)
	res, err := collection.InsertOne(context.TODO(), schedule)
	return res, err
}

func UpdateSchedule(db *mongo.Database, schedule *Schedule) (*mongo.UpdateResult, error) {
	collection := db.Collection(ScheduleCollection)
	filter := bson.M{"_id": schedule.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        schedule.Name,
			"description": schedule.Description,
			"events":      schedule.Events,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func DeleteSchedule(db *mongo.Database, scheduleID bson.ObjectID) (*mongo.DeleteResult, error) {
	collection := db.Collection(ScheduleCollection)
	filter := bson.M{"_id": scheduleID}
	res, err := collection.DeleteOne(context.TODO(), filter)
	return res, err
}

func GetScheduleByID(db *mongo.Database, scheduleID bson.ObjectID) (*Schedule, error) {
	collection := db.Collection(ScheduleCollection)
	filter := bson.M{"_id": scheduleID}
	var schedule Schedule
	err := collection.FindOne(context.TODO(), filter).Decode(&schedule)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func AddEventToSchedule(db *mongo.Database, scheduleID, eventID bson.ObjectID) (*mongo.UpdateResult, error) {
	collection := db.Collection(ScheduleCollection)
	filter := bson.M{"_id": scheduleID}
	update := bson.M{
		"$addToSet": bson.M{
			"events": eventID,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func RemoveEventFromSchedule(db *mongo.Database, scheduleID, eventID bson.ObjectID) (*mongo.UpdateResult, error) {
	collection := db.Collection(ScheduleCollection)
	filter := bson.M{"_id": scheduleID}
	update := bson.M{
		"$pull": bson.M{
			"events": eventID,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func ClearEventsFromSchedule(db *mongo.Database, scheduleID bson.ObjectID) (*mongo.UpdateResult, error) {
	collection := db.Collection(ScheduleCollection)
	filter := bson.M{"_id": scheduleID}
	update := bson.M{
		"$set": bson.M{
			"events": []bson.ObjectID{},
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func GetSchedulesByEvent(db *mongo.Database, eventID bson.ObjectID) ([]Schedule, error) {
	collection := db.Collection(ScheduleCollection)
	filter := bson.M{"events": eventID}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var schedules []Schedule
	for cursor.Next(context.TODO()) {
		var schedule Schedule
		if err := cursor.Decode(&schedule); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return schedules, nil
}
