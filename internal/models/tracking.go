package models

type TrackingTestRequest struct {
	TestID      string `json:"test_id" bson:"test_id"`
	TestName    string `json:"test_name" bson:"test_name"`
	NormalValue string `json:"normal_value" bson:"normal_value"`
	Unit        string `json:"unit" bson:"unit"`
	Order       int    `json:"order" bson:"order"`
}

type TrackingTestData struct {
	TestID      string `json:"test_id" bson:"test_id"`
	TestName    string `json:"test_name" bson:"test_name"`
	NormalValue string `json:"normal_value" bson:"normal_value"`
	Unit        string `json:"unit" bson:"unit"`
	Order       int    `json:"order" bson:"order"`
}

type Tracking struct {
	ID    string             `json:"id" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Tests []TrackingTestData `json:"tests" bson:"tests"`
}

type CreateTrackingRequest struct {
	Name  string                `json:"tracking_name"`
	Tests []TrackingTestRequest `json:"tests"`
}

type TrackingQueryOptions struct {
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

func NewTracking(name string, testRequests []TrackingTestRequest) Tracking {
	var tests []TrackingTestData
	for _, test := range testRequests {
		tests = append(tests, TrackingTestData{
			TestID:      test.TestID,
			TestName:    test.TestName,
			NormalValue: test.NormalValue,
			Unit:        test.Unit,
			Order:       test.Order, // Order starts from 1
		})
	}

	return Tracking{
		ID:    GenerateComboID(name), // Use tracking name for ID
		Name:  name,
		Tests: tests,
	}
}
