package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

// fetch は指定URLから中身を取得して返すヘルパー
func fetch(url string) ([]byte, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func startServer() {
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		http.HandleFunc("/healthz", healthzHandler)
		http.HandleFunc("/", handler)
		http.ListenAndServe(":"+port, nil)
	}()
	time.Sleep(500 * time.Millisecond) // サーバ起動待ち
}

func TestProxyMatchesUpstream(t *testing.T) {
	startServer()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	localURL := fmt.Sprintf("http://localhost:%s/280blocker.txt", port)
	upstreamURL := fmt.Sprintf("https://280blocker.net/files/280blocker_adblock_%04d%02d.txt",
		time.Now().Year(), time.Now().Month())

	t.Logf("Comparing local proxy:\n  %s\nwith upstream:\n  %s", localURL, upstreamURL)

	local, err := fetch(localURL)
	if err != nil {
		t.Fatalf("failed to fetch from proxy: %v", err)
	}

	upstream, err := fetch(upstreamURL)
	if err != nil {
		t.Fatalf("failed to fetch from upstream: %v", err)
	}

	if !bytes.Equal(local, upstream) {
		t.Errorf("❌ Mismatch detected (first 200 bytes)")
		t.Errorf("Local:\n%s", string(local[:min(200, len(local))]))
		t.Errorf("Upstream:\n%s", string(upstream[:min(200, len(upstream))]))
		t.FailNow()
	}

	t.Log("✅ Proxy content matches upstream")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
