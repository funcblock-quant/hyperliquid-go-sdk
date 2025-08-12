package examples

import (
	"testing"

	sdk "github.com/funcblock-quant/hyperliquid-go-sdk"
)

const walletAddress = "youWalletAddress"

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
	orders, err := info.OpenOrders(walletAddress)
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
	orders, err := info.FrontendOpenOrders(walletAddress)
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
	txs, err := info.UserDepositWithdrawTxs(walletAddress, nil, nil)
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
	txs, err := info.UserFills(walletAddress)
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
	txs, err := info.UserFundingHistory(walletAddress, 0, nil)
	if err != nil {
		t.Fatalf("Failed to fetch user funding history: %v", err)
	}
	t.Logf("User funding history: %v", txs)
}

func TestUserPortfolio(t *testing.T) {
	info, err := sdk.NewInfo(sdk.TestnetAPIURL)
	if err != nil {
		t.Fatalf("Failed to create sdk.Info for all_tokens_test: %v", err)
	}
	portfolio, err := info.UserPortfolio(walletAddress)
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
	fees, err := info.UserFees(walletAddress)
	if err != nil {
		t.Fatalf("Failed to fetch user fees: %v", err)
	}
	t.Logf("User fees: %v", fees)
}

func TestUserStatus(t *testing.T) {
	info, err := sdk.NewInfo(sdk.TestnetAPIURL)
	if err != nil {
		t.Fatalf("Failed to create sdk.Info for all_tokens_test: %v", err)
	}
	userState, err := info.UserState(walletAddress)
	if err != nil {
		t.Fatalf("Failed to fetch user state: %v", err)
	}
	t.Logf("User state: %v", userState)
}
