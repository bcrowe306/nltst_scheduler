package models

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const EventCollection = "events"

// Reminder interval enuns 15min, 30min, 1hr, 2hr, 3hr,	 1day, 2day, 1week. in minutes as time.Duration
const (
	Reminder15Min = 15 * time.Minute
	Reminder30Min = 30 * time.Minute
	Reminder1Hr   = 1 * time.Hour
	Reminder2Hr   = 2 * time.Hour
	Reminder3Hr   = 3 * time.Hour
	Reminder1Day  = 24 * time.Hour
	Reminder2Day  = 48 * time.Hour
	Reminder1Week = 7 * 24 * time.Hour
)

type Event struct {
	ID                  string               `bson:"_id,omitempty" json:"_id" query:"_id" form:"_id"`
	Name                string               `bson:"name" json:"name" query:"name" form:"name"`
	Description         string               `bson:"description" json:"description" query:"description" form:"description"`
	Template            string               `bson:"template" json:"template" query:"template" form:"template"`
	StartTime           string               `bson:"startTime" json:"startTime" query:"startTime" form:"startTime"`
	EndTime             string               `bson:"endTime" json:"endTime" query:"endTime" form:"endTime"`
	Date                time.Time            `bson:"date" json:"date" query:"date" form:"date"`
	ReminderInterval    time.Duration        `bson:"reminderInterval" json:"reminderInterval" query:"reminderInterval" form:"reminderInterval"`
	ReminderEnabled     bool                 `bson:"reminderEnabled" json:"reminderEnabled" query:"reminderEnabled" form:"reminderEnabled"`
	TeamID              string               `bson:"teamId" json:"teamId" query:"teamId" form:"teamId"`
	PositionAssignments []PositionAssignment `bson:"positionAssignments,omitempty" json:"positionAssignments" query:"positionAssignments" form:"positionAssignments"`
	CreatedAt           time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt           time.Time            `json:"updatedAt" bson:"updatedAt"`
}

type PositionAssignment struct {
	ID           string   `bson:"_id,omitempty" json:"_id" query:"_id" form:"_id"`
	PositionName string   `bson:"positionName" json:"positionName" query:"positionName" form:"positionName"`
	Description  string   `bson:"description" json:"description" query:"description" form:"description"`
	MemberID     string   `bson:"memberId" json:"memberId" query:"memberId" form:"memberId"`
	Members      []string `bson:"members,omitempty" json:"members" query:"members" form:"members"`
}

type PositionAssignmentWithMember struct {
	ID           string   `bson:"_id,omitempty" json:"_id" query:"_id" form:"_id"`
	PositionName string   `bson:"positionName" json:"positionName" query:"positionName" form:"positionName"`
	Description  string   `bson:"description" json:"description" query:"description" form:"description"`
	Member       Member   `bson:"member" json:"member" query:"member" form:"member"`
	Members      []Member `bson:"members,omitempty" json:"members,omitempty" query:"members,omitempty" form:"members,omitempty"`
}

type EventWithMemberDetails struct {
	ID                  string                         `bson:"_id,omitempty" json:"_id" query:"_id" form:"_id"`
	Name                string                         `bson:"name" json:"name" query:"name" form:"name"`
	Description         string                         `bson:"description" json:"description" query:"description" form:"description"`
	Template            string                         `bson:"template" json:"template" query:"template" form:"template"`
	StartTime           string                         `bson:"startTime" json:"startTime" query:"startTime" form:"startTime"`
	EndTime             string                         `bson:"endTime" json:"endTime" query:"endTime" form:"endTime"`
	Date                time.Time                      `bson:"date" json:"date" query:"date" form:"date"`
	ReminderInterval    time.Duration                  `bson:"reminderInterval" json:"reminderInterval" query:"reminderInterval" form:"reminderInterval"`
	ReminderEnabled     bool                           `bson:"reminderEnabled" json:"reminderEnabled" query:"reminderEnabled" form:"reminderEnabled"`
	TeamID              string                         `bson:"teamId" json:"teamId" query:"teamId" form:"teamId"`
	CreatedAt           time.Time                      `json:"createdAt" bson:"createdAt"`
	UpdatedAt           time.Time                      `json:"updatedAt" bson:"updatedAt"`
	PositionAssignments []PositionAssignmentWithMember `bson:"positionAssignments" json:"positionAssignments" query:"positionAssignments" form:"positionAssignments"`
}

type EventByPosition struct {
	ID          string    `bson:"_id,omitempty" json:"_id" query:"_id" form:"_id"`
	Name        string    `bson:"name" json:"name" query:"name" form:"name"`
	Description string    `bson:"description" json:"description" query:"description" form:"description"`
	Template    string    `bson:"template" json:"template" query:"template" form:"template"`
	StartTime   string    `bson:"startTime" json:"startTime" query:"startTime" form:"startTime"`
	EndTime     string    `bson:"endTime" json:"endTime" query:"endTime" form:"endTime"`
	Date        time.Time `bson:"date" json:"date" query:"date" form:"date"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
	Member      Member    `bson:"member" json:"member" query:"member" form:"member"`
}

type EventsGroupedByPositionName struct {
	PositionName string            `bson:"_id" json:"_id" query:"_id" form:"_id"`
	Events       []EventByPosition `bson:"events" json:"events" query:"events" form:"events"`
}

// Types for member schedule view
type MemberScheduleEvent struct {
	ID               string        `bson:"_id,omitempty" json:"_id" query:"_id" form:"_id"`
	Name             string        `bson:"name" json:"name" query:"name" form:"name"`
	Description      string        `bson:"description" json:"description" query:"description" form:"description"`
	Template         string        `bson:"template" json:"template" query:"template" form:"template"`
	StartTime        string        `bson:"startTime" json:"startTime" query:"startTime" form:"startTime"`
	EndTime          string        `bson:"endTime" json:"endTime" query:"endTime" form:"endTime"`
	Date             time.Time     `bson:"date" json:"date" query:"date" form:"date"`
	ReminderInterval time.Duration `bson:"reminderInterval" json:"reminderInterval" query:"reminderInterval" form:"reminderInterval"`
	ReminderEnabled  bool          `bson:"reminderEnabled" json:"reminderEnabled" query:"reminderEnabled" form:"reminderEnabled"`
	TeamID           string        `bson:"teamId" json:"teamId" query:"teamId" form:"teamId"`
	PositionName     string        `bson:"positionName,omitempty" json:"positionName" query:"positionName" form:"positionName"`
	CreatedAt        time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt" bson:"updatedAt"`
}

type MemberSchedule struct {
	ID          string                `bson:"_id" json:"_id"`
	FirstName   string                `json:"firstName" bson:"firstName" query:"firstName" form:"firstName"`
	LastName    string                `json:"lastName" bson:"lastName" query:"lastName" form:"lastName"`
	Email       string                `json:"email" bson:"email" query:"email" form:"email"`
	PhoneNumber string                `json:"phoneNumber" bson:"phoneNumber" query:"phoneNumber" form:"phoneNumber"`
	CreatedAt   time.Time             `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt" bson:"updatedAt"`
	Events      []MemberScheduleEvent `bson:"events" json:"events" query:"events" form:"events"`
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

func GetEventsSortedByDate(db *mongo.Database) ([]Event, error) {
	collection := db.Collection(EventCollection)
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: 1}})
	cursor, err := collection.Find(context.TODO(), bson.M{}, opts)
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
	event.ID = uuid.NewString()
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()
	event.PositionAssignments = make([]PositionAssignment, 0)
	res, err := collection.InsertOne(context.TODO(), event)
	return res, err
}

func UpdateEvent(db *mongo.Database, event *Event) (*mongo.UpdateResult, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"_id": event.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        event.Name,
			"description": event.Description,
			"template":    event.Template,
			"startTime":   event.StartTime,
			"endTime":     event.EndTime,
			"date":        event.Date,
			"teamId":      event.TeamID,
			"updatedAt":   time.Now(),
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

func AddPosition(db *mongo.Database, eventID string, position PositionAssignment) (*mongo.UpdateResult, error) {
	collection := db.Collection(EventCollection)
	position.ID = uuid.NewString()
	filter := bson.M{"_id": eventID}
	update := bson.M{
		"$addToSet": bson.M{
			"positionAssignments": position,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func RemovePosition(db *mongo.Database, eventID string, positionID string) (*mongo.UpdateResult, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"_id": eventID}
	update := bson.M{
		"$pull": bson.M{
			"positionAssignments": bson.M{"_id": positionID},
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func AssignPositionToMember(db *mongo.Database, eventID string, positionID string, memberID string) (*mongo.UpdateResult, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"_id": eventID, "positionAssignments._id": positionID}
	update := bson.M{
		"$set": bson.M{
			"positionAssignments.$.memberId": memberID,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func UnassignPositionFromMember(db *mongo.Database, eventID string, positionID string) (*mongo.UpdateResult, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"_id": eventID, "positionAssignments._id": positionID}
	update := bson.M{
		"$set": bson.M{
			"positionAssignments.$.memberId": "",
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
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
	filter := bson.M{"positionAssignments.memberId": memberID}
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

func RemovePositionAssignment(db *mongo.Database, eventID string, positionID string) (*mongo.UpdateResult, error) {
	collection := db.Collection(EventCollection)
	filter := bson.M{"_id": eventID}
	update := bson.M{
		"$pull": bson.M{
			"positionAssignments": bson.M{"_id": positionID},
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

// Special queries and aggregations

func GetEventsWithMemberDetails(db *mongo.Database) ([]EventWithMemberDetails, error) {
	collection := db.Collection(EventCollection)
	pipeline := mongo.Pipeline{
		{
			{Key: "$unwind", Value: bson.M{"path": "$positionAssignments", "preserveNullAndEmptyArrays": true}},
		},
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "members",
				"localField":   "positionAssignments.memberId",
				"foreignField": "_id",
				"as":           "memberDetails",
			}},
		},
		{
			{Key: "$unwind", Value: bson.M{"path": "$memberDetails", "preserveNullAndEmptyArrays": true}},
		},
		{
			{Key: "$group", Value: bson.M{
				"_id":         "$_id",
				"name":        bson.M{"$first": "$name"},
				"description": bson.M{"$first": "$description"},
				"teamId":      bson.M{"$first": "$teamId"},
				"date":        bson.M{"$first": "$date"},
				"startTime":   bson.M{"$first": "$startTime"},
				"endTime":     bson.M{"$first": "$endTime"},
				"createdAt":   bson.M{"$first": "$createdAt"},
				"updatedAt":   bson.M{"$first": "$updatedAt"},
				"template":    bson.M{"$first": "$template"},
				"positionAssignments": bson.M{
					"$push": bson.M{
						"_id":          "$positionAssignments._id",
						"description":  "$positionAssignments.description",
						"positionName": "$positionAssignments.positionName",
						"member":       "$memberDetails",
					},
				},
			}}},
		{{Key: "$sort", Value: bson.M{"date": 1}}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []EventWithMemberDetails
	for cursor.Next(context.TODO()) {
		var event EventWithMemberDetails
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		results = append(results, event)
	}
	if err := cursor.Err(); err != nil {
		log.Print(err)
		return nil, err
	}

	return results, nil
}

func GetEventsGroupedByPositionName(db *mongo.Database) ([]EventsGroupedByPositionName, error) {
	collection := db.Collection(EventCollection)
	pipeline := mongo.Pipeline{
		{
			{Key: "$unwind", Value: "$positionAssignments"},
		},
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "members",
				"localField":   "positionAssignments.memberId",
				"foreignField": "_id",
				"as":           "member",
			}},
		},
		{
			{Key: "$unwind", Value: bson.M{"path": "$member", "preserveNullAndEmptyArrays": true}},
		},
		{
			{Key: "$group", Value: bson.M{
				"_id": "$positionAssignments.positionName",
				"events": bson.M{
					"$push": bson.M{
						"_id":         "$_id",
						"name":        "$name",
						"description": "$description",
						"date":        "$date",
						"startTime":   "$startTime",
						"endTime":     "$endTime",
						"template":    "$template",
						"createdAt":   "$createdAt",
						"updatedAt":   "$updatedAt",
						"member":      "$member",
					},
				},
			}},
		},
		{
			{Key: "$sort", Value: bson.M{"events.date": 1, "events.startTime": 1}},
		},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []EventsGroupedByPositionName
	for cursor.Next(context.TODO()) {
		var group EventsGroupedByPositionName
		if err := cursor.Decode(&group); err != nil {
			return nil, err
		}
		results = append(results, group)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// GetMemberSchedule retrieves the schedule for a specific member within a date range.
func GetMemberSchedule(db *mongo.Database, memberID string, startDate, endDate time.Time) (*MemberSchedule, error) {
	collection := db.Collection("members")

	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.M{"_id": memberID}},
		},
		{
			{Key: "$lookup", Value: bson.M{
				"from": "events",
				"let":  bson.M{"memberId": "$_id"},
				"pipeline": mongo.Pipeline{
					{
						{Key: "$match", Value: bson.M{
							"$expr": bson.M{
								"$in": []interface{}{
									"$$memberId",
									bson.M{
										"$map": bson.M{
											"input": "$positionAssignments",
											"as":    "pa",
											"in":    "$$pa.memberId",
										},
									},
								},
							},
						}},
					},
					{
						{Key: "$match", Value: bson.M{
							"date": bson.M{"$gte": startDate, "$lte": endDate},
						}},
					},
					{
						{Key: "$addFields", Value: bson.M{
							"positionName": bson.M{
								"$let": bson.M{
									"vars": bson.M{
										"matchedPA": bson.M{
											"$arrayElemAt": []interface{}{
												bson.M{
													"$filter": bson.M{
														"input": "$positionAssignments",
														"as":    "pa",
														"cond":  bson.M{"$eq": []interface{}{"$$pa.memberId", "$$memberId"}},
													},
												},
												0,
											},
										},
									},
									"in": "$$matchedPA.positionName",
								},
							},
						}},
					},
					{
						{Key: "$project", Value: bson.M{
							"positionAssignments": 0,
						}},
					},
					{
						{Key: "$sort", Value: bson.M{"date": 1, "startTime": 1}},
					},
				},
				"as": "events",
			}},
		},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []MemberSchedule
	for cursor.Next(context.TODO()) {
		var schedule MemberSchedule
		if err := cursor.Decode(&schedule); err != nil {
			return nil, err
		}
		results = append(results, schedule)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(results) > 0 {
		return &results[0], nil
	}
	return nil, nil
}

func GetMembersSchedules(db *mongo.Database, startDate, endDate time.Time) ([]MemberSchedule, error) {
	collection := db.Collection("members")
	pipeline := mongo.Pipeline{
		{
			{Key: "$lookup", Value: bson.M{
				"from": "events",
				"let":  bson.M{"memberId": "$_id"},
				"pipeline": mongo.Pipeline{
					{
						{Key: "$match", Value: bson.M{
							"$expr": bson.M{
								"$in": []interface{}{
									"$$memberId",
									bson.M{
										"$map": bson.M{
											"input": "$positionAssignments",
											"as":    "pa",
											"in":    "$$pa.memberId",
										},
									},
								},
							},
						}},
					},
					{
						{Key: "$match", Value: bson.M{
							"date": bson.M{"$gte": startDate, "$lte": endDate},
						}},
					},
					{
						{Key: "$addFields", Value: bson.M{
							"positionName": bson.M{
								"$let": bson.M{
									"vars": bson.M{
										"matchedPA": bson.M{
											"$arrayElemAt": []interface{}{
												bson.M{
													"$filter": bson.M{
														"input": "$positionAssignments",
														"as":    "pa",
														"cond":  bson.M{"$eq": []interface{}{"$$pa.memberId", "$$memberId"}},
													},
												},
												0,
											},
										},
									},
									"in": "$$matchedPA.positionName",
								},
							},
						}},
					},
					{
						{Key: "$project", Value: bson.M{
							"positionAssignments": 0,
						}},
					},
					{
						{Key: "$sort", Value: bson.M{"date": 1, "startTime": 1}},
					},
				},
				"as": "events",
			}},
		},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []MemberSchedule
	for cursor.Next(context.TODO()) {
		var schedule MemberSchedule
		if err := cursor.Decode(&schedule); err != nil {
			return nil, err
		}
		results = append(results, schedule)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
