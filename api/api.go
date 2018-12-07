package api

import (
	"../common"
	"../common/rediscli"
	"../config"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

func submitTask(c *gin.Context) {
	project := c.Param("project")
	url := c.Query("url")
	if strings.Trim(url, " ") == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "invalid url",
		})
	} else {
		header := c.Request.Header
		contentType := header.Get("Content-Type")
		body, _ := ioutil.ReadAll(c.Request.Body)
		message := common.NewMessage(&header, body, project, url)
		defer c.Request.Body.Close()
		rediscli.LPush(fmt.Sprintf("project:%s", project), message.ToJson())
		c.JSON(http.StatusOK, gin.H{
			"project":     project,
			"url":         url,
			"message_id":  message.MeessgeId,
			"contentType": contentType,
		})
	}
}

func info(c *gin.Context) {
	var total_stock int64
	projects := rediscli.Len(rediscli.ListKey("project")...)
	for _, v := range projects {
		total_stock += v
	}
	c.JSON(http.StatusOK, gin.H{
		"total_stock": total_stock,
		"project":     projects,
		"redis_stats": rediscli.Info(),
	})
}

func callback(c *gin.Context) {
	header := c.Request.Header
	contentType := header.Get("Content-Type")
	body, _ := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	fmt.Println("body:", string(body), contentType)
	c.String(http.StatusOK, "1")
}

func Run() {
	router := gin.Default()

	router.GET("/info", info)
	router.POST("/callback", callback)
	router.POST("/task/:project/", submitTask)
	router.GET("/task/:project/", submitTask)
	listen := fmt.Sprintf(":%d", config.API_PROT)
	router.Run(listen)
}
