package autodebit

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zpmep/hmacutil"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"zalopay-api/common"
)

type Infor struct {
	AppId    string `form:"app_id" json:"app_id" require:"true`
	AppUser  string `form:"app_user" json:"app_user" require:"true`
	Amount   string `form:"amount" json:"amount" require:"true`
	Item     string `form:"item" json:"item" require:"true`
	Embedata string `form:"embed_data" json:"embed_data" require:"true"`
	Mackey   string `form:"mac_key" json:"mac_key" require:"true`
}

func CreateOrder(c *gin.Context) {
	var infor Infor
	if err := c.Bind(&infor); err != nil {
		c.JSON(http.StatusBadRequest, "Bingding Fail")
		return
	}
	appTime := common.GetAppTime()
	appTransId := common.GetTransID(appTime)
	params := make(url.Values)
	params.Add("app_id", infor.AppId)
	params.Add("amount", infor.Amount)
	params.Add("app_user", infor.AppUser)
	params.Add("item", infor.Item)
	params.Add("embed_data", infor.Embedata)
	params.Add("mac_key", infor.Mackey)
	params.Add("app_time", strconv.FormatInt(appTime, 10))
	params.Add("app_trans_id", appTransId)
	params.Add("description", "ZaloPay-Payment for order:"+strconv.Itoa(common.DefaultId))
	log.Println("Input request: ", params)

	data := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", params.Get("app_id"), params.Get("app_trans_id"), params.Get("app_user"),
		params.Get("amount"), params.Get("app_time"), params.Get("embed_data"), params.Get("item"))
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, infor.Mackey, data))

	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/create", params)

	if err != nil {
		// log.Fatal(err)
		log.Println("http POST error: ", err.Error())
		c.JSON(http.StatusBadRequest, "CREATE ORDER ERROR")
		return
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// log.Fatal(err)
		log.Println("GetBinding.Unmarshal data error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot unmarshal data",
		})
		return
	}

	log.Println("response data: ", string(body))
	c.JSON(200, result)
}
