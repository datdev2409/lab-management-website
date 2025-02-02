package models

import "go.mongodb.org/mongo-driver/v2/bson"

type TestResult struct {
	TestID     string `json:"test_id" bson:"test_id"`
	Result     string `json:"result" bson:"result"`
	ResultText string `json:"result_text" bson:"result_text"`
}

type Record struct {
	ID          string       `json:"id" bson:"_id"`
	ComboName   string       `json:"combo_name" bson:"combo_name"`
	PatientID   string       `json:"patient_id" bson:"patient_id"`
	TestResults []TestResult `json:"test_results" bson:"test_results"`
	Status      string       `json:"status" bson:"status"`
}

func (r *Record) MarshalBSON() ([]byte, error) {
	if r.TestResults == nil {
		r.TestResults = []TestResult{}
	}

	type my Record
	return bson.Marshal((*my)(r))
}
