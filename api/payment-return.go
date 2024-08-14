package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

const UserMessage = "Thank you for using us!"

type PaymentReturn struct {
	UserMessage     string `json:"userMessage"`
	Status          string `form:"status" binding:"required" json:"status"`
	OrderID         string `form:"orderID" binding:"required" json:"orderID"`
	MerchantOrderID string `form:"merchantOrderID" binding:"required" json:"merchantOrderID"`
	ErrorMessage    string `form:"errorMessage" json:"errorMessage"`
}

func (s *Server) paymentReturn(ctx *gin.Context) {
	var req PaymentReturn
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	req.UserMessage = UserMessage

	message, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, string(message))
}
