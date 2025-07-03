package controller

import (
	"net/http"
	db "retialops/DB"
	"retialops/helper"
	"retialops/models"
	"retialops/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func LoginPage(c *gin.Context){
	c.HTML(http.StatusOK,"login.html",nil)

}

func Login(c *gin.Context){

	var input struct{
		EmpID		string 	`form:"employeeId" binding:"required"`
		Password 	string 	`form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil{
		c.HTML(http.StatusBadRequest,"login.html",gin.H{"error":"Failed to get employeeID or password, please try again later."})
		return 
	}

	if strings.TrimSpace(input.EmpID) == "" || strings.TrimSpace(input.Password) == "" {
		c.HTML(http.StatusBadRequest,"login.html",gin.H{"error":"Invalid employeeID or password"})
		return 
	}

	var user models.User
	if err := db.Db.Where("emp_id = ?",input.EmpID).First(&user).Error; err != nil{
		c.HTML(http.StatusBadRequest,"login.html",gin.H{"error":"User not found"})
		return
	}

	if user.Status != "Active"{
		c.HTML(http.StatusUnauthorized,"login.html",gin.H{"error":"Account blocked contact administrator"})
		return 
	}

	if !utils.CheckPasswordHash(input.Password,user.Password){
		c.HTML(http.StatusUnauthorized,"login.html",gin.H{"error":"Password incorrect, please enter correct password"})
		return 
	}

	token, err := helper.CreateToken(user.EmpID,user.SuperUser)

	if err != nil{
		c.HTML(http.StatusInternalServerError,"login.html",gin.H{"error":"Failed to create token,please try again later."})
		return 
	}

	c.SetCookie("JWT-User",token,3600,"/","",false,true)
	
	c.Header("Cache-Control","no-store, no-cache, must-revalidate, max-age=0")

	c.Redirect(http.StatusFound,"/")
}