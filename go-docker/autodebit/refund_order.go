package autodebit

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zpmep/hmacutil" // go get github.com/zpmep/hmacutil
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"zalopay-api/common"
)

type Refund struct {
	AppId       string `form:"app_id" json:"app_id" require:"true`
	ZptransID   string `form:"zp_trans_id" json:"zptrans_id" require:"true`
	Amount      string `form:"amount" json:"amount" require:"true`
	Description string `form:"description" json:"description" require:"true"`
	Mackey      string `form:"mac_key" json:"mackey" require:"true`
}

func RefundOrder(c *gin.Context) {
	var refund Refund
	if err := c.Bind(&refund); err != nil {
		c.JSON(http.StatusBadRequest, "Binding Fail")
		return
	}
	appTime := common.GetAppTime()
	mRefundId := common.GetMRefundId(refund.AppId, appTime)
	params := make(url.Values)
	params.Add("timestamp", strconv.FormatInt(appTime, 10))
	params.Add("app_id", refund.AppId)
	params.Add("amount", refund.Amount)
	params.Add("zp_trans_id", refund.ZptransID)
	params.Add("description", refund.Description)
	params.Add("mac_key", refund.Mackey)
	params.Add("m_refund_id", mRefundId)
	log.Println("Input request:", params)
	// app_id|zp_trans_id|amount|description|timestamp
	data := fmt.Sprintf("%v|%v|%v|%v|%v", refund.AppId, params.Get("zp_trans_id"), params.Get("amount"), params.Get("description"), params.Get("timestamp"))
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, refund.Mackey, data))

	// Content-Type: application/x-www-form-urlencoded
	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/refund", params)

	if err != nil {
		// log.Fatal(err)
		log.Println("http GET error: ", err.Error())
		c.JSON(http.StatusBadRequest, "REFUND ORDER ERROR")
		return
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("GetBinding.Unmarshal data error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot unmarshal data",
		})
		return
	}

	log.Println("response data: ", string(body))
	c.JSON(200, result)

}
