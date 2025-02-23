package storage

import (
	"fmt"
	"time"
)

func GetCurrentTime() time.Time {
	hcm, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		fmt.Println("Error:", err)
		return time.Now()
	}

	return time.Now().In(hcm)
}
