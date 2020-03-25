package coinHistory_crud

import (
	"context"
	"errors"
	"gift/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type CoinHistory struct {
	UserId    string    `bson:"user_id"`
	GetUserId string    `bson:"get_user_id"`
	HistoryId string    `bson:"history_id"`
	Coin      float32   `bson:"coin"`
	Time      time.Time `bson:"create_time"`
}

type LogMgr struct {
	client     *mongo.Client
	collection *mongo.Collection
}

var (
	GLogMgr *LogMgr
	Client  *mongo.Client
)

type Ret struct {
	Code        int
	Param       string
	Msg         string
	TotalCount  int64
	NowPageNo   int64
	NowPageSize int64
	Data        []CoinHistory
}

func init() {

	var (
		ctx        context.Context
		opts       *options.ClientOptions
		client     *mongo.Client
		err        error
		collection *mongo.Collection
	)
	// 1.连接数据库
	ctx, _ = context.WithTimeout(context.Background(), time.Duration(2000)*time.Millisecond) // ctx
	opts = options.Client().ApplyURI("mongodb://localhost:27017").SetMaxPoolSize(10)         // opts
	if client, err = mongo.Connect(ctx, opts); err != nil {
		errors.New(err.Error())
		return
	}
	//2.链接数据库和表
	collection = client.Database("test").Collection("coin_history")
	//3.赋值单例
	GLogMgr = &LogMgr{
		client:     client,
		collection: collection,
	}
}

//保存数据
func (logMgr *LogMgr) SaveCoinHistory(history *CoinHistory) (err error) {
	if _, err = logMgr.collection.InsertOne(context.TODO(), &history); err != nil {
		return errors.New(common.InsertDataError + err.Error())
	}
	return
}

//查询数据
func (logMgr *LogMgr) SelectMongodb(userId string, pageNo int64, pageSize int64) (ret Ret) {
	var (
		cur         *mongo.Cursor
		ctx         context.Context
		err         error
		coinHistory *CoinHistory
	)
	defer cur.Close(ctx)
	ctx = context.TODO()
	ret.NowPageNo = pageNo
	ret.NowPageSize = pageSize
	if cur, err = logMgr.collection.Find(ctx, bson.M{"get_user_id": userId},
		options.Find().SetSort(bson.M{"create_time": -1}).SetLimit(pageSize).SetSkip(pageNo-1)); err != nil {
		return
	}

	for cur.Next(ctx) {
		coinHistory = &CoinHistory{}
		if err = cur.Decode(coinHistory); err != nil {
			return
		}
		ret.Msg = "success"
		ret.Data = append(ret.Data, *coinHistory)
	}
	return
}
