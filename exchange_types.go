package sdk

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

//go:generate easyjson -all types.go

const (
	TifAlo string = "Alo" // ALO (add liquidity only, i.e. "post only") will be canceled instead of immediately matching.
	TifIoc string = "Ioc" // IOC (immediate or cancel) will have the unfilled part canceled instead of resting.
	TifGtc string = "Gtc" // GTC (good til canceled) orders have no special behavior.

	TakeProfit string = "tp"
	StopLose   string = "sl"

	GroupingNa   string = "na"
	GroupingTpsl string = "tpsl"
)

type AssetInfo struct {
	Name       string `json:"name"`
	SzDecimals int    `json:"szDecimals"`
}

type Meta struct {
	Universe []AssetInfo `json:"universe"`
}

type SpotAssetInfo struct {
	Name        string `json:"name"`
	Tokens      []int  `json:"tokens"`
	Index       int    `json:"index"`
	IsCanonical bool   `json:"isCanonical"`
}

type SpotTokenInfo struct {
	Name        string          `json:"name"`
	SzDecimals  int             `json:"szDecimals"`
	WeiDecimals int             `json:"weiDecimals"`
	Index       int             `json:"index"`
	TokenID     string          `json:"tokenId"`
	IsCanonical bool            `json:"isCanonical"`
	EvmContract json.RawMessage `json:"evmContract"`
	FullName    *string         `json:"fullName"`
}

type SpotMeta struct {
	Universe []SpotAssetInfo `json:"universe"`
	Tokens   []SpotTokenInfo `json:"tokens"`
}

type SpotAssetCtx struct {
	DayNtlVlm         string  `json:"dayNtlVlm"`
	MarkPx            string  `json:"markPx"`
	MidPx             *string `json:"midPx"`
	PrevDayPx         string  `json:"prevDayPx"`
	CirculatingSupply string  `json:"circulatingSupply"`
	Coin              string  `json:"coin"`
}

// WebSocket message types

type WsMsg struct {
	Channel string         `json:"channel"`
	Data    map[string]any `json:"data"`
}

// Exchange Request types

type OrderRequest struct {
	Coin       string
	IsBuy      bool
	Size       float64
	LimitPx    float64
	OrderType  OrderType
	ReduceOnly bool
	Cloid      *string
}

type ModifyRequest struct {
	Oid          any
	OrderRequest OrderRequest
}

type CancelRequest struct {
	Coin string
	Oid  uint64
}

type MarketRequest struct {
	Coin        string
	IsBuy       bool
	ReduceOnly  bool
	Size        float64
	MarketPrice float64
	Slippage    float64
	Cloid       *string
}

type CancelByCloidRequest struct {
	Coin  string
	Cloid string
}

type OrderType struct {
	Limit   *LimitOrderType   `json:"limit,omitempty" msgpack:"limit,omitempty"`
	Trigger *TriggerOrderType `json:"trigger,omitempty" msgpack:"trigger,omitempty"`
}

type LimitOrderType struct {
	Tif string `json:"tif" msgpack:"tif"` // "Alo", "Ioc", "Gtc"
}

type TriggerOrderType struct {
	TriggerPx string `json:"triggerPx" msgpack:"triggerPx"`
	IsMarket  bool   `json:"isMarket" msgpack:"isMarket"`
	Tpsl      string `json:"tpsl" msgpack:"tpsl"` // "tp" or "sl"
}

type BuilderInfo struct {
	Builder string `json:"b" msgpack:"b"`
	Fee     int    `json:"f" msgpack:"f"`
}

type OrderWire struct {
	Asset      int           `json:"a" msgpack:"a"`
	IsBuy      bool          `json:"b" msgpack:"b"`
	LimitPx    string        `json:"p" msgpack:"p"`
	Size       string        `json:"s" msgpack:"s"`
	ReduceOnly bool          `json:"r" msgpack:"r"`
	OrderType  OrderTypeWire `json:"t" msgpack:"t"`
	Cloid      string        `json:"c,omitempty" msgpack:"c,omitempty"`
}

type ModifyWire struct {
	Oid   any       `json:"oid" msgpack:"oid"`
	Order OrderWire `json:"order" msgpack:"order"`
}

type CancelWire struct {
	Asset int    `json:"a" msgpack:"a"`
	Oid   uint64 `json:"o" msgpack:"o"`
}

type CancelByCloidWire struct {
	Asset int    `json:"asset" msgpack:"asset"`
	Cloid string `json:"cloid" msgpack:"cloid"`
}

type TriggerOrderTypeWire struct {
	TriggerPx string `json:"triggerPx" msgpack:"triggerPx"`
	IsMarket  bool   `json:"isMarket" msgpack:"isMarket"`
	Tpsl      string `json:"tpsl" msgpack:"tpsl"` // "tp" or "sl"
}

type OrderTypeWire struct {
	Limit   *LimitOrderType       `json:"limit,omitempty" msgpack:"limit,omitempty"`
	Trigger *TriggerOrderTypeWire `json:"trigger,omitempty" msgpack:"trigger,omitempty"`
}

type ExchangeRequest struct {
	Action       Action          `json:"action"`
	Nonce        uint64          `json:"nonce"`
	Signature    *Signature      `json:"signature"`
	VaultAddress *common.Address `json:"vaultAddress,omitempty"`
	ExpiresAfter *uint64         `json:"expiresAfter,omitempty"`
}

func (req *OrderRequest) ToWire(asset, assetDec int) OrderWire {
	wire := OrderWire{
		Asset:      asset,
		IsBuy:      req.IsBuy,
		LimitPx:    FloatToString(adjustPrice(req.LimitPx, asset, assetDec)),
		ReduceOnly: req.ReduceOnly,
		OrderType:  req.OrderType.ToWire(),
		Size:       FloatToString(RoundToDecimal(req.Size, assetDec)),
	}
	if req.Cloid != nil {
		wire.Cloid = *req.Cloid
	}

	return wire
}

func (req *ModifyRequest) ToWire(asset, assetDec int) ModifyWire {
	return ModifyWire{
		Oid:   req.Oid,
		Order: req.OrderRequest.ToWire(asset, assetDec),
	}
}

func (req *CancelRequest) ToWire(asset int) CancelWire {
	return CancelWire{
		Asset: asset,
		Oid:   req.Oid,
	}
}

func (req *CancelByCloidRequest) ToWire(asset int) CancelByCloidWire {
	return CancelByCloidWire{
		Asset: asset,
		Cloid: req.Cloid,
	}
}

func (tp *OrderType) ToWire() OrderTypeWire {
	wire := OrderTypeWire{}

	if tp.Limit != nil {
		wire.Limit = &LimitOrderType{
			Tif: tp.Limit.Tif,
		}
	} else if tp.Trigger != nil {
		wire.Trigger = &TriggerOrderTypeWire{
			TriggerPx: tp.Trigger.TriggerPx,
			IsMarket:  tp.Trigger.IsMarket,
			Tpsl:      tp.Trigger.Tpsl,
		}
	}

	return wire
}

func (tp *TriggerOrderType) ToWire() TriggerOrderTypeWire {
	return TriggerOrderTypeWire{
		TriggerPx: tp.TriggerPx,
		IsMarket:  tp.IsMarket,
		Tpsl:      tp.Tpsl,
	}
}

// Exchange action types

type Action interface {
	Tp() string
}

type OrderAction struct {
	Type     string       `json:"type" msgpack:"type"`
	Orders   []OrderWire  `json:"orders" msgpack:"orders"`
	Grouping string       `json:"grouping" msgpack:"grouping"`
	Builder  *BuilderInfo `json:"builder,omitempty" msgpack:"builder,omitempty"`
}

func (action *OrderAction) Tp() string {
	return action.Type
}

type ModifyAction struct {
	Type     string       `json:"type" msgpack:"type"`
	Modifies []ModifyWire `json:"modifies" msgpack:"modifies"`
}

func (action *ModifyAction) Tp() string {
	return action.Type
}

type CancelAction struct {
	Type    string       `json:"type" msgpack:"type"`
	Cancels []CancelWire `json:"cancels" msgpack:"cancels"`
}

func (action *CancelAction) Tp() string {
	return action.Type
}

type CancelByCloidAction struct {
	Type    string              `json:"type" msgpack:"type"`
	Cancels []CancelByCloidWire `json:"cancels" msgpack:"cancels"`
}

func (action *CancelByCloidAction) Tp() string {
	return action.Type
}

type UpdateLeverageAction struct {
	Type     string `json:"type" msgpack:"type"`
	Asset    int    `json:"asset" msgpack:"asset"`
	IsCross  bool   `json:"isCross" msgpack:"isCross"`
	Leverage int    `json:"leverage" msgpack:"leverage"`
}

func (action *UpdateLeverageAction) Tp() string {
	return action.Type
}

type UpdateIsolatedMarginAction struct {
	Type   string `json:"type" msgpack:"type"`
	Asset  int    `json:"asset" msgpack:"asset"`
	IsBuy  bool   `json:"isBuy" msgpack:"isBuy"`
	Amount int    `json:"ntli" msgpack:"ntli"`
}

func (action *UpdateIsolatedMarginAction) Tp() string {
	return action.Type
}

type VaultUsdTransferAction struct {
	Type         string `json:"type" msgpack:"type"`
	VaultAddress string `json:"vaultAddress" msgpack:"vaultAddress"`
	IsDeposit    bool   `json:"isDeposit" msgpack:"isDeposit"`
	Usd          int    `json:"usd" msgpack:"usd"` // Amount in USD cents
}

func (action *VaultUsdTransferAction) Tp() string {
	return action.Type
}

// Response related

type ExchangeRestingOrder struct {
	Oid uint64 `json:"oid"`
}

type ExchangeFilledOrder struct {
	Oid       uint64 `json:"oid"`
	TotalSize string `json:"totalSz"`
	AveragePx string `json:"avgPx"`
}

type ExchangeDataStatusObject struct {
	Error   *string               `json:"error"`
	Resting *ExchangeRestingOrder `json:"resting"`
	Filled  *ExchangeFilledOrder  `json:"filled"`
}

type ExchangeDataStatus struct {
	String *string
	Object *ExchangeDataStatusObject
}

func (s *ExchangeDataStatus) UnmarshalJSON(data []byte) error {
	str := new(string)
	if err := json.Unmarshal(data, str); err == nil {
		s.String = str
		return nil
	}

	obj := new(ExchangeDataStatusObject)
	if err := json.Unmarshal(data, &obj); err == nil {
		s.Object = obj
		return nil
	} else {
		return err
	}
}

// Parse returns the parsed value of the status.
// possible result types:
// 1. string
// 2. *ExchangeRestingOrder
// 3. *ExchangeFilledOrder
// 4. error
func (s *ExchangeDataStatus) Parse() any {
	if s.String != nil {
		return *s.String
	}
	if s.Object != nil {
		if s.Object.Error != nil {
			return errors.New(*s.Object.Error)
		}
		if s.Object.Resting != nil {
			return s.Object.Resting
		}
		if s.Object.Filled != nil {
			return s.Object.Filled
		}
	}
	return fmt.Errorf("invalid exchange status: %+v", s)
}

type ExchangeDataStatuses struct {
	Statuses []ExchangeDataStatus `json:"statuses"`
}

type ExchangeSuccessResponse struct {
	Type string                `json:"type"`
	Data *ExchangeDataStatuses `json:"data"`
}

type ExchangeResponse struct {
	Error   *string
	Success *ExchangeSuccessResponse
}

func (e *ExchangeResponse) UnmarshalJSON(data []byte) error {
	errMsg := new(string)
	if err := json.Unmarshal(data, errMsg); err == nil {
		e.Error = errMsg
		return nil
	}

	success := new(ExchangeSuccessResponse)
	if err := json.Unmarshal(data, success); err == nil {
		e.Success = success
		return nil
	} else {
		return err
	}
}

type ExchangeResponsesStatus struct {
	Status   string           `json:"status"`
	Response ExchangeResponse `json:"response"`
}

func (e *ExchangeResponsesStatus) Parse() (*ExchangeSuccessResponse, error) {
	if e.Status == "ok" {
		return e.Response.Success, nil
	} else {
		return nil, errors.New(*e.Response.Error)
	}
}

func adjustPrice(price float64, asset, assetDec int) float64 {
	// 6 decimals for perps, 8 decimals for spot
	decimals := 6
	// spot assets start at 10000
	if asset >= 10_000 {
		decimals = 8
	}
	// Adjust by asset-specific decimal offset
	decimals -= assetDec
	// Format to 5 significant figures
	return RoundToSignificantAndDecimal(price, 5, decimals)
}
