package api

import (
	"KeywordCollection/constant"
	"KeywordCollection/dao/table"
	"KeywordCollection/global"
	"KeywordCollection/server/taskQueue"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// RegiserService ...
func RegiserService(engine *gin.Engine) {
	api := engine.Group("/api")
	api.GET("/query", queryHander)
}

func queryHander(c *gin.Context) {
	clientId := c.Query("clientID")
	if clientId == "" {
		errText := "非法查询请求"
		log.WithFields(log.Fields{"clientID": clientId}).Error(errText)
		c.JSON(http.StatusOK, gin.H{"c": 0, "m": errText, "d": gin.H{}})
		return
	}

	client, ok := global.QueryClients.Load(clientId)
	if !ok {
		errText := "非法查询请求"
		log.WithFields(log.Fields{"clientID": clientId}).Error(errText)
		c.JSON(http.StatusOK, gin.H{"c": 0, "m": errText, "d": gin.H{}})
		return
	}

	session := client.(*melody.Session)
	if session.IsClosed() {
		errText := "客户端已断开"
		log.WithFields(log.Fields{"clientID": clientId}).Error(errText)
		c.JSON(http.StatusOK, gin.H{"c": 0, "m": errText, "d": gin.H{}})
		return
	}

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{"c": 0, "m": "id is empty", "d": gin.H{}})
		return
	}

	engine := c.DefaultQuery("engine", constant.EngineBaidu)
	if engine == "" {
		c.JSON(http.StatusOK, gin.H{"c": 0, "m": "engine is empty", "d": gin.H{}})
		return
	}

	device := c.DefaultQuery("device", constant.DevicePc)
	if device == "" {
		c.JSON(http.StatusOK, gin.H{"c": 0, "m": "device is empty", "d": gin.H{}})
		return
	}

	keyword := c.Query("wd")
	if keyword == "" {
		c.JSON(http.StatusOK, gin.H{"c": 0, "m": "wd is empty", "d": gin.H{}})
		return
	}

	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusOK, gin.H{"c": 0, "m": "domain is empty", "d": gin.H{}})
		return
	}

	page := c.DefaultQuery("page", "1")
	if page == "" {
		c.JSON(http.StatusOK, gin.H{"c": 0, "m": "page is empty", "d": gin.H{}})
		return
	}

	keywordId, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, gin.H{"c": 0, "m": "id error", "d": gin.H{}})
		return
	}

	ipage, err := strconv.Atoi(page)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, gin.H{"c": 0, "m": "page error", "d": gin.H{}})
		return
	}

	qctx := &table.QueryContext{
		Priority:     1,
		KeywordId:    keywordId,
		Engine:       engine,
		Type:         device,
		Keyword:      keyword,
		Page:         ipage,
		ClientId:     clientId,
		SerialNumber: global.GetIdWorker().String(),
	}

	err = qtasks.ContentRetrievalServer().Put(qctx)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, gin.H{"c": 0, "m": "idomainId error", "d": gin.H{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"c": 1, "m": "查询提交成功", "d": gin.H{"status": 1, "id": id}})
}
