package managers

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/models"
	"github.com/Dadard29/go-core/repositories"
	"math/rand"
	"strconv"
	"time"
)

func subsGenerateAccessToken() string {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	i := rand.Int()

	keyByte := []byte(strconv.Itoa(i))

	hash := sha256.New()

	hash.Write(keyByte)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func sendNotificationQuotaReached(subDb models.Subscription, quota int, profile models.Profile) error {
	if profile.BeNotified {
		if handler, check := recoveryByMapping[profile.RecoverBy]; check {
			if subDb.RequestCount == quota {

				text := fmt.Sprintf("Your quota (%d) for the API `%s`"+
					" is reached\n", quota, subDb.ApiName)

				if err := handler.sendNotification(profile, text); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func sendNotificationQuotaAlmostReached(subDb models.Subscription, quota int, profile models.Profile) error {
	if profile.BeNotified {
		if handler, check := recoveryByMapping[profile.RecoverBy]; check {
			step := float32(0.9)
			// send notification if quota almost reached (90% exactly of maximum quota)
			if float32(subDb.RequestCount) == (step*float32(quota)) && subDb.RequestCount < quota {

				text := fmt.Sprintf("Your quota (%d) for the API `%s`"+
					" is about be reached\n", quota, subDb.ApiName)

				if err := handler.sendNotification(profile, text); err != nil {
					return err
				}

			}
		}
	}

	return nil
}

func getQuota(profile models.Profile) (int, error) {
	var quotaStr string
	var err error
	if profile.Silver {
		quotaStr, err = api.Api.Config.GetValueFromFile(
			config.Profile,
			config.ProfileQuota,
			config.ProfileQuotaSilver)

		if err != nil {
			return 0, err
		}
	} else {
		quotaStr, err = api.Api.Config.GetValueFromFile(
			config.Profile,
			config.ProfileQuota,
			config.ProfileQuotaNotSilver)

		if err != nil {
			return 0, err
		}
	}

	quota, err := strconv.Atoi(quotaStr)
	if err != nil {
		return 0, err
	}

	return quota, nil
}

func SubsManagerGetFromToken(subToken string) (models.SubscriptionJson, string, error) {
	var s models.SubscriptionJson

	subDb, msg, err := repositories.SubsGetFromPkAndToken(subToken)
	if err != nil {
		return s, msg, err
	}

	a, msg, err := repositories.ApiGet(subDb.ApiName)
	if err != nil {
		return s, msg, err
	}

	if subDb.FromEchoSlam {
		return models.NewSubscriptionJson(subDb, a, -1), "subs for echo slam retrieved", nil
	}

	profile, msg, err := repositories.ProfileGetFromKey(subDb.ProfileKey)
	if err != nil {
		return s, msg, err
	}

	quota, err := getQuota(profile)
	if err != nil {
		return s, "error getting quota", err
	}

	if subDb.RequestCount >= quota {
		msg := fmt.Sprintf("quota (%d) reached", quota)
		return s, msg, errors.New(msg)
	}

	newSubDb, msg, err := repositories.SubsUpdateRequestCount(&subDb)
	if err != nil {
		return s, msg, err
	}

	// send notif asynchronously to avoid request blocking
	go func() {
		if err := sendNotificationQuotaAlmostReached(*newSubDb, quota, profile); err != nil {
			logger.Warning(err.Error())
		}

		if err := sendNotificationQuotaReached(*newSubDb, quota, profile); err != nil {
			logger.Warning(err.Error())
		}
	}()

	return models.NewSubscriptionJson(*newSubDb, a, quota), "sub checked", nil
}

func SubsManagerGetFromApiName(apiName string, profileKey string, fromEchoSlam bool) (models.SubscriptionJson, string, error) {
	var s models.SubscriptionJson

	a, msg, err := repositories.ApiGet(apiName)
	if err != nil {
		return s, msg, err
	}

	profile, msg, err := repositories.ProfileGetFromKey(profileKey)
	if err != nil {
		return s, msg, err
	}

	quota, err := getQuota(profile)
	if err != nil {
		return s, "error getting quota", err
	}

	subDb, msg, err := repositories.SubsGetFromApiName(apiName, profileKey, fromEchoSlam)
	if err != nil {
		return s, msg, err
	}

	return models.NewSubscriptionJson(subDb, a, quota), "sub checked", nil
}

func SubsManagerList(profileKey string, fromEchoSlam bool) ([]models.SubscriptionJson, string, error) {
	var s []models.SubscriptionJson

	p, msg, err := repositories.ProfileGetFromKey(profileKey)
	if err != nil {
		return s, msg, err
	}

	quota, err := getQuota(p)
	if err != nil {
		return s, "error getting quota", err
	}

	subDb, msg, err := repositories.SubsListFromProfile(p)
	if err != nil {
		return s, msg, err
	}

	var subListJson = make([]models.SubscriptionJson, 0)
	for _, sub := range subDb {

		// filter
		if sub.FromEchoSlam != fromEchoSlam {
			continue
		}

		a, _, err := repositories.ApiGet(sub.ApiName)
		if err != nil {
			continue
		}

		subListJson = append(subListJson, models.NewSubscriptionJson(sub, a, quota))
	}

	return subListJson, msg, nil
}

func SubsManagerCreate(profileKey string, apiName string, fromEchoSlam bool) (models.SubscriptionJson, string, error) {
	var s models.SubscriptionJson

	p, msg, err := repositories.ProfileGetFromKey(profileKey)
	if err != nil {
		return s, msg, err
	}

	quota, err := getQuota(p)
	if err != nil {
		return s, "error getting quota", err
	}

	a, msg, err := repositories.ApiGet(apiName)
	if err != nil {
		return s, msg, err
	}

	if p.Silver != true && a.Restricted {
		msg := "you need to be silver labelled to access restricted APIs"
		return s, msg, errors.New(msg)
	}

	subDb, msg, err := repositories.SubsCreate(models.Subscription{
		AccessToken:    subsGenerateAccessToken(),
		ProfileKey:     p.ProfileKey,
		ApiName:        a.Name,
		DateSubscribed: time.Now(),
		RequestCount:   0,
		FromEchoSlam:   fromEchoSlam,
	})
	if err != nil {
		return s, msg, err
	}

	subJson := models.NewSubscriptionJson(subDb, a, quota)

	return subJson, msg, nil
}

func SubsManagerDelete(profileKey string, apiName string) (models.SubscriptionJson, string, error) {
	var s models.SubscriptionJson

	p, msg, err := repositories.ProfileGetFromKey(profileKey)
	if err != nil {
		return s, msg, err
	}

	quota, err := getQuota(p)
	if err != nil {
		return s, "error getting quota", err
	}

	a, msg, err := repositories.ApiGet(apiName)
	if err != nil {
		return s, msg, err
	}

	subDb, msg, err := repositories.SubsDelete(p.ProfileKey, a.Name)
	if err != nil {
		return s, msg, err
	}

	subJson := models.NewSubscriptionJson(subDb, a, quota)

	return subJson, msg, nil
}

func SubsManagerUpdate(profileKey string, apiName string) (models.SubscriptionJson, string, error) {
	var s models.SubscriptionJson

	p, msg, err := repositories.ProfileGetFromKey(profileKey)
	if err != nil {
		return s, msg, err
	}

	quota, err := getQuota(p)
	if err != nil {
		return s, "error getting quota", err
	}

	a, msg, err := repositories.ApiGet(apiName)
	if err != nil {
		return s, msg, err
	}

	subUpdated, msg, err := repositories.SubsRegenerateToken(profileKey, apiName, subsGenerateAccessToken())
	if err != nil {
		return s, msg, err
	}

	subJson := models.NewSubscriptionJson(subUpdated, a, quota)

	return subJson, "token regenerated", nil
}

// warning: admin operation
func SubsManagerResetRequestCount() (string, error) {
	allSubs := repositories.SubsGetAll()
	for _, s := range allSubs {
		_, _, err := repositories.SubsResetRequestCount(&s)
		if err != nil {
			logger.Warning(err.Error())
		}
	}

	return "all requests count reset", nil
}
