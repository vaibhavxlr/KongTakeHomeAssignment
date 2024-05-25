package internal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vaibhavxlr/KongTakeHomeAssignment/internal/DTOs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockMongoColl struct{}

type MockCursor struct {
	Data []DTOs.Service
	Ind  int
}

func (m *MockCursor) Decode(val interface{}) error {
	if m.Ind < len(m.Data) {
		*val.(*DTOs.Service) = m.Data[m.Ind]
		m.Ind++
	}
	return nil
}

func (*MockCursor) Close(ctx context.Context) error {
	return nil
}

func (m *MockCursor) Next(ctx context.Context) bool {
	return m.Ind < len(m.Data)
}

func (*MockMongoColl) Find(ctx context.Context, filter interface{}, dbOptions ...*options.FindOptions) (Cursor, error) {
	services := []DTOs.Service{{
		ID:            "d3b0738a-d9f1-11eb-ba80-0242ac130004",
		Name:          "Konnect API 7",
		Info:          "lorem ipsum lorem ipsum",
		VersionsCount: 4,
	},
		{
			ID:            "d3b0738b-d9f1-11eb-ba80-0242ac130004",
			Name:          "Konnect API 8",
			Info:          "lorem ipsum lorem ipsum",
			VersionsCount: 3,
		},
		{
			ID:            "d3b0738c-d9f1-11eb-ba80-0242ac130004",
			Name:          "Konnect API 9",
			Info:          "lorem ipsum lorem ipsum",
			VersionsCount: 2,
		},
		{
			ID:            "d3b0738d-d9f1-11eb-ba80-0242ac130004",
			Name:          "Konnect API 10",
			Info:          "lorem ipsum lorem ipsum",
			VersionsCount: 1,
		},
		{
			ID:            "d3b0738e-d9f1-11eb-ba80-0242ac130004",
			Name:          "Konnect API 11",
			Info:          "lorem ipsum lorem ipsum",
			VersionsCount: 5,
		},
		{
			ID:            "d3b0738f-d9f1-11eb-ba80-0242ac130004",
			Name:          "Konnect API 12",
			Info:          "lorem ipsum lorem ipsum",
			VersionsCount: 2,
		},
	}
	return &MockCursor{Data: services}, nil
}

func TestListServices(t *testing.T) {
	mockColl := &MockMongoColl{}
	r, err := http.NewRequest("GET", "/services?curr=1&count=2&sortOrder=0", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	listServices(w, r, mockColl)
	expected := `{"services":[{"id":"d3b0738a-d9f1-11eb-ba80-0242ac130004","name":"Konnect API 7","info":"lorem ipsum lorem ipsum","versionsCount":4},{"id":"d3b0738b-d9f1-11eb-ba80-0242ac130004","name":"Konnect API 8","info":"lorem ipsum lorem ipsum","versionsCount":3}],"sortOrder":{"A-Z":0,"Z-A":1},"pageDetails":{"curr":1,"total":3,"count":2}}`
	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
	}
}
