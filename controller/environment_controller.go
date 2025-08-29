package controller

import (
	"group_ten_server/dao"
	"group_ten_server/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateEnvironment 新增环境
func CreateEnvironment(c *gin.Context) {
	var env model.Environment
	if err := c.ShouldBindJSON(&env); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := dao.CreateEnvironment(&env); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "data": env})
}

// GetEnvironmentByName 查询单个环境
func GetEnvironmentByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusOK, gin.H{"error": "name不能为空"})
		return
	}
	env, err := dao.GetEnvironmentByName(name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "未找到"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": env})
}

// GetAllEnvironments 查询所有环境
func GetAllEnvironments(c *gin.Context) {
	envs, err := dao.GetAllEnvironments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": envs})
}

// UpdateEnvironmentByName 根据name更新环境
func UpdateEnvironmentByName(c *gin.Context) {
	name := c.Param("name")
	var update model.Environment
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := dao.UpdateEnvironmentByName(name, &update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": update})
}

// DeleteEnvironmentByName 根据name删除环境
func DeleteEnvironmentByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusOK, gin.H{"error": "name不能为空"})
		return
	}
	if err := dao.DeleteEnvironmentByName(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
