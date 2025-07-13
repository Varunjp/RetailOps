package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         	uint   		`gorm:"primarykey;autoIncrement"`
	Username   	string 		`gorm:"not null"`
	EmpID      	string 		`gorm:"not null;unique;"`
	Phone      	string
	Password   	string 		`gorm:"not null"`
	SuperUser  	bool		`gorm:"default:false"`
	Created_at 	time.Time
	Status		string
	Attendance	[]Attendance	`gorm:"constraint:ONDELETE:CASCADE; foreignKey:UserID"`
}


type Attendance	struct{
	ID			uint 		`gorm:"primaryKey;autoIncrement"`
	UserID		uint		`gorm:"index"`
	Date		time.Time
	Status 		string 
	OTamount	float64	
	OTStatus	bool 
	User		User		`gorm:"constraint:ONDELETE:CASCADE; foreignKey:UserID"`
}

type Product struct{
	ID			uint					`gorm:"primaryKey;autoIncrement"`
	ItemID		string					`gorm:"not null"`
	ItemName	string 
	Status 		string 		
	Created_at  time.Time 
	ProductDetail []ProductDetail 		`gorm:"constraint:ONDELETE:CASCADE;foreignKey: ProductID"`
	CounterSale   []CounterSale 		`gorm:"constraint:ONDELETE:CASCADE;foreignKey: ProductID"`
	LineSale 	[]LineSale  			`gorm:"constraint:ONDELETE:CASCADE;foreignKey: ProductID"`
	Deleted_at  gorm.DeletedAt
}

type ProductDetail struct{
	ID				uint				`gorm:"primaryKey;autoIncrement"`
	ProductID		uint 				`gorm:"index"`
	Rate			float64 			`gorm:"not null"`
	Stock			int	
	Product 		Product				`gorm:"foreignKey:ProductID;constraint:ONDELETE:CASCADE"`	
	Deleted_at  	gorm.DeletedAt	
}

type CounterSale struct{
	ID 				uint			`gorm:"primaryKey;autoIncrement"`
	ProductID		uint			`gorm:"index"`
	ItemName 		string 		
	Rate			float64
	TotalAmount		float64
	Cash			float64
	Account			float64
	Discount		float64
	Balance			float64
	Created_at		time.Time
	Product			Product 		`gorm:"constraint:ONDELETE:CASCADE;foreignKey:ProductID"`
	Deleted_at  	gorm.DeletedAt
}

type LineSale struct{
	ID 				uint					`gorm:"primaryKey;autoIncrement"`
	EmpID			string 					`gorm:"not null"`
	ProductID 		uint					`gorm:"index"`
	Vehicle 		string 					
	ItemName 		string 		
	Rate			float64
	StockIn			int 
	StockOut		int
	Damage			int 
	Created_at 		time.Time
	Status 			bool 									
	Product			Product 				`gorm:"constraint:ONDELETE:CASCADE;foreignKey:ProductID"`
	VyaparSale 		VyaparSale				`gorm:"foreignKey:SaleID;constraint:OnDelete:CASCADE"`
	Deleted_at  	gorm.DeletedAt
}

type LineSaleClosing struct {
	ID 				uint 					`gorm:"primaryKey;autoIncrement"`
	EmpID 			string 					`gorm:"not null"`
	SaleAmount		float64
	Discount 		float64
	Vehicle 		string 
	ActualSale 		float64
	Cash 			float64
	AccountPayment	float64
	Balance 		float64
	Status 			bool 					`gorm:"default:false"`
	Created_at   	time.Time
	LineSaleExpenses []LineSaleExpenses 	`gorm:"constraint:ONDELETE:CASCADE;foreignKey:LineSaleID"`
}

type LineSaleExpenses struct{
	ID 				uint 					`gorm:"primaryKey;autoIncrement"`
	LineSaleID		uint 					`gorm:"not null"`
	Type 			string
	Amount 			float64
	Created_at 		time.Time
	LineSale		LineSaleClosing				`gorm:"constraint:ONDELETE:CASCADE;foreignKey:LineSaleID"`
}

type LineSaleEdit struct{
	ID 				uint 				`gorm:"primaryKey;autoIncrement"`
	LineSaleID 		uint 				`gorm:"not null"`
	StockOut 		int 				
	Created_at		time.Time
}

type VyaparSale struct {
	ID 				uint 			`gorm:"primaryKey;autoIncrement"`
	SaleID 			uint 			`gorm:"index"`
	StockOut 		int 			
	Created_at  	time.Time 
}

type Batch struct{
	ID 				uint		`gorm:"primaryKey;autoIncrement"`
	BatchID 		string		`gorm:"unique;not null"`
	TotalCount 		float64
	TotalWeight		float64
	Created_at		time.Time 
	BatchProduction []BatchProduction `gorm:"foreignKey:BID;constraint:OnDelete:CASCADE"`
}

type BatchProduction struct{
	ID 					uint  		`gorm:"primaryKey;autoIncrement"`
	BID   				uint  		`gorm:"index"`
	RawMaterials 		float64
	ItemProduced		float64
	Bottle				int
	Pouch				int 
	Created_at 			time.Time
	Batch				Batch			`gorm:"foreignKey:BID;constraint:OnDelete:CASCADE"`
	Deleted_at 			gorm.DeletedAt
}

type Credits struct{
	ID 			uint 		`gorm:"primaryKey;autoIncrement"`
	Type 		string		`gorm:"not null;unique"`
	Balance 	float64
	Created_at	time.Time
}

type Vehicle struct{
	ID			uint 		`gorm:"primaryKey;autoIncrement"`
	Name		string		`gorm:"not null"`
	License 	string		`gorm:"unique"`
	Created_at 	time.Time
	Delete_at 	gorm.DeletedAt
}