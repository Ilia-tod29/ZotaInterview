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

	c := client.GetZotaClient(s.config.SecretKey)
	respBody, err := s.callZotaDepositAPI(ctx, c, &req)
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

func (s *Server) callZotaDepositAPI(ctx *gin.Context, client *client.ZotaClient, req *CreateDepositRequest) (string, error) {
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	createDepositRes, err := client.Post(ctx, fmt.Sprintf("api/v1/deposit/request/%s", s.config.EndpointId), jsonBytes)
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
	req.Signature = util.GenerateSignature(s.config.EndpointId + req.MerchantOrderID + req.OrderAmount + req.CustomerEmail + s.config.SecretKey)
	s.generateCheckoutUrl(req)

	return nil
}

// TODO: WHY is it needed? https://doc.zota.com/deposit/1.0/#deposit-request
func (s *Server) generateCheckoutUrl(req *CreateDepositRequest) {
	req.CheckoutUrl = s.config.EndpointId
}
