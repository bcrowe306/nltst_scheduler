package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const MemberCollection = "members"

type Member struct {
	ID          string    `bson:"_id" json:"_id"`
	FirstName   string    `json:"firstName" bson:"firstName" query:"firstName" form:"firstName"`
	LastName    string    `json:"lastName" bson:"lastName" query:"lastName" form:"lastName"`
	Email       string    `json:"email" bson:"email" query:"email" form:"email"`
	PhoneNumber string    `json:"phoneNumber" bson:"phoneNumber" query:"phoneNumber" form:"phoneNumber"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (m *Member) FullName() string {
	return m.FirstName + " " + m.LastName
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

func GetMemberByID(db *mongo.Database, id string) (*Member, error) {
	collection := db.Collection(MemberCollection)
	var member Member
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func GetMembersByIDs(db *mongo.Database, ids []string) ([]Member, error) {
	collection := db.Collection(MemberCollection)
	cursor, err := collection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": ids}})
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
	member.ID = uuid.NewString()
	member.CreatedAt = time.Now()
	member.UpdatedAt = time.Now()
	collection := db.Collection(MemberCollection)
	res, err := collection.InsertOne(context.TODO(), member)
	return res, err
}

func UpdateMember(db *mongo.Database, id string, member *Member) (*mongo.UpdateResult, error) {
	collection := db.Collection(MemberCollection)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"firstName":   member.FirstName,
			"lastName":    member.LastName,
			"email":       member.Email,
			"phoneNumber": member.PhoneNumber,
			"updatedAt":   time.Now(),
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func DeleteMember(db *mongo.Database, id string) (*mongo.DeleteResult, error) {
	collection := db.Collection(MemberCollection)
	res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return res, err
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
