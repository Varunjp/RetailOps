package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         	uint   		`gorm:"primarykey;autoIncrement"`
	Username   	string 		`gorm:"not null"`
	EmpID      	string 		`gorm:"not null;unique; index"`
	Phone      	string
	Password   	string 		`gorm:"not null"`
	SuperUser  	bool
	Created_at 	time.Time
	Status		string
	Attendance	[]Attendance	`gorm:"constraint:ONDELETE:CASCADE; foreignKey:EmpID"`
}


type Attendance	struct{
	ID			uint 		`gorm:"primaryKey;autoIncrement"`
	EmpID		uint		`gorm:"index"`
	Date		time.Time
	Status 		string 
	OTamount	float64	
	OTStatus	bool 
	User		User		`gorm:"constraint:ONDELETE:CASCADE; foreignKey:EmpID"`
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
	ID 				uint		`gorm:"primaryKey;autoIncrement"`
	EmpID			string 
	ProductID 		uint		`gorm:"index"`
	ItemName 		string 		
	Rate			float64
	StockIn			int 
	StockOut		int
	Damage			int 
	SaleAmount		float64
	Discount		float64
	ActualSale		float64
	Cash			float64
	Balance			float64
	Created_at 		time.Time
	Product			Product 	`gorm:"constraint:ONDELETE:CASCADE;foreignKey: ProductID"`
	Deleted_at  	gorm.DeletedAt
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