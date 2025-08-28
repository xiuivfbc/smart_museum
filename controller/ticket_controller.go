package controller

import (
	"group_ten_server/dao"
	"group_ten_server/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTicket(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier"` // 邮箱或电话
		Time       string `json:"time"`       // 时间字符串
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	userID, err := dao.GetUserIDByEmailOrPhone(req.Identifier)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}
	// 解析时间
	t, err := time.Parse(time.RFC3339, req.Time)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "时间格式错误，需为RFC3339格式"})
		return
	}
	ticket := model.Ticket{
		ID:   userID,
		Time: t,
	}
	if err := dao.CreateTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "ticket": ticket})
}

func ListTickets(c *gin.Context) {
	identifier := c.Query("identifier")
	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少identifier参数"})
		return
	}
	userID, err := dao.GetUserIDByEmailOrPhone(identifier)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}
	tickets, err := dao.ListTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	// 只返回该用户的ticket
	var userTickets []model.Ticket
	for _, t := range tickets {
		if t.ID == userID {
			userTickets = append(userTickets, t)
		}
	}
	c.JSON(http.StatusOK, gin.H{"tickets": userTickets})
}

func DeleteTicket(c *gin.Context) {
	identifier := c.Query("identifier")
	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少identifier参数"})
		return
	}
	userID, err := dao.GetUserIDByEmailOrPhone(identifier)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}
	if err := dao.DeleteTicketByID(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func UpdateTicket(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier"`
		Time       string `json:"time"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	userID, err := dao.GetUserIDByEmailOrPhone(req.Identifier)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}
	t, err := time.Parse(time.RFC3339, req.Time)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "时间格式错误，需为RFC3339格式"})
		return
	}
	ticket := model.Ticket{
		ID:   userID,
		Time: t,
	}
	if err := dao.UpdateTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "ticket": ticket})
}
