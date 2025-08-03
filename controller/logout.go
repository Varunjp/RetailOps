package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {

	// Clear the session
	c.SetCookie("JWT-User", "", -1, "/", "", false, true)

	// Redirect to the login page
	c.Redirect(http.StatusSeeOther, "/login")
}