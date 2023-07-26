package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestConvertHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/convert", convertHandler)

	t.Run("InvalidInput", func(t *testing.T) {
		w := performRequest(router, "/convert")
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid input")
	})

	t.Run("InvalidAmount", func(t *testing.T) {
		w := performRequest(router, "/convert?source=USD&target=JPY&amount=abc")
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid amount")
	})

	t.Run("ConvertSuccess", func(t *testing.T) {
		w := performRequest(router, "/convert?source=USD&target=TWD&amount=$1,525")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"msg":"success"`)
		assert.Contains(t, w.Body.String(), `"amount":"$46,427.10"`)
	})
}

func performRequest(router *gin.Engine, url string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	router.ServeHTTP(w, req)
	return w
}

func TestAddCommasToNumber(t *testing.T) {
	assert.Equal(t, "123,456.78", addCommasToNumber("123456.78"))
	assert.Equal(t, "1,234,567,890.12", addCommasToNumber("1234567890.12"))
	assert.Equal(t, "0.00", addCommasToNumber("0.00"))
	assert.Equal(t, "1,234.56", addCommasToNumber("1234.56"))
}

func TestGetExchangeRate(t *testing.T) {
	// Test valid exchange rate
	assert.Equal(t, 111.801, getExchangeRate("USD", "JPY"))
	assert.Equal(t, 30.444, getExchangeRate("USD", "TWD"))
	assert.Equal(t, 3.669, getExchangeRate("TWD", "JPY"))

	// Test invalid exchange rate
	assert.Equal(t, 0.0, getExchangeRate("USD", "EUR"))
	assert.Equal(t, 0.0, getExchangeRate("TWD", "EUR"))
}

func TestConversionAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/convert", convertHandler)

	t.Run("ConvertUSDToTWD", func(t *testing.T) {
		w := performRequest(router, "/convert?source=USD&target=TWD&amount=$1,525")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"msg":"success"`)
		assert.Contains(t, w.Body.String(), `"amount":"$46,427.10"`)
	})

	t.Run("ConvertJPYToTWD", func(t *testing.T) {
		w := performRequest(router, "/convert?source=JPY&target=TWD&amount=¥1,000")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"msg":"success"`)
		assert.Contains(t, w.Body.String(), `"amount":"$269.56"`)
	})
}
