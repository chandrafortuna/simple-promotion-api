package main

import (
	"log"
	"net/http"

	p "github.com/chandrafortuna/simple-promotion-api/domain/promotion"
	h "github.com/chandrafortuna/simple-promotion-api/handler"
	"github.com/gorilla/mux"
)

func main() {
	//register service
	promoRepository := p.NewRepository([]*p.Promotion{})
	promoService := p.NewService(promoRepository)
	handler := h.NewHandler(promoService)

	router := mux.NewRouter()
	router.HandleFunc("/promo", handler.CreatePromo).Methods("POST")
	router.HandleFunc("/promo", handler.GetAvailablePromo).Methods("GET")
	router.HandleFunc("/promo/apply", handler.ApplyPromo).Methods("POST")
	router.HandleFunc("/promo/distribute", handler.PromoDistribution).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
