package adapter

import (
	"github.com/gin-gonic/gin"

	"github.com/zpmep/hmacutil"
	"log"
	"net/url"
)

type CallbackOrder struct {
	Data string `json:"data"`
	Mac  string `json:"mac"`
	Type string `json:"type"`
	Key2 string `json:"key2"`
}

func HandleCallback(c *gin.Context) {
	var cbOrder *CallbackOrder
	params := make(url.Values)
	if err := c.Bind(&cbOrder); err == nil {
		params.Add("data", cbOrder.Data)
		params.Add("mac", cbOrder.Mac)
		params.Add("type", cbOrder.Type)
		params.Add("key2", cbOrder.Key2)
	} else {
		log.Println("binding fail")
	}

	rqmac := hmacutil.HexStringEncode(hmacutil.SHA256, cbOrder.Key2, cbOrder.Data)
	result := make(map[string]interface{})
	if rqmac != cbOrder.Mac {
		// callback lỗi
		result["return_code"] = -1
		result["return_message"] = "mac not equal"

	} else {
		//callback hợp lệ
		result["return_code"] = 1
		result["return_message"] = "success"
	}

	c.JSON(200, result)

}
