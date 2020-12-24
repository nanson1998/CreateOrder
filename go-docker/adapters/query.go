package adapter

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
	"net/http"
	"github.com/gin-gonic/gin"

    "github.com/zpmep/hmacutil" // go get github.com/zpmep/hmacutil
)

var (
    appTransID = "201224_112763" // Input your app trans id
)


func QueryOrder(c *gin.Context) {
    data := fmt.Sprintf("%v|%s|%s", appID, appTransID, key1) // appid|apptransid|key1
    params := map[string]interface{}{
        "app_id":       appID,
        "app_trans_id": appTransID,
        "mac":          hmacutil.HexStringEncode(hmacutil.SHA256, key1, data),
    }

    jsonStr, err := json.Marshal(params)
    if err != nil {
        log.Fatal(err)
    }

    res, err := http.Post("https://sb-openapi.zalopay.vn/v2/query", "application/json", bytes.NewBuffer(jsonStr))

    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()

    body, _ := ioutil.ReadAll(res.Body)

    var result map[string]interface{}

    if err := json.Unmarshal(body, &result); err != nil {
        log.Fatal(err)
    }

    c.JSON(200, result)
    
}