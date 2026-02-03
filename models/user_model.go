package models

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
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
	ID            string    `bson:"_id" json:"_id"`
	Name          string    `json:"name" bson:"name"`
	Email         string    `json:"email" bson:"email"`
	PhoneNumber   string    `json:"phoneNumber" bson:"phoneNumber"`
	PasswordHash  string    `json:"passwordHash" bson:"passwordHash"`
	Enabled       bool      `json:"enabled" bson:"enabled"`
	EmailVerified bool      `json:"emailVerified" bson:"emailVerified"`
	PhoneVerified bool      `json:"phoneVerified" bson:"phoneVerified"`
	CreatedAt     time.Time `json:"created" bson:"createTime"`
	UpdatedAt     time.Time `json:"updated" bson:"updateTime"`
	LastLogin     time.Time `json:"lastLogin" bson:"lastLogin"`
}

func CreateUser(db *mongo.Database, name string, email string, password string, phoneNumber string) (*mongo.InsertOneResult, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := User{
		Name:          name,
		Email:         email,
		PhoneNumber:   phoneNumber,
		PasswordHash:  passwordHash,
		Enabled:       true,
		EmailVerified: false,
		PhoneVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		LastLogin:     time.Time{},
	}

	user.ID = uuid.NewString()
	res, err := insertUser(db, &user)
	return res, err
}

func insertUser(db *mongo.Database, user *User) (*mongo.InsertOneResult, error) {
	collection := db.Collection(UserCollection)
	res, err := collection.InsertOne(context.TODO(), user)
	return res, err
}

func FindUserByID(db *mongo.Database, id string) (*User, error) {
	collection := db.Collection(UserCollection)
	var user User
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByIDString(db *mongo.Database, idStr string) (*User, error) {
	return FindUserByID(db, idStr)
}

func FindUserByEmailPassword(db *mongo.Database, email string, password string) (*User, error) {
	log.Printf("Attempting to find user by email: %s", email)
	user, err := FindUserByEmail(db, email)
	if err != nil {
		return nil, err
	}
	log.Printf("User found: %s", user.Email)
	if !checkPasswordHash(password, user.PasswordHash) {
		return nil, mongo.ErrNoDocuments
	}

	return user, nil
}

func UpdateUserLoginTime(db *mongo.Database, userID string) error {
	collection := db.Collection(UserCollection)
	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{
			"lastLogin": time.Now(),
		}},
	)
	return err
}

func FindUserByEmail(db *mongo.Database, email string) (*User, error) {
	collection := db.Collection(UserCollection)
	var user User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	log.Printf("Found user with email: %s", user.Email)
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

func UpdateUser(db *mongo.Database, userID string, user *User) error {
	collection := db.Collection(UserCollection)
	// Remove ID from the update document to avoid immutable field error
	// Use bson.M to set the fields to be updated

	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{
			"email":       user.Email,
			"phoneNumber": user.PhoneNumber,
			"name":        user.Name,
			"updatedAt":   time.Now(),
			"enabled":     user.Enabled,
		}},
	)
	return err
}

func DeleteUser(db *mongo.Database, id string) error {
	collection := db.Collection(UserCollection)
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func DeleteUserByIDString(db *mongo.Database, idStr string) error {

	return DeleteUser(db, idStr)
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
