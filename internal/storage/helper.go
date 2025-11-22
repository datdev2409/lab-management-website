package storage

import (
	"fmt"
	"time"
	_ "time/tzdata"
)

func GetCurrentTime() time.Time {
	hcm, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		fmt.Println("Error:", err)
		return time.Now()
	}

	return time.Now().In(hcm)
}
