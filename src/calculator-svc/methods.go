package calculatorsvc

import "fmt"

func calculate(req CalculationRequest) (float64, error) {
	switch req.Operation {
	case "add":
		return req.A + req.B, nil
	case "subtract":
		return req.A - req.B, nil
	case "multiply":
		return req.A * req.B, nil
	case "divide":
		if req.B == 0 {
			return 0, fmt.Errorf("division by zero is not allowed")
		}
		return req.A / req.B, nil
	default:
		return 0, fmt.Errorf("invalid operation specified")
	}
}
