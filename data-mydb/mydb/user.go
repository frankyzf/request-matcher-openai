package mydb

import (
	"time"

	"request-matcher-openai/gocommon/util"

	"gorm.io/gorm"
)

type User struct {
	ID                 string         `gorm:"primary_key;index;type:char(255);not null" json:"id"`
	CreatedAt          time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	ExpiredAt          *time.Time     `sql:"index" json:"expired_at"`
	CountryCode        string         `json:"country_code" gorm:"default:'65'"` //singapore
	Name               string         `json:"name"`
	Password           string         `json:"password"`
	Phone              string         `json:"phone"`
	Email              string         `json:"email"`
	PasswordUpdateTime time.Time      `json:"password_update_time" gorm:"default:current_timestamp"`
	Birthday           string         `json:"birthday"`
	Nationality        string         `json:"nationality"`
	PostCode           string         `json:"post_code"`
	UnitNumber         string         `json:"unit_number"`
	Address            string         `json:"address"`
	NeedChangePassword bool           `json:"need_change_password" gorm:"default:true"`
	Enable             bool           `json:"enable" gorm:"default:true"`
	Approve            bool           `json:"approve"  gorm:"default:true"`
	Remark             string         `json:"remark"`
	Questions          string         `json:"questions" gorm:"type:text"`
}

func (User) TableName() string {
	return "user"
}

func (User) GetType() string {
	return "user"
}

func (p User) GetID() string {
	return p.ID
}

func (p *User) SetID(id string) {
	p.ID = id
}

func (p User) GetName() string {
	return ""
}

func (p User) GetUpdateTimestamp() time.Time {
	return p.UpdatedAt
}

func (p User) GetCreateTimestamp() time.Time {
	return p.CreatedAt
}

func (p *User) SetUpdateTimestamp(t time.Time) {
	p.UpdatedAt = t
}

func (p User) GetDeleteTimestamp() gorm.DeletedAt {
	return p.DeletedAt
}

func (p User) IsDeleted() bool {
	return p.DeletedAt.Valid
}

func (p *User) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = util.GetUUID()
	}
	return nil
}

type UserShort struct {
	ID                 string         `gorm:"primary_key;index;type:char(255);not null" json:"id"`
	CreatedAt          time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	ExpiredAt          *time.Time     `sql:"index" json:"expired_at"`
	CountryCode        string         `json:"country_code" gorm:"default:'65'"` //singapore
	Name               string         `json:"name"`
	Phone              string         `json:"phone"`
	Email              string         `json:"email"`
	Birthday           string         `json:"birthday"`
	Nationality        string         `json:"nationality"`
	PostCode           string         `json:"post_code"`
	UnitNumber         string         `json:"unit_number"`
	Address            string         `json:"address"`
	PasswordUpdateTime time.Time      `json:"password_update_time" gorm:"default:current_timestamp"`
	Token              string         `json:"token"`
	ContactEmail       string         `json:"contact_email"`
	NeedChangePassword bool           `json:"need_change_password" gorm:"default:true"`
	Enable             bool           `json:"enable" gorm:"default:true"`
	Approve            bool           `json:"approve"  gorm:"default:false"`
	AccountType        string         `json:"account_type"`
	Remark             string         `json:"remark"`
	Questions          string         `json:"questions" gorm:"type:text"`
	APIKey             string         `json:"api_key"`
}

func (UserShort) TableName() string {
	return "user_short"
}

func (UserShort) GetType() string {
	return "user_short"
}

func (p UserShort) GetID() string {
	return p.ID
}

func (p *UserShort) SetID(id string) {
	p.ID = id
}

func (p UserShort) GetName() string {
	return ""
}

func (p UserShort) GetCreateTimestamp() time.Time {
	return p.CreatedAt
}

func (p UserShort) GetUpdateTimestamp() time.Time {
	return p.UpdatedAt
}

func (p *UserShort) SetUpdateTimestamp(t time.Time) {
	p.UpdatedAt = t
}

func (p UserShort) GetDeleteTimestamp() gorm.DeletedAt {
	return p.DeletedAt
}

func (p UserShort) IsDeleted() bool {
	return false
}

type UserLoginMessage struct {
	Email                  string `json:"email"`
	Phone                  string `json:"phone"`
	Password               string `json:"password"`
	Captcha                string `json:"captcha"`
	Code                   string `json:"code,omitempty"`
	ForceSkipCaptureVerify bool   `json:"force_skip_capture_verify"`
}

type UserSignupMessage struct {
	CountryCode string `json:"country_code" gorm:"default:'65'"` //singapore
	Name        string `json:"name"`
	Password    string `json:"password,omitempty"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Birthday    string `json:"birthday"`
	Nationality string `json:"nationality"`
	PostCode    string `json:"post_code"`
	UnitNumber  string `json:"unit_number"`
	Address     string `json:"address"`
	Token       string `json:"token,omitempty" gorm:"-"`
	Captcha     string `json:"captcha,omitempty" gorm:"-"`
	ExpiredAt   string `json:"expired_at"`
}

type UserUpdateMessage struct {
	CountryCode string  `json:"country_code" gorm:"default:'65'"` //singapore
	UserID      string  `json:"user_id"`
	Name        string  `json:"name"`
	Phone       *string `json:"phone"`
	Email       string  `json:"email"`
	Remark      string  `json:"remark,omitempty"`
	ExpiredAt   string  `json:"expired_at"`
	Birthday    string  `json:"birthday"`
	Nationality string  `json:"nationality"`
	PostCode    string  `json:"post_code"`
	UnitNumber  string  `json:"unit_number"`
	Address     string  `json:"address"`
}

type UserMessage struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func ConvertMessageToUser(msg UserMessage) User {
	data := User{
		Name:  msg.Name,
		Phone: msg.Phone,
		Email: msg.Email,
	}
	return data
}

type UserUpdatePasswordMessage struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	Code        string `json:"code"`
}

type UserSelfDeleteMessage struct {
	Code string `json:"code"`
}

func ConvertUserToBaseAccount(user User) BaseAccount {
	ba := BaseAccount{
		ID:                 user.ID,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
		DeletedAt:          user.DeletedAt,
		Name:               user.Name,
		Password:           user.Password,
		Phone:              user.Phone,
		Email:              user.Email,
		PasswordUpdateTime: user.PasswordUpdateTime,
		NeedChangePassword: user.NeedChangePassword,
		Enable:             user.Enable,
		Approve:            user.Approve,
		AccountType:        "user",
	}
	return ba
}

func ConvertUserSignupMessageToUser(msg UserSignupMessage, loc *time.Location) User {
	if msg.ExpiredAt == "" {
		msg.ExpiredAt = "2049-12-31 23:59:59"
	}
	expiredAt, _ := time.ParseInLocation("2006-01-02 15:04:05", msg.ExpiredAt, loc)

	return User{
		CountryCode: msg.CountryCode,
		Name:        msg.Name,
		Password:    msg.Password,
		Email:       msg.Email,
		Phone:       msg.Phone,
		Birthday:    msg.Birthday,
		Nationality: msg.Nationality,
		PostCode:    msg.PostCode,
		UnitNumber:  msg.UnitNumber,
		ExpiredAt:   &expiredAt,
		Address:     msg.Address,
	}
}

func ConvertUserToUserShort(msg User) UserShort {
	return UserShort{
		ID:                 msg.ID,
		CountryCode:        msg.CountryCode,
		Name:               msg.Name,
		Birthday:           msg.Birthday,
		Nationality:        msg.Nationality,
		PostCode:           msg.PostCode,
		UnitNumber:         msg.UnitNumber,
		Address:            msg.Address,
		Phone:              msg.Phone,
		Email:              msg.Email,
		NeedChangePassword: msg.NeedChangePassword,
		Enable:             msg.Enable,
		Approve:            msg.Approve,
		ExpiredAt:          msg.ExpiredAt,
		PasswordUpdateTime: msg.PasswordUpdateTime,
		Remark:             msg.Remark,
	}
}

func ConvertUserUpdateMessageToUser(msg UserUpdateMessage, loc *time.Location) User {
	if msg.ExpiredAt == "" {
		msg.ExpiredAt = "2049-12-31 23:59:59"
	}
	expiredAt, _ := time.ParseInLocation("2006-01-02 15:04:05", msg.ExpiredAt, loc)
	phone := ""
	if msg.Phone != nil {
		phone = *msg.Phone
	}
	user := User{
		ID:          msg.UserID,
		Name:        msg.Name,
		Email:       msg.Email,
		Phone:       phone,
		Remark:      msg.Remark,
		Birthday:    msg.Birthday,
		Nationality: msg.Nationality,
		PostCode:    msg.PostCode,
		UnitNumber:  msg.UnitNumber,
		Address:     msg.Address,
		ExpiredAt:   &expiredAt,
	}
	return user
}
