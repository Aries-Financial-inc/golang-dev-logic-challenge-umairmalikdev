package controllers

import (
	"math"
	"github.com/gofiber/fiber/v2"
	"github.com/umairmalik/fiber-options-analysis/models"
)

// XYValue represents a point on the profit/loss chart.
type XYValue struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// AnalysisResponse encapsulates the response data for the options analysis.
type AnalysisResponse struct {
	XYValues        []XYValue `json:"xy_values"`
	MaxProfit       float64   `json:"max_profit"`
	MaxLoss         float64   `json:"max_loss"`
	BreakEvenPoints []float64 `json:"break_even_points"`
}

// AnalyzeOptionsContracts handles the analysis of options contracts.
func AnalyzeOptionsContracts(c *fiber.Ctx) error {
	var contracts []models.OptionsContract

	// Parse the request body into the contracts slice.
	if err := c.BodyParser(&contracts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validate each contract's type and position.
	for _, contract := range contracts {
		if contract.Type != "Call" && contract.Type != "Put" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid option type"})
		}
		if contract.LongShort != "long" && contract.LongShort != "short" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid position type"})
		}
	}

	// Perform calculations.
	xyValues := calculateXYValues(contracts)
	maxProfit := getMaxProfit(xyValues)
	maxLoss := getMaxLoss(xyValues)
	breakEvenPoints := findBreakEvenPoints(xyValues)

	// Prepare the response.
	response := AnalysisResponse{
		XYValues:        xyValues,
		MaxProfit:       maxProfit,
		MaxLoss:         maxLoss,
		BreakEvenPoints: breakEvenPoints,
	}

	// Return the response as JSON.
	return c.Status(fiber.StatusOK).JSON(response)
}

// calculateXYValues computes the profit/loss for a range of prices.
func calculateXYValues(contracts []models.OptionsContract) []XYValue {
	var xyValues []XYValue

	// Iterate over a range of prices from 50 to 150.
	for price := 50.0; price <= 150.0; price += 1.0 {
		totalProfitLoss := 0.0

		// Calculate the profit/loss for each contract at the current price.
		for _, contract := range contracts {
			profitLoss := 0.0

			switch contract.Type {
			case "Call":
				profitLoss = calculateCallProfitLoss(contract, price)
			case "Put":
				profitLoss = calculatePutProfitLoss(contract, price)
			}

			totalProfitLoss += profitLoss
		}

		xyValues = append(xyValues, XYValue{X: price, Y: totalProfitLoss})
	}

	return xyValues
}

// calculateCallProfitLoss computes the profit/loss for a call option at a given price.
func calculateCallProfitLoss(contract models.OptionsContract, price float64) float64 {
	if contract.LongShort == "long" {
		return math.Max(0, price-contract.StrikePrice) - contract.Ask
	}
	return contract.Bid - math.Max(0, price-contract.StrikePrice)
}

// calculatePutProfitLoss computes the profit/loss for a put option at a given price.
func calculatePutProfitLoss(contract models.OptionsContract, price float64) float64 {
	if contract.LongShort == "long" {
		return math.Max(0, contract.StrikePrice-price) - contract.Ask
	}
	return contract.Bid - math.Max(0, contract.StrikePrice-price)
}

// getMaxProfit finds the maximum profit from the list of XY values.
func getMaxProfit(xyValues []XYValue) float64 {
	maxProfit := math.Inf(-1)
	for _, xy := range xyValues {
		if xy.Y > maxProfit {
			maxProfit = xy.Y
		}
	}
	return maxProfit
}

// getMaxLoss finds the maximum loss from the list of XY values.
func getMaxLoss(xyValues []XYValue) float64 {
	maxLoss := math.Inf(1)
	for _, xy := range xyValues {
		if xy.Y < maxLoss {
			maxLoss = xy.Y
		}
	}
	return maxLoss
}

// findBreakEvenPoints determines the prices where profit/loss is zero.
func findBreakEvenPoints(xyValues []XYValue) []float64 {
	var breakEvenPoints []float64

	for i := 1; i < len(xyValues); i++ {
		// Check for sign change between consecutive points.
		if (xyValues[i-1].Y <= 0 && xyValues[i].Y >= 0) || (xyValues[i-1].Y >= 0 && xyValues[i].Y <= 0) {
			breakEvenPoints = append(breakEvenPoints, xyValues[i].X)
		}
	}

	return breakEvenPoints
}