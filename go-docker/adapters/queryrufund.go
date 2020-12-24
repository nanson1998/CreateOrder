// go version go1.11.1 linux/amd64
package adapter

import (
	"net/url"
	"bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    
    "strconv"
	"github.com/gin-gonic/gin"
    "github.com/zpmep/hmacutil" // go get github.com/zpmep/hmacutil
)



func QueryRefund(c *gin.Context) {
    params := make(url.Values)
    params.Add("app_id", appID)
    params.Add("m_refund_id", appTransID)
    params.Add("timestamp",  strconv.Itoa(12345678910))
  
    data := fmt.Sprintf("%v|%v|%v", params.Get("app_id"), params.Get("m_refund_id"), params.Get("timestamp")) 
  params.Add("mac",hmacutil.HexStringEncode(hmacutil.SHA256, key1, data))

  jsonStr, err := json.Marshal(params)
    if err != nil {
        log.Fatal(err)
  }

  res, err := http.Post("https://sb-openapi.zalopay.vn/v2/query_refund", "application/json", bytes.NewBuffer(jsonStr))

    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()

    body, _ := ioutil.ReadAll(res.Body)

    var result map[string]interface{}

    if err := json.Unmarshal(body, &result); err != nil {
        log.Fatal(err)
    }

    c.JSON(200,result)
    
}