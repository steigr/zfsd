package util

import "strconv"

func QuotaToInt(quota string) int64 {
	q1, err := strconv.ParseInt(quota,10,64)
	if err != nil {
		return 0
	}
	return q1
}
