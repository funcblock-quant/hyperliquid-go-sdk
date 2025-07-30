package exchange_api

import (
	"fmt"
	sdk "hyperliquid-go-sdk"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type ApproveBuilderFeeRequest struct {
	Nonce            uint64         `json:"nonce"`
	Builder          common.Address `json:"builder"`
	MaxFeeRate       string         `json:"maxFeeRate"`
	HyperliquidChain string         `json:"hyperliquidChain"` // e.g., "Mainnet" or "Testnet"
	SignatureChainId string         `json:"signatureChainId"` // e.g., "0xa4b1" for Arbitrum Mainnet
}

type ApproveBuilderFeeAction struct {
	ApproveBuilderFeeRequest
	Type string `json:"type"`
}

func (action *ApproveBuilderFeeAction) Tp() string {
	return action.Type
}

func FromBuilderFeeReq(req *ApproveBuilderFeeRequest) *ApproveBuilderFeeAction {
	return &ApproveBuilderFeeAction{
		Type: "approveBuilderFee",
		ApproveBuilderFeeRequest: ApproveBuilderFeeRequest{
			Builder:          req.Builder,
			MaxFeeRate:       req.MaxFeeRate,
			HyperliquidChain: req.HyperliquidChain,
			SignatureChainId: req.SignatureChainId,
			Nonce:            req.Nonce,
		},
	}
}

func ApproveBuilderFee(e *sdk.Exchange, req ApproveBuilderFeeRequest) (any, error) {
	nonce := e.NextNonce()
	req.Nonce = nonce
	action := FromBuilderFeeReq(&req)
	actionT := map[string]interface{}{
		"nonce":            new(big.Int).SetUint64(nonce),
		"builder":          req.Builder.Hex(),
		"maxFeeRate":       req.MaxFeeRate,
		"hyperliquidChain": req.HyperliquidChain,
		"signatureChainId": req.SignatureChainId,
	}

	// Note: The signature is generated based on the *action* content and the *request* nonce.
	sig, err := signApproveBuilderFee(e, actionT)
	if err != nil {
		return nil, fmt.Errorf("failed to sign approveBuilderFee action: %w", err)
	}

	// The response structure for approveBuilderFee might be simpler (e.g., just status confirmation)
	// Adjust parsing if needed based on actual API response.
	respType, statuses, err := e.PostActionAndParseResponse(action, sig, nonce)
	if err != nil {
		return nil, fmt.Errorf("approveBuilderFee request failed: %w", err)
	}

	// Assuming the response format is similar to others, return the first status element if available.
	// This might need adjustment based on the actual API response for approveBuilderFee.
	if len(statuses) > 0 {
		return statuses[0], nil
	}

	// If no specific status data is returned, return the response type (e.g., "ok")
	return respType, nil
}

var approveBuilderFeePrimaryType = "HyperliquidTransaction:ApproveBuilderFee"

func signApproveBuilderFeePayloadTypes() []apitypes.Type {
	return []apitypes.Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "maxFeeRate", Type: "string"},
		{Name: "builder", Type: "address"},
		{Name: "nonce", Type: "uint64"},
	}
}

func signApproveBuilderFee(e *sdk.Exchange, action apitypes.TypedDataMessage) (*sdk.Signature, error) {
	payload := signApproveBuilderFeePayloadTypes()
	return sdk.SignUserSignedAction(e.Signer(), action, payload, approveBuilderFeePrimaryType)
}

func ApproveBuilderFeeTypedData(req ApproveBuilderFeeRequest) (*apitypes.TypedData, error) {
	action := map[string]interface{}{
		"nonce":            new(big.Int).SetUint64(req.Nonce),
		"builder":          req.Builder.Hex(),
		"maxFeeRate":       req.MaxFeeRate,
		"hyperliquidChain": req.HyperliquidChain,
		"signatureChainId": req.SignatureChainId,
	}
	payloadTypes := signApproveBuilderFeePayloadTypes()
	primaryType := approveBuilderFeePrimaryType

	signData, err := sdk.UserSignedPayload(primaryType, payloadTypes, action)
	return signData, err
}
