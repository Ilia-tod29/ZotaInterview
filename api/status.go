package api

import (
	"ZotaInterview/client"
	"ZotaInterview/util"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type StatusRequest struct {
	OrderID         string `form:"orderID" binding:"required"`
	MerchantOrderID string `form:"merchantOrderID" binding:"required"`
	MerchantID      string `form:"merchantID"`
	Timestamp       string `form:"timestamp"`
	Signature       string `form:"signature"`
}

func (s *Server) checkDepositStatus(ctx *gin.Context) {
	var req StatusRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	s.generateStatusFields(&req)

	c := client.GetZotaClient(s.config.SecretKey)
	respBody, err := s.CallZotaStatusAPI(ctx, c, &req)
	if err != nil {
		if code, errConv := strconv.Atoi(respBody); errConv == nil {
			ctx.JSON(code, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, respBody)
}

func (s *Server) CallZotaStatusAPI(ctx *gin.Context, client *client.ZotaClient, req *StatusRequest) (string, error) {
	qParams := url.Values{
		"merchantID":      {req.MerchantID},
		"merchantOrderID": {req.MerchantOrderID},
		"orderID":         {req.OrderID},
		"timestamp":       {req.Timestamp},
		"signature":       {req.Signature},
	}
	createDepositRes, err := client.Get(ctx, "api/v1/query/order-status", qParams)
	if err != nil {
		return strconv.Itoa(createDepositRes.StatusCode), err
	}

	defer createDepositRes.Body.Close()
	body, readErr := io.ReadAll(createDepositRes.Body)
	if readErr != nil {
		return "", err
	}
	return string(body), nil
}

func (s *Server) generateStatusFields(req *StatusRequest) {
	req.MerchantID = s.config.MerchantId
	req.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	req.Signature = util.GenerateSignature(req.MerchantID + req.MerchantOrderID + req.OrderID + req.Timestamp + s.config.SecretKey)
}
