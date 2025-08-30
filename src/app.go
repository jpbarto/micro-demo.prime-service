package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

type PrimeResponse struct {
	Hostname string `json:"hostname"`
	Primes   []int  `json:"primes"`
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func primesUpTo(n int) []int {
	primes := []int{}
	for i := 2; i <= n; i++ {
		if isPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err == nil {
		return hostname
	}
	// fallback to local IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "unknown"
}

func primeHandler(w http.ResponseWriter, r *http.Request) {
	numStr := r.URL.Query().Get("number")
	n, err := strconv.Atoi(numStr)
	if err != nil || n < 2 {
		http.Error(w, "Invalid or missing 'number' parameter", http.StatusBadRequest)
		return
	}

	primes := primesUpTo(n)
	resp := PrimeResponse{
		Hostname: getHostname(),
		Primes:   primes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/primes", primeHandler)
	http.HandleFunc("/health", healthHandler)

	port := 8080
	log.Printf("Prime service running on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
