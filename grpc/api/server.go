package main

import (
	"context"
	pb "currency/proto/currency"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	httpServer *http.Server
	grpcClient pb.CurrencyServiceClient
}

type Currencies struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Converted struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

func NewServer(grpcClient pb.CurrencyServiceClient) *Server {
	mux := http.NewServeMux()
	s := &Server{
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
		grpcClient: grpcClient,
	}

	mux.HandleFunc("/currencies", s.GetCurrencies)
	mux.HandleFunc("/currencies/convert", s.Convert)

	return s
}

func (s *Server) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	currencyList, err := s.grpcClient.GetAllCurrencies(ctx, &pb.Empty{})

	if err != nil {
		http.Error(w, "Failed to fetch currencies", http.StatusInternalServerError)
		return
	}

	var currencies []Currencies
	for _, currency := range currencyList.Currencies {
		currencies = append(currencies, Currencies{
			Code: currency.Code,
			Name: currency.Name,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currencies)
}

func (s *Server) Convert(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	from := strings.ToLower(r.URL.Query().Get("from"))
	to := strings.ToLower(r.URL.Query().Get("to"))
	amount := r.URL.Query().Get("amount")

	req := &pb.ConvertRequest{
		Amount:       amount,
		FromCurrency: from,
		ToCurrency:   to,
	}

	converted, err := s.grpcClient.ConvertCurrency(ctx, req)

	if err != nil {
		http.Error(w, "Failed to convert", http.StatusInternalServerError)
		return
	}

	res := &Converted{
		Currency: converted.Currency,
		Amount:   converted.ConvertedAmount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
