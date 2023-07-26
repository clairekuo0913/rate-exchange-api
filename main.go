package main

import (
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var exchangeRates = map[string]map[string]float64{
	"TWD": {
		"TWD": 1,
		"JPY": 3.669,
		"USD": 0.03281,
	},
	"JPY": {
		"TWD": 0.26956,
		"JPY": 1,
		"USD": 0.00885,
	},
	"USD": {
		"TWD": 30.444,
		"JPY": 111.801,
		"USD": 1,
	},
}

func main() {
	router := gin.Default()

	router.GET("/convert", convertHandler)

	router.Run(":8080")
}

func convertHandler(c *gin.Context) {
	source := c.DefaultQuery("source", "")
	target := c.DefaultQuery("target", "")
	amount := c.DefaultQuery("amount", "")

	if source == "" || target == "" || amount == "" {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Extract the money value from the amount using regular expression
	re := regexp.MustCompile(`[0-9,.]+`)
	matches := re.FindAllString(amount, -1)
	if len(matches) == 0 {
		c.JSON(400, gin.H{"error": "Invalid amount"})
		return
	}
	println(matches)
	amountStr := strings.ReplaceAll(matches[0], ",", "") // Remove any commas in the amount string
	amountVal, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid amount"})
		return
	}
	println("val = ", amountStr)

	rate := getExchangeRate(source, target)
	convertedAmount := amountVal * rate

	// Round the converted amount to 2 decimal places
	convertedAmount = math.Round(convertedAmount*100) / 100

	// Format the converted amount with commas as thousand separators
	convertedAmountStr := strconv.FormatFloat(convertedAmount, 'f', 2, 64)
	convertedAmountStr = addCommasToNumber(convertedAmountStr)

	c.JSON(200, gin.H{
		"msg":    "success",
		"amount": "$" + convertedAmountStr,
	})
}

func getExchangeRate(source, target string) float64 {
	if rate, ok := exchangeRates[source][target]; ok {
		return rate
	}
	return 0
}

func addCommasToNumber(numStr string) string {
	parts := strings.Split(numStr, ".")
	intPart := parts[0]
	var formattedNum string

	for i := len(intPart) - 1; i >= 0; i-- {
		formattedNum = string(intPart[i]) + formattedNum
		if (len(intPart)-i)%3 == 0 && i != 0 {
			formattedNum = "," + formattedNum
		}
	}

	if len(parts) == 2 {
		formattedNum += "." + parts[1]
	}

	return formattedNum
}
