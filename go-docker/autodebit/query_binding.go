package autodebit

import (
	"bytes"
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

type Query struct {
	AppId      string `form:"app_id" json:"app_id" require:"true`
	ApptransId string `form:"app_trans_id" json:"apptrans_id" require:"true`
	MacKey     string `form:"mackey" json:"mac_key" require:"true`
}

func QueryBinding(c *gin.Context) {
	var query Query
	if err := c.Bind(&query); err != nil {
		c.JSON(http.StatusBadRequest, "Binding Fail")
		return
	}
	params := make(url.Values)
	params.Add("app_id", query.AppId)
	params.Add("app_trans_id", query.ApptransId)
	params.Add("mac_key", query.MacKey)
	log.Println("Input request: ", params)
	reqDate := common.GetAppTime()
	macInput := fmt.Sprintf("%v|%v|%v", query.AppId, query.ApptransId, reqDate)
	appID, err := strconv.Atoi(query.AppId)
	pr := map[string]interface{}{
		"req_date":     reqDate,
		"app_id":       appID,
		"app_trans_id": query.ApptransId,
		"mac":          hmacutil.HexStringEncode(hmacutil.SHA256, query.MacKey, macInput),
	}
	jsonStr, err := json.Marshal(pr)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.Post("https://sb-openapi.zalopay.vn/v2/agreement/query", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Println("http POST error: ", err.Error())
		c.JSON(http.StatusBadRequest, "QUERY BINDING ERROR")
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
