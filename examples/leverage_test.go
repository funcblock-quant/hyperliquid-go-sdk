package examples

import "testing"

func TestUpdateLeverage(t *testing.T) {
	exchange := getTestExchange(t)
	coin := "SOL"
	leverage := 10

	t.Run("isolated margin", func(t *testing.T) {
		err := exchange.UpdateLeverage(coin, false, leverage)
		if err != nil {
			t.Fatalf("Failed to update leverage: %v", err)
		}
		t.Log("Update leverage successfully")
	})

	t.Run("cross margin", func(t *testing.T) {
		err := exchange.UpdateLeverage(coin, true, leverage)
		if err != nil {
			t.Fatalf("Failed to update leverage: %v", err)
		}
		t.Log("Update leverage successfully")
	})
}

func TestUpdateIsolatedMargin(t *testing.T) {
	exchange := getTestExchange(t)

	amount := 15.0 // Amount in USD
	coin := "BTC"

	err := exchange.UpdateIsolatedMargin(coin, amount)
	if err != nil {
		t.Fatalf("Failed to update isolated margin: %v", err)
	}
	t.Log("Update isolated margin successfully")
}
