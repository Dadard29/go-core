package models

import (
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

type Profile struct {
	ProfileKey      string    `gorm:"type:varchar(70);index:profile_key;primary_key"`
	Username        string    `gorm:"type:varchar(70);index:username"`
	PasswordEncrypt string    `gorm:"type:varchar(70);index:password_encrypt"`
	DateCreated     time.Time `gorm:"type:date;index:date_created"`
}

type ProfileChangePassword struct {
	NewPassword string
}

type ProfileJson struct {
	ProfileKey  string
	Username    string
	DateCreated time.Time
}

func NewProfileJson(p Profile) ProfileJson {
	return ProfileJson{
		ProfileKey:  p.ProfileKey,
		Username:    p.Username,
		DateCreated: p.DateCreated,
	}
}

func (Profile) TableName() string {
	return "profile"
}

func (Profile) NewProfileKey() string {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	randInt := rand.Int()
	randBytes := []byte(strconv.Itoa(randInt))
	hash := sha256.New()
	hash.Write(randBytes)
	key := fmt.Sprintf("%x", hash.Sum(nil))
	return key
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Subscription struct {
	AccessToken    string    `gorm:"type:varchar(70);index:access_token;primary_key"`
	ProfileKey     string    `gorm:"type:varchar(70);index:profile_key"`
	Profile        Profile   `gorm:"foreignkey:ProfileKey"`
	ApiName        string    `gorm:"type:varchar(70);index:api_name"`
	Api            ApiModel  `gorm:"foreignkey:ApiName"`
	DateSubscribed time.Time `gorm:"type:date;index:date_subscribed"`
}

type SubscriptionJson struct {
	AccessToken    string
	Profile        ProfileJson
	Api            ApiModel
	DateSubscribed time.Time
}

func (Subscription) TableName() string {
	return "subscription"
}
