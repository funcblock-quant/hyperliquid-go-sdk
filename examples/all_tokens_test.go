package examples

import (
	sdk "hyperliquid-go-sdk"
	"testing"
)

func TestGetAllTokens(t *testing.T) {
	info, err := sdk.NewInfo(sdk.MainnetAPIURL)
	if err != nil {
		t.Fatalf("Failed to create sdk.Info for all_tokens_test: %v", err)
	}
	perpCoins := info.PerpCoins()
	spotCoins := info.SpotCoins()

	t.Logf("Perp coins (len: %d): %v", len(perpCoins), perpCoins)
	t.Logf("Spot coins (len: %d): %v", len(spotCoins), spotCoins)
}
