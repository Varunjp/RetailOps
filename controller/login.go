package controller

import (
	"net/http"

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


}