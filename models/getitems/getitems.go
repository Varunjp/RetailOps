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

type expenseItem struct {
	Item   string  `json:"item"`
	Amount float64 `json:"amount"`
}

type ExpenseRequest struct {
	Transaction []expenseItem `json:"expenses"`
	VehicleID   string        `json:"vehicleId"`
}

type Linetable struct {
	Name     string  `json:"name"`
	StockIn  int     `json:"stock_in"`
	StockOut int     `json:"stock_out"`
	Rate     float64 `json:"rate"`
}

type Revenue struct {
	TotalAmount float64 `json:"total_amount"`
	Discount    float64 `json:"discount"`
	ActualSale  float64 `json:"actual_sale"`
	Cash        float64 `json:"cash"`
	Account     float64 `json:"account"`
}

type Expenses struct {
	Item   string  `json:"item"`
	Amount float64 `json:"amount"`
}