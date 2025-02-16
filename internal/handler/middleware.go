package handler

import (
	"net/http"
	"strings"

	"github.com/437d5/merch-store/pkg/token"
	"github.com/gin-gonic/gin"
)

const bearerPrefix = "Bearer "

func (h *Handler) AuthMiddleware(c *gin.Context) {
	const op = "/internal/handler/middleware/AuthMiddleware"

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.logger.Warn("missing token", "op", op)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		h.logger.Warn("invalid token format", "op", op)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
		c.Abort()
		return
	}

	t := strings.TrimPrefix(authHeader, bearerPrefix)

	userId, ok, err := token.ValidateToken(t, h.cfg.JWT.Secret)
	if err != nil {
		h.logger.Error("failed to check token", "op", op, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed parse token"})
		c.Abort()
		return
	}

	if !ok {
		h.logger.Warn("invalid token", "op", op)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	c.Set("user_id", userId)
	c.Next()
}
