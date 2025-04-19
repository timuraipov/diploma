package client

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// Client — структура API клиента
type AccrualClient struct {
	baseURL string
	resty   *resty.Client
}

// NewClient — конструктор нового клиента
func NewClient(baseURL string) *AccrualClient {
	return &AccrualClient{
		baseURL: baseURL,
		resty:   resty.New(),
	}
}

type OrderResponse struct {
	OrderNumber string  `json:"order"`
	Status      string  `json:"status"`
	Amount      float64 `json:"amount"`
}

// GetOrder — делает GET запрос на /api/orders/{number}
func (c *AccrualClient) GetOrder(number string) (*OrderResponse, int, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, number)

	resp, err := c.resty.R().
		SetResult(&OrderResponse{}).
		Get(url)

	if err != nil {
		return nil, 0, err // ошибка сетевого уровня
	}

	statusCode := resp.StatusCode()

	if resp.IsError() {
		// тут можешь вернуть nil + статус и nil, если хочешь обработать ошибку выше
		return nil, statusCode, nil
	}

	order, ok := resp.Result().(*OrderResponse)
	if !ok {
		return nil, statusCode, fmt.Errorf("не удалось распарсить ответ")
	}

	return order, statusCode, nil
}
