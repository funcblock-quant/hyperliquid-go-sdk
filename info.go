package sdk

import (
	"context" // Import the context package
	"encoding/json"
	"fmt"
)

type Info struct {
	client         *Client
	coinToAsset    map[string]int
	assetToDecimal map[int]int
	perpCoins      []string
	spotCoins      []string
}

func NewInfo(apiBaseURL string) (*Info, error) {
	info := &Info{
		client:         NewClient(context.Background(), apiBaseURL),
		coinToAsset:    make(map[string]int),
		assetToDecimal: make(map[int]int),
	}

	// Always attempt to fetch meta and spotMeta as skipMeta is effectively false
	// and meta/spotMeta are effectively nil at this point.
	var meta *Meta
	var spotMeta *SpotMeta
	var err error

	meta, err = info.Meta()
	if err != nil {
		return nil, fmt.Errorf("error getting meta info: %w", err)
	}

	spotMeta, err = info.SpotMeta()
	if err != nil {
		return nil, fmt.Errorf("error getting spot meta info: %w", err)
	}

	// Map perp assets
	if meta != nil {
		for asset, assetInfo := range meta.Universe {
			info.coinToAsset[assetInfo.Name] = asset
			info.assetToDecimal[asset] = assetInfo.SzDecimals
			info.perpCoins = append(info.perpCoins, assetInfo.Name)
		}
	}

	// Map spot assets starting at 10000
	if spotMeta != nil {
		for _, spotInfo := range spotMeta.Universe {
			asset := spotInfo.Index + 10000
			info.coinToAsset[spotInfo.Name] = asset
			info.assetToDecimal[asset] = spotMeta.Tokens[spotInfo.Tokens[0]].SzDecimals
		}

		for _, spotInfo := range spotMeta.Tokens {
			info.spotCoins = append(info.spotCoins, spotInfo.Name)
		}
	}

	return info, nil
}

func (i *Info) ApiBaseUrl() string {
	return i.client.baseURL
}

func (i *Info) PerpCoins() []string {
	return i.perpCoins
}

func (i *Info) SpotCoins() []string {
	return i.spotCoins
}

func (i *Info) Meta() (*Meta, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "meta",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch meta: %w", err)
	}

	var meta Meta
	if err := json.Unmarshal(resp, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal meta response: %w", err)
	}

	return &meta, nil
}

func (i *Info) SpotMeta() (*SpotMeta, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "spotMeta",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch spot meta: %w", err)
	}

	var spotMeta SpotMeta
	if err := json.Unmarshal(resp, &spotMeta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spot meta response: %w", err)
	}

	return &spotMeta, nil
}

func (i *Info) CoinToAsset(coin string) (int, error) {
	asset, exist := i.coinToAsset[coin]
	if !exist {
		return 0, fmt.Errorf("coin %s not found", coin)
	}
	return asset, nil
}

func (i *Info) AssetToDecimal(asset int) (int, error) {
	decimal, exist := i.assetToDecimal[asset]
	if !exist {
		return 0, fmt.Errorf("asset %d not found", asset)
	}
	return decimal, nil
}

func (i *Info) UserState(address string) (*UserState, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "clearinghouseState",
		"user": address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user state: %w", err)
	}

	var result UserState
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user state: %w", err)
	}
	return &result, nil
}

func (i *Info) SpotUserState(address string) (*SpotState, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "spotClearinghouseState",
		"user": address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch spot user state: %w", err)
	}

	var result SpotState
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spot user state: %w", err)
	}
	return &result, nil
}

func (i *Info) OpenOrders(address string) ([]OpenOrder, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "openOrders",
		"user": address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch open orders: %w", err)
	}

	var result []OpenOrder
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal open orders: %w", err)
	}
	return result, nil
}

func (i *Info) FrontendOpenOrders(address string) ([]FrontendOpenOrder, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "frontendOpenOrders",
		"user": address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch frontend open orders: %w", err)
	}

	var frontendOrders []FrontendOpenOrder
	if err := json.Unmarshal(resp, &frontendOrders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal frontend open orders: %w", err)
	}

	// Convert frontend orders to standard open orders
	return frontendOrders, nil
}

func (i *Info) UserDepositWithdrawTxs(address string, startTime, endTime *int64) ([]DepositWithdrawTx, error) {
	payload := map[string]any{
		"type":      "userNonFundingLedgerUpdates",
		"user":      address,
		"startTime": 0,
	}
	if endTime != nil {
		payload["endTime"] = *endTime
	}
	if startTime != nil {
		payload["startTime"] = *startTime
	}

	resp, err := i.client.post("/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user deposit & withdraw txs: %w", err)
	}

	var result []DepositWithdrawTx
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user deposit & withdraw txs: %w", err)
	}
	return result, nil
}

func (i *Info) UserPortfolio(address string) ([]PortFolioTimeRangeItem, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "portfolio",
		"user": address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user portfolio: %w", err)
	}

	var result []PortFolioTimeRangeItem
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user portfolio: %w", err)
	}
	return result, nil
}

func (i *Info) AllMids() (map[string]string, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "allMids",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all mids: %w", err)
	}

	var result map[string]string
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal all mids: %w", err)
	}
	return result, nil
}

func (i *Info) UserFills(address string) ([]Fill, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "userFills",
		"user": address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user fills: %w", err)
	}

	var result []Fill
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user fills: %w", err)
	}
	return result, nil
}

func (i *Info) UserFillsByTime(address string, startTime int64, endTime *int64) ([]Fill, error) {
	payload := map[string]any{
		"type":            "userFillsByTime",
		"user":            address,
		"startTime":       startTime,
		"aggregateByTime": true,
	}
	if endTime != nil {
		payload["endTime"] = *endTime
	}

	resp, err := i.client.post("/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user fills by time: %w", err)
	}

	var result []Fill
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user fills by time: %w", err)
	}
	return result, nil
}

func (i *Info) MetaAndAssetCtxs() (map[string]any, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "metaAndAssetCtxs",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch meta and asset contexts: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal meta and asset contexts: %w", err)
	}
	return result, nil
}

func (i *Info) SpotMetaAndAssetCtxs() (map[string]any, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "spotMetaAndAssetCtxs",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch spot meta and asset contexts: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spot meta and asset contexts: %w", err)
	}
	return result, nil
}

func (i *Info) FundingHistory(
	coin string,
	startTime int64,
	endTime *int64,
) ([]FundingHistory, error) {

	payload := map[string]any{
		"type":      "fundingHistory",
		"coin":      coin,
		"startTime": startTime,
	}
	if endTime != nil {
		payload["endTime"] = *endTime
	}

	resp, err := i.client.post("/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch funding history: %w", err)
	}

	var result []FundingHistory
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal funding history: %w", err)
	}
	return result, nil
}

func (i *Info) UserFundingHistory(
	user string,
	startTime int64,
	endTime *int64,
) ([]UserFundingHistory, error) {
	payload := map[string]any{
		"type":      "userFunding",
		"user":      user,
		"startTime": startTime,
	}
	if endTime != nil {
		payload["endTime"] = *endTime
	}

	resp, err := i.client.post("/info", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user funding history: %w", err)
	}

	var result []UserFundingHistory
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user funding history: %w", err)
	}
	return result, nil
}

func (i *Info) L2Snapshot(coin string) (*L2Book, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "l2Book",
		"coin": coin,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch L2 snapshot: %w", err)
	}

	var result L2Book
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal L2 snapshot: %w", err)
	}
	return &result, nil
}

func (i *Info) CandlesSnapshot(coin, interval string, startTime, endTime int64) ([]Candle, error) {
	req := map[string]any{
		"coin":      coin,
		"interval":  interval,
		"startTime": startTime,
		"endTime":   endTime,
	}

	resp, err := i.client.post("/info", map[string]any{
		"type": "candleSnapshot",
		"req":  req,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch candles snapshot: %w", err)
	}

	var result []Candle
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal candles snapshot: %w", err)
	}
	return result, nil
}

func (i *Info) UserFees(address string) (*UserFees, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "userFees",
		"user": address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user fees: %w", err)
	}

	var result UserFees
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user fees: %w", err)
	}
	return &result, nil
}

func (i *Info) UserStakingSummary(address string) (*StakingSummary, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "delegatorSummary",
		"user": address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch staking summary: %w", err)
	}

	var result StakingSummary
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal staking summary: %w", err)
	}
	return &result, nil
}

func (i *Info) UserStakingDelegations(address string) ([]StakingDelegation, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "delegations",
		"user": address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch staking delegations: %w", err)
	}

	var result []StakingDelegation
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal staking delegations: %w", err)
	}
	return result, nil
}

func (i *Info) UserStakingRewards(address string) ([]StakingReward, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "delegatorRewards",
		"user": address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch staking rewards: %w", err)
	}

	var result []StakingReward
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal staking rewards: %w", err)
	}
	return result, nil
}

func (i *Info) QueryOrderByOid(user string, oid int64) (*OpenOrder, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "orderStatus",
		"user": user,
		"oid":  oid,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order status: %w", err)
	}

	var result OpenOrder
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order status: %w", err)
	}
	return &result, nil
}

func (i *Info) QueryOrderByCloid(user string, cloid string) (*OpenOrder, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "orderStatus",
		"user": user,
		"oid":  cloid,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order status by cloid: %w", err)
	}

	var result OpenOrder
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order status: %w", err)
	}
	return &result, nil
}

func (i *Info) QueryReferralState(user string) (*ReferralState, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "referral",
		"user": user,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch referral state: %w", err)
	}

	var result ReferralState
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal referral state: %w", err)
	}
	return &result, nil
}

func (i *Info) QuerySubAccounts(user string) ([]SubAccount, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "subAccounts",
		"user": user,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sub accounts: %w", err)
	}

	var result []SubAccount
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sub accounts: %w", err)
	}
	return result, nil
}

func (i *Info) QueryUserToMultiSigSigners(multiSigUser string) ([]MultiSigSigner, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "userToMultiSigSigners",
		"user": multiSigUser,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch multi-sig signers: %w", err)
	}

	var result []MultiSigSigner
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal multi-sig signers: %w", err)
	}
	return result, nil
}

func (i *Info) ExtraAgents(user string) ([]ExtraAgent, error) {
	resp, err := i.client.post("/info", map[string]any{
		"type": "extraAgents",
		"user": user,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch extra agents: %w", err)
	}

	var result []ExtraAgent
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal extra agents: %w", err)
	}
	return result, nil
}
