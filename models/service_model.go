package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const ServiceCollection = "services"

type Position struct {
	Name        string `bson:"name" json:"name" query:"name" form:"name"`
	Description string `bson:"description" json:"description" query:"description" form:"description"`
}

type Service struct {
	ID          string     `bson:"_id,omitempty" json:"_id" query:"_id" form:"_id"`
	Name        string     `bson:"name" json:"name" query:"name" form:"name"`
	Description string     `bson:"description" json:"description" query:"description" form:"description"`
	Positions   []Position `bson:"positions,omitempty" json:"positions" query:"positions" form:"positions"`
	CreatedAt   time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time  `bson:"updatedAt" json:"updatedAt"`
}

func GetAllServices(db *mongo.Database) ([]Service, error) {
	collection := db.Collection(ServiceCollection)
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var services []Service
	for cursor.Next(context.TODO()) {
		var service Service
		if err := cursor.Decode(&service); err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return services, nil
}

func InsertService(db *mongo.Database, service *Service) (*mongo.InsertOneResult, error) {
	service.ID = uuid.NewString()
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	service.Positions = make([]Position, 0)
	collection := db.Collection(ServiceCollection)
	res, err := collection.InsertOne(context.TODO(), service)
	return res, err
}

func UpdateService(db *mongo.Database, service *Service) (*mongo.UpdateResult, error) {
	collection := db.Collection(ServiceCollection)
	filter := bson.M{"_id": service.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        service.Name,
			"description": service.Description,
			"updatedAt":   time.Now(),
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func DeleteService(db *mongo.Database, serviceID string) (*mongo.DeleteResult, error) {
	collection := db.Collection(ServiceCollection)
	filter := bson.M{"_id": serviceID}
	res, err := collection.DeleteOne(context.TODO(), filter)
	return res, err
}

func GetServiceByID(db *mongo.Database, serviceID string) (*Service, error) {
	collection := db.Collection(ServiceCollection)
	var service Service
	err := collection.FindOne(context.TODO(), bson.M{"_id": serviceID}).Decode(&service)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func GetServiceByName(db *mongo.Database, name string) (*Service, error) {
	collection := db.Collection(ServiceCollection)
	var service Service
	err := collection.FindOne(context.TODO(), bson.M{"name": name}).Decode(&service)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func GetServiceByIDString(db *mongo.Database, idStr string) (*Service, error) {
	return GetServiceByID(db, idStr)
}

func DeleteServiceByIDString(db *mongo.Database, idStr string) (*mongo.DeleteResult, error) {

	return DeleteService(db, idStr)
}

func AddPositionToService(db *mongo.Database, serviceID string, position Position) (*mongo.UpdateResult, error) {
	collection := db.Collection(ServiceCollection)
	filter := bson.M{"_id": serviceID}
	// Use $addToSet to avoid duplicate positions. Make sure positions is array using $ifNull if null.

	update := bson.M{
		"$addToSet": bson.M{
			"positions": position,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func RemovePositionFromService(db *mongo.Database, serviceID string, positionName string) (*mongo.UpdateResult, error) {
	collection := db.Collection(ServiceCollection)
	filter := bson.M{"_id": serviceID}
	update := bson.M{
		"$pull": bson.M{
			"positions": bson.M{"name": positionName},
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func UpdatePositionInService(db *mongo.Database, serviceID string, position Position) (*mongo.UpdateResult, error) {
	collection := db.Collection(ServiceCollection)
	filter := bson.M{"_id": serviceID, "positions.name": position.Name}
	update := bson.M{
		"$set": bson.M{
			"positions.$.description": position.Description,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}
