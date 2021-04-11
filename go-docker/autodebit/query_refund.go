package autodebit

import (
	"bytes"
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

type GetRefund struct {
	AppId     string `form:"app_id" json:"app_id" require:"true`
	MacKey    string `form:"mac_key" json:"mac_key" require:"true`
	MRefundId string `form:"m_refund_id" json:"m_refund_id" require:"true`
}

func QueryRefund(c *gin.Context) {
	var getrefund GetRefund
	if err := c.Bind(&getrefund); err != nil {
		c.JSON(http.StatusBadRequest, "binding Fail")
		return
	}
	timeStamp := common.GetAppTime()
	params := make(url.Values)
	params.Add("app_id", getrefund.AppId)
	params.Add("m_refund_id", getrefund.MRefundId)
	params.Add("mac_key", getrefund.MacKey)
	appID, err := strconv.Atoi(getrefund.AppId)
	log.Println("Input request:", params)
	pr := map[string]interface{}{
		"app_id":      appID,
		"m_refund_id": getrefund.MRefundId,
		"timestamp":   timeStamp, // miliseconds
	}

	data := fmt.Sprintf("%v|%v|%v", appID, getrefund.MRefundId, timeStamp)
	pr["mac"] = hmacutil.HexStringEncode(hmacutil.SHA256, getrefund.MacKey, data)
	jsonStr, err := json.Marshal(pr)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post("https://sb-openapi.zalopay.vn/v2/query_refund", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		// log.Fatal(err)
		log.Println("http GET error: ", err.Error())
		c.JSON(http.StatusBadRequest, "QUERY REFUND ERROR")
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
