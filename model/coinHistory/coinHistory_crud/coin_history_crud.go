package coinHistory_crud

import (
	"context"
	"errors"
	"gift/common"
	"gift/util"
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

type Ret struct {
	Code        int
	Param       string
	Msg         string
	TotalCount  int64
	NowPageNo   int64
	NowPageSize int64
	Data        []CoinHistory
}

var GLogMgr *util.LogMgr

func init() {
	GLogMgr = util.Connect("test", "coin_history")
}

//保存数据
func SaveCoinHistory(history *CoinHistory) (err error) {
	defer GLogMgr.Client.Disconnect(context.TODO())
	if _, err = GLogMgr.Collection.InsertOne(context.TODO(), &history); err != nil {
		return errors.New(common.INSERT_DATA_ERROR + err.Error())
	}
	return
}

//查询数据
func SelectMongodb(userId string, pageNo int64, pageSize int64) (ret Ret) {
	var (
		cur         *mongo.Cursor
		ctx         context.Context
		err         error
		coinHistory *CoinHistory
	)
	ctx = context.TODO()
	defer cur.Close(ctx)
	defer GLogMgr.Client.Disconnect(context.TODO())
	ret.NowPageNo = pageNo
	ret.NowPageSize = pageSize
	if cur, err = GLogMgr.Collection.Find(ctx, bson.M{"get_user_id": userId},
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
