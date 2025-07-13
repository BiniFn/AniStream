package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"
)

var (
	addr      = flag.String("addr", ":1234", "Address to listen on")
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/58.0.3029.110 Safari/537.3"
	allowedExts = map[string]bool{
		".m3u8": true, ".ts": true, ".png": true,
		".jpg": true, ".webp": true, ".ico": true,
		".html": true, ".js": true, ".css": true,
		".txt": true,
	}
)

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/proxy", proxyHandler)

	srv := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	errChan := make(chan error, 1)
	go func() {
		log.Printf("📡 proxy listening on %s", *addr)

		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	select {
	case err := <-errChan:
		log.Fatalf("Error starting server: %v", err)
	case sig := <-stop:
		log.Printf("🛑 caught signal %s, shutting down...", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("graceful shutdown failed: %v", err)
		}
		log.Println("✅ server stopped gracefully")
		os.Exit(0)
	}
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	p := r.URL.Query().Get("p")
	if p == "" {
		http.Error(w, "`p` is required", http.StatusBadRequest)
		return
	}
	rawURL, err := base64.URLEncoding.DecodeString(p)
	if err != nil {
		http.Error(w, "Invalid `p` parameter", http.StatusBadRequest)
		return
	}
	targetURL := string(rawURL)

	serverName := r.URL.Query().Get("s")
	isHianime := strings.HasPrefix(strings.ToLower(serverName), "hd")

	headers := http.Header{}
	headers.Set("Accept", "*/*")
	if isHianime {
		headers.Set("Referer", "https://megacloud.blog/")
		headers.Set("Origin", "https://megacloud.blog")
	} else {
		headers.Set("Referer", "https://megaplay.buzz/")
		headers.Set("Origin", "https://megaplay.buzz")
		headers.Set("User-Agent", userAgent)
	}

	req, err := http.NewRequestWithContext(ctx, r.Method, targetURL, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		http.Error(w, "bad target URL", http.StatusBadRequest)
		return
	}

	req.Header = headers

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error fetching upstream:", err)
		http.Error(w, "upstream fetch failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	u, _ := url.Parse(targetURL)
	ext := path.Ext(u.Path)

	if ext == ".m3u8" {
		w.Header().Del("Content-Length")
	}

	w.WriteHeader(resp.StatusCode)

	encodeProxyURL := func(next string) string {
		full := next
		if !strings.HasPrefix(next, "http") {
			// relative in playlist
			base := targetURL[:strings.LastIndex(targetURL, "/")+1]
			full = base + next
		}
		pEnc := base64.URLEncoding.EncodeToString([]byte(full))
		return "/proxy?p=" + pEnc + "&s=" + url.QueryEscape(serverName)
	}

	if ext == ".m3u8" {
		scanner := bufio.NewScanner(resp.Body)
		flusher, _ := w.(http.Flusher)

		for scanner.Scan() {
			line := scanner.Text()
			out := line

			if !strings.HasPrefix(line, "#") {
				e := strings.ToLower(path.Ext(line))
				if _, ok := allowedExts[e]; ok {
					out = encodeProxyURL(line)
				}
			}

			io.WriteString(w, out+"\n")
			if flusher != nil {
				flusher.Flush()
			}
		}
		if err := scanner.Err(); err != nil {
			log.Println("⚠️ playlist scan error:", err)
		}
	} else {
		// static content (ts, images, etc.)—just pipe directly
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		io.Copy(w, resp.Body)
	}

	log.Printf("🔗 proxied %s %s -> %s", r.Method, r.RemoteAddr, targetURL)
}
