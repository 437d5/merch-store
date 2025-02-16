package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/437d5/merch-store/internal/config"
	"github.com/437d5/merch-store/internal/service"
	"github.com/437d5/merch-store/pkg/token"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService *service.UserService
	marketService *service.MarketService
	transactionService *service.TransactionService
	logger *slog.Logger
	cfg config.Config
}

func NewHandler(
	userService *service.UserService,
	marketService *service.MarketService,
	transactionService *service.TransactionService,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		userService: userService,
		marketService: marketService,
		transactionService: transactionService,
		logger: logger,
	}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	
	api.GET("/info", h.AuthMiddleware, h.GetUserInfo)
	api.POST("/sendCoin", h.AuthMiddleware, h.SendCoin)
	api.GET("/buy/:item", h.AuthMiddleware, h.BuyItem)
	api.POST("/auth", h.Auth)
}

func (h *Handler) GetUserInfo(c *gin.Context) {
	const op = "/internal/handler/handlers/GetUserInfo"

	userId := c.GetInt("user_id")
	u, err := h.userService.UserInfo(c.Request.Context(), userId)
	if err != nil {
		h.logger.Error("failed get user info", "op", op, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	
	tList, err := h.transactionService.GetTransactionsByUser(c.Request.Context(), userId)
	if err != nil {
		h.logger.Error("failed get transaction list", "op", op, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	response := gin.H{
		"coins": u.Coins,
		"inventory": formatInventory(u.Inventory),
		"coinHistory": formatTranscations(tList, userId),
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) SendCoin(c *gin.Context) {
	const op = "/internal/handler/handlers/SendCoin"
	
	userId := c.GetInt("user_id")

	var req struct {
		ToUsername string `json:"toUser" binding:"required"`
		Amount int `json:"amount" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		h.logger.Warn("invalid request", "op", op, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	err := h.transactionService.TransferCoins(
		c.Request.Context(), userId, req.Amount, req.ToUsername,
	)
	if err != nil {
		if errors.Is(err, service.ErrNotEnoughCoins) {
			h.logger.Warn("not enough coins", "op", op, "id", userId)
			c.JSON(http.StatusBadRequest, gin.H{"errors": "not enough coins"})
			return
		} else if errors.Is(err, service.ErrInvalidAmount) {
			h.logger.Warn("invalid amount", "op", op, "amount", req.Amount)
			c.JSON(http.StatusBadRequest, gin.H{"errors": "invalid amount"})
			return
		} else {
			h.logger.Error("failed transfer coins", "op", op, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"errors": err})
			return
		}
	}

	c.Status(http.StatusOK)
}

func (h *Handler) BuyItem(c *gin.Context) {
	const op = "/internal/handler/handlers/BuyItem"
	
	userId := c.GetInt("user_id")

	err := h.marketService.BuyMerch(c.Request.Context(), userId, c.Param("item"))
	if err != nil {
		if errors.Is(err, service.ErrNotEnoughCoins) {
			h.logger.Warn("not enough coins to buy merch", "op", op, "id", userId)
			c.JSON(http.StatusBadRequest, gin.H{"errors": "not enough money"})
			return
		} else {
			h.logger.Error("failed buy item", "op", op, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "failed buy merch"})
			return
		}
	}

	c.Status(http.StatusOK)
}

func (h *Handler) Auth(c *gin.Context) {
	const op = "/internal/handler/handlers/Auth"

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request", "op", op, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	u, err := h.userService.AuthUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.logger.Warn("authentification failed", "op", op, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentification failed"})
		return
	}

	t, err := token.CreateToken(u.Id, u.Name, h.cfg.JWT.Secret, token.JWTExpAt)
	if err != nil {
		h.logger.Error("failed to create JWT", "op", op, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	h.logger.Info("user authentificated", "op", op, "username", req.Username)

	c.Header("Authorization", bearerPrefix+t)
	c.JSON(http.StatusOK, gin.H{"token": t})
}