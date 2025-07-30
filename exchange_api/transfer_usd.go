package exchange_api

import (
	"fmt"
	sdk "hyperliquid-go-sdk"
	"math/big"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type TransferUSDRequest struct {
	Nonce            uint64 `json:"time"`
	Amount           string `json:"amount"`
	Destination      string `json:"destination"`
	HyperliquidChain string `json:"hyperliquidChain"` // e.g., "Mainnet" or "Testnet"
	SignatureChainId string `json:"signatureChainId"` // e.g., "0xa4b1" for Arbitrum Mainnet
}

type TransferUSDAction struct {
	TransferUSDRequest
	Type string `json:"type"`
}

func (action *TransferUSDAction) Tp() string {
	return action.Type
}

func FromBuilderTransferUSDReq(req *TransferUSDRequest) *TransferUSDAction {
	return &TransferUSDAction{
		Type: "usdSend",
		TransferUSDRequest: TransferUSDRequest{
			Amount:           req.Amount,
			Destination:      req.Destination,
			HyperliquidChain: req.HyperliquidChain,
			SignatureChainId: req.SignatureChainId,
			Nonce:            req.Nonce,
		},
	}
}

func TansferUSD(e *sdk.Exchange, req TransferUSDRequest) (any, error) {
	nonce := e.NextNonce()
	req.Nonce = nonce
	action := FromBuilderTransferUSDReq(&req)
	actionT := map[string]interface{}{
		"destination":      req.Destination,
		"amount":           req.Amount,
		"time":             new(big.Int).SetUint64(nonce),
		"hyperliquidChain": req.HyperliquidChain,
		"signatureChainId": req.SignatureChainId,
	}

	// Note: The signature is generated based on the *action* content and the *request* nonce.
	sig, err := signTransferUSD(e, actionT)
	if err != nil {
		return nil, fmt.Errorf("failed to sign Transfer USD action: %w", err)
	}

	// The response structure for transfer USD might be simpler (e.g., just status confirmation)
	// Adjust parsing if needed based on actual API response.
	respType, statuses, err := e.PostActionAndParseResponse(action, sig, nonce)
	if err != nil {
		return nil, fmt.Errorf("transferUSDC request failed: %w", err)
	}

	// Assuming the response format is similar to others, return the first status element if available.
	// This might need adjustment based on the actual API response for approveBuilderFee.
	if len(statuses) > 0 {
		return statuses[0], nil
	}

	// If no specific status data is returned, return the response type (e.g., "ok")
	return respType, nil
}

var transferUSDPrimaryType = "HyperliquidTransaction:UsdSend"

func signTransferUSDPayload() []apitypes.Type {
	return []apitypes.Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "destination", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "time", Type: "uint64"},
	}
}

func signTransferUSD(e *sdk.Exchange, action apitypes.TypedDataMessage) (*sdk.Signature, error) {
	payload := signTransferUSDPayload()
	return sdk.SignUserSignedAction(e.Signer(), action, payload, transferUSDPrimaryType)
}
