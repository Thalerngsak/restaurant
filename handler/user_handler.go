package handler

import (
	"fmt"
	"github.com/Thalerngsak/restaurant/model"
	"github.com/Thalerngsak/restaurant/service"
	"github.com/Thalerngsak/restaurant/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type userHandler struct {
	userService service.UserService
	tokenMaker  token.JWTMaker
}

func NewUserHandler(userService service.UserService, tokenMaker token.JWTMaker) *userHandler {
	return &userHandler{userService: userService, tokenMaker: tokenMaker}
}

func (h userHandler) Login(c *gin.Context) {
	var user model.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.userService.GetByUserName(user.Username)
	fmt.Printf("user Id : %v", results.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := h.tokenMaker.GenerateToken(results.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}
