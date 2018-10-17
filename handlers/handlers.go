package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/y14636/itshome-claims/claims"
	"github.com/y14636/itshome-claims/modifiedclaims"
)

// GetClaimsResultsHandler returns all current claim items
// func GetClaimsResultsHandler(c *gin.Context) {
// 	c.JSON(http.StatusOK, claims.GetResults())
// }

// GetClaimsListHandler returns all current claim items
func GetClaimsListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, claims.Get())
}

// GetModifiedClaimsListHandler returns all current claim items
func GetModifiedClaimsListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, modifiedclaims.GetModifiedClaims())
}

// AddClaimsHandler adds a new claim to the claims list
func AddClaimsHandler(c *gin.Context) {
	claimsItem, statusCode, err := convertHTTPBodyToClaims(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(statusCode, gin.H{"id": claims.Add(claimsItem.ClaimType, claimsItem.ServiceId, claimsItem.ReceiptDate, claimsItem.FromDate, claimsItem.ToDate, claimsItem.PlaceOfService, claimsItem.ProviderId,
		claimsItem.ProviderType, claimsItem.ProviderSpecialty, claimsItem.ProcedureCode, claimsItem.DiagnosisCode,
		claimsItem.NetworkIndicator, claimsItem.SubscriberId, claimsItem.PatientAccountNumber, claimsItem.SccfNumber,
		claimsItem.RevenueCode, claimsItem.BillType, claimsItem.Modifier, claimsItem.PlanCode, claimsItem.SfMessageCode,
		claimsItem.PricingMethod, claimsItem.PricingRule, claimsItem.DeliveryMethod, claimsItem.InputDate, claimsItem.FileName)})
}

// AddModifiedClaimsHandler adds a new claim to the modified claims list
func AddModifiedClaimsHandler(c *gin.Context) {
	claimsItem, statusCode, err := convertHTTPBodyToClaims(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(statusCode, gin.H{"id": modifiedclaims.AddModifiedClaim(claimsItem.ClaimType, claimsItem.FromDate, claimsItem.ToDate, claimsItem.PlaceOfService, claimsItem.ProviderId,
		claimsItem.ProviderType, claimsItem.ProviderSpecialty, claimsItem.ProcedureCode, claimsItem.DiagnosisCode,
		claimsItem.NetworkIndicator, claimsItem.SubscriberId, claimsItem.PatientAccountNumber, claimsItem.SccfNumber,
		claimsItem.RevenueCode, claimsItem.BillType, claimsItem.Modifier, claimsItem.PlanCode, claimsItem.SfMessageCode,
		claimsItem.PricingMethod, claimsItem.PricingRule, claimsItem.DeliveryMethod, claimsItem.InputDate, claimsItem.FileName)})
}

// DeleteClaimsHandler will delete a specified claim based on user http input
func DeleteClaimsHandler(c *gin.Context) {
	claimsID := c.Param("id")
	if err := modifiedclaims.Delete(claimsID); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

func convertHTTPBodyToClaims(httpBody io.ReadCloser) (claims.Claims, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return claims.Claims{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToClaims(body)
}

func convertJSONBodyToClaims(jsonBody []byte) (claims.Claims, int, error) {
	var claimsItem claims.Claims
	err := json.Unmarshal(jsonBody, &claimsItem)
	if err != nil {
		return claims.Claims{}, http.StatusBadRequest, err
	}
	return claimsItem, http.StatusOK, nil
}
