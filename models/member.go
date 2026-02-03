package models

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const MemberCollection = "members"

type Member struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

func GetAllMembers(db *mongo.Database) ([]Member, error) {
	collection := db.Collection(MemberCollection)
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var members []Member
	for cursor.Next(context.TODO()) {
		var member Member
		if err := cursor.Decode(&member); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

func InsertMember(db *mongo.Database, member *Member) (*mongo.InsertOneResult, error) {
	collection := db.Collection(MemberCollection)
	res, err := collection.InsertOne(context.TODO(), member)
	return res, err
}

func UpdateMember(db *mongo.Database, member *Member) (*mongo.UpdateResult, error) {
	collection := db.Collection(MemberCollection)
	filter := bson.M{"_id": member.ID}
	update := bson.M{
		"$set": bson.M{
			"firstName": member.FirstName,
			"lastName":  member.LastName,
			"email":     member.Email,
			"phone":     member.Phone,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func DeleteMember(db *mongo.Database, id bson.ObjectID) (*mongo.DeleteResult, error) {
	collection := db.Collection(MemberCollection)
	res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return res, err
}

func DeleteMemberByIDString(db *mongo.Database, idStr string) (*mongo.DeleteResult, error) {
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, err
	}
	return DeleteMember(db, id)
}

func GetMemberByEmail(db *mongo.Database, email string) (*Member, error) {
	collection := db.Collection(MemberCollection)
	var member Member
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func GetMemberByPhone(db *mongo.Database, phone string) (*Member, error) {
	collection := db.Collection(MemberCollection)
	var member Member
	err := collection.FindOne(context.TODO(), bson.M{"phone": phone}).Decode(&member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}
