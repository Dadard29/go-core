package models

import (
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

// Entity
type Profile struct {
	ProfileKey      string    `gorm:"type:varchar(70);index:profile_key;primary_key"`
	Username        string    `gorm:"type:varchar(70);index:username"`
	PasswordEncrypt string    `gorm:"type:varchar(70);index:password_encrypt"`
	DateCreated     time.Time `gorm:"type:date;index:date_created"`
	Silver          bool      `gorm:"type:bool;index:silver"`
	RecoverBy       string    `gorm:"type:varchar(70);index:recover_by"`
	Contact         string    `gorm:"type:varchar(70);index:contact"`
	BeNotified      bool      `gorm:"type:bool;index:be_notified"`
}

func (Profile) TableName() string {
	return "profile"
}

type TempProfile struct {
	ConfirmationCode string    `gorm:"type:varchar(30);index:confirmation_code;primary_key"`
	Username         string    `gorm:"type:varchar(70);index:username"`
	PasswordEncrypt  string    `gorm:"type:varchar(70);index:password_encrypt"`
	ExpirationTime   time.Time `gorm:"type:datetime;index:expiration_time"`
}

func (TempProfile) TableName() string {
	return "temp_profile"
}

// DTO
type ProfileJson struct {
	ProfileKey  string
	Username    string
	DateCreated time.Time
	Silver      bool
	RecoverBy   string
	Contact     string
	BeNotified  bool
}

func NewProfileJson(p Profile) ProfileJson {
	return ProfileJson{
		ProfileKey:  p.ProfileKey,
		Username:    p.Username,
		DateCreated: p.DateCreated,
		Silver:      p.Silver,
		RecoverBy:   p.RecoverBy,
		Contact:     p.Contact,
		BeNotified:  p.BeNotified,
	}
}

// methods
func NewProfileKey() string {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	randInt := rand.Int()
	randBytes := []byte(strconv.Itoa(randInt))
	hash := sha256.New()
	hash.Write(randBytes)
	key := fmt.Sprintf("%x", hash.Sum(nil))
	return key
}

type ProfileChangePassword struct {
	NewPassword string
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
