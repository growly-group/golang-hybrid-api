package calculatorsvc

type CalculationRequest struct {
	A         float64 `json:"a"`
	B         float64 `json:"b"`
	Operation string  `json:"operation"`
}
