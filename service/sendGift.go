package service

import (
	"bytes"
	"errors"
	"gift/common"
	"gift/model/coinHistory/coinHistory_crud"
	"gift/model/user/user_crud"
	"gift/util"
	"github.com/garyburd/redigo/redis"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SendGift struct {
	UserId    string  `json:"userId"`
	GetUserId string  `json:"getUserId"`
	Coin      float32 `json:"coin"`
}
type a []UserCoin
type Ret struct {
	Code        int
	Param       string
	Msg         string
	TotalCount  int64
	NowPageNo   int64
	NowPageSize int64
	Data        a
}
type UserCoin struct {
	UserName string
	CostCoin float32
	UserId   string
}

func SendGiftService(send *SendGift) (ret Ret, sendError error) {
	// 从池里获取连接
	rc := util.RedisClient.Get()
	// 用完后将连接放回连接池
	defer rc.Close()
	pidStr := strconv.Itoa(os.Getpid())
	//分布式锁
	lockName := getString(common.Lock, send.UserId, send.GetUserId)
	_, err := redis.String(rc.Do("GET", lockName))
	if err == nil {
		//fmt.Println(err)
		return
	}
	rc.Do("SETNX", lockName, pidStr)
	rc.Do("EXPIRE", lockName, 10)
	//1.查询送礼用户的信息，并查看这个金额是否大于他自己的金额
	error, user := user_crud.SelectMongodbByUserId(send.UserId)
	if error != nil {
		sendError = errors.New("user does not exist")
	}
	err, toUser := user_crud.SelectMongodbByUserId(send.GetUserId)
	if err != nil {
		sendError = errors.New("toUser does not exist")
	}

	//1.1.如果大于，则扣减需要送礼的用户的金币，否则提示失败
	if send.Coin > user.UserCoin {
		sendError = errors.New("您的金币不够哟!请充值后再赠送")
	}
	//1.2用户赠送，删减金币数
	err1 := user_crud.UpdateCoin(*user, -send.Coin)
	if err1 != nil {
		sendError = errors.New("用户删除金币失败")
	}
	//2.增加金币数
	err2 := user_crud.UpdateCoin(*toUser, send.Coin)
	if err2 != nil {
		user_crud.UpdateCoin(*user, send.Coin)
	}

	//2.5生成对应的redis key，规则是gift_主播Id
	key := strings.Join([]string{common.GIFT, toUser.UserId}, common.UNDERLINE)

	//3.增加流水记录
	coinHistory := new(coinHistory_crud.CoinHistory)
	coinHistory.UserId = user.UserId
	coinHistory.Coin = send.Coin
	coinHistory.GetUserId = toUser.UserId
	//历史记录ID:来源+用户ID+日期时间戳
	coinHistory.HistoryId = getString(common.GIFT, user.UserId, strconv.FormatInt(time.Now().Unix(), 10))
	coinHistory.Time = time.Now()
	err3 := coinHistory_crud.SaveCoinHistory(coinHistory)

	if err3 != nil {
		user_crud.UpdateCoin(*toUser, -send.Coin)
		user_crud.UpdateCoin(*user, send.Coin)
	}
	rc.Do("Del", common.Lock+send.UserId+send.GetUserId)
	//4.增加redis zset数据
	sign := getString(user.UserName, common.UNDERLINE, user.UserId)

	if _, e := rc.Do("ZREVRANK", key, user.UserId); e != nil {
		rc.Do("zadd", key, send.Coin, sign)
	} else {
		rc.Do("ZIncrBy", key, send.Coin, sign)
	}
	ret.Code = 0
	ret.Msg = "success"
	return ret, sendError
}

func GetSortGift(userId string) (ret Ret, sortError error) {
	var userCoin *UserCoin
	var buffer bytes.Buffer
	//得到redis的key
	buffer.WriteString(common.GIFT)
	buffer.WriteString(common.UNDERLINE)
	buffer.WriteString(userId)
	key := buffer.String()
	//获取连接
	rc := util.RedisClient.Get()
	defer rc.Close()
	userMap, err := redis.StringMap(rc.Do("zrevrange", key, 0, -1, "withscores"))
	if err != nil {
		sortError = errors.New("redis get failed")
	}
	userCoin = &UserCoin{}
	var i int64 = 0
	//将map循环放入切片中
	for user := range userMap {
		userCoin.UserName = strings.Split(user, common.UNDERLINE)[0]
		userCoin.UserId = strings.Split(user, common.UNDERLINE)[1]
		v1, _ := strconv.ParseFloat(userMap[user], 32)
		userCoin.CostCoin = float32(v1)
		ret.Data = append(ret.Data, *userCoin)
		i++
	}
	ret.TotalCount = i
	sort.Stable(ret.Data)
	return ret, sortError
}

func GiftList(userId string, pageNo int64, pageSize int64) interface{} {
	return coinHistory_crud.SelectMongodb(userId, pageNo, pageSize)
}

func (s a) Len() int { return len(s) }

func (s a) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s a) Less(i, j int) bool { return s[i].CostCoin > s[j].CostCoin }

func getString(a string, b string, c string) string {
	var buffer bytes.Buffer
	buffer.WriteString(a)
	buffer.WriteString(b)
	buffer.WriteString(c)
	return buffer.String()
}
