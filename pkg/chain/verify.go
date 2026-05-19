package chain

import "time"

func IsValidChain(hash string, custodyCount int) bool {
	return hash != "" && custodyCount >= 1
}

func HasMonotonicCustodyTimestamps(times []time.Time) bool {
	if len(times) <= 1 {
		return true
	}
	for i := 1; i < len(times); i++ {
		if times[i].Before(times[i-1]) {
			return false
		}
	}
	return true
}

func HasRequiredCustodyFields(actor, action string) bool {
	return actor != "" && action != ""
}
