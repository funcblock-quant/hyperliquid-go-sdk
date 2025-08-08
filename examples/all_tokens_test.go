package examples

import (
	"testing"
	
	sdk "github.com/funcblock-quant/hyperliquid-go-sdk"
)

func TestGetAllTokens(t *testing.T) {
	info, err := sdk.NewInfo(sdk.TestnetAPIURL)
	if err != nil {
		t.Fatalf("Failed to create sdk.Info for all_tokens_test: %v", err)
	}
	perpCoins := info.PerpCoins()
	spotCoins := info.SpotCoins()

	t.Logf("Perp coins (len: %d): %v", len(perpCoins), perpCoins)
	t.Logf("Spot coins (len: %d): %v", len(spotCoins), spotCoins)
}

// 查看用户订单
func TestOpenOrders(t *testing.T) {
	info, err := sdk.NewInfo(sdk.TestnetAPIURL)
	if err != nil {
		t.Fatalf("Failed to create sdk.Info for all_tokens_test: %v", err)
	}
	orders, err := info.OpenOrders("0xae61Feb5A2D59aDF2291f509FF87e449E7fa5790")
	if err != nil {
		t.Fatalf("Failed to fetch open orders: %v", err)
	}
	t.Logf("Open orders: %v", orders)
}

// 查看用户前端订单
func TestFrontendOpenOrders(t *testing.T) {
	info, err := sdk.NewInfo(sdk.TestnetAPIURL)
	if err != nil {
		t.Fatalf("Failed to create sdk.Info for all_tokens_test: %v", err)
	}
	orders, err := info.FrontendOpenOrders("0xae61Feb5A2D59aDF2291f509FF87e449E7fa5790")
	if err != nil {
		t.Fatalf("Failed to fetch frontend open orders: %v", err)
	}
	t.Logf("Frontend open orders: %v", orders)
}

func TestUserDepositWithdrawTxs(t *testing.T) {
	info, err := sdk.NewInfo(sdk.TestnetAPIURL)
	if err != nil {
		t.Fatalf("Failed to create sdk.Info for all_tokens_test: %v", err)
	}
	txs, err := info.UserDepositWithdrawTxs("0xae61Feb5A2D59aDF2291f509FF87e449E7fa5790", nil, nil)
	if err != nil {
		t.Fatalf("Failed to fetch user deposit & withdraw txs: %v", err)
	}
	t.Logf("User deposit & withdraw txs: %v", txs)
}

// 查看用户交易记录
func TestUserFills(t *testing.T) {
	info, err := sdk.NewInfo(sdk.TestnetAPIURL)
	if err != nil {
		t.Fatalf("Failed to create sdk.Info for all_tokens_test: %v", err)
	}
	txs, err := info.UserFills("0xae61Feb5A2D59aDF2291f509FF87e449E7fa5790")
	if err != nil {
		t.Fatalf("Failed to fetch user fills: %v", err)
	}
	t.Logf("User fills: %v", txs)
}

// 获取用户资金费用历史
func TestUserFundingHistory(t *testing.T) {
	info, err := sdk.NewInfo(sdk.TestnetAPIURL)
	if err != nil {
		t.Fatalf("Failed to create sdk.Info for all_tokens_test: %v", err)
	}
	txs, err := info.UserFundingHistory("0xae61Feb5A2D59aDF2291f509FF87e449E7fa5790", 1700000000000, nil)
	if err != nil {
		t.Fatalf("Failed to fetch user funding history: %v", err)
	}
	
	t.Logf("Total funding history records: %d", len(txs))
	for i, tx := range txs {
		if i >= 5 { // 只显示前5条记录
			break
		}
		t.Logf("Funding [%d]: Coin=%s, FundingRate=%s, Szi=%s, Usdc=%s, Hash=%s, Time=%d", 
			i, tx.Delta.Coin, tx.Delta.FundingRate, tx.Delta.Szi, tx.Delta.Usdc, tx.Hash, tx.Time)
	}
	
	// 如果有数据，显示总的资金费用
	if len(txs) > 0 {
		totalFunding := 0.0
		for _, tx := range txs {
			if fundingValue, err := strconv.ParseFloat(tx.Delta.Usdc, 64); err == nil {
				totalFunding += fundingValue
			}
		}
		t.Logf("Total funding amount (USDC): %.8f", totalFunding)
	}
}

func TestUserPortfolio(t *testing.T) {
	info, err := sdk.NewInfo(sdk.TestnetAPIURL)
	if err != nil {
		t.Fatalf("Failed to create sdk.Info for all_tokens_test: %v", err)
	}
	portfolio, err := info.UserPortfolio("0xae61Feb5A2D59aDF2291f509FF87e449E7fa5790")
	if err != nil {
		t.Fatalf("Failed to fetch user portfolio: %v", err)
	}
	t.Logf("User portfolio: %v", portfolio)
}

// 获取用户手续费结构（相对复杂，暂时用不到）
func TestUserFees(t *testing.T) {
	info, err := sdk.NewInfo(sdk.TestnetAPIURL)
	if err != nil {
		t.Fatalf("Failed to create sdk.Info for all_tokens_test: %v", err)
	}
	fees, err := info.UserFees("0xae61Feb5A2D59aDF2291f509FF87e449E7fa5790")
	if err != nil {
		t.Fatalf("Failed to fetch user fees: %v", err)
	}
	t.Logf("User fees: %v", fees)
}