package responsemodel

type LineSaleClosing struct {
	ID          uint
	Date        string
	Sale_amount float64
	Discount    float64
	Actual_sale float64
	Cash        float64
	Account     float64
}

type LineSaleStockOut struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	StockIn  int    `json:"stockin"`
	StockOut int    `json:"stockout"`
}

type Vehicle struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	License string `json:"license"`
}

type StockUpdate struct {
	ItemID   string `json:"item_id"`
	ItemName string `json:"item_name"`
	StockIn  int    `json:"stock_in"`
	StockOut int    `json:"stock_out"`
}

type StockUpdateRequest struct {
	VehicleID    string        `json:"vehicle_id"`
	StockUpdates []StockUpdate `json:"stock_updates"`
}

type VyaparResponse struct {
	ItemName   string `json:"item_name"`
	Difference int    `json:"difference"`
	Status     string `json:"status"`
}

type VyaparUpdate struct {
	ItemID   string `json:"item_id"`
	ItemName string `json:"item_name"`
	StockOut int    `json:"stock_out"`
}

type VyaparUpdateRequest struct {
	StockUpdates []VyaparUpdate `json:"stock_updates"`
}