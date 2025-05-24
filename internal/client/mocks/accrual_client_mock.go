package client_mocks

import (
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/timuraipov/diploma/internal/client"
)

// Client — структура API клиента
type accrualClient struct {
	baseURL string
	resty   *resty.Client
}

// NewClient — конструктор нового клиента
func NewMockClient(baseURL string) client.AccrualClient {
	return &accrualClient{
		baseURL: baseURL,
		resty:   resty.New(),
	}
}

// GetOrder — делает GET запрос на /api/orders/{number}
func (c *accrualClient) GetOrder(number string) (*client.OrderResponse, int, error) {
	order := client.OrderResponse{
		OrderNumber: number,
		Status:      "PROCESSED",
		Amount:      0.0,
	}
	return &order, http.StatusOK, nil
}
