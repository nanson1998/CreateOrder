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

type object map[string]interface{}

type CreateBindingRequest struct {
	AppId            string `form:"app_id" json:"app_id" require:"true`
	Billing          string `form:"billing" json:"billing"`
	BindingType      string `form:"binding_type" json:"binding_type" require:"true`
	Identifier       string `form:"identifier" json:"identifier" require:"true`
	MaxAmount        string `form:"max_amount" json:"max_amount" require:"true"`
	RedirectUrl      string `form:"redirect_url"`
	RedirectDeepLink string `form:"redirect_deep_link"`
	CallbackUrl      string `form:"callback_url"`
	MacKey           string `form:"mac_key" json:"mac_key" require:"true"`
}

func CreateBinding(c *gin.Context) {
	var create CreateBindingRequest
	if err := c.Bind(&create); err != nil {
		c.JSON(http.StatusBadRequest, "Binding Fail")
		return
	}
	appTime := common.GetAppTime()
	appTransId := common.GetTransID(appTime)
	binding_data, _ := json.Marshal(object{})
	params := make(url.Values)
	params.Add("binding_data", string(binding_data))
	params.Add("req_date", strconv.FormatInt(appTime, 10))
	params.Add("app_trans_id", appTransId)
	params.Add("app_id", create.AppId)
	params.Add("max_amount", create.MaxAmount)
	params.Add("binding_type", create.BindingType)
	params.Add("identifier", create.Identifier)
	params.Add("billing", create.Billing)
	log.Println("Input request:", params)

	macInput := fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", params.Get("app_id"), params.Get("app_trans_id"), params.Get("binding_data"),
		params.Get("binding_type"), params.Get("identifier"), params.Get("max_amount"), params.Get("req_date"))
	params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, create.MacKey, macInput))

	res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/agreement/bind", params)

	if err != nil {
		// log.Fatal(err)
		log.Println("http POST error: ", err.Error())
		c.JSON(http.StatusBadRequest, "CREATE BINDING ERROR")
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
