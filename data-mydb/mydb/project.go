package mydb

import (
	"request-matcher-openai/gocommon/util"
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID           string         `gorm:"primary_key;index;type:char(255);not null" json:"id"`
	CreatedAt    time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name         string         `json:"name"`
	URL          string         `json:"url" gorm:"url"`
	ContactEmail string         `json:"contact_email"`
	ContactPhone string         `json:"contact_phone"`
	Logo         string         `json:"logo"`
	Criterial    string         `json:"criterial"`
	Details      string         `json:"details"`
}

func (Project) TableName() string {
	return "project"
}

func (Project) GetType() string {
	return "project"
}

func (p Project) GetID() string {
	return p.ID
}

func (p *Project) SetID(id string) {
	p.ID = id
}

func (p Project) GetName() string {
	return ""
}

func (p Project) GetUpdateTimestamp() time.Time {
	return p.UpdatedAt
}

func (p Project) GetCreateTimestamp() time.Time {
	return p.CreatedAt
}

func (p *Project) SetUpdateTimestamp(t time.Time) {
	p.UpdatedAt = t
}

func (p Project) GetDeleteTimestamp() gorm.DeletedAt {
	return p.DeletedAt
}

func (p Project) IsDeleted() bool {
	return p.DeletedAt.Valid
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = util.GetUUID()
	}
	return nil
}

type ProjectShort struct {
	ID           string         `gorm:"primary_key;index;type:char(255);not null" json:"id"`
	CreatedAt    time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name         string         `json:"name"`
	URL          string         `json:"url" gorm:"url"`
	ContactEmail string         `json:"contact_email"`
	ContactPhone string         `json:"contact_phone"`
	Logo         string         `json:"logo"`
	Criterial    string         `json:"criterial"`
	Details      string         `json:"details"`
}

func (ProjectShort) TableName() string {
	return "project_short"
}

func (ProjectShort) GetType() string {
	return "project_short"
}

func (p ProjectShort) GetID() string {
	return p.ID
}

func (p *ProjectShort) SetID(id string) {
	p.ID = id
}

func (p ProjectShort) GetName() string {
	return ""
}

func (p ProjectShort) GetCreateTimestamp() time.Time {
	return p.CreatedAt
}

func (p ProjectShort) GetUpdateTimestamp() time.Time {
	return p.UpdatedAt
}

func (p *ProjectShort) SetUpdateTimestamp(t time.Time) {
	p.UpdatedAt = t
}

func (p ProjectShort) GetDeleteTimestamp() gorm.DeletedAt {
	return p.DeletedAt
}

func (p ProjectShort) IsDeleted() bool {
	return false
}

type ProjectMessage struct {
	Name         string `json:"name"`
	URL          string `json:"url" gorm:"url"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`
	Logo         string `json:"logo"`
	Criterial    string `json:"criterial"`
	Details      string `json:"details"`
}

func ConvertMessageToProject(msg ProjectMessage) Project {
	data := Project{
		Name:         msg.Name,
		URL:          msg.URL,
		ContactEmail: msg.ContactEmail,
		ContactPhone: msg.ContactPhone,
		Logo:         msg.Logo,
		Criterial:    msg.Criterial,
		Details:      msg.Details,
	}
	return data
}
