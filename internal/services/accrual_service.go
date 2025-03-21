package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/grnsv/gophermart/internal/models"
)

type accrualService struct {
	address string
}

type AccrualResponse struct {
	Order   int           `json:"order,string"`
	Status  models.Status `json:"status"`
	Accrual float64       `json:"accrual"`
}

func NewAccrualService(address string) AccrualService {
	return &accrualService{address: address}
}

func (s *accrualService) GetAccrual(ctx context.Context, order *models.Order) (*models.Order, error) {
	url := fmt.Sprintf("%s/api/orders/%d", s.address, order.ID)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	client := http.DefaultClient

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled")

		case <-ticker.C:
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				return nil, fmt.Errorf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusOK:
				var accrualResp AccrualResponse
				if err := json.NewDecoder(resp.Body).Decode(&accrualResp); err != nil {
					return nil, fmt.Errorf("failed to decode response: %v", err)
				}

				switch accrualResp.Status {
				case models.StatusRegistered, models.StatusProcessing:
					continue
				case models.StatusInvalid, models.StatusProcessed:
					order.Status = accrualResp.Status
					order.Accrual = accrualResp.Accrual
					return order, nil
				default:
					return nil, fmt.Errorf("unknown accrual status: %s", accrualResp.Status)
				}

			case http.StatusNoContent:
				continue

			case http.StatusTooManyRequests:
				retryAfter := resp.Header.Get("Retry-After")
				if retryAfter == "" {
					retryAfter = "60"
				}

				duration, err := time.ParseDuration(retryAfter + "s")
				if err != nil {
					duration = 60 * time.Second
				}

				select {
				case <-ctx.Done():
					return nil, fmt.Errorf("context canceled")
				case <-time.After(duration):
					continue
				}

			case http.StatusInternalServerError:
				body, _ := io.ReadAll(resp.Body)
				return nil, fmt.Errorf("accrual system internal error: %s", body)

			default:
				return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			}
		}
	}
}
