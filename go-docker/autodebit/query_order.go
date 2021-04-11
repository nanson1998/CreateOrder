package autodebit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/zpmep/hmacutil" // go get github.com/zpmep/hmacutil
)

type GetOrder struct {
	AppId      string `form:"app_id" json:"app_id" require:"true`
	ApptransId string `form:"app_trans_id" json:"apptrans_id" require:"true`
	MacKey     string `form:"mac_key"json:"mac_key" require:"true`
}

func QueryOrder(c *gin.Context) {
	var getorder GetOrder
	if err := c.Bind(&getorder); err != nil {
		c.JSON(http.StatusBadRequest, "Binding Fail")
		return
	}
	params := make(url.Values)
	params.Add("app_id", getorder.AppId)
	params.Add("app_trans_id", getorder.ApptransId)
	params.Add("mac_key", getorder.MacKey)
	log.Println("Input request:", params)
	appID, err := strconv.Atoi(getorder.AppId)
	data := fmt.Sprintf("%v|%s|%s", appID, getorder.ApptransId, getorder.MacKey) // appid|apptransid|key1
	pr := map[string]interface{}{
		"app_id":       appID,
		"app_trans_id": getorder.ApptransId,
		"mac":          hmacutil.HexStringEncode(hmacutil.SHA256, getorder.MacKey, data),
	}

	jsonStr, err := json.Marshal(pr)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post("https://sb-openapi.zalopay.vn/v2/query", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		// log.Fatal(err)
		log.Println("http GET error: ", err.Error())
		c.JSON(http.StatusBadRequest, "QUERY ORDER ERROR")
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
