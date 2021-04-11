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

type Check struct {
	AppId      string `form:"app_id" json:"app_id" require:"true`
	Identifier string `form:"identifier" json:"identifier" require:"true`
	Paytoken   string `form:"pay_token" json:"paytoken" require:"true`
	Amount     string `form:"amount" json:"amount" require:"true`
	MacKey     string `form:"mac_key" json:"mac_key" require:"true`
}

func CheckBalance(c *gin.Context) {
	var check Check
	if err := c.Bind(&check); err != nil {
		c.JSON(http.StatusBadRequest, "Binding Fail")
		return
	}
	appTime := common.GetAppTime()
	params := make(url.Values)
	params.Add("req_date", strconv.FormatInt(appTime, 10))
	params.Add("app_id", check.AppId)
	params.Add("amount", check.Amount)
	params.Add("mac_key", check.MacKey)
	params.Add("pay_token", check.Paytoken)
	params.Add("identifier", check.Identifier)
	log.Println("Input request", params)

	macInput := fmt.Sprintf("%v|%v|%v|%v|%v", params.Get("app_id"), params.Get("pay_token"), params.Get("identifier"),
		params.Get("amount"), params.Get("req_date"))
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, check.MacKey, macInput))

	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/agreement/balance", params)
	if err != nil {
		log.Println("http POST error: ", err.Error())
		c.JSON(http.StatusBadRequest, "CHECK BALANCE ERROR")
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
