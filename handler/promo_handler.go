package handler

import (
	"encoding/json"
	"log"
	"net/http"

	domainPromo "github.com/chandrafortuna/simple-promotion-api/domain/promotion"
	uuid "github.com/satori/go.uuid"
)

type Handler struct {
	service domainPromo.Service
}

func NewHandler(s domainPromo.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) CreatePromo(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req domainPromo.PromoRequest
	if err := decoder.Decode(&req); err != nil {
		log.Println("Error:", err)
		Error(w, http.StatusBadRequest, err, err.Error())
		return
	}

	err := req.Validate()
	if err != nil {
		Error(w, http.StatusBadRequest, err, err.Error())
		return
	}

	uid, err := uuid.NewV4()
	if err != nil {
		Error(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	promo, err := req.ToPromo(uid)
	promotion, err := h.service.CreatePromotion(promo)
	if err != nil {
		Error(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	JSON(w, http.StatusCreated, promotion)
}

func (h *Handler) ApplyPromo(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req domainPromo.ApplyPromoRequest
	if err := decoder.Decode(&req); err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	res, err := h.service.ApplyPromotion(req)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, res)
}

func (h *Handler) PromoDistribution(w http.ResponseWriter, r *http.Request) {
	err := h.service.PromoDistribution()
	if err != nil {
		Error(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, true)
}

func (h *Handler) GetAvailablePromo(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetAvailablePromo()
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, res)
}
