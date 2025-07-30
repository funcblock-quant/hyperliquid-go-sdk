package sdk

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func RoundToDecimal(value float64, decimals int) float64 {
	factor := math.Pow10(decimals)
	return math.Round(value*factor) / factor
}

// RoundToSignificantAndDecimal rounds a value to a specified number of significant figures
// and then ensures it has at most the specified number of decimal places.
func RoundToSignificantAndDecimal(value float64, sigFigs, maxDecimals int) float64 {
	// Format the number with specified significant figures
	formattedValue := fmt.Sprintf("%.*g", sigFigs, value)
	// Convert back to float64
	rounded, _ := strconv.ParseFloat(formattedValue, 64)
	return RoundToDecimal(rounded, maxDecimals)
}

// FloatToString converts a float64 to a string format suitable for hashing
func FloatToString(x float64) string {
	// Format with fixed number of decimal places
	formatted := fmt.Sprintf("%.8f", x)
	// Remove trailing zeros
	for strings.HasSuffix(formatted, "0") {
		formatted = formatted[:len(formatted)-1]
	}
	// Remove decimal point if it's the last character
	formatted = strings.TrimSuffix(formatted, ".")
	// Convert -0 to 0
	if formatted == "-0" {
		return "0"
	}
	return formatted
}

// FloatToInt converts a float64 to an int, considering the specified power of 10
func FloatToInt(x float64, power int) int {
	// Multiply the float by 10^power to shift decimal places
	withDecimals := x * math.Pow10(power)
	// Round to nearest integer and return the result
	return int(math.Round(withDecimals))
}

func FloatToUsdInt(x float64) int {
	return FloatToInt(x, 6)
}
