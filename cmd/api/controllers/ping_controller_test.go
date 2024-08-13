package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	assert := assert.New(t)
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	r.GET("/ping", NewPingController().Ping)
	c.Request, _ = http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, c.Request)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	var response struct {
		Code        int    `json:"code"`
		Message     string `json:"message"`
	}

	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(200, response.Code)
	assert.Equal("pong", response.Message)
}
