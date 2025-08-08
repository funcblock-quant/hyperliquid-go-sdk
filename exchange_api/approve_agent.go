package exchange_api

import (
	"fmt"
	sdk "github.com/funcblock-quant/hyperliquid-go-sdk"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// ApproveAgentRequest defines the parameters for the approveAgent action.
type ApproveAgentRequest struct {
	Nonce            uint64         `json:"nonce"`
	AgentAddress     common.Address `json:"agentAddress"`
	AgentName        string         `json:"agentName"`
	HyperliquidChain string         `json:"hyperliquidChain"` // e.g., "Mainnet" or "Testnet"
	SignatureChainId string         `json:"signatureChainId"` // e.g., "0xa4b1" for Arbitrum Mainnet
}

// ApproveAgentAction represents the structure for the approveAgent action sent to the API.
type ApproveAgentAction struct {
	ApproveAgentRequest
	Type string `json:"type"` // Should be "approveAgent"
}

func (action *ApproveAgentAction) Tp() string {
	return action.Type
}

func FromApproveAgentReq(req *ApproveAgentRequest) *ApproveAgentAction {
	return &ApproveAgentAction{
		Type: "approveAgent",
		ApproveAgentRequest: ApproveAgentRequest{
			AgentAddress:     req.AgentAddress,
			AgentName:        req.AgentName,
			HyperliquidChain: req.HyperliquidChain,
			SignatureChainId: req.SignatureChainId,
			Nonce:            req.Nonce,
		},
	}
}

// ApproveAgent sends a request to approve an agent wallet for the main account.
// All parameters are required as per the API specification.
func ApproveAgent(e *sdk.Exchange, req ApproveAgentRequest) (any, error) {
	nonce := e.NextNonce()
	req.Nonce = nonce
	action := FromApproveAgentReq(&req)
	actionT := map[string]interface{}{
		"nonce":            new(big.Int).SetUint64(nonce),
		"agentAddress":     req.AgentAddress.Hex(),
		"agentName":        req.AgentName,
		"hyperliquidChain": req.HyperliquidChain,
		"signatureChainId": req.SignatureChainId,
	}

	// Note: The signature is generated based on the *action* content and the *request* nonce.
	sig, err := signAgent(e, actionT)
	if err != nil {
		return nil, fmt.Errorf("failed to sign approveAgent action: %w", err)
	}

	// The response structure for approveAgent might be simpler (e.g., just status confirmation)
	// Adjust parsing if needed based on actual API response.
	respType, statuses, err := e.PostActionAndParseResponse(action, sig, nonce)
	if err != nil {
		return nil, fmt.Errorf("approveAgent request failed: %w", err)
	}

	// Assuming the response format is similar to others, return the first status element if available.
	// This might need adjustment based on the actual API response for approveAgent.
	if len(statuses) > 0 {
		return statuses[0], nil
	}

	// If no specific status data is returned, return the response type (e.g., "ok")
	return respType, nil
}

var signAgentPrimaryType = "HyperliquidTransaction:ApproveAgent"

func signAgentPayloadTypes() []apitypes.Type {
	return []apitypes.Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "agentAddress", Type: "address"},
		{Name: "agentName", Type: "string"},
		{Name: "nonce", Type: "uint64"},
	}
}

// Sign for Approve Agent
func signAgent(e *sdk.Exchange, action apitypes.TypedDataMessage) (*sdk.Signature, error) {
	payload := signAgentPayloadTypes()
	primaryType := signAgentPrimaryType
	return sdk.SignUserSignedAction(e.Signer(), action, payload, primaryType)
}

func ApproveAgentTypedData(req ApproveAgentRequest) (*apitypes.TypedData, error) {
	action := map[string]interface{}{
		"nonce":            new(big.Int).SetUint64(req.Nonce),
		"agentAddress":     req.AgentAddress.Hex(),
		"agentName":        req.AgentName,
		"hyperliquidChain": req.HyperliquidChain,
		"signatureChainId": req.SignatureChainId,
	}
	payloadTypes := signAgentPayloadTypes()
	primaryType := signAgentPrimaryType

	signData, err := sdk.UserSignedPayload(primaryType, payloadTypes, action)
	return signData, err
}
