// go version go1.11.1 linux/amd64
package adapter

import (
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


func GetListBank(c *gin.Context) {
    params := make(url.Values)
    params.Add("appid", "2553")
    params.Add("reqtime", strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)) // miliseconds

    data := fmt.Sprintf("%v|%v", params.Get("appid"), params.Get("reqtime")) //appid|reqtime
    params.Add("mac", hmacutil.HexStringEncode(hmacutil.SHA256, key1, data))

    res, err := http.Get("https://sbgateway.zalopay.vn/api/getlistmerchantbanks?" + params.Encode())

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