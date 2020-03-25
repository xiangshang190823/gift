package user_crud

import (
	"context"
	"errors"
	"gift/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	UserName string  `bson:"user_name"`
	UserId   string  `bson:"user_id"`
	Status   int     `bson:"status"` //0 主播  1 用户
	UserCoin float32 `bson:"user_coin"`
}

type LogMgr struct {
	client     *mongo.Client
	collection *mongo.Collection
}

var (
	GLogMgr *LogMgr
	Client  *mongo.Client
)

func init() {

	var (
		ctx        context.Context
		opts       *options.ClientOptions
		client     *mongo.Client
		err        error
		collection *mongo.Collection
	)
	// 连接数据库
	ctx, _ = context.WithTimeout(context.Background(), time.Duration(2000)*time.Millisecond) // ctx
	opts = options.Client().ApplyURI("mongodb://localhost:27017").SetMaxPoolSize(10)         // opts
	if client, err = mongo.Connect(ctx, opts); err != nil {
		return
	}

	//链接数据库和表
	collection = client.Database("test").Collection("user")

	//单例
	GLogMgr = &LogMgr{
		client:     client,
		collection: collection,
	}
}

//保存数据
func (logMgr *LogMgr) SaveMongodb() (err error) {
	var (
		insetRest *mongo.InsertOneResult
		_         interface{}
		users     []interface{}
	)
	//userId := bson.NewObjectId().String()
	user := User{"主播", "1", 0, 0}

	if insetRest, err = logMgr.collection.InsertOne(context.TODO(), &user); err != nil {
		return
	}
	_ = insetRest.InsertedID

	users = append(users, &User{UserName: "用户1", UserId: "10001", Status: 1, UserCoin: 5102},
		&User{UserName: "用户2", UserId: "10002", Status: 1, UserCoin: 84263},
		&User{UserName: "用户3", UserId: "10003", Status: 1, UserCoin: 1520},
		&User{UserName: "用户4", UserId: "10004", Status: 1, UserCoin: 358})
	if _, err = logMgr.collection.InsertMany(context.TODO(), users); err != nil {
		errors.New(common.InsertDataError)
		return
	}
	return
}

//查询数据
func (logMgr *LogMgr) SelectMongodbByUserId(userId string) (err error, user *User) {
	if err = logMgr.collection.FindOne(context.TODO(), bson.M{"user_id": userId}).Decode(&user); err != nil {
		return
	}
	return
}

//更新数据
func (logMgr *LogMgr) UpdateCoin(user User, coin float32) (err error) {
	var (
		ctx context.Context
	)
	if singleResult := logMgr.collection.FindOneAndUpdate(ctx, bson.M{"user_id": user.UserId, "user_coin": user.UserCoin},
		bson.M{"$inc": bson.M{"user_coin": coin}}); singleResult.Err() != nil {
		errors.New(common.UpdateDataError)
		return
	}
	return
}
