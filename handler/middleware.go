package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/model"
	"strings"
)

func (h *Handler) CorsMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

func (h *Handler) userIdentity(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		h.logger.Errorf("userIdentity:empty auth header")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Message: "empty auth header"})
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		h.logger.Errorf("userIdentity:invalid auth header")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Message: "invalid auth header"})
		return
	}
	if len(headerParts[1]) == 0 {
		h.logger.Errorf("userIdentity:token is empty")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Message: "token is empty"})
		return
	}
	userPerms, err := h.service.AppUser.ParseToken(headerParts[1])
	if err != nil {
		h.logger.Errorf("userIdentity:%s", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.Set("perms", userPerms.Permissions)
	ctx.Set("role", userPerms.Role)
	ctx.Set("userId", userPerms.UserId)
}
