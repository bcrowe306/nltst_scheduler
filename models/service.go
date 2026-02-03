package models

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const ServiceCollection = "services"

type Position struct {
	Name        string
	Description string
}

type Service struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Name        string
	Description string
	Positions   []Position
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
			"positions":   service.Positions,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func DeleteService(db *mongo.Database, serviceID bson.ObjectID) (*mongo.DeleteResult, error) {
	collection := db.Collection(ServiceCollection)
	filter := bson.M{"_id": serviceID}
	res, err := collection.DeleteOne(context.TODO(), filter)
	return res, err
}

func GetServiceByID(db *mongo.Database, serviceID bson.ObjectID) (*Service, error) {
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
	serviceID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, err
	}
	return GetServiceByID(db, serviceID)
}

func DeleteServiceByIDString(db *mongo.Database, idStr string) (*mongo.DeleteResult, error) {
	serviceID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, err
	}
	return DeleteService(db, serviceID)
}

func AddPositionToService(db *mongo.Database, serviceID bson.ObjectID, position Position) (*mongo.UpdateResult, error) {
	collection := db.Collection(ServiceCollection)
	filter := bson.M{"_id": serviceID}
	update := bson.M{
		"$push": bson.M{
			"positions": position,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func RemovePositionFromService(db *mongo.Database, serviceID bson.ObjectID, positionName string) (*mongo.UpdateResult, error) {
	collection := db.Collection(ServiceCollection)
	filter := bson.M{"_id": serviceID}
	update := bson.M{
		"$pull": bson.M{
			"positions": bson.M{"name": positionName},
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	return res, err
}

func UpdatePositionInService(db *mongo.Database, serviceID bson.ObjectID, position Position) (*mongo.UpdateResult, error) {
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
