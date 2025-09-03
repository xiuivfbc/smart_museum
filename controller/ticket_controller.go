package controller

import (
	"group_ten_server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTicket(c *gin.Context) {
	var req struct {
		ID   int    `json:"id"`
		Time string `json:"time"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	service.CreateTicketService(c, req)
}

func ListTicket(c *gin.Context) {
	var req struct {
		ID int `json:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	service.ListTicketService(c, req)
}

func DeleteTicket(c *gin.Context) {
	var req struct {
		ID int `json:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	service.DeleteTicketService(c, req)
}

func UpdateTicket(c *gin.Context) {
	var req struct {
		ID   int    `json:"id"`
		Time string `json:"time"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	service.UpdateTicketService(c, req)
}

func GetTodayEntryCount(c *gin.Context) {
	service.GetTodayEntryCountService(c)
}

func GetAllEntryCounts(c *gin.Context) {
	service.GetAllEntryCountsService(c)

}

func UseTicket(c *gin.Context) {
	service.UseTicketService(c)
}
