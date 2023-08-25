package absensi

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetFirstLastDateCurrentMonth() (firstOfMonth, lastOfMonth time.Time) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth = firstOfMonth.AddDate(0, 1, -1)

	fmt.Println(firstOfMonth)
	fmt.Println(lastOfMonth)
	return

}

func GetTimestampFromObjectID(objectID primitive.ObjectID) time.Time {
	timestamp := objectID.Timestamp()
	return timestamp
}

func ConvertTimestampToJkt(waktu time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return waktu.In(loc)
}
