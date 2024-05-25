package mydb

import (
	"request-matcher-openai/gocommon/util"
	"time"

	"gorm.io/gorm"
)

type QualificationMatch struct {
	ID        string         `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UserID    string         `json:"user_id" gorm:"user_id"`
	ProjectID string         `json:"project_id" gorm:"project_id"`
}

func (QualificationMatch) TableName() string {
	return "qualification_match"
}

func (QualificationMatch) GetType() string {
	return "qualification_match"
}

func (p QualificationMatch) GetID() string {
	return p.ID
}

func (p *QualificationMatch) SetID(id string) {
	p.ID = id
}

func (p QualificationMatch) GetName() string {
	return ""
}

func (p QualificationMatch) GetUpdateTimestamp() time.Time {
	return p.UpdatedAt
}

func (p QualificationMatch) GetCreateTimestamp() time.Time {
	return p.CreatedAt
}

func (p *QualificationMatch) SetUpdateTimestamp(t time.Time) {
	p.UpdatedAt = t
}

func (p QualificationMatch) GetDeleteTimestamp() gorm.DeletedAt {
	return p.DeletedAt
}

func (p QualificationMatch) IsDeleted() bool {
	return p.DeletedAt.Valid
}

func (p *QualificationMatch) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = util.GetUUID()
	}
	return nil
}

type QualificationMatchShort struct {
	ID        string         `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UserID    string         `json:"user_id" gorm:"user_id"`
	ProjectID string         `json:"project_id" gorm:"project_id"`
}

func (QualificationMatchShort) TableName() string {
	return "qualification_match_short"
}

func (QualificationMatchShort) GetType() string {
	return "qualification_match_short"
}

func (p QualificationMatchShort) GetID() string {
	return p.ID
}

func (p *QualificationMatchShort) SetID(id string) {
	p.ID = id
}

func (p QualificationMatchShort) GetName() string {
	return ""
}

func (p QualificationMatchShort) GetCreateTimestamp() time.Time {
	return p.CreatedAt
}

func (p QualificationMatchShort) GetUpdateTimestamp() time.Time {
	return p.UpdatedAt
}

func (p *QualificationMatchShort) SetUpdateTimestamp(t time.Time) {
	p.UpdatedAt = t
}

func (p QualificationMatchShort) GetDeleteTimestamp() gorm.DeletedAt {
	return p.DeletedAt
}

func (p QualificationMatchShort) IsDeleted() bool {
	return false
}
