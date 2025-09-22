package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	authsvc "github.com/growly-group/golang-hybrid-api/src/auth-svc"
	calculatorsvc "github.com/growly-group/golang-hybrid-api/src/calculator-svc"
	gatewaysvc "github.com/growly-group/golang-hybrid-api/src/gateway-svc"
	pdfsvc "github.com/growly-group/golang-hybrid-api/src/pdf-svc"
	"github.com/joho/godotenv"
)

type entrypointFunc func()

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found or failed to load")
	}
	targetServices := os.Getenv("TARGET_SERVICES")
	if targetServices == "" {
		fmt.Println("TARGET_SERVICES environment variable not set")
		os.Exit(1)
	}

	entrypoints := map[string]entrypointFunc{
		"auth-svc":       authsvc.Entrypoint,
		"calculator-svc": calculatorsvc.Entrypoint,
		"gateway-svc":    gatewaysvc.Entrypoint,
		"pdf-svc":        pdfsvc.Entrypoint,
	}

	services := strings.Split(targetServices, ",")
	var wg sync.WaitGroup

	for _, serviceName := range services {
		trimmedServiceName := strings.TrimSpace(serviceName)
		if trimmedServiceName == "" {
			continue
		}

		entrypoint, ok := entrypoints[trimmedServiceName]
		if !ok {
			fmt.Printf("No entrypoint found for service: %s\n", trimmedServiceName)
			continue
		}

		wg.Add(1)
		go func(ep entrypointFunc, name string) {
			defer wg.Done()
			fmt.Printf("Starting service: %s\n", name)
			ep()
			fmt.Printf("Service finished: %s\n", name)
		}(entrypoint, trimmedServiceName)
	}

	wg.Wait()
	fmt.Println("All services have finished.")
}
