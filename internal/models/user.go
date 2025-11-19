package models

type User struct {
	BaseModel
	FirstName string `gorm:"type:varchar;not null" json:"first_name"`
	LastName  string `gorm:"type:varchar;not null" json:"last_name"`
	Email     string `gorm:"type:varchar;not null" json:"email"`
	Password  string `gorm:"type:varchar;not null" json:"-"`
}

type CreateUserDTO struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"email,required"`
	Password  string `json:"password" validate:"required"`
}
