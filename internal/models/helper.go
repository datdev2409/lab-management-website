package models

import (
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ToBSONDocument(data interface{}) (bson.M, error) {
	doc, err := bson.Marshal(data)
	if err != nil {
		return nil, err
	}
	var result bson.M
	if err := bson.Unmarshal(doc, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func NormalizeString(s string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(strings.ToLower(s), "")
}

func GeneratePatientID(phone, name string) string {
	return NormalizeString(phone) + NormalizeString(name)
}

func GenerateRecordID(patientID string, createdAt time.Time) string {
	return patientID + createdAt.Format("060102") // YYMMDD
}

func GenerateTestID(testName string) string {
	return NormalizeString(testName)
}

func GenerateComboID(comboName string) string {
	return NormalizeString(comboName)
}

func GenerateTrackingID(patientID string, createdAt time.Time) string {
	return patientID + createdAt.Format("060102")
}
