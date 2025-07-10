package getitems

type inventoryItem struct {
	ItemID   string  `json:"itemId"`
	ItemName string  `json:"itemName"`
	Rate     float64 `json:"rate"`
	StockIn  float64 `json:"stockIn"`
}

type InventoryRequest struct {
	Items []inventoryItem `json:"items"`
}