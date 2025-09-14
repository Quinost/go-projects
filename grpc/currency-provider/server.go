package main

import (
	"context"
	pb "currency/proto/currency"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	currencyAPIURL = "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies.min.json"
	convertAPIURL  = "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/{code}.min.json"
)

type Server struct {
	pb.UnimplementedCurrencyServiceServer
	httpClient *http.Client
}

func (s *Server) ConvertCurrency(ctx context.Context, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
	url := strings.Replace(convertAPIURL, "{code}", req.FromCurrency, 1)

	response, err := s.httpClient.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch currency %v", err)
	}
	defer response.Body.Close()

	var root map[string]json.RawMessage
	json.NewDecoder(response.Body).Decode(&root)
	var rates map[string]decimal.Decimal
	json.Unmarshal(root[req.FromCurrency], &rates)

	decimal, _ := decimal.NewFromString(req.Amount)
	amountTo := decimal.Mul(rates[req.ToCurrency])

	return &pb.ConvertResponse{
		ConvertedAmount: amountTo.StringFixed(2),
		Currency:        req.ToCurrency,
	}, nil
}

func (s *Server) GetAllCurrencies(ctx context.Context, req *pb.Empty) (*pb.CurrencyList, error) {
	response, err := s.httpClient.Get(currencyAPIURL)

	if err != nil || response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch currencies %v", err)
	}
	defer response.Body.Close()

	var currenciesMap map[string]string
	if err := json.NewDecoder(response.Body).Decode(&currenciesMap); err != nil {
		return nil, fmt.Errorf("failed to decode response %v", err)
	}

	var currencies []*pb.Currency
	for code, name := range currenciesMap {
		currencies = append(currencies, &pb.Currency{
			Code: code,
			Name: name,
		})
	}

	return &pb.CurrencyList{
		Currencies: currencies,
	}, nil
}
