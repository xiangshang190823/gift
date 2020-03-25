package util

import (
"context"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
"time"
)

type LogMgr struct {
	Collection *mongo.Collection
	Client     *mongo.Client
}

var (
	GLogMgr *LogMgr
	Client  *mongo.Client
)

var (
	ctx        context.Context
	opts       *options.ClientOptions
	client     *mongo.Client
	err        error
	collection *mongo.Collection
)

func init() {
	// 连接数据库
	ctx, _ = context.WithTimeout(context.Background(), time.Duration(2000)*time.Millisecond) // ctx
	opts = options.Client().ApplyURI("mongodb://localhost:27017").SetMaxPoolSize(10)         // opts
	if client, err = mongo.Connect(ctx, opts); err != nil {
		return
	}
}

func Connect(db string, collectionName string) (GLogMgr *LogMgr) {
	//链接数据库和表
	collection = client.Database(db).Collection(collectionName)
	//单例
	GLogMgr = &LogMgr{
		Client:     client,
		Collection: collection,
	}
	return GLogMgr
}
