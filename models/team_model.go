package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const TeamCollection = "teams"

type Team struct {
	ID          string    `bson:"_id,omitempty" json:"_id" query:"_id" form:"_id"`
	Name        string    `json:"name" query:"name" form:"name"`
	Description string    `json:"description" query:"description" form:"description"`
	Members     []string  `json:"members" query:"members" form:"members"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

type TeamView struct {
	ID          string    `bson:"_id,omitempty" json:"_id" query:"_id" form:"_id"`
	Name        string    `json:"name" query:"name" form:"name"`
	Description string    `json:"description" query:"description" form:"description"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
	Members     []Member  `json:"members" query:"members" form:"members"`
}

func GetAllTeams(db *mongo.Database) ([]TeamView, error) {
	collection := db.Collection(TeamCollection)

	//
	// 	db.teams.aggregate([
	// {
	// 	$lookup: {
	// 	from: 'members',
	// 	localField: 'members',
	// 	foreignField: '_id',
	// 	as: 'memberDetails'
	// 	}
	// }
	// ]);
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "members"},
			{Key: "localField", Value: "members"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "members"},
		}}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var teams []TeamView
	for cursor.Next(context.TODO()) {
		var team TeamView
		if err := cursor.Decode(&team); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return teams, nil

}

func GetTeamMembers(db *mongo.Database, teamID string) ([]Member, error) {
	team, err := GetTeamByID(db, teamID)
	if err != nil {
		return nil, err
	}
	return team.Members, nil
}

func InsertTeam(db *mongo.Database, team *Team) (*mongo.InsertOneResult, error) {
	team.ID = uuid.NewString()
	team.CreatedAt = time.Now()
	team.UpdatedAt = time.Now()
	collection := db.Collection(TeamCollection)
	res, err := collection.InsertOne(context.TODO(), team)
	return res, err
}

func UpdateTeam(db *mongo.Database, id string, team *Team) (*mongo.UpdateResult, error) {
	collection := db.Collection(TeamCollection)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"name":        team.Name,
			"description": team.Description,
			"updatedAt":   time.Now(),
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func DeleteTeam(db *mongo.Database, teamID string) (*mongo.DeleteResult, error) {
	collection := db.Collection(TeamCollection)
	filter := bson.M{"_id": teamID}
	res, err := collection.DeleteOne(context.TODO(), filter)
	return res, err
}

func GetTeamByID(db *mongo.Database, teamID string) (*TeamView, error) {
	// Use aggregation to lookup members and return TeamView by teamID
	collection := db.Collection(TeamCollection)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "_id", Value: teamID},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "members"},
			{Key: "localField", Value: "members"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "members"},
		}}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if cursor.Next(context.TODO()) {
		var team TeamView
		if err := cursor.Decode(&team); err != nil {
			return nil, err
		}
		return &team, nil
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return nil, mongo.ErrNoDocuments
}

func GetTeamByIDString(db *mongo.Database, idStr string) (*TeamView, error) {
	return GetTeamByID(db, idStr)
}

func GetTeamsByMemberID(db *mongo.Database, memberID string) ([]Team, error) {
	collection := db.Collection(TeamCollection)
	cursor, err := collection.Find(context.TODO(), bson.M{"members": memberID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var teams []Team
	for cursor.Next(context.TODO()) {
		var team Team
		if err := cursor.Decode(&team); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return teams, nil
}

func AddMemberToTeam(db *mongo.Database, teamID, memberID string) (*mongo.UpdateResult, error) {
	collection := db.Collection(TeamCollection)
	filter := bson.M{"_id": teamID}
	update := bson.M{
		"$addToSet": bson.M{
			"members": memberID,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func RemoveMemberFromTeam(db *mongo.Database, teamID, memberID string) (*mongo.UpdateResult, error) {
	collection := db.Collection(TeamCollection)
	filter := bson.M{"_id": teamID}
	update := bson.M{
		"$pull": bson.M{
			"members": memberID,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func DeleteTeamByIDString(db *mongo.Database, idStr string) (*mongo.DeleteResult, error) {
	return DeleteTeam(db, idStr)
}
