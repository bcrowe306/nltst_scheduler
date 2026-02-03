package models

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

const UserCollection = "users"

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func checkPasswordHash(plain_password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain_password))
	return err == nil
}

type User struct {
	ID            bson.ObjectID `bson:"_id,omitempty"`
	Email         string
	PhoneNumber   string
	PasswordHash  string
	Enabled       bool
	EmailVerified bool
	PhoneVerified bool
}

func CreateUser(db *mongo.Database, email string, password string) (*mongo.InsertOneResult, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := User{
		Email:         email,
		PasswordHash:  passwordHash,
		Enabled:       true,
		EmailVerified: false,
		PhoneVerified: false,
	}
	res, err := insertUser(db, &user)
	return res, err
}

func insertUser(db *mongo.Database, user *User) (*mongo.InsertOneResult, error) {
	collection := db.Collection(UserCollection)
	res, err := collection.InsertOne(context.TODO(), user)
	return res, err
}

func FindUserByID(db *mongo.Database, id bson.ObjectID) (*User, error) {
	collection := db.Collection(UserCollection)
	var user User
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByIDString(db *mongo.Database, idStr string) (*User, error) {
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, err
	}
	return FindUserByID(db, id)
}

func FindUserByEmailPassword(db *mongo.Database, email string, password string) (*User, error) {
	user, err := FindUserByEmail(db, email)
	if err != nil {
		return nil, err
	}
	if !checkPasswordHash(password, user.PasswordHash) {
		return nil, mongo.ErrNoDocuments
	}
	return user, nil
}

func FindUserByEmail(db *mongo.Database, email string) (*User, error) {
	collection := db.Collection(UserCollection)
	var user User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByPhoneNumber(db *mongo.Database, phoneNumber string) (*User, error) {
	collection := db.Collection(UserCollection)
	var user User
	err := collection.FindOne(context.TODO(), bson.M{"phoneNumber": phoneNumber}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(db *mongo.Database, user *User) error {
	collection := db.Collection(UserCollection)
	// Remove ID from the update document to avoid immutable field error
	// Use bson.M to set the fields to be updated

	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"email":       user.Email,
			"phoneNumber": user.PhoneNumber,
		}},
	)
	return err
}

func DeleteUser(db *mongo.Database, id bson.ObjectID) error {
	collection := db.Collection(UserCollection)
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func DeleteUserByIDString(db *mongo.Database, idStr string) error {
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return err
	}
	return DeleteUser(db, id)
}

func DeleteUserByEmail(db *mongo.Database, email string) error {
	collection := db.Collection(UserCollection)
	_, err := collection.DeleteOne(context.TODO(), bson.M{"email": email})
	return err
}

func GetAllUsers(db *mongo.Database) ([]User, error) {
	collection := db.Collection(UserCollection)
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var users []User
	for cursor.Next(context.TODO()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
