package api

import (
	"log"
	"sync"
	"time"
	"github.com/gin-gonic/gin"
)

type bruteForceEntry struct {
	Attempts      int
	FirstAttempt  time.Time
	BanUntil      time.Time
	BanLevel      int // 0: нет, 1: 1м, 2: 10м, 3: 1ч, 4: сутки
}

var (
	bruteForceMap = make(map[string]*bruteForceEntry)
	bruteForceMu  sync.Mutex
)

func getBruteForceKey(c *gin.Context) string {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		apiKey = c.GetHeader("Authorization")
		if len(apiKey) > 7 && apiKey[:7] == "Bearer " {
			apiKey = apiKey[7:]
		}
	}
	ip := c.ClientIP()
	return ip+":"+apiKey
}

func BruteForceCheck(c *gin.Context) (bool, string) {
	key := getBruteForceKey(c)
	bruteForceMu.Lock()
	entry, ok := bruteForceMap[key]
	if !ok {
		entry = &bruteForceEntry{FirstAttempt: time.Now()}
		bruteForceMap[key] = entry
	}
	now := time.Now()
	if entry.BanUntil.After(now) {
		msg := "Слишком много попыток ввести корректный API ключ. Если утрачен доступ к ключу, перевыпустите его в личном кабинете. Следующее окно возможностей откроется через "
		switch entry.BanLevel {
		case 1:
			msg += "1 минуту."
		case 2:
			msg += "10 минут."
		case 3:
			msg += "1 час."
		case 4:
			msg += "24 часа."
		default:
			msg += entry.BanUntil.Sub(now).Round(time.Second).String()
		}
		log.Printf("BRUTEFORCE: BLOCKED %s (level %d) до %s", key, entry.BanLevel, entry.BanUntil)
		bruteForceMu.Unlock()
		return true, msg
	}
	bruteForceMu.Unlock()
	c.Set("bruteforce_key", key)
	return false, ""
}

func RegisterBruteForceFailure(c *gin.Context) {
	key, ok := c.Get("bruteforce_key")
	if !ok {
		keyStr := getBruteForceKey(c)
		key = keyStr
		c.Set("bruteforce_key", keyStr)
	}
	bruteForceMu.Lock()
	defer bruteForceMu.Unlock()
	entry := bruteForceMap[key.(string)]
	now := time.Now()
	if now.Sub(entry.FirstAttempt) > time.Minute {
		entry.Attempts = 1
		entry.FirstAttempt = now
		entry.BanLevel = 0
		entry.BanUntil = time.Time{}
		log.Printf("BRUTEFORCE: RESET %s", key)
		return
	}
	entry.Attempts++
	log.Printf("BRUTEFORCE: ATTEMPT %s #%d", key, entry.Attempts)
	if entry.Attempts > 10 {
		switch entry.BanLevel {
		case 0:
			entry.BanLevel = 1
			entry.BanUntil = now.Add(time.Minute)
			log.Printf("BRUTEFORCE: BAN 1m %s", key)
		case 1:
			entry.BanLevel = 2
			entry.BanUntil = now.Add(10 * time.Minute)
			log.Printf("BRUTEFORCE: BAN 10m %s", key)
		case 2:
			entry.BanLevel = 3
			entry.BanUntil = now.Add(time.Hour)
			log.Printf("BRUTEFORCE: BAN 1h %s", key)
		case 3:
			entry.BanLevel = 4
			entry.BanUntil = now.Add(24 * time.Hour)
			log.Printf("BRUTEFORCE: BAN 24h %s", key)
		}
	}
} 