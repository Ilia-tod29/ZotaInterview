package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	UserMessage          = "Thank you for using us!"
	PossibleStatusesLink = "https://doc.zota.com/deposit/1.0/#order-statuses"
)

type PaymentReturn struct {
	UserMessage      string `json:"userMessage"`
	Status           string `form:"status" binding:"required" json:"status"`
	OrderID          string `form:"orderID" binding:"required" json:"orderID"`
	MerchantOrderID  string `form:"merchantOrderID" binding:"required" json:"merchantOrderID"`
	ErrorMessage     string `form:"errorMessage" json:"errorMessage"`
	PossibleStatuses string `json:"possibleStatuses"`
}

// paymentReturn is the handler being called when we hit the defined by the application /payment-return endpoint
func (s *Server) paymentReturn(ctx *gin.Context) {
	var req PaymentReturn
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	req.UserMessage = UserMessage

	if req.ErrorMessage != "" {
		req.PossibleStatuses = PossibleStatusesLink
	}

	ctx.JSON(http.StatusOK, req)
}
