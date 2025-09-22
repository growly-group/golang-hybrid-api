# Temperature Service (temperature-svc) Specification

## Overview
The Temperature Service provides comprehensive temperature conversion utilities, supporting multiple temperature scales, unit conversions, and temperature-related calculations. It integrates seamlessly with the hybrid API ecosystem supporting both in-process and HTTP-based communication patterns.

## Port Assignment
- **Default Port**: 8084
- **Health Check**: 8084/health

## Core Features

### 1. Temperature Conversions
- Convert between Celsius, Fahrenheit, Kelvin, Rankine, Réaumur
- Bulk conversion operations
- Precision control for conversions
- Formula display and calculations

### 2. Temperature Calculations
- Temperature difference calculations
- Average temperature calculations
- Temperature range validations
- Heat index and wind chill calculations

### 3. Unit Utilities
- Supported unit information
- Conversion formula display
- Temperature scale comparisons
- Historical temperature data validation

### 4. Weather Integrations
- Heat index calculations
- Dew point calculations
- Apparent temperature
- Thermal comfort indices

## API Endpoints

### Basic Conversions

#### POST /temperature/convert
Convert temperature between different scales.

**Request Body:**
```json
{
  "value": 25.5,
  "fromUnit": "celsius",
  "toUnit": "fahrenheit",
  "precision": 2
}
```

**Response (200):**
```json
{
  "success": true,
  "input": {
    "value": 25.5,
    "unit": "celsius",
    "symbol": "°C"
  },
  "output": {
    "value": 77.9,
    "unit": "fahrenheit",
    "symbol": "°F"
  },
  "formula": "°F = (°C × 9/5) + 32",
  "calculation": "°F = (25.5 × 9/5) + 32 = 77.9",
  "precision": 2
}
```

#### POST /temperature/convert/bulk
Convert multiple temperature values.

**Request Body:**
```json
{
  "conversions": [
    {
      "value": 0,
      "fromUnit": "celsius",
      "toUnit": "fahrenheit"
    },
    {
      "value": 100,
      "fromUnit": "celsius",
      "toUnit": "kelvin"
    },
    {
      "value": 32,
      "fromUnit": "fahrenheit",
      "toUnit": "celsius"
    }
  ],
  "precision": 2
}
```

**Response (200):**
```json
{
  "success": true,
  "results": [
    {
      "input": { "value": 0, "unit": "celsius", "symbol": "°C" },
      "output": { "value": 32.0, "unit": "fahrenheit", "symbol": "°F" }
    },
    {
      "input": { "value": 100, "unit": "celsius", "symbol": "°C" },
      "output": { "value": 373.15, "unit": "kelvin", "symbol": "K" }
    },
    {
      "input": { "value": 32, "unit": "fahrenheit", "symbol": "°F" },
      "output": { "value": 0.0, "unit": "celsius", "symbol": "°C" }
    }
  ],
  "totalConversions": 3,
  "processingTime": "2ms"
}
```

### Temperature Calculations

#### POST /temperature/calculations/difference
Calculate temperature difference between two values.

**Request Body:**
```json
{
  "temperature1": {
    "value": 25,
    "unit": "celsius"
  },
  "temperature2": {
    "value": 77,
    "unit": "fahrenheit"
  },
  "outputUnit": "celsius"
}
```

**Response (200):**
```json
{
  "success": true,
  "temperature1": { "value": 25, "unit": "celsius", "symbol": "°C" },
  "temperature2": { "value": 25, "unit": "celsius", "symbol": "°C" },
  "difference": {
    "value": 0,
    "unit": "celsius",
    "symbol": "°C"
  },
  "calculation": "25°C - 25°C = 0°C"
}
```

#### POST /temperature/calculations/average
Calculate average temperature from multiple readings.

**Request Body:**
```json
{
  "temperatures": [
    { "value": 20, "unit": "celsius" },
    { "value": 68, "unit": "fahrenheit" },
    { "value": 295, "unit": "kelvin" }
  ],
  "outputUnit": "celsius",
  "precision": 2
}
```

**Response (200):**
```json
{
  "success": true,
  "temperatures": [
    { "value": 20, "unit": "celsius", "normalized": 20 },
    { "value": 20, "unit": "celsius", "normalized": 20 },
    { "value": 21.85, "unit": "celsius", "normalized": 21.85 }
  ],
  "average": {
    "value": 20.62,
    "unit": "celsius",
    "symbol": "°C"
  },
  "count": 3,
  "range": {
    "min": { "value": 20, "unit": "celsius" },
    "max": { "value": 21.85, "unit": "celsius" }
  }
}
```

#### POST /temperature/calculations/range
Validate temperature range and provide statistics.

**Request Body:**
```json
{
  "temperatures": [
    { "value": 18, "unit": "celsius" },
    { "value": 22, "unit": "celsius" },
    { "value": 25, "unit": "celsius" },
    { "value": 19, "unit": "celsius" }
  ],
  "unit": "celsius"
}
```

**Response (200):**
```json
{
  "success": true,
  "statistics": {
    "count": 4,
    "min": { "value": 18, "unit": "celsius" },
    "max": { "value": 25, "unit": "celsius" },
    "average": { "value": 21, "unit": "celsius" },
    "median": { "value": 20.5, "unit": "celsius" },
    "standardDeviation": { "value": 2.94, "unit": "celsius" },
    "range": { "value": 7, "unit": "celsius" }
  },
  "distribution": {
    "cold": 1,    // < 20°C
    "moderate": 2, // 20-24°C
    "warm": 1     // > 24°C
  }
}
```

### Weather-Related Calculations

#### POST /temperature/weather/heat-index
Calculate heat index (apparent temperature).

**Request Body:**
```json
{
  "temperature": {
    "value": 85,
    "unit": "fahrenheit"
  },
  "humidity": 85,
  "outputUnit": "fahrenheit"
}
```

**Response (200):**
```json
{
  "success": true,
  "temperature": { "value": 85, "unit": "fahrenheit" },
  "humidity": 85,
  "heatIndex": {
    "value": 108,
    "unit": "fahrenheit",
    "symbol": "°F"
  },
  "category": "extreme_caution",
  "description": "Heat exhaustion and heat cramps are possible with prolonged exposure",
  "formula": "HI = c1 + c2*T + c3*RH + c4*T*RH + c5*T² + c6*RH² + c7*T²*RH + c8*T*RH² + c9*T²*RH²"
}
```

#### POST /temperature/weather/wind-chill
Calculate wind chill factor.

**Request Body:**
```json
{
  "temperature": {
    "value": 10,
    "unit": "fahrenheit"
  },
  "windSpeed": 15,
  "windUnit": "mph",
  "outputUnit": "fahrenheit"
}
```

**Response (200):**
```json
{
  "success": true,
  "temperature": { "value": 10, "unit": "fahrenheit" },
  "windSpeed": { "value": 15, "unit": "mph" },
  "windChill": {
    "value": -7,
    "unit": "fahrenheit",
    "symbol": "°F"
  },
  "category": "cold",
  "description": "Frostbite possible in 30 minutes",
  "formula": "WC = 35.74 + 0.6215*T - 35.75*V^0.16 + 0.4275*T*V^0.16"
}
```

#### POST /temperature/weather/dew-point
Calculate dew point temperature.

**Request Body:**
```json
{
  "temperature": {
    "value": 25,
    "unit": "celsius"
  },
  "humidity": 60,
  "outputUnit": "celsius"
}
```

**Response (200):**
```json
{
  "success": true,
  "temperature": { "value": 25, "unit": "celsius" },
  "humidity": 60,
  "dewPoint": {
    "value": 16.7,
    "unit": "celsius",
    "symbol": "°C"
  },
  "comfort": "comfortable",
  "description": "Pleasant humidity level"
}
```

### Unit Information

#### GET /temperature/units
List all supported temperature units.

**Response (200):**
```json
{
  "units": [
    {
      "name": "celsius",
      "symbol": "°C",
      "fullName": "Celsius",
      "description": "Metric temperature scale where water freezes at 0° and boils at 100°",
      "absoluteZero": -273.15,
      "commonRange": { "min": -40, "max": 50 },
      "scientificUse": true
    },
    {
      "name": "fahrenheit",
      "symbol": "°F",
      "fullName": "Fahrenheit",
      "description": "Imperial temperature scale where water freezes at 32° and boils at 212°",
      "absoluteZero": -459.67,
      "commonRange": { "min": -40, "max": 120 },
      "scientificUse": false
    },
    {
      "name": "kelvin",
      "symbol": "K",
      "fullName": "Kelvin",
      "description": "Absolute temperature scale used in scientific applications",
      "absoluteZero": 0,
      "commonRange": { "min": 233.15, "max": 323.15 },
      "scientificUse": true
    }
  ]
}
```

#### GET /temperature/units/:unit/info
Get detailed information about a specific temperature unit.

#### GET /temperature/conversions/formulas
Get conversion formulas between temperature units.

**Response (200):**
```json
{
  "formulas": [
    {
      "from": "celsius",
      "to": "fahrenheit",
      "formula": "°F = (°C × 9/5) + 32",
      "inverse": "°C = (°F - 32) × 5/9"
    },
    {
      "from": "celsius",
      "to": "kelvin",
      "formula": "K = °C + 273.15",
      "inverse": "°C = K - 273.15"
    },
    {
      "from": "fahrenheit",
      "to": "kelvin",
      "formula": "K = (°F - 32) × 5/9 + 273.15",
      "inverse": "°F = (K - 273.15) × 9/5 + 32"
    }
  ]
}
```

### Validation and Utilities

#### POST /temperature/validate
Validate temperature values and ranges.

**Request Body:**
```json
{
  "temperature": {
    "value": -300,
    "unit": "celsius"
  },
  "context": "weather" // "weather", "scientific", "cooking", "medical"
}
```

**Response (200):**
```json
{
  "valid": false,
  "temperature": { "value": -300, "unit": "celsius" },
  "errors": [
    {
      "code": "BELOW_ABSOLUTE_ZERO",
      "message": "Temperature is below absolute zero (-273.15°C)",
      "severity": "error"
    }
  ],
  "warnings": [],
  "suggestions": [
    "Check if the value should be in a different unit",
    "Verify the temperature reading is correct"
  ],
  "context": {
    "name": "weather",
    "validRange": { "min": -89, "max": 58 },
    "unit": "celsius"
  }
}
```

#### POST /temperature/compare
Compare temperatures across different scales.

**Request Body:**
```json
{
  "temperature": {
    "value": 100,
    "unit": "celsius"
  },
  "compareUnits": ["fahrenheit", "kelvin", "rankine"]
}
```

**Response (200):**
```json
{
  "original": { "value": 100, "unit": "celsius", "symbol": "°C" },
  "comparisons": [
    { "value": 212, "unit": "fahrenheit", "symbol": "°F" },
    { "value": 373.15, "unit": "kelvin", "symbol": "K" },
    { "value": 671.67, "unit": "rankine", "symbol": "°R" }
  ],
  "significantPoints": {
    "waterBoiling": true,
    "waterFreezing": false,
    "absoluteZero": false,
    "bodyTemperature": false
  }
}
```

## Data Models

### Temperature Model
```go
type Temperature struct {
    Value     float64 `json:"value"`
    Unit      string  `json:"unit"`
    Symbol    string  `json:"symbol"`
    Precision int     `json:"precision"`
}

type TemperatureRange struct {
    Min Temperature `json:"min"`
    Max Temperature `json:"max"`
}
```

### Conversion Request/Response Models
```go
type ConversionRequest struct {
    Value     float64 `json:"value" validate:"required"`
    FromUnit  string  `json:"fromUnit" validate:"required,oneof=celsius fahrenheit kelvin rankine reaumur"`
    ToUnit    string  `json:"toUnit" validate:"required,oneof=celsius fahrenheit kelvin rankine reaumur"`
    Precision int     `json:"precision" validate:"min=0,max=10"`
}

type ConversionResponse struct {
    Success     bool        `json:"success"`
    Input       Temperature `json:"input"`
    Output      Temperature `json:"output"`
    Formula     string      `json:"formula"`
    Calculation string      `json:"calculation"`
    Precision   int         `json:"precision"`
}

type BulkConversionRequest struct {
    Conversions []ConversionRequest `json:"conversions" validate:"required,min=1,max=1000"`
    Precision   int                 `json:"precision" validate:"min=0,max=10"`
}
```

### Weather Calculation Models
```go
type HeatIndexRequest struct {
    Temperature Temperature `json:"temperature" validate:"required"`
    Humidity    float64     `json:"humidity" validate:"required,min=0,max=100"`
    OutputUnit  string      `json:"outputUnit" validate:"oneof=celsius fahrenheit"`
}

type WindChillRequest struct {
    Temperature Temperature `json:"temperature" validate:"required"`
    WindSpeed   float64     `json:"windSpeed" validate:"required,min=0"`
    WindUnit    string      `json:"windUnit" validate:"oneof=mph kmh ms"`
    OutputUnit  string      `json:"outputUnit" validate:"oneof=celsius fahrenheit"`
}

type DewPointRequest struct {
    Temperature Temperature `json:"temperature" validate:"required"`
    Humidity    float64     `json:"humidity" validate:"required,min=0,max=100"`
    OutputUnit  string      `json:"outputUnit" validate:"oneof=celsius fahrenheit"`
}
```

### Unit Information Model
```go
type TemperatureUnit struct {
    Name         string          `json:"name"`
    Symbol       string          `json:"symbol"`
    FullName     string          `json:"fullName"`
    Description  string          `json:"description"`
    AbsoluteZero float64         `json:"absoluteZero"`
    CommonRange  TemperatureRange `json:"commonRange"`
    ScientificUse bool           `json:"scientificUse"`
    Countries    []string        `json:"countries"`
    History      UnitHistory     `json:"history"`
}

type UnitHistory struct {
    InventedYear int    `json:"inventedYear"`
    Inventor     string `json:"inventor"`
    Origin       string `json:"origin"`
}
```

## SDK Interface

### Direct Mode (In-Process)
```go
type TemperatureSdk struct {
    Convert           func(value float64, fromUnit, toUnit string, precision int) (*ConversionResponse, error)
    ConvertBulk       func(conversions []ConversionRequest, precision int) (*BulkConversionResponse, error)
    CalculateDifference func(temp1, temp2 Temperature, outputUnit string) (*DifferenceResponse, error)
    CalculateAverage  func(temperatures []Temperature, outputUnit string, precision int) (*AverageResponse, error)
    CalculateHeatIndex func(temp Temperature, humidity float64, outputUnit string) (*HeatIndexResponse, error)
    CalculateWindChill func(temp Temperature, windSpeed float64, windUnit, outputUnit string) (*WindChillResponse, error)
    CalculateDewPoint func(temp Temperature, humidity float64, outputUnit string) (*DewPointResponse, error)
    ValidateTemperature func(temp Temperature, context string) (*ValidationResponse, error)
    GetUnits          func() ([]TemperatureUnit, error)
    GetFormulas       func() ([]ConversionFormula, error)
    CompareTemperature func(temp Temperature, compareUnits []string) (*ComparisonResponse, error)
}

func NewTemperatureSdk(mode string) *TemperatureSdk
```

### HTTP Mode (Network)
Standard HTTP JSON API calls to the temperature service endpoints.

## Environment Variables

```env
# Server Configuration
TEMPERATURE_SERVICE_PORT=8084
TEMPERATURE_SERVICE_HOST=0.0.0.0

# Calculation Configuration
TEMPERATURE_DEFAULT_PRECISION=2
TEMPERATURE_MAX_PRECISION=10
TEMPERATURE_MAX_BULK_CONVERSIONS=1000

# Validation Configuration
TEMPERATURE_STRICT_VALIDATION=true
TEMPERATURE_ALLOW_EXTREME_VALUES=false

# Weather API Integration (optional)
WEATHER_API_ENABLED=false
WEATHER_API_KEY=your-weather-api-key
WEATHER_API_URL=https://api.weather.com

# Caching Configuration
TEMPERATURE_CACHE_ENABLED=true
TEMPERATURE_CACHE_TTL=1h
TEMPERATURE_CACHE_SIZE=10000

# Performance Configuration
TEMPERATURE_MAX_CONCURRENT_REQUESTS=1000
TEMPERATURE_REQUEST_TIMEOUT=30s

# Monitoring
TEMPERATURE_METRICS_ENABLED=true
TEMPERATURE_LOG_LEVEL=info
```

## Conversion Algorithms

### Core Conversion Functions
```go
// Celsius conversions
func celsiusToFahrenheit(c float64) float64 {
    return (c * 9.0 / 5.0) + 32.0
}

func celsiusToKelvin(c float64) float64 {
    return c + 273.15
}

func celsiusToRankine(c float64) float64 {
    return (c * 9.0 / 5.0) + 491.67
}

func celsiusToReaumur(c float64) float64 {
    return c * 4.0 / 5.0
}

// Weather calculations
func calculateHeatIndex(tempF, humidity float64) float64 {
    // Rothfusz regression equation
    c1 := -42.379
    c2 := 2.04901523
    c3 := 10.14333127
    c4 := -0.22475541
    c5 := -0.00683783
    c6 := -0.05481717
    c7 := 0.00122874
    c8 := 0.00085282
    c9 := -0.00000199
    
    T := tempF
    RH := humidity
    
    hi := c1 + (c2*T) + (c3*RH) + (c4*T*RH) + (c5*T*T) + (c6*RH*RH) + 
          (c7*T*T*RH) + (c8*T*RH*RH) + (c9*T*T*RH*RH)
    
    return hi
}

func calculateWindChill(tempF, windMph float64) float64 {
    // Only valid for temperatures <= 50°F and wind speeds > 3 mph
    if tempF > 50 || windMph <= 3 {
        return tempF
    }
    
    wc := 35.74 + (0.6215 * tempF) - (35.75 * math.Pow(windMph, 0.16)) + 
          (0.4275 * tempF * math.Pow(windMph, 0.16))
    
    return wc
}
```

## Performance Requirements

- Conversion response time: < 1ms for single conversions
- Bulk conversion: < 100ms for 1000 conversions
- Concurrent requests: 1000+ simultaneous
- Memory usage: < 100MB base memory
- CPU usage: < 10% under normal load

## Error Handling

### HTTP Status Codes
- 200: Success
- 400: Bad Request (invalid parameters)
- 422: Unprocessable Entity (invalid temperature values)
- 429: Too Many Requests (rate limited)
- 500: Internal Server Error

### Error Response Format
```json
{
  "error": {
    "code": "INVALID_TEMPERATURE_UNIT",
    "message": "Unsupported temperature unit: celsius2",
    "details": {
      "provided": "celsius2",
      "supported": ["celsius", "fahrenheit", "kelvin", "rankine", "reaumur"]
    },
    "suggestions": [
      "Use 'celsius' instead of 'celsius2'",
      "Check the list of supported units at /temperature/units"
    ]
  }
}
```

### Common Error Codes
- `INVALID_TEMPERATURE_UNIT`: Unsupported temperature unit
- `BELOW_ABSOLUTE_ZERO`: Temperature below absolute zero
- `INVALID_PRECISION`: Precision out of range
- `BULK_LIMIT_EXCEEDED`: Too many conversions in bulk request
- `INVALID_HUMIDITY`: Humidity outside 0-100% range
- `INVALID_WIND_SPEED`: Negative wind speed

## Dependencies

```go
// Add to go.mod
require (
    github.com/gin-gonic/gin v1.10.1
    github.com/go-playground/validator/v10 v10.15.5
    github.com/prometheus/client_golang v1.17.0
    github.com/patrickmn/go-cache v2.1.0+incompatible
)
```

## Testing Requirements

### Unit Tests
- All conversion algorithms
- Weather calculation functions
- Input validation
- Error handling

### Integration Tests
- Complete API endpoint testing
- Bulk conversion functionality
- Weather calculation accuracy
- Performance benchmarks

### Accuracy Tests
- Known temperature conversions
- Weather calculation validation
- Precision testing
- Edge case handling

## Monitoring & Logging

### Metrics to Track
```
# Conversion metrics
temperature_conversions_total{from_unit, to_unit, status}
temperature_conversion_duration_seconds{from_unit, to_unit}
temperature_bulk_conversions_total{status}

# Weather calculation metrics
temperature_weather_calculations_total{type, status}
temperature_weather_duration_seconds{type}

# Error metrics
temperature_errors_total{error_type}
temperature_validation_failures_total{reason}
```

### Performance Logging
- Conversion accuracy tracking
- Response time monitoring
- Cache hit/miss rates
- Request volume patterns

## Migration Path

### Phase 1: Core Conversions
- Basic temperature conversions
- Unit information API
- Input validation

### Phase 2: Advanced Calculations
- Weather-related calculations
- Statistical functions
- Bulk operations

### Phase 3: Integration Features
- Weather API integration
- Historical data support
- Advanced analytics
- Caching optimization