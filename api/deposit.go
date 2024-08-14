package api

import (
	"ZotaInterview/client"
	"ZotaInterview/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

type CreateDepositRequest struct {
	MerchantOrderID           string `json:"merchantOrderID"`
	MerchantOrderDesc         string `json:"merchantOrderDesc" binding:"required"`
	OrderAmount               string `json:"orderAmount" binding:"required"`
	OrderCurrency             string `json:"orderCurrency" binding:"required"`
	CustomerEmail             string `json:"customerEmail" binding:"required"`
	CustomerFirstName         string `json:"customerFirstName" binding:"required"`
	CustomerLastName          string `json:"customerLastName" binding:"required"`
	CustomerAddress           string `json:"customerAddress" binding:"required"`
	CustomerCountryCode       string `json:"customerCountryCode" binding:"required"`
	CustomerCity              string `json:"customerCity" binding:"required"`
	CustomerState             string `json:"customerState"`
	CustomerZipCode           string `json:"customerZipCode" binding:"required"`
	CustomerPhone             string `json:"customerPhone" binding:"required"`
	CustomerIP                string `json:"customerIP" binding:"required"`
	CustomerPersonalID        string `json:"customerPersonalID"`
	CustomerBankCode          string `json:"customerBankCode"`
	CustomerBankAccountNumber string `json:"customerBankAccountNumber"`
	RedirectUrl               string `json:"redirectUrl" binding:"required"`
	CallbackUrl               string `json:"callbackUrl"`
	CheckoutUrl               string `json:"checkoutUrl"`
	CustomParam               string `json:"customParam"`
	Language                  string `json:"language"`
	Signature                 string `json:"signature"`
}

// depositMoney is the handler being called when we hit the defined by the application /deposit endpoint
func (s *Server) depositMoney(ctx *gin.Context) {
	var req CreateDepositRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.validateDepositRequest(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	respBody, err := s.callZotaDepositAPI(ctx, s.zotaClient, &req)
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

// callZotaDepositAPI performs the POST call in order to create a new deposit
func (s *Server) callZotaDepositAPI(ctx *gin.Context, client client.ZotaClientInterface, req *CreateDepositRequest) (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	createDepositRes, err := client.Post(ctx, fmt.Sprintf("api/v1/deposit/request/%s", s.config.EndpointId), jsonBytes)
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

	var responseMap map[string]interface{}
	if err = json.Unmarshal(body, &responseMap); err != nil {
		return nil, err
	}
	return responseMap, nil
}

// validateDepositRequest validates the presented data and generates the needed data for performing the request, but not provided by the user
func (s *Server) validateDepositRequest(req *CreateDepositRequest) error {
	_, err := strconv.ParseFloat(req.OrderAmount, 64)
	if err != nil {
		return err
	}
	_, err = strconv.Atoi(req.CustomerZipCode)
	if err != nil {
		return err
	}

	req.MerchantOrderID = util.GenerateNonce(10)
	s.generateCheckoutUrl(req)
	req.Signature = util.GenerateSignature(s.config.EndpointId + req.MerchantOrderID + req.OrderAmount + req.CustomerEmail + s.config.SecretKey)

	return nil
}

// generateCheckoutUrl sets the checkout url
func (s *Server) generateCheckoutUrl(req *CreateDepositRequest) {
	req.CheckoutUrl = s.config.HTTPServerAddress + "/deposit"
}
