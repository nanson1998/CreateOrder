package autodebit

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func GetBinding(c *gin.Context) {
	bindingToken := c.Param("binding_token")
	if bindingToken == "" {
		log.Println("Missing Token")
	}
	fmt.Println("binding_token:", bindingToken)
	res, err := http.Get("https://sb-openapi.zalopay.vn/v2/agreement/binding/" + bindingToken)
	if err != nil {
		// log.Fatal(err)
		log.Println("http GET error: ", err.Error())
		c.JSON(http.StatusBadRequest, "GET BINDING ERROR")
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
