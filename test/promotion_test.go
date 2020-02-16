package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	p "github.com/chandrafortuna/simple-promotion-api/domain/promotion"
	h "github.com/chandrafortuna/simple-promotion-api/handler"
	"github.com/chandrafortuna/simple-promotion-api/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

var (
	id, _    = uuid.NewV4()
	start, _ = utils.ParseTimeFromString("2020-02-15 00:00:00")
	end, _   = utils.ParseTimeFromString("2020-02-22 23:00:00")
	promo    = &p.Promotion{
		ID:          id,
		Balance:     10,
		Code:        "PROMOTEST123",
		Title:       "Test Promo",
		Percentage:  null.NewInt(10, true),
		Qty:         20,
		Redeem:      2,
		StartDate:   null.NewTime(start, true),
		EndDate:     null.NewTime(end, true),
		Status:      int64(1),
		MinNight:    null.NewInt(2, true),
		MinRoom:     null.NewInt(2, true),
		CheckinDays: null.NewString("Sunday", true),
	}

	rooms1 = &p.RoomRequest{
		Date:  "2020-02-16 10:00:00",
		Night: null.NewInt(2, true),
		Price: float64(150000),
		Qty:   null.NewInt(2, true),
		Room:  "Room1",
	}

	rooms2 = &p.RoomRequest{
		Date:  "2020-02-25 10:00:00",
		Night: null.NewInt(1, true),
		Price: float64(150000),
		Qty:   null.NewInt(2, true),
		Room:  "Room2",
	}

	rooms3 = &p.RoomRequest{
		Date:  "2020-02-17 10:00:00",
		Night: null.NewInt(2, true),
		Price: float64(150000),
		Qty:   null.NewInt(2, true),
		Room:  "Room3",
	}

	rooms4 = &p.RoomRequest{
		Date:  "2020-02-16 10:00:00",
		Night: null.NewInt(2, true),
		Price: float64(150000),
		Qty:   null.NewInt(1, true),
		Room:  "Room4",
	}

	applyPromoRequest = p.ApplyPromoRequest{
		Rooms:      []*p.RoomRequest{rooms1, rooms2},
		Code:       "PROMOTEST123",
		TotalPrice: float64(150000),
	}

	repo          = p.NewRepository([]*p.Promotion{promo})
	service       = p.NewService(repo)
	handler       = h.NewHandler(service)
	promotionBody = `
	{
		"title": "test promo",
		"code": "PROMOTEST01",
		"quota": 3,
		"percentage": 10,
		"amount": "",
		"startDate": "2020-02-16 00:04:05",
		"endDate": "2020-02-16 23:59:05",
		"minRoom": 2,
		"minNight": 1,
		"checkinDays": "Sunday",
		"bookingDays": "null",
		"bookingHourStart": null,
		"bookingHourEnd": null
	}`

	applyBody = `
	{
		"rooms":[
			{
				"date":"2020-02-15 12:00:00",
				"room": "Room A - Standard",
				"price": 120000,
				"night": null,
				"qty": 3
			},
			{
				"date":"2020-02-15 12:00:00",
				"room": "Room A - Standard",
				"price": 150000,
				"night": 2,
				"qty": null
			},
			{
				"date":"2020-02-15 12:00:00",
				"room": "Room A - Standard",
				"price": 180000,
				"night": 1,
				"qty": 2
			},
			{
				"date":"2020-02-16 12:00:00",
				"room": "Room A - Standard",
				"price": 220000,
				"night": 1,
				"qty": 2
			},
			{
				"date":"2020-02-17 12:00:00",
				"room": "Room A - Standard",
				"price": 220000,
				"night": 1,
				"qty": 2
			}],
		"totalPrice":2000000,
		"code":"PROMOTEST123"
	}`
)

func TestHTTPCreatePromo(t *testing.T) {
	reqBody := strings.NewReader(promotionBody)
	req, err := http.NewRequest("POST", "/promo", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.CreatePromo)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

func TestHTTPApplyPromo(t *testing.T) {
	reqBody := strings.NewReader(applyBody)
	req, err := http.NewRequest("POST", "/promo/apply", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.ApplyPromo)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestFunctionApplyPromo(t *testing.T) {
	res, err := service.ApplyPromotion(applyPromoRequest)

	fmt.Sprintf("res: %v", res)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, len(res.Rooms), len(applyPromoRequest.Rooms))

	for _, room := range res.Rooms {
		if room.Room == "Room1" {
			assert.Equal(t, room.Price, float64(150000))
			assert.Equal(t, room.PromoPrice, (float64(150000) - float64(15000)))
			assert.Equal(t, room.Message, "")
		}

		if room.Room == "Room2" {
			assert.Equal(t, room.Price, float64(150000))
			assert.Equal(t, room.PromoPrice, float64(150000))
			assert.Contains(t, room.Message, "Min Night rule is failed")
		}

		if room.Room == "Room3" {
			assert.Equal(t, room.Price, float64(150000))
			assert.Equal(t, room.PromoPrice, float64(150000))
			assert.Contains(t, room.Message, "checkin time rule is failed")
		}
		if room.Room == "Room4" {
			assert.Equal(t, room.Price, float64(150000))
			assert.Equal(t, room.PromoPrice, float64(150000))
			assert.Equal(t, room.Message, "This promo apply min night")
		}
	}
}
