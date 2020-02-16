package promotion

import (
	"errors"
	"log"

	uuid "github.com/satori/go.uuid"
)

// Repository represent interface of promotion repository
type Repository interface {
	GetPromotionByCode(code string) (*Promotion, error)
	GetPromotionByID(id uuid.UUID) (*Promotion, error)
	ExistsByCode(code string) (bool, error)
	GetAllAvailable() ([]*Promotion, error)
	Save(*Promotion) error
	Update(*Promotion) error
}

var ErrPromoNotFound = errors.New("Promo Not Found")
var ErrPromoDistributionNotFound = errors.New("Promo Distribution Not Found")

// TempRepository represent temporary repository of promo
type TempRepository struct {
	promoCollection []*Promotion
}

// GetPromotionByCode represent get promotion by code
func (r *TempRepository) GetPromotionByCode(code string) (*Promotion, error) {
	for _, promo := range r.promoCollection {
		if promo.Code == code {
			return promo, nil
		}
	}
	return nil, ErrPromoNotFound
}

// GetPromotionByID represent get promotion by ID
func (r *TempRepository) GetPromotionByID(id uuid.UUID) (*Promotion, error) {
	for _, promo := range r.promoCollection {
		if promo.ID == id {
			return promo, nil
		}
	}
	return nil, ErrPromoNotFound
}

// ExistsByCode represent get existance promotion code
func (r *TempRepository) ExistsByCode(code string) (bool, error) {
	for _, promo := range r.promoCollection {
		log.Println("promo.Code", promo.Code)
		log.Println("code", code)
		log.Println("promo.Code == code", promo.Code == code)
		if promo.Code == code {
			return true, nil
		}
	}
	return false, nil
}

// GetAllAvailable represent get all promotion
func (r *TempRepository) GetAllAvailable() ([]*Promotion, error) {
	res := []*Promotion{}
	for _, promo := range r.promoCollection {
		if promo.Status == int64(1) {
			res = append(res, promo)
		}
	}
	return res, nil
}

// Save represent save promotion repository
func (r *TempRepository) Save(p *Promotion) error {
	r.promoCollection = append(r.promoCollection, p)
	return nil
}

func (r *TempRepository) Update(p *Promotion) error {
	for _, promo := range r.promoCollection {
		if promo.ID == p.ID {
			promo = p
			return nil
		}
	}
	return ErrPromoNotFound
}

// NewRepository initiate Repository
func NewRepository(p []*Promotion) (r Repository) {
	r = &TempRepository{p}
	return
}
