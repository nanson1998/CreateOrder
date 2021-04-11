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

type Unbind struct {
	AppId      string `form:"app_id" json:"app_id" require:"true`
	BindingId  string `form:"binding_id" json:"binding_id" require:"true`
	Identifier string `form:"identifier" json:"identifier" require:"true`
	MacKey     string `form:"mac_key" json:"mac_key" require:"true"`
}

func Unbinding(c *gin.Context) {
	var unbind Unbind
	if err := c.Bind(&unbind); err != nil {
		c.JSON(http.StatusBadRequest, "Binding Fail")
		return
	}
	appTime := common.GetAppTime()
	params := make(url.Values)
	params.Add("app_id", unbind.AppId)
	params.Add("binding_id", unbind.BindingId)
	params.Add("identifier", unbind.Identifier)
	params.Add("mac_key", unbind.MacKey)
	params.Add("req_date", strconv.FormatInt(appTime, 10))
	log.Println("Input request: ", params)
	macInput := fmt.Sprintf("%v|%v|%v|%v", params.Get("app_id"), params.Get("identifier"),
		params.Get("binding_id"), params.Get("req_date"))
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, unbind.MacKey, macInput))
	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/agreement/unbind", params)
	if err != nil {
		// log.Fatal(err)
		log.Println("http GET error: ", err.Error())
		c.JSON(http.StatusBadRequest, "UNBINDING ERROR")
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
