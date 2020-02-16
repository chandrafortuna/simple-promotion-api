package promotion

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/chandrafortuna/simple-promotion-api/utils"
)

// Service represent promotion service
type Service struct {
	repo Repository
}

// NewService represent promotion service constructor
func NewService(r Repository) Service {
	return Service{
		repo: r,
	}
}

// ApplyPromotion represent apply promotion of the service
func (s *Service) ApplyPromotion(req ApplyPromoRequest) (pr *ApplyPromoResponse, err error) {
	//getPromoByCode
	promoExists, err := s.repo.ExistsByCode(req.Code)
	if err != nil {
		return nil, errors.New("Failed to get existance of promo")
	}

	if !promoExists {
		return nil, errors.New("Promo is Not Exists")
	}

	promo, err := s.repo.GetPromotionByCode(req.Code)
	if err != nil {
		return nil, errors.New("Failed to get promo")
	}

	bookingTime := time.Now()
	var rooms []*RoomResponse
	totalPromo := float64(0)
	totalPrice := float64(0)
	for _, room := range req.Rooms {
		parsedDate, err := utils.ParseTimeFromString(room.Date)
		promoPrice := float64(0)
		message := ""
		if err != nil {
			return nil, errors.New("Invalid Date")
		}
		err = promo.ApplyRule(parsedDate, bookingTime, room.Night, room.Qty)
		if err != nil {
			message = fmt.Sprintf("Promo not applied: %v", err)
		} else {
			promoPrice, err = promo.CalculatePromo(room.Price)
			if err != nil {
				return nil, fmt.Errorf("Promo calculation failed: %v", err)
			}
		}

		if promoPrice == 0 {
			promoPrice = room.Price
		}

		res := &RoomResponse{
			Date:       parsedDate,
			Room:       room.Room,
			Price:      room.Price,
			Night:      room.Night,
			Qty:        room.Qty,
			PromoPrice: promoPrice,
			Saving:     room.Price - promoPrice,
			Message:    message,
		}

		totalPromo += (room.Price - promoPrice)
		totalPrice += room.Price
		rooms = append(rooms, res)
	}

	pr = &ApplyPromoResponse{
		Rooms:         rooms,
		PromoPrice:    totalPromo,
		FinalPrice:    totalPrice - totalPromo,
		OriginalPrice: totalPrice,
	}

	return pr, nil
}

// PromoDistribution represent promotion distribution of the promotion service
func (s *Service) PromoDistribution() error {
	promos, err := s.repo.GetAllAvailable()
	if err != nil {
		fmt.Errorf("Failed to get available promo: %v", err)
	}

	for _, promo := range promos {
		p := s.distribute(promo)
		err = s.repo.Update(p)
		if err != nil {
			return fmt.Errorf("Failed to update promo: %v", err)
		}
	}

	return nil
}

// CreatePromotion represent create promotion of the promotion service
func (s *Service) CreatePromotion(promotion *Promotion) (*Promotion, error) {

	codeIsExists, err := s.repo.ExistsByCode(promotion.Code)
	if err != nil {
		return nil, errors.New("Failed to get existance promo code")
	}

	if codeIsExists {
		return nil, errors.New("Duplicated Promo code")
	}

	err = s.repo.Save(s.distribute(promotion))
	if err != nil {
		return nil, errors.New("Failed to Save")
	}
	return promotion, nil
}

// GetAvailablePromo represent get ll available promotion
func (s *Service) GetAvailablePromo() ([]*Promotion, error) {
	promotions, err := s.repo.GetAllAvailable()
	if err != nil {
		return nil, errors.New("Failed to Save")
	}
	return promotions, nil
}

// GetAvailablePromo represent get ll available promotion
func (s *Service) distribute(p *Promotion) *Promotion {
	dayRange := int64(1)
	now := time.Now()

	if p.StartDate.Valid && p.EndDate.Valid {
		duration := p.EndDate.Time.Sub(p.StartDate.Time)
		if now.After(p.StartDate.Time) && now.Before(p.EndDate.Time) {
			duration = p.EndDate.Time.Sub(now)
		}
		dayRange = int64(math.Ceil(duration.Hours() / 24))
	}

	if dayRange < 1 {
		dayRange = 1
	}

	qtyPerDay := int64(1)
	if p.Balance > dayRange {
		qtyPerDay = int64(math.Round(float64(p.Balance) / float64(dayRange)))
	}

	if qtyPerDay > p.Balance {
		qtyPerDay = p.Balance
	}

	pd := &PromoDistribution{
		PromoID: p.ID,
		Qty:     qtyPerDay,
		Redeem:  0,
		Balance: qtyPerDay,
	}

	p.Distribution = pd

	return p
}
