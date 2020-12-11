package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zpmep/hmacutil"
)

type object map[string]interface{}

var (
	app_id = "2553"
	key1   = "PcY4iZIKFCIdgZvA6ueMcMHHUbRLYjPL"
	key2   = "kLtgPl8HHhfvMuDHPwKfgfsY4Ydm9eIz"
)

func CreateOrder(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	transID := rand.Intn(1000000) // Generate random trans id
	embedData, _ := json.Marshal(object{})
	items, _ := json.Marshal([]object{})
	// request data
	params := make(url.Values)
	params.Add("app_id", app_id)
	params.Add("amount", "50000")
	params.Add("app_user", "user123")
	params.Add("embed_data", string(embedData))
	params.Add("item", string(items))
	params.Add("description", "LAZADA - Payment for the order #"+strconv.Itoa(transID))
	params.Add("bank_code", "zalopayapp")

	now := time.Now()
	params.Add("app_time", strconv.FormatInt(now.UnixNano()/int64(time.Millisecond), 10)) // miliseconds

	params.Add("app_trans_id", fmt.Sprintf("%02d%02d%02d_%v", now.Year()%100, int(now.Month()), now.Day(), transID))

	// appid|app_trans_id|appuser|amount|apptime|embeddata|item
	data := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", params.Get("app_id"), params.Get("app_trans_id"), params.Get("app_user"),
		params.Get("amount"), params.Get("app_time"), params.Get("embed_data"), params.Get("item"))
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, key1, data))

	// Content-Type: application/x-www-form-urlencoded
	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/create", params)

	// parse response
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
