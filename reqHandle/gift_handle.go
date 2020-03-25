package reqHandle

import (
	"encoding/json"
	"errors"
	"fmt"
	"gift/common"
	"gift/service"
	"io/ioutil"
	"net/http"
	"strconv"
)

/*送礼*/
func SendGift(writer http.ResponseWriter, request *http.Request) {
	s, _ := ioutil.ReadAll(request.Body)
	var gift service.SendGift
	err := json.Unmarshal(s, &gift)
	// 根据请求body创建一个json解析器实例
	if err != nil {
		fmt.Println(err)
	}
	ret, _ := service.SendGiftService(&gift)
	s, _ = json.Marshal(ret)
	writer.Write(s)
}

/*送礼排行榜*/
func GetSortGift(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	userId := query.Get("userId")
	if len(userId) < 1 {
		errors.New(common.ParameterError)
	}
	var s = make([]byte, 5, 10)
	map1, _ := service.GetSortGift(userId)
	s, _ = json.Marshal(map1)
	writer.Write(s)
}

/*资金流水*/
func GiftList(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	userId := query.Get("userId")
	pageNo := query.Get("pageNo")
	pageSize := query.Get("pageSize")
	var pn int64
	var ps int64
	if len(userId) < 1 {
		errors.New(common.ParameterError)
	}
	if len(pageNo) < 1 {
		pn = 1
	} else {
		pn, _ = strconv.ParseInt(pageNo, 10, 64)
	}
	if len(pageSize) < 1 {
		ps = 10
	} else {
		ps, _ = strconv.ParseInt(pageSize, 10, 64)
	}
	var s1 = make([]byte, 10, 20)
	s1, _ = json.Marshal(service.GiftList(userId, pn, ps))
	writer.Write(s1)
}
