package datastore

import (
	"context"
	"fmt"
	"io"

	"github.com/puipuipartpicker/kbpartpicker/api/pkg/di"
	"github.com/puipuipartpicker/kbpartpicker/api/pkg/env"
	"github.com/puipuipartpicker/kbpartpicker/api/pkg/logging"
	"github.com/puipuipartpicker/kbpartpicker/api/pkg/retry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

const (
	envAppEnv           env.VarName = "APP_ENV"
	envDatabaseURI      env.VarName = "DB_URI"
	envDatabaseUsername env.VarName = "DB_USERNAME"
	envDatabasePassword env.VarName = "DB_PASSWORD"
	envDatabaseName     env.VarName = "DB_NAME"
)

var (
	db *mongo.Database
)

// NewBaseRepo returns a base repository.
func NewBaseRepo(db *mongo.Database) *BaseRepo {
	return &BaseRepo{db: db}
}

type BaseRepo struct {
	db *mongo.Database
}

type wrapClient struct {
	client *mongo.Client
}

func (w *wrapClient) Close() error {
	return w.client.Disconnect(context.Background())
}

// GetClient returns mongo database client.
func GetClient(ctx context.Context) (*mongo.Client, string, error) {
	if GetAppEnv().IsTest() {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		return client, "test", err
	}

	database := env.String(envDatabaseName)

	clientOpts := options.Client().ApplyURI(env.String(envDatabaseURI))
	clientOpts.SetReadPreference(readpref.Primary())
	clientOpts.SetAuth(options.Credential{
		AuthSource: database,
		Username:   env.String(envDatabaseUsername),
		Password:   env.String(envDatabasePassword),
	})

	client, err := mongo.Connect(ctx, clientOpts)

	return client, database, err
}

// GetDatabase returns mongo database connection.
func GetDatabase() *mongo.Database {
	if db != nil {
		return db
	}

	ctx := context.Background()

	client, database, err := GetClient(ctx)
	if err != nil {
		di.LogInitFatal("failed to connect mongo", err)
	}

	l := di.GetLogger().Named("repository")

	if err := pingDatabase(ctx, l, client); err != nil {
		di.LogInitFatal("failed to ping database", err)
	}

	var c io.Closer = &wrapClient{
		client: client,
	}

	di.RegisterCloser("MongoDB connection", c)

	db = client.Database(database)

	return db
}

// createFilter returns basic mongo bson.M filter (include option to use soft delete or not).
func createFilter(ignoreDeletedDocument bool) bson.M {
	filter := bson.M{}
	if !ignoreDeletedDocument {
		filter = bson.M{
			"$or": []bson.M{
				{"deleted_at": bson.M{"$exists": false}},
				{"deleted_at": nil},
			},
		}
	}

	return filter
}

func GetAppEnv() env.AppEnv {
	return env.Env(envAppEnv, env.EnvTest)
}

func pingDatabase(ctx context.Context, logger logging.Logger, client *mongo.Client) error {
	// use default configuration
	retrier, err := retry.New()
	if err != nil {
		return fmt.Errorf("failed to initialize retrier: %w", err)
	}

	count := 1

	err = retrier.Do(ctx, func(ctx context.Context) error {
		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			logger.Warn(fmt.Sprintf("failed to ping to database, try count %d", count), zap.Error(err))
			count++

			return fmt.Errorf("failed to ping to database, try count %d: %w", count, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to establish database connection: %w", err)
	}

	return nil
}

