package sdk

import (
	"encoding/json"
	"fmt"
)

//go:generate easyjson -all models.go

type L2Book struct {
	Coin   string    `json:"coin"`
	Levels [][]Level `json:"levels"`
	Time   int64     `json:"time"`
}

type Level struct {
	N  int     `json:"n"`
	Px float64 `json:"px,string"`
	Sz float64 `json:"sz,string"`
}

type AssetPosition struct {
	Position Position `json:"position"`
	Type     string   `json:"type"`
}

type Position struct {
	Coin           string     `json:"coin"`
	CumFunding     CumFunding `json:"cumFunding"`
	EntryPx        *string    `json:"entryPx"`
	Leverage       Leverage   `json:"leverage"`
	LiquidationPx  *string    `json:"liquidationPx"`
	MarginUsed     string     `json:"marginUsed"`
	PositionValue  string     `json:"positionValue"`
	ReturnOnEquity string     `json:"returnOnEquity"`
	Szi            string     `json:"szi"`
	UnrealizedPnl  string     `json:"unrealizedPnl"`
	MaxLeverage    int        `json:"maxLeverage"`
}

type CumFunding struct {
	AllTime     string `json:"allTime"`
	SinceChange string `json:"sinceChange"`
	SinceOpen   string `json:"sinceOpen"`
}

type Leverage struct {
	Type   string  `json:"type"`
	Value  int     `json:"value"`
	RawUsd *string `json:"rawUsd,omitempty"`
}

type UserState struct {
	AssetPositions             []AssetPosition `json:"assetPositions"`
	CrossMarginSummary         MarginSummary   `json:"crossMarginSummary"`
	MarginSummary              MarginSummary   `json:"marginSummary"`
	CrossMaintenanceMarginUsed string          `json:"crossMaintenanceMarginUsed"`
	Withdrawable               string          `json:"withdrawable"`
	Time                       int64           `json:"time"` // in ms
}

type MarginSummary struct {
	AccountValue    string `json:"accountValue"`
	TotalMarginUsed string `json:"totalMarginUsed"`
	TotalNtlPos     string `json:"totalNtlPos"`
	TotalRawUsd     string `json:"totalRawUsd"`
}

type OpenOrder struct {
	Coin      string  `json:"coin"`
	LimitPx   float64 `json:"limitPx,string"`
	Oid       int64   `json:"oid"`
	Side      string  `json:"side"`
	Size      float64 `json:"sz,string"`
	Timestamp int64   `json:"timestamp"`
}

// FrontendOpenOrder represents an open order with additional frontend-specific fields
type FrontendOpenOrder struct {
	Coin             string  `json:"coin"`
	IsPositionTpsl   bool    `json:"isPositionTpsl"`
	IsTrigger        bool    `json:"isTrigger"`
	LimitPx          float64 `json:"limitPx,string"`
	Oid              int64   `json:"oid"`
	OrderType        string  `json:"orderType"`
	OrigSz           float64 `json:"origSz,string"`
	ReduceOnly       bool    `json:"reduceOnly"`
	Side             string  `json:"side"`
	Size             float64 `json:"sz,string"`
	Timestamp        int64   `json:"timestamp"`
	TriggerCondition string  `json:"triggerCondition"`
	TriggerPx        float64 `json:"triggerPx,string"`
}

type Fill struct {
	ClosedPnl     string `json:"closedPnl"`
	Coin          string `json:"coin"`
	Crossed       bool   `json:"crossed"`
	Dir           string `json:"dir"`
	Hash          string `json:"hash"`
	Oid           int64  `json:"oid"`
	Price         string `json:"px"`
	Side          string `json:"side"`
	StartPosition string `json:"startPosition"`
	Size          string `json:"sz"`
	Time          int64  `json:"time"`
	Fee           string `json:"fee"`
	FeeToken      string `json:"feeToken"`
	BuilderFee    string `json:"builderFee"`
	Tid           int64  `json:"tid"`
}

type DepositWithdrawTx struct {
	Time   int64      `json:"time"`
	Hash   string     `json:"hash"`
	Action ActionType `json:"delta"`
}

func (t *DepositWithdrawTx) UnmarshalJSON(data []byte) error {
	type Alias DepositWithdrawTx
	aux := &struct {
		Action json.RawMessage `json:"delta"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var actionType struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(aux.Action, &actionType); err != nil {
		return err
	}

	switch actionType.Type {
	case "accountClassTransfer":
		var a USDClassTransferAction
		if err := json.Unmarshal(aux.Action, &a); err != nil {
			return err
		}
		t.Action = a
	case "withdraw":
		var a WithdrawAction
		if err := json.Unmarshal(aux.Action, &a); err != nil {
			return err
		}
		t.Action = a
	case "deposit":
		var a DepositAction
		if err := json.Unmarshal(aux.Action, &a); err != nil {
			return err
		}
		t.Action = a
	case "transfer", "internalTransfer":
		var a TransferAction
		if err := json.Unmarshal(aux.Action, &a); err != nil {
			return err
		}
		t.Action = a
	default:
		var a OtherAction
		a.Type = actionType.Type
		t.Action = a
	}

	return nil
}

type ActionType interface {
	actionType() string
}

type USDClassTransferAction struct {
	Type   string `json:"type"`
	Amount string `json:"usdc"`
	ToPerp bool   `json:"toPerp"`
}

func (a USDClassTransferAction) actionType() string {
	return a.Type
}

type WithdrawAction struct {
	Type   string `json:"type"`
	Amount string `json:"usdc"`
	Fee    string `json:"fee"`
}

func (a WithdrawAction) actionType() string {
	return a.Type
}

type DepositAction struct {
	Type   string `json:"type"`
	Amount string `json:"usdc"`
}

func (a DepositAction) actionType() string {
	return a.Type
}

type TransferAction struct {
	Type        string `json:"type"`
	Amount      string `json:"usdc"`
	Destination string `json:"destination"`
}

func (a TransferAction) actionType() string {
	return a.Type
}

type OtherAction struct {
	Type string `json:"type"`
}

func (a OtherAction) actionType() string {
	return a.Type
}

type HistoryEntry struct {
	Timestamp int64
	Value     string
}

type TimeRangeData struct {
	AccountValueHistory []HistoryEntry `json:"accountValueHistory"`
	PnlHistory          []HistoryEntry `json:"pnlHistory"`
	Volume              string         `json:"vlm"`
}

type PortFolioTimeRangeItem struct {
	RangeName string
	Data      TimeRangeData
}

func (tri *PortFolioTimeRangeItem) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) != 2 {
		return fmt.Errorf("invalid time range item format")
	}

	if err := json.Unmarshal(raw[0], &tri.RangeName); err != nil {
		return err
	}

	if err := json.Unmarshal(raw[1], &tri.Data); err != nil {
		return err
	}

	return nil
}

func (he *HistoryEntry) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) != 2 {
		return fmt.Errorf("invalid history entry format")
	}

	if err := json.Unmarshal(raw[0], &he.Timestamp); err != nil {
		return err
	}

	if err := json.Unmarshal(raw[1], &he.Value); err != nil {
		return err
	}

	return nil
}

type FundingHistory struct {
	Coin        string `json:"coin"`
	FundingRate string `json:"fundingRate"`
	Premium     string `json:"premium"`
	Time        int64  `json:"time"`
}

type UserFundingHistory struct {
	Delta FundingDelta `json:"delta"`
	Hash  string       `json:"hash"`
	Time  int64        `json:"time"`
}

type FundingDelta struct {
	Coin        string `json:"coin"`
	FundingRate string `json:"fundingRate"`
	Szi         string `json:"szi"`
	Type        string `json:"type"`
	Usdc        string `json:"usdc"`
}

type Candle struct {
	Timestamp int64  `json:"T"`
	Close     string `json:"c"`
	High      string `json:"h"`
	Interval  string `json:"i"`
	Low       string `json:"l"`
	Number    int    `json:"n"`
	Open      string `json:"o"`
	Symbol    string `json:"s"`
	Time      int64  `json:"t"`
	Volume    string `json:"v"`
}

type UserFees struct {
	ActiveReferralDiscount string       `json:"activeReferralDiscount"`
	DailyUserVolume        []UserVolume `json:"dailyUserVlm"`
	FeeSchedule            FeeSchedule  `json:"feeSchedule"`
	UserAddRate            string       `json:"userAddRate"`
	UserCrossRate          string       `json:"userCrossRate"`
}

type UserVolume struct {
	Date      string `json:"date"`
	Exchange  string `json:"exchange"`
	UserAdd   string `json:"userAdd"`
	UserCross string `json:"userCross"`
}

type FeeSchedule struct {
	Add              string `json:"add"`
	Cross            string `json:"cross"`
	ReferralDiscount string `json:"referralDiscount"`
	Tiers            Tiers  `json:"tiers"`
}

type Tiers struct {
	MM  []MMTier  `json:"mm"`
	VIP []VIPTier `json:"vip"`
}

type MMTier struct {
	Add                 string `json:"add"`
	MakerFractionCutoff string `json:"makerFractionCutoff"`
}

type VIPTier struct {
	Add       string `json:"add"`
	Cross     string `json:"cross"`
	NtlCutoff string `json:"ntlCutoff"`
}

type StakingSummary struct {
	Delegated              string `json:"delegated"`
	Undelegated            string `json:"undelegated"`
	TotalPendingWithdrawal string `json:"totalPendingWithdrawal"`
	NPendingWithdrawals    int    `json:"nPendingWithdrawals"`
}

type StakingDelegation struct {
	Validator            string `json:"validator"`
	Amount               string `json:"amount"`
	LockedUntilTimestamp int64  `json:"lockedUntilTimestamp"`
}

type StakingReward struct {
	Time        int64  `json:"time"`
	Source      string `json:"source"`
	TotalAmount string `json:"totalAmount"`
}

type ReferralState struct {
	ReferralCode string   `json:"referralCode"`
	Referrer     string   `json:"referrer"`
	Referred     []string `json:"referred"`
}

type SubAccount struct {
	Name        string   `json:"name"`
	User        string   `json:"user"`
	Permissions []string `json:"permissions"`
}

type MultiSigSigner struct {
	User      string `json:"user"`
	Threshold int    `json:"threshold"`
}

type Trade struct {
	Coin  string   `json:"coin"`
	Side  string   `json:"side"`
	Px    string   `json:"px"`
	Sz    string   `json:"sz"`
	Time  int64    `json:"time"`
	Hash  string   `json:"hash"`
	Tid   int64    `json:"tid"`
	Users []string `json:"users"`
}

// UniqueKey returns a unique key for the trade
func (t *Trade) UniqueKey() string {
	return fmt.Sprintf("%s-%d-%d", t.Coin, t.Time, t.Tid)
}

type ExtraAgent struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	ValidUntil int64  `json:"validUntil"`
}

type SpotStateBalance struct {
	Coin     string `json:"coin"`
	Token    int64  `json:"token"`
	Total    string `json:"total"`
	Hold     string `json:"hold"`
	EntryNtl string `json:"entryNtl"`
}
type SpotState struct {
	Balances []SpotStateBalance `json:"balances"`
}
