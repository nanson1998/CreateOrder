package adapter

import (
	"math/rand"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/zpmep/hmacutil" // go get github.com/zpmep/hmacutil
)
func Refund(c *gin.Context) {
    params := make(url.Values)
    params.Add("app_id", appID)
    params.Add("zp_trans_id", "190508000000017")
    params.Add("amount", "50000")
    params.Add("description", "ZaloPay Refund Demo")

    now := time.Now()
    timestamp := now.UnixNano() / int64(time.Millisecond) // Miliseconds
    params.Add("timestamp", strconv.FormatInt(timestamp, 10))

    uid := fmt.Sprintf("%d%d", timestamp, 111+rand.Intn(888))

    params.Add("m_refund_id", fmt.Sprintf("%02d%02d%02d_%v_%v", now.Year()%100, int(now.Month()), now.Day(), appID, uid))

    // app_id|zp_trans_id|amount|description|timestamp
    data := fmt.Sprintf("%v|%v|%v|%v|%v", appID, params.Get("zp_trans_id"), params.Get("amount"), params.Get("description"), params.Get("timestamp"))
    params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, key1, data))

    

    // Content-Type: application/x-www-form-urlencoded
    res, err := http.PostForm("https://sb-openapi.zalopay.vn/v2/refund", params)

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