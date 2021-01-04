package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zpmep/hmacutil" // go get github.com/zpmep/hmacutil
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func Refund(c *gin.Context) {
	params := make(url.Values)
	params.Add("app_id", "15111")
	params.Add("zp_trans_id", "201230000000632")
	params.Add("amount", "1000000")
	params.Add("description", "ZaloPay Refund Demo")

	now := time.Now()
	timestamp := now.UnixNano() / int64(time.Millisecond) // Miliseconds
	params.Add("timestamp", strconv.FormatInt(timestamp, 10))

	uid := fmt.Sprintf("%d%d", timestamp, 111+rand.Intn(888))

	params.Add("m_refund_id", fmt.Sprintf("%02d%02d%02d_%v_%v", now.Year()%100, int(now.Month()), now.Day(), appID, uid))

	// app_id|zp_trans_id|amount|description|timestamp
	data := fmt.Sprintf("%v|%v|%v|%v|%v", appID, params.Get("zp_trans_id"), params.Get("amount"), params.Get("description"), params.Get("timestamp"))
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, key1, data))

	log.Printf("%+v", params)
	time := fmt.Sprintf("%v", params.Get("timestamp"))
	log.Printf(time)

	// Content-Type: application/x-www-form-urlencoded
	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/refund", params)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var result map[string]interface{}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}

	c.JSON(200, result)
}
