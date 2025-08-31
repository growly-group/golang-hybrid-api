package calculatorsvc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type CalculatorSdk struct {
	Calculate func(req CalculationRequest) (float64, error)
}

func NewCalculatorSdk(mode string) *CalculatorSdk {
	if mode == "http" {
		return newHttpSdk()
	}

	return &CalculatorSdk{
		Calculate: calculate,
	}
}

func newHttpSdk() *CalculatorSdk {
	return &CalculatorSdk{
		Calculate: func(req CalculationRequest) (float64, error) {
			url := os.Getenv("CALCULADOR_SERVICE_URL")
			if url == "" {
				return 0, fmt.Errorf("CALCULADOR_SERVICE_URL environment variable not set")
			}

			body, err := json.Marshal(req)
			if err != nil {
				return 0, fmt.Errorf("failed to marshal request: %w", err)
			}

			resp, err := http.Post(url+"/calculator", "application/json", bytes.NewBuffer(body))
			if err != nil {
				return 0, fmt.Errorf("failed to call calculator service: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return 0, fmt.Errorf("calculator service returned non-200 status: %d", resp.StatusCode)
			}

			var result struct {
				Data float64 `json:"data"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return 0, fmt.Errorf("failed to decode response: %w", err)
			}

			return result.Data, nil
		},
	}
}
