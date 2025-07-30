package examples

import (
	ex "hyperliquid-go-sdk/exchange_api"
	"testing"
)

func TestTransferUSD(t *testing.T) {
	exchange := getTestExchange(t)

	req := ex.TransferUSDRequest{
		Destination:      "0x1d4c6be5659e3d801c68ec6fc0fdf88b588b70a2",
		Amount:           "2",
		HyperliquidChain: "Mainnet",
		SignatureChainId: "0x66eee",
	}

	res, err := ex.TansferUSD(exchange, req)
	if err != nil {
		t.Fatalf("Transfer USDC failed: %v", err)
	}

	t.Logf("Transfer USDC completed successfully: %+v", res)
}
