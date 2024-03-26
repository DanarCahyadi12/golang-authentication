package entity

type User struct {
	Id        int       `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email;unique"`
	Password  string    `gorm:"column:password"`
	CreatedAt uint8     `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt uint8     `gorm:"column:updated_at;autCreateTime:milli;autoUpdateTime:milli"`
	Products  []Product `gorm:"foreignKey:user_id;references:id"`
}
