package core

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/noerw/osem_notify/utils"
)

/**
 * in memory + yaml persisted cache for check results, ensuring we don't resend
 * notifications on every check

 * TODO: reminder functionality: extract additional results with Status ERR
 * from cache with time.Since(lastNotifyDate) > remindAfter.
 * would require to serialize the full result..
 */

var cache = viper.New()

func init() {
	fileName := utils.GetConfigFile("osem_notify_cache")

	cache.SetConfigType("yaml")
	cache.SetConfigFile(fileName)

	if _, err := os.Stat(fileName); err == nil {
		err := cache.ReadInConfig()
		if err != nil {
			log.Error("Error when reading cache file:", err)
		}
	}
}

func (results BoxCheckResults) filterChangedFromCache() BoxCheckResults {
	remaining := BoxCheckResults{}

	for box, boxResults := range results {
		// get results from cache. they are indexed by an event ID per boxId
		// filter, so that only changed result.Status remain
		remaining[box] = []CheckResult{}
		for _, result := range boxResults {
			cached := cache.GetStringMap(fmt.Sprintf("watchcache.%s.%s", box.Id, result.EventID()))
			if result.Status != cached["laststatus"] {
				remaining[box] = append(remaining[box], result)
			}
		}
	}

	return remaining
}

func updateCache(box *Box, results []CheckResult) error {
	for _, result := range results {
		key := fmt.Sprintf("watchcache.%s.%s", box.Id, result.EventID())
		cache.Set(key+".laststatus", result.Status)
	}
	return cache.WriteConfig()
}

func ClearCache() error {
	fileName := utils.GetConfigFile("osem_notify_cache")
	_, err := os.Stat(fileName)
	if err != nil {
		return nil
	}
	return os.Remove(fileName)
}

func PrintCache() {
	for key, val := range cache.AllSettings() {
		log.Infof("%20s: %v", key, val)
	}
}
