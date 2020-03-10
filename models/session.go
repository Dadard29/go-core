package models

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Session struct {
	AccessToken string        `gorm:"type:varchar(30);index:access_token;primary_key"`
	Duration    time.Duration `gorm:"type:int;index:duration"`
	CreatedAt   time.Time     `gorm:"type:date;index:created_at"`
}

func (Session) NewAccessToken() string {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	randInt := rand.Int()
	randBytes := []byte(strconv.Itoa(randInt))
	hash := sha256.New()
	hash.Write(randBytes)
	key := fmt.Sprintf("%x", hash.Sum(nil))
	return key[:8]
}
