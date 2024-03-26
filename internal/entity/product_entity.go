package entity

type Product struct {
	Id     int    `gorm:"column:id;primaryKey"`
	Name   string `gorm:"column:name"`
	Stock  int32  `gorm:"column:stock"`
	Price  int32  `gorm:"column:price"`
	UserId int    `gorm:"column:user_id"`
	User   User   `gorm:"foreignKey:user_id;references:id"`
}
