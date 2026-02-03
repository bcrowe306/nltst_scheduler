package models

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const TeamCollection = "teams"

type Team struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Name        string
	Description string
	Members     []bson.ObjectID
}

func GetAllTeams(db *mongo.Database) ([]Team, error) {
	collection := db.Collection(TeamCollection)
	cursor, err := collection.Find(context.TODO(), bson.M{})
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

func InsertTeam(db *mongo.Database, team *Team) (*mongo.InsertOneResult, error) {
	collection := db.Collection(TeamCollection)
	res, err := collection.InsertOne(context.TODO(), team)
	return res, err
}

func UpdateTeam(db *mongo.Database, team *Team) (*mongo.UpdateResult, error) {
	collection := db.Collection(TeamCollection)
	filter := bson.M{"_id": team.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        team.Name,
			"description": team.Description,
			"members":     team.Members,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func DeleteTeam(db *mongo.Database, teamID bson.ObjectID) (*mongo.DeleteResult, error) {
	collection := db.Collection(TeamCollection)
	filter := bson.M{"_id": teamID}
	res, err := collection.DeleteOne(context.TODO(), filter)
	return res, err
}

func GetTeamByID(db *mongo.Database, teamID bson.ObjectID) (*Team, error) {
	collection := db.Collection(TeamCollection)
	var team Team
	err := collection.FindOne(context.TODO(), bson.M{"_id": teamID}).Decode(&team)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func GetTeamByIDString(db *mongo.Database, idStr string) (*Team, error) {
	teamID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, err
	}
	return GetTeamByID(db, teamID)
}

func GetTeamsByMemberID(db *mongo.Database, memberID bson.ObjectID) ([]Team, error) {
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

func AddMemberToTeam(db *mongo.Database, teamID, memberID bson.ObjectID) (*mongo.UpdateResult, error) {
	collection := db.Collection(TeamCollection)
	filter := bson.M{"_id": teamID}
	update := bson.M{
		"$addToSet": bson.M{
			"members": memberID,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func RemoveMemberFromTeam(db *mongo.Database, teamID, memberID bson.ObjectID) (*mongo.UpdateResult, error) {
	collection := db.Collection(TeamCollection)
	filter := bson.M{"_id": teamID}
	update := bson.M{
		"$pull": bson.M{
			"members": memberID,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func DeleteTeamByIDString(db *mongo.Database, idStr string) (*mongo.DeleteResult, error) {
	teamID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, err
	}
	return DeleteTeam(db, teamID)
}
