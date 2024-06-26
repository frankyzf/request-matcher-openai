package mydb

import (
	"request-matcher-openai/gocommon/util"
	"time"

	"gorm.io/gorm"
)

// Generated by https://quicktype.io

type Account struct {
	ID                 string         `gorm:"primary_key;index;type:char(255);not null" json:"id"`
	CreatedAt          time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	CountryCode        string         `json:"country_code" gorm:"default:'65'"`
	Name               string         `json:"name"`
	Password           string         `json:"password"`
	Phone              string         `json:"phone"`
	Email              string         `json:"email" gorm:"not null"`
	NeedChangePassword bool           `json:"need_change_password" gorm:"default:true"`
	Enable             bool           `json:"enable" gorm:"default:true"`
	LastLoginTime      time.Time      `json:"last_login_time" gorm:"default:current_timestamp"`
	PasswordUpdateTime time.Time      `json:"password_update_time" gorm:"default:current_timestamp"`
	ContactEmail       string         `json:"contact_email"`
}

func (Account) TableName() string {
	return "account_user"
}

func (Account) GetType() string {
	return "account_user"
}

func (p Account) GetID() string {
	return p.ID
}

func (p *Account) SetID(id string) {
	p.ID = id
}

func (p Account) GetName() string {
	return p.Name
}

func (p Account) GetOwnerID() string {
	return p.ID
}

func (p Account) GetCreateTimestamp() time.Time {
	return p.CreatedAt
}

func (p Account) GetUpdateTimestamp() time.Time {
	return p.UpdatedAt
}

func (p *Account) SetUpdateTimestamp(t time.Time) {
	p.UpdatedAt = t
}

func (p Account) GetDeleteTimestamp() gorm.DeletedAt {
	return p.DeletedAt
}

func (p Account) IsDeleted() bool {
	return p.DeletedAt.Valid
}

func (p *Account) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = util.GetUUID()
	}
	return nil
}

type AccountMessage struct {
	CountryCode  string  `json:"country_code"`
	Name         string  `json:"name"`
	Password     string  `json:"password"`
	Email        string  `json:"email"`
	ContactEmail *string `json:"contact_email"`
	Phone        string  `json:"phone"`
	Token        string  `json:"token" `
	Captcha      string  `json:"captcha" `
}

type AccountShort struct {
	ID                 string         `json:"id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	CountryCode        string         `json:"country_code"`
	Name               string         `json:"name"`
	Phone              string         `json:"phone"`
	Email              string         `json:"email"`
	AccountType        string         `json:"account_type"`
	NeedChangePassword bool           `json:"need_change_password"`
	Enable             bool           `json:"enable"`
	PasswordUpdateTime time.Time      `json:"password_update_time" gorm:"default:current_timestamp"`
	ContactEmail       string         `json:"contact_email"`
}

func (AccountShort) TableName() string {
	return ""
}

func (AccountShort) GetType() string {
	return "account_user_short"
}

func (p AccountShort) GetID() string {
	return p.ID
}

func (p *AccountShort) SetID(id string) {
	p.ID = id
}

func (p AccountShort) GetName() string {
	return p.Name
}

func (p AccountShort) GetOwnerID() string {
	return p.ID
}

func (p AccountShort) IsDeleted() bool {
	return false //normally it is not
}

func (p AccountShort) GetCreateTimestamp() time.Time {
	return p.CreatedAt
}

func (p AccountShort) GetUpdateTimestamp() time.Time {
	return p.UpdatedAt
}

func (p *AccountShort) SetUpdateTimestamp(t time.Time) {
	p.UpdatedAt = t
}

func (p AccountShort) GetDeleteTimestamp() gorm.DeletedAt {
	return p.DeletedAt
}

func ConvertBaseAccountToUserShort(ba BaseAccount) UserShort {
	item := UserShort{
		ID:                 ba.ID,
		CountryCode:        ba.CountryCode,
		Name:               ba.Name,
		Phone:              ba.Phone,
		Email:              ba.Email,
		ContactEmail:       ba.ContactEmail,
		NeedChangePassword: ba.NeedChangePassword,
		Enable:             ba.Enable,
		Approve:            ba.Approve,
		AccountType:        ba.AccountType,
		PasswordUpdateTime: ba.PasswordUpdateTime,
	}
	return item
}

func ConvertAccountToBaseAccount(account Account) BaseAccount {
	ba := BaseAccount{
		ID:                 account.ID,
		CreatedAt:          account.CreatedAt,
		UpdatedAt:          account.UpdatedAt,
		DeletedAt:          account.DeletedAt,
		CountryCode:        account.CountryCode,
		Name:               account.Name,
		Password:           account.Password,
		Phone:              account.Phone,
		Email:              account.Email,
		ContactEmail:       account.ContactEmail,
		NeedChangePassword: account.NeedChangePassword,
		Enable:             account.Enable,
		PasswordUpdateTime: account.PasswordUpdateTime,
		Approve:            true,
		AccountType:        "account",
	}
	return ba
}

func ConvertBaseAccountToAccount(ba BaseAccount) Account {
	item := Account{
		ID:                 ba.ID,
		CreatedAt:          ba.CreatedAt,
		UpdatedAt:          ba.UpdatedAt,
		DeletedAt:          ba.DeletedAt,
		CountryCode:        ba.CountryCode,
		Name:               ba.Name,
		Password:           ba.Password,
		Phone:              ba.Phone,
		Email:              ba.Email,
		NeedChangePassword: ba.NeedChangePassword,
		Enable:             ba.Enable,
		PasswordUpdateTime: ba.PasswordUpdateTime,
	}
	return item
}

type AccountUpdatePasswordMessage struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type AccountProfileUpdateMessage struct {
	Name          string `json:"name"`
	PreferPrinter string `json:"prefer_printer"`
}
