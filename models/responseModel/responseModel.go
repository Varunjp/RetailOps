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