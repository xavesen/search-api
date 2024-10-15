package storage

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/xavesen/search-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		log.Errorf("Error converting userId string %s to object id while checking user rights for index with name %s: %s", userId, indexName, err.Error())
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

func (s *MongoStorage) AddIndexToUser(ctx context.Context, userId string, indexName string) error {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Errorf("Error converting userId string %s to object id while adding index with name %s to user: %s", userId, indexName, err.Error())
		return err
	}

	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "indexes", Value: indexName},
		}},
	}

	result, err := s.usersCollection.UpdateByID(ctx, oid, update)
	if err != nil {
		log.Errorf("Error pushing index with name %s to indexes array in users document with id %s: %s", indexName, userId, err)
		return err
	} else if result.MatchedCount == 0 {
		log.Errorf("Error pushing index with name %s to indexes array in users document with id %s: user doesn't exist", indexName, userId)
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *MongoStorage) GetUserInfo(ctx context.Context, login string) (*models.User, error) {
	var user *models.User
	filter := bson.D{
		{Key: "login", Value: login},
	}

	if err := s.usersCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Warningf("Tried to find in db non-existent user with login %s ", login)
		} else {
			log.Errorf("Error searching for user with login %s in db: %s", login, err.Error())
		}
		return nil, err
	}

	return user, nil
}

func (s *MongoStorage) SetRefreshToken(ctx context.Context, userId string, refreshToken string) error {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Errorf("Error converting userId string %s to object id while adding refresh token %s to user: %s", userId, refreshToken, err.Error())
		return err
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "refreshToken", Value: refreshToken},
		}},
	}

	result, err := s.usersCollection.UpdateByID(ctx, oid, update)
	if err != nil {
		log.Errorf("Error setting refresh token %s for user %s: %s", refreshToken, userId, err)
		return err
	} else if result.MatchedCount < 1 {
		log.Errorf("Error setting refresh token %s for user %s: No user with such id", refreshToken, userId)
		return mongo.ErrNoDocuments
	}

	return nil
}