package service

import (
	"fmt"
	"group_ten_server/config"
	"group_ten_server/dao"
	"group_ten_server/model"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func CreateTicketService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		ID   int    `json:"id"`
		Time string `json:"time"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}
	if reqStruct.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "缺少id参数"})
		return
	}
	t, err := time.Parse(time.RFC3339, reqStruct.Time)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "时间格式错误，需为RFC3339格式"})
		return
	}
	saveDir := config.Conf.GetString("server.path") + "/qrcodes"
	content := "http://" + config.Conf.GetString("gin.ip") + ":" + config.Conf.GetString("gin.port") + "/ticket/use?ticket_id=" + strconv.Itoa(reqStruct.ID)
	fileName := fmt.Sprintf("qrcode_%d.png", reqStruct.ID)
	path := filepath.Join(saveDir, fileName)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		log.Println("创建目录失败:", err)
		return
	}
	if err := qrcode.WriteFile(content, qrcode.Medium, 256, path); err != nil {
		log.Println("生成二维码失败:", err)
		return
	}
	ticket := model.Ticket{
		ID:   reqStruct.ID,
		Time: t,
		Path: fileName,
	}
	if err := dao.CreateTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "ticket": ticket, "path": ticket.Path})
}

func ListTicketService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		ID int `json:"id"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}
	log.Println("ListTicket request:", reqStruct)
	if reqStruct.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "缺少id参数"})
		return
	}
	ticket, err := dao.GetTicketByID(reqStruct.ID)
	if err != nil || ticket == nil {
		c.JSON(http.StatusOK, gin.H{"error": "未找到该ticket"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ticket": ticket})
}

func DeleteTicketService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		ID int `json:"id"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}
	if reqStruct.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "缺少id参数"})
		return
	}
	if err := dao.DeleteTicketByID(reqStruct.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func UpdateTicketService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		ID   int    `json:"id"`
		Time string `json:"time"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}
	if reqStruct.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "缺少id参数"})
		return
	}
	t, err := time.Parse(time.RFC3339, reqStruct.Time)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "时间格式错误，需为RFC3339格式"})
		return
	}
	fileName := fmt.Sprintf("qrcode_%d.png", reqStruct.ID)
	ticket := model.Ticket{
		ID:   reqStruct.ID,
		Time: t,
		Path: fileName,
	}
	if err := dao.UpdateTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "ticket": ticket})
}

func GetTodayEntryCountService(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"count": config.EntryNum})
}

func GetAllEntryCountsService(c *gin.Context) {
	counts, err := dao.GetAllEntryCounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": counts})
}

func UseTicketService(c *gin.Context) {
	id := c.Query("ticket_id")
	ticket_id, _ := strconv.Atoi(id)
	// 先查票据是否存在
	ticket, err := dao.GetTicketByID(ticket_id)
	if err != nil || ticket == nil {
		// 返回 damie.png
		damiePath := filepath.Join(config.Conf.GetString("filepath2"))
		c.File(damiePath)
		return
	}
	// 删除二维码图片
	saveDir := config.Conf.GetString("server.path") + "/qrcodes"
	fileName := fmt.Sprintf("qrcode_%d.png", ticket_id)
	qrPath := filepath.Join(saveDir, fileName)
	if err := os.Remove(qrPath); err != nil {
		log.Println("删除二维码失败:", err)
		// 不阻断后续票据删除
	}
	if err := dao.DeleteTicketByID(ticket_id); err != nil {
		log.Println("删除票失败:", err)
		return
	}
	config.EntryNum++
	c.File(config.Conf.GetString("filepath"))
}
