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

type Paytoken struct {
	AppId        string `form:"app_id" json:"app_id" require:"true`
	PayToken     string `form:"pay_token" json:"pay_token" require:"true`
	Identifier   string `form:"identifier" json:"identifier" require:"true`
	Zptranstoken string `form:"zp_trans_token" json:"zptranstoken" require:"true"`
	MacKey       string `form:"mac_key" json:"mac_key" require:"true"`
}

func Pay(c *gin.Context) {
	var pay Paytoken
	if err := c.Bind(&pay); err != nil {
		c.JSON(http.StatusBadRequest, "Binding Fail")
		return
	}
	appTime := common.GetAppTime()
	params := make(url.Values)
	params.Add("app_id", pay.AppId)
	params.Add("pay_token", pay.PayToken)
	params.Add("identifier", pay.Identifier)
	params.Add("zp_trans_token", pay.Zptranstoken)
	params.Add("mac_key", pay.MacKey)
	params.Add("req_date", strconv.FormatInt(appTime, 10))
	log.Println("Input request: ", params)
	macInput := fmt.Sprintf("%v|%v|%v|%v|%v", params.Get("app_id"), params.Get("identifier"), params.Get("zp_trans_token"),
		params.Get("pay_token"), params.Get("req_date"))
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, pay.MacKey, macInput))

	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/agreement/pay", params)

	if err != nil {
		// log.Fatal(err)
		log.Println("http GET error: ", err.Error())
		c.JSON(http.StatusBadRequest, "PAY ERROR")
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
