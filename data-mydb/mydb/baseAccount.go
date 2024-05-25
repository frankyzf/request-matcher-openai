package mydb

import (
	"gorm.io/gorm"
	"time"
)

type BaseAccount struct {
	ID                 string         `json:"id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name               string         `json:"name"`
	Password           string         `json:"password"`
	Phone              string         `json:"phone"`
	Email              string         `json:"email"`
	ContactEmail       string         `json:"contact_email"`
	AccountType        string         `json:"account_type"`
	PasswordUpdateTime time.Time      `json:"password_update_time" gorm:"default:current_timestamp"`
	CountryCode        string         `json:"country_code" gorm:"default:'65'"`
	NeedChangePassword bool           `json:"need_change_password" gorm:"default:true"`
	Enable             bool           `json:"enable" gorm:"default:true"`
	Approve            bool           `json:"approve"`
}
