package api

import (
	"ZotaInterview/client"
	"ZotaInterview/util"
	"encoding/json"
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

	respBody, err := s.CallZotaStatusAPI(ctx, s.zotaClient, &req)
	if err != nil {
		if codeVal, ok := respBody["StatusCode"]; ok {
			if code, ok := codeVal.(int); ok {
				ctx.JSON(code, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, respBody)
}

func (s *Server) CallZotaStatusAPI(ctx *gin.Context, client client.ZotaClientInterface, req *StatusRequest) (map[string]interface{}, error) {
	qParams := url.Values{
		"merchantID":      {req.MerchantID},
		"merchantOrderID": {req.MerchantOrderID},
		"orderID":         {req.OrderID},
		"timestamp":       {req.Timestamp},
		"signature":       {req.Signature},
	}
	createDepositRes, err := client.Get(ctx, "api/v1/query/order-status", qParams)
	if err != nil {
		return map[string]interface{}{
			"StatusCode": createDepositRes.StatusCode,
		}, err
	}

	defer createDepositRes.Body.Close()
	body, readErr := io.ReadAll(createDepositRes.Body)
	if readErr != nil {
		return map[string]interface{}{}, err
	}
	// Unmarshal the response body into a map
	var responseMap map[string]interface{}
	if err = json.Unmarshal(body, &responseMap); err != nil {
		return nil, err
	}
	return responseMap, nil
}

func (s *Server) generateStatusFields(req *StatusRequest) {
	req.MerchantID = s.config.MerchantId
	req.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	req.Signature = util.GenerateSignature(req.MerchantID + req.MerchantOrderID + req.OrderID + req.Timestamp + s.config.SecretKey)
}
