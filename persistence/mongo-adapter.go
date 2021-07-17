package persistence

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoOptions options.ClientOptions

type MongoAdapter struct {
	dialInfo           *options.ClientOptions
	dbName             string
	db                 *mongo.Database
	client             *mongo.Client
	MessagesCollection *mongo.Collection
}

type Index struct {
	Name       string
	Collection string
	Keys       interface{}
	Unique     bool
	Sparse     bool
	Expires    time.Duration
}

var mongoSingleton *MongoAdapter

func NewPersistenceLayer(urls []string, user, password, dbName string, extraArgs map[string]interface{}) (*MongoAdapter, error) {
	if mongoSingleton == nil {
		var err error
		clientOptions := &options.ClientOptions{}
		mode := readpref.SecondaryPreferredMode
		if val, k := extraArgs["read_pref"]; k {
			mode, err = readpref.ModeFromString(val.(string))
			if err != nil {
				return nil, err
			}
		}
		if val, k := extraArgs["timeout"]; k {
			timeout, k := val.(time.Duration)
			if !k {
				return nil, fmt.Errorf("missing or invalid timeout value")
			}
			clientOptions.ConnectTimeout = &timeout
		}

		if val, k := extraArgs["replicaset"]; k && val != nil {
			val := val.(string)
			clientOptions.Hosts = urls
			clientOptions.SetReplicaSet(val)
			readPref, err := readpref.New(mode)
			if err != nil {
				return nil, err
			}
			clientOptions.ReadPreference = readPref
		}
		if user != "" && password != "" {
			auth := options.Credential{
				AuthSource: "admin",
				Username:   user,
				Password:   password,
			}
			clientOptions.SetAuth(auth)
		}

		if len(urls) == 1 {
			clientOptions.SetDirect(true)
			clientOptions.ApplyURI(urls[0])
		}
		if dbName == "" {
			dbName = "notificationaway"
		}

		mongoSingleton = &MongoAdapter{
			dialInfo: clientOptions,
			dbName:   dbName,
		}

	}
	return mongoSingleton, nil
}

func (m *MongoAdapter) Connect(ctx context.Context) error {
	var err error
	m.client, err = mongo.Connect(ctx, m.dialInfo)
	if err != nil {
		return err
	}
	err = m.client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	m.db = m.client.Database(m.dbName)
	m.MessagesCollection = m.db.Collection(MessagesDB{}.CollectionName())
	return m.createIndexesForCollection(MessagesDB{}.Indexes())
}

func (m *MongoAdapter) createIndexesForCollection(indexes []Index) error {
	var err error
	var indexModels []mongo.IndexModel

	for _, index := range indexes {
		indexOptions := options.Index().SetUnique(index.Unique).SetName(index.Name)
		if index.Sparse {
			indexOptions.SetSparse(index.Sparse)
		}
		if index.Expires > 0 {
			indexOptions.SetExpireAfterSeconds(int32(index.Expires.Seconds()))
		}
		indexModel := mongo.IndexModel{
			Keys:    index.Keys,
			Options: indexOptions,
		}
		indexModels = append(indexModels, indexModel)
	}

	indexView := m.db.Collection(indexes[0].Collection).Indexes()
	result, indexErr := indexView.CreateMany(context.Background(), indexModels)
	if result == nil || indexErr != nil {
		err = errors.Wrapf(indexErr, "unable to add indexes for %s collection", indexes[0].Collection)
	}
	return err
}

func (m *MongoAdapter) CloseConnection(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
