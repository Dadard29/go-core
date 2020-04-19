package models

import "time"

type Subscription struct {
	AccessToken    string    `gorm:"type:varchar(70);index:access_token;primary_key"`
	ProfileKey     string    `gorm:"type:varchar(70);index:profile_key"`
	Profile        Profile   `gorm:"foreignkey:ProfileKey"`
	ApiName        string    `gorm:"type:varchar(70);index:api_name"`
	Api            ApiModel  `gorm:"foreignkey:ApiName"`
	DateSubscribed time.Time `gorm:"type:date;index:date_subscribed"`
	RequestCount   int       `gorm:"type:int;index:request_count"`
	FromEchoSlam bool `gorm:"type:bool;index:from_echo_slam"`
}

func NewSubscriptionJson(s Subscription, a ApiModel, quota int) SubscriptionJson {
	return SubscriptionJson{
		AccessToken:    s.AccessToken,
		Api:            a,
		DateSubscribed: s.DateSubscribed,
		RequestCount:   s.RequestCount,
		Quota:          quota,
	}
}

type SubscriptionJson struct {
	AccessToken    string
	Api            ApiModel
	DateSubscribed time.Time
	RequestCount   int
	Quota          int
}

func (Subscription) TableName() string {
	return "subscription"
}

