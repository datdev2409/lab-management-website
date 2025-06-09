package models

import "go.mongodb.org/mongo-driver/v2/bson"

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
