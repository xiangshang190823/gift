package user_crud

import (
	"context"
	"errors"
	"gift/common"
	"gift/util"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	UserName string  `bson:"user_name"`
	UserId   string  `bson:"user_id"`
	Status   int     `bson:"status"` //0 主播  1 用户
	UserCoin float32 `bson:"user_coin"`
}

var GLogMgr *util.LogMgr

func init() {
	GLogMgr = util.Connect("test", "user")
}

//保存数据
func SaveMongodb() (err error) {
	var (
		insetRest *mongo.InsertOneResult
		_         interface{}
		users     []interface{}
	)
	user := User{"主播", "1", 0, 0}
	defer GLogMgr.Client.Disconnect(context.TODO())
	if insetRest, err = GLogMgr.Collection.InsertOne(context.TODO(), &user); err != nil {
		return
	}
	_ = insetRest.InsertedID

	users = append(users, &User{UserName: "用户1", UserId: "10001", Status: 1, UserCoin: 5102},
		&User{UserName: "用户2", UserId: "10002", Status: 1, UserCoin: 84263},
		&User{UserName: "用户3", UserId: "10003", Status: 1, UserCoin: 1520},
		&User{UserName: "用户4", UserId: "10004", Status: 1, UserCoin: 358})
	if _, err = GLogMgr.Collection.InsertMany(context.TODO(), users); err != nil {
		err = errors.New(common.InsertDataError)
		return
	}
	return
}

//查询数据
func SelectMongodbByUserId(userId string) (err error, user *User) {
	defer GLogMgr.Client.Disconnect(context.TODO())
	if err = GLogMgr.Collection.FindOne(context.TODO(), bson.M{"user_id": userId}).Decode(&user); err != nil {
		return
	}
	return
}

//更新数据
func UpdateCoin(user User, coin float32) (err error) {
	var ctx context.Context
	defer GLogMgr.Client.Disconnect(context.TODO())
	if singleResult := GLogMgr.Collection.FindOneAndUpdate(ctx, bson.M{"user_id": user.UserId, "user_coin": user.UserCoin},
		bson.M{"$inc": bson.M{"user_coin": coin}}); singleResult.Err() != nil {
		errors.New(common.UpdateDataError)
		return
	}
	return
}
