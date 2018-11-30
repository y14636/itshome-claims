package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/y14636/itshome-claims/claims"
	"github.com/y14636/itshome-claims/logging"
	"github.com/y14636/itshome-claims/model"
	"github.com/y14636/itshome-claims/modifiedclaims"
	"github.com/y14636/itshome-claims/searchclaims"
)

type LogMessage struct {
	Additional []string `json:"additional"`
	FileName   string   `json:"fileName"`
	Level      int      `json:"level"`
	LineNumber string   `json:"lineNumber"`
	Message    string   `json:"message"`
	Timestamp  string   `json:"timestamp"`
}

//GetClaimsResultsHandler returns claim items from search
func GetClaimsResultsHandler(c *gin.Context) {
	searchString := c.Param("search")
	c.JSON(http.StatusOK, searchclaims.GetResults(searchString))
}

func AddClaimsHandler(c *gin.Context) {
	claimsDataString := c.Param("claimsData")
	if err := modifiedclaims.AddClaims(claimsDataString); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

// GetClaimsListHandler returns all current claim items
func GetClaimsListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, claims.Get())
}

// GetClaimsListHandler returns all current claim items
func GetClaimsListByIdHandler(c *gin.Context) {
	idString := c.Param("claimsId")
	c.JSON(http.StatusOK, claims.GetListById(idString))
}

// GetModifiedClaimsListHandler returns all current claim items
func GetModifiedClaimsListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, modifiedclaims.GetModifiedClaims())
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

func LogWebMessages(c *gin.Context) {
	message, statusCode, err := convertHTTPBodyToLogging(c.Request.Body)
	fmt.Println("LogWebMessages=", message)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(statusCode, gin.H{"messages": logging.LogWebMessages("ui-" + message.Message + message.Additional[0])})
}

func convertHTTPBodyToLogging(httpBody io.ReadCloser) (LogMessage, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return LogMessage{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToLogging(body)
}

func convertJSONBodyToLogging(jsonBody []byte) (LogMessage, int, error) {
	var message LogMessage
	err := json.Unmarshal(jsonBody, &message)
	if err != nil {
		return LogMessage{}, http.StatusBadRequest, err
	}
	return message, http.StatusOK, nil
}

func convertHTTPBodyToClaims(httpBody io.ReadCloser) (structs.Claims, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return structs.Claims{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToClaims(body)
}

func convertJSONBodyToClaims(jsonBody []byte) (structs.Claims, int, error) {
	var claimsItem structs.Claims
	err := json.Unmarshal(jsonBody, &claimsItem)
	if err != nil {
		return structs.Claims{}, http.StatusBadRequest, err
	}
	return claimsItem, http.StatusOK, nil
}
