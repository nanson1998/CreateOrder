package adapter

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

const (
	appID = "15111"
	key1  = "1IViCAnn06yX4arMLDMscNRn1lIOIoW2"
	key2  = "1e7NjauDRo1NXgjS32NhQoN8NZN86lrz"
)

type Infor struct {
	Amount   string `form:"amount" require:"true"`
	BankCode string `form:"bank_code"`
	AppUser  string `form:"app_user"`
	Phone    string `form:"phone"`
	Email    string `form:"email"`
	Address  string `form:"address"`
}

func CreateOrder(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	transID := rand.Intn(1000000)
	embedData, _ := json.Marshal(object{})
	//embedData, _ := json.Marshal(object{"bankgroup":"ATM"})
	//items, _ := json.Marshal([]object{})

	params := make(url.Values)
	params.Add("app_id", appID)
	params.Add("embed_data", string(embedData))
	params.Add("item", "[{\"itemid\":\"knb\",\"itemname\":\"kim nguyen bao\",\"itemprice\":198400,\"itemquantity\":1}]")
	params.Add("description", "Payment for order:"+strconv.Itoa(transID))
	var infor Infor

	// if c.ShouldBind(&infor) == nil{
	// 	params.Add("amount",infor.Amount)
	// 	params.Add("bank_code",infor.BankCode)
	// 	params.Add("app_user",infor.AppUser)
	// 	params.Add("phone",infor.Phone)
	// }
	if err := c.Bind(&infor); err == nil {

		params.Add("amount", infor.Amount)
		params.Add("bank_code", infor.BankCode)
		params.Add("app_user", infor.AppUser)
		params.Add("phone", infor.Phone)
		params.Add("email", infor.Email)
		params.Add("address", infor.Address)
		//params.Add("item",items)

	} else {
		log.Println("binding fail")
	}
	now := time.Now()
	params.Add("app_time", strconv.FormatInt(now.UnixNano()/int64(time.Millisecond), 10))
	params.Add("app_trans_id", fmt.Sprintf("%02d%02d%02d_%v", now.Year()%100, int(now.Month()), now.Day(), transID))

	data := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", params.Get("app_id"), params.Get("app_trans_id"), params.Get("app_user"),
		params.Get("amount"), params.Get("app_time"), params.Get("embed_data"), params.Get("item"))
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, key1, data))

	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/create", params)

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
