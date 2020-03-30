package models

import "time"

// url are set as string to manage the DB transactions easier
type ApiModel struct {
	Name             string     `gorm:"type:varchar(30);index:name;primary_key"`
	DocumentationUrl string     `gorm:"type:varchar(70);index:documentation_url"`
	Hostname         string     `gorm:"type:varchar(70);index:hostname"`
	VCSUrl           string     `gorm:"type:varchar(70);index:vcs_url"`
	BuildUrl         string     `gorm:"type:varchar(70);index:build_url"`
	IconUrl          string     `gorm:"type:varchar(100);index:icon_url"`
	Image            string     `gorm:"type:varchar(70);index:image"`
	CreationDate     *time.Time `gorm:"type:date;index:creation_date"`
	IsStandard       bool       `gorm:"type:bool;index:is_standard"`
	Restricted       bool       `gorm:"type:bool;index:restricted"`
}

func (ApiModel) TableName() string {
	return "api"
}
