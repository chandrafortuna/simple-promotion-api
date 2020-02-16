package promotion

import (
	"errors"
	"log"
	"time"

	"github.com/chandrafortuna/simple-promotion-api/utils"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v3"
)

// Promotion represent entity of the promo
type Promotion struct {
	ID               uuid.UUID          `db:"id" json:"id"`
	Title            string             `db:"title" json:"title"`
	Code             string             `db:"code" json:"code"`
	StartDate        null.Time          `db:"start_date" json:"startDate"`
	EndDate          null.Time          `db:"end_date" json:"endDate"`
	Percentage       null.Int           `db:"percentage" json:"percentage"`
	Amount           null.Float         `db:"amount" json:"amount"`
	Qty              int64              `db:"qty" json:"qty"`
	Redeem           int64              `db:"reedem" json:"redeem"`
	Balance          int64              `db:"balance" json:"balance"`
	Status           int64              `db:"status" json:"status"`
	MinNight         null.Int           `db:"min_night" json:"minNight"`
	MinRoom          null.Int           `db:"min_room" json:"minRoom"`
	CheckinDays      null.String        `db:"checkin_day" json:"checkinDays"`
	BookingDays      null.String        `db:"booking_day" json:"bookingDays"`
	BookingHourStart null.Int           `db:"booking_hour_start" json:"bookingHourStart"`
	BookingHourEnd   null.Int           `db:"booking_hour_end" json:"bookingHourEnd"`
	Distribution     *PromoDistribution `db:"distribution" json:"distribution"`
}

func (p *Promotion) CalculatePromo(price float64) (float64, error) {
	if p.Percentage.Valid {
		return p.calculatePromoPercentage(price)
	}

	if p.Amount.Valid {
		return p.calculatePromoAmount(price)
	}

	return 0, errors.New("Invalid Promotion Percentage/Amount")

}

func (p *Promotion) calculatePromoPercentage(price float64) (float64, error) {
	promoPercentage := p.Percentage.Int64
	discount := (price * float64(promoPercentage) / float64(100))
	newPrice := price - discount
	log.Println("promoPercentage:", promoPercentage)
	log.Println("discount:", discount)
	log.Println("calculatePromoPercentage:", newPrice)
	return newPrice, nil
}

func (p *Promotion) calculatePromoAmount(price float64) (float64, error) {
	promoAmount := p.Amount.Float64
	newPrice := price - promoAmount
	log.Println("calculatePromoPercentage:", newPrice)
	return newPrice, nil
}

func (p *Promotion) hasBalance() bool {
	log.Println("p.Balance", p.Balance)
	if p.Distribution == nil {
		return p.Balance > 0
	}

	log.Println("p.Distribution.Balance", p.Distribution.Balance)
	return p.Distribution.Balance > 0
}

func (p *Promotion) dateRangeRule(t time.Time) error {
	if !p.StartDate.Valid || !p.EndDate.Valid {
		return nil
	}

	log.Println("BookingTime:", t)
	if t.After(p.StartDate.Time) && t.Before(p.EndDate.Time) {
		return nil
	}

	return errors.New("Promo is not started or has been ended")
}

func (p *Promotion) checkinRule(t time.Time) error {
	if !p.CheckinDays.Valid {
		return nil
	}

	checkinDayRule, err := utils.ParseWeekday(p.CheckinDays.String)
	if err != nil {
		return errors.New("Promo Checkin Day is Invalid")
	}

	if checkinDayRule != t.Weekday() {
		return errors.New("checkin time rule is failed")
	}

	return nil
}

func (p *Promotion) bookingDayRule(t time.Time) error {
	if !p.BookingDays.Valid {
		return nil
	}

	bookingDayRule, err := utils.ParseWeekday(p.BookingDays.String)
	if err != nil {
		return errors.New("Promo Booking Day is Invalid")
	}

	if bookingDayRule != t.Weekday() {
		return errors.New("booking time rule is failed")
	}

	return nil
}

func (p *Promotion) bookingHourRule(t time.Time) error {
	if !p.BookingHourStart.Valid || !p.BookingHourEnd.Valid {
		return nil
	}

	bookingHour := int64(t.Hour())
	if bookingHour >= p.BookingHourStart.Int64 && bookingHour <= p.BookingHourEnd.Int64 {
		return nil
	}

	return errors.New("booking hour rule is failed")
}

func (p *Promotion) minNightRule(n null.Int) error {
	if !p.MinNight.Valid {
		return nil
	}

	if !n.Valid {
		return errors.New("This promo apply min night")
	}

	if n.Int64 < p.MinNight.Int64 {
		return errors.New("Min Night rule is failed")
	}

	return nil
}

func (p *Promotion) minRoomRule(n null.Int) error {
	log.Println("p.MinRoom.Valid:", p.MinRoom.Valid)
	if !p.MinRoom.Valid {
		return nil
	}

	if !n.Valid {
		return errors.New("This promo apply min room")
	}

	log.Println("p.MinRoom.Int64:", p.MinRoom.Int64)
	if n.Int64 < p.MinRoom.Int64 {
		return errors.New("Min Room rule is failed")
	}

	return nil
}

func (p *Promotion) ApplyRule(checkinTime time.Time, bookingTime time.Time, night null.Int, room null.Int) error {

	if !p.hasBalance() {
		return errors.New("Promo is not available")
	}

	if err := p.dateRangeRule(bookingTime); err != nil {
		return err
	}

	if err := p.minNightRule(night); err != nil {
		return err
	}

	if err := p.minRoomRule(room); err != nil {
		return err
	}

	if err := p.checkinRule(checkinTime); err != nil {
		return err
	}

	if err := p.bookingDayRule(bookingTime); err != nil {
		return err
	}

	if err := p.bookingHourRule(bookingTime); err != nil {
		return err
	}

	return nil
}

// PromoDistribution represent entity of the promo
type PromoDistribution struct {
	PromoID uuid.UUID `db:"promo_id" json:"promoId"`
	Qty     int64     `db:"qty" json:"qty"`
	Redeem  int64     `db:"reedem" json:"redeem"`
	Balance int64     `db:"balance" json:"balance"`
}

// PromoRequest represent entity of the Promotion Request
type PromoRequest struct {
	Title            string      `json:"title"`
	Code             string      `json:"code"`
	StartDate        null.String `json:"startDate"`
	EndDate          null.String `json:"endDate"`
	Percentage       null.Int    `json:"percentage"`
	Amount           null.Float  `json:"amount"`
	Quota            int64       `json:"quota"`
	MinNight         null.Int    `json:"minNight"`
	MinRoom          null.Int    `json:"minRoom"`
	CheckinDays      null.String `json:"checkinDays"`
	BookingDays      null.String `json:"bookingDays"`
	BookingHourStart null.Int    `json:"bookingHourStart"`
	BookingHourEnd   null.Int    `json:"bookingHourEnd"`
}

func (req *PromoRequest) ToPromo(id uuid.UUID) (*Promotion, error) {
	var startDate, endDate null.Time
	if req.StartDate.Valid && req.EndDate.Valid {
		_startDate, err := utils.ParseTimeFromString(utils.NullStringToString(req.StartDate))
		if err != nil {
			return nil, err
		}
		startDate = null.TimeFrom(_startDate)

		_endDate, err := utils.ParseTimeFromString(utils.NullStringToString(req.EndDate))
		if err != nil {
			return nil, err
		}
		endOfDay := time.Date(_endDate.Year(), _endDate.Month(), _endDate.Day(), 23, 59, 59, 0, _endDate.Location())
		endDate = null.TimeFrom(endOfDay)

		if _startDate.After(_endDate) {
			err = errors.New("End Date must greather than Start Date")
			return nil, err
		}
	}

	promo := &Promotion{
		ID:               id,
		Title:            req.Title,
		Code:             req.Code,
		StartDate:        startDate,
		EndDate:          endDate,
		Percentage:       req.Percentage,
		Amount:           req.Amount,
		Qty:              req.Quota,
		Redeem:           0,
		Balance:          req.Quota,
		Status:           1,
		MinNight:         req.MinNight,
		MinRoom:          req.MinRoom,
		CheckinDays:      req.CheckinDays,
		BookingDays:      req.BookingDays,
		BookingHourStart: req.BookingHourStart,
		BookingHourEnd:   req.BookingHourEnd,
		Distribution:     nil,
	}

	return promo, nil
}

func (promoReq *PromoRequest) Validate() error {
	if !promoReq.Percentage.Valid && !promoReq.Amount.Valid {
		return errors.New("Either Percentage or Amount must be filled")
	}

	if promoReq.Percentage.Valid && promoReq.Amount.Valid {
		return errors.New("You have to fill only one either Percentage or Amount")
	}

	if promoReq.StartDate.Valid || promoReq.EndDate.Valid {
		if promoReq.StartDate.String == "" || promoReq.EndDate.String == "" {
			return errors.New("Start Date and End Date range required")
		}
	}

	return nil
}

// RoomRequest represent entity of the Room Request params
type RoomRequest struct {
	Date  string   `json:"date"`
	Room  string   `json:"room"`
	Price float64  `json:"price"`
	Night null.Int `json:"night"`
	Qty   null.Int `json:"qty"`
}

// ApplyPromoRequest represent entity of the PromoRequest params
type ApplyPromoRequest struct {
	Rooms      []*RoomRequest `json:"rooms"`
	TotalPrice float64        `json:"totalPrice"`
	Code       string         `json:"code"`
}

// RoomResponse represent entity of the Room Response
type RoomResponse struct {
	Date       time.Time `json:"date"`
	Room       string    `json:"room"`
	Price      float64   `json:"price"`
	Night      null.Int  `json:"night"`
	Qty        null.Int  `json:"qty"`
	PromoPrice float64   `json:"promoPrice"`
	Saving     float64   `json:"saving"`
	Message    string    `json:"message"`
}

// ApplyPromoResponse represent entity of the Promo response
type ApplyPromoResponse struct {
	Rooms         []*RoomResponse `json:"rooms"`
	PromoPrice    float64         `json:"promoPrice"`
	FinalPrice    float64         `json:"finalPrice"`
	OriginalPrice float64         `json:"originalPrice"`
}
