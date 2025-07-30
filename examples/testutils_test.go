package examples

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	sdk "hyperliquid-go-sdk"
)

var (
	testPrivateKey *ecdsa.PrivateKey //lint:ignore U1000 unused
	testExchange   *sdk.Exchange
)

func TestMain(m *testing.M) {
	// Setup test environment
	var err error
	privKeyHex := os.Getenv("HL_PRIVATE_KEY")
	if privKeyHex == "" {
		// Skip all tests if no private key is provided
		fmt.Println("Skipping all tests: HL_PRIVATE_KEY environment variable not set")
		os.Exit(0)
	}

	testSigner, err := sdk.NewLocalSignerFromHex(privKeyHex)
	if err != nil {
		panic("failed to create local signer: " + err.Error())
	}
	fmt.Printf("test signer address: %s\n", testSigner.Address())

	vault := os.Getenv("HL_VAULT_ADDRESS")
	var vaultAddress *common.Address
	if vault != "" {
		addr := common.HexToAddress(vault)
		vaultAddress = &addr
	}

	info, err := sdk.NewInfo(sdk.MainnetAPIURL)
	if err != nil {
		panic("failed to create sdk.Info for testutils: " + err.Error())
	}
	meta, err := info.Meta()
	if err != nil {
		panic("failed to fetch meta data: " + err.Error())
	}
	// Initialize test exchange
	testExchange = sdk.NewExchange(sdk.MainnetAPIURL, vaultAddress, meta, testSigner)

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func getTestExchange(t *testing.T) *sdk.Exchange {
	t.Helper()
	if testExchange == nil {
		t.Fatal("test exchange not initialized")
	}
	return testExchange
}

//lint:ignore U1000 unused
func skipIfNoPrivateKey(t *testing.T) {
	t.Helper()
	if os.Getenv("HL_PRIVATE_KEY") == "" {
		t.Skip("Skipping test: no private key provided")
	}
}
