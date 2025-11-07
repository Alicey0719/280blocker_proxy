package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func handler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	targetURL := fmt.Sprintf("https://280blocker.net/files/280blocker_adblock_%04d%02d.txt",
		now.Year(), now.Month())

	clientIP := getClientIP(r)
	log.Printf("[%s] %s %s %s -> %s",
		time.Now().Format("2006-01-02 15:04:05"),
		clientIP, r.Method, r.URL.Path, targetURL)

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("User-Agent", "NABARI")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(status)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "34165"
	}

	http.HandleFunc("/280blocker.txt", handler)
	http.HandleFunc("/healthz", healthzHandler)
	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
