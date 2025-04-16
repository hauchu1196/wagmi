package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func GetLiveProxy(timeout time.Duration) ([]string, error) {
	// Fetch list of proxies
	proxies, err := FetchProxies()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch proxies: %v", err)
	}

	// Filter for HTTPS proxies
	var httpsProxies []Proxy
	for _, p := range proxies {
		if p.HTTPS == "yes" {
			httpsProxies = append(httpsProxies, p)
		}
	}

	if len(httpsProxies) == 0 {
		return nil, fmt.Errorf("no HTTPS proxies found")
	}

	// Channel to receive working proxies
	proxyChan := make(chan string, len(httpsProxies))
	// Channel to signal when all goroutines are done
	doneChan := make(chan struct{})
	var wg sync.WaitGroup

	// Launch goroutines to test proxies concurrently
	for _, proxy := range httpsProxies {
		wg.Add(1)
		go func(proxy Proxy) {
			defer wg.Done()

			// Create transport with proxy
			transport := &http.Transport{
				Proxy: http.ProxyURL(&url.URL{
					Scheme: "http",
					Host:   fmt.Sprintf("%s:%s", proxy.IP, proxy.Port),
				}),
			}
			client := &http.Client{
				Timeout:   timeout,
				Transport: transport,
			}

			// Test the proxy
			resp, err := client.Get("https://api.ipify.org?format=text")
			if err != nil {
				return
			}
			resp.Body.Close()

			// If we get here, the proxy is working
			proxyChan <- proxy.GetProxyURL()
		}(proxy)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(proxyChan)
		doneChan <- struct{}{}
	}()

	// Collect working proxies
	var workingProxies []string
	for proxyURL := range proxyChan {
		workingProxies = append(workingProxies, proxyURL)
	}
	<-doneChan

	if len(workingProxies) == 0 {
		return nil, fmt.Errorf("no working proxy found")
	}

	return workingProxies, nil
}

func makeClientWithProxy(proxy string, timeout time.Duration) *http.Client {
	proxyFunc := http.ProxyURL(&url.URL{Scheme: "http", Host: proxy[len("http://"):]})
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy: proxyFunc,
		},
	}
}

func MakeClient(proxyURL string, timeout time.Duration) *http.Client {
	return makeClientWithProxy(proxyURL, timeout)
}
