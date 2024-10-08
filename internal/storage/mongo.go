package storage

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

type MongoStorage struct {
	client 				*mongo.Client
	database 			*mongo.Database
	usersCollection		*mongo.Collection
}

func NewMongoStorage(ctx context.Context, addr string, db string, user string, password string) (*MongoStorage, error) {
	log.Infof("Initializing client and connecting mongo db %s on %s with user %s", db, addr, user)

	clientCreds := options.Credential{
		Username: user,
		Password: password,
		AuthSource: db,
	}
	clientOpts := options.Client()
	clientOpts.SetAuth(clientCreds)
	clientOpts.SetHosts([]string{addr})

	log.Debug("Initializing mongo client")
	newClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Errorf("Error while initializing mongo client for db %s on %s with user %s: %s", db, addr, user, err.Error())
		return nil, err
	}

	log.Debug("Connecting mongo db")
	if err = newClient.Ping(ctx, nil); err != nil {
		log.Errorf("Error while connecting mongo db %s on %s with user %s: %s", db, addr, user, err.Error())
		return nil, err
	}

	log.Debug("Initializing db and collections")
	appDb := newClient.Database(db)
	usersCol := appDb.Collection("users")

	newStorage := &MongoStorage{
		client: newClient,
		database: appDb,
		usersCollection: usersCol,
	}

	log.Info("Successfully initialized and connected mongo db")
	return newStorage, nil
}

func (s *MongoStorage) CheckUserIndexRights(ctx context.Context, userId string, indexName string) (bool, error) {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Warningf("Error converting userId string %s to object id while checking user rights for index with name %s: %s", userId, indexName, err.Error())
		return false, err
	}

	filter := bson.D{
		{Key: "$and", Value: bson.A{
				bson.D{{Key: "_id", Value: oid}}, // user object id equals userId argument
				bson.D{{Key: "indexes", Value: indexName},}, // argument indexName is present in indexes field inside user document
			},
		},
	}

	count, err := s.usersCollection.CountDocuments(ctx, filter)
	if err != nil {
		log.Errorf("Error finding in db user with id %s and %s index name in indexes field: %s", userId, indexName, err)
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}