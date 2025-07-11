package getitems

type inventoryItem struct {
	ItemID   string  `json:"itemId"`
	ItemName string  `json:"itemName"`
	Rate     float64 `json:"rate"`
	StockIn  float64 `json:"stockIn"`
}

type InventoryRequest struct {
	Transactions struct {
		Vehicle string `json:"vehicle"`
	} `json:"transaction"`
	Items []inventoryItem `json:"items"`
}
