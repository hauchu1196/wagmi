package proxy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Proxy represents a proxy server
type Proxy struct {
	IP        string
	Port      string
	Code      string
	Country   string
	Anonymity string
	Google    string
	HTTPS     string
	LastCheck string
}

// FetchProxies retrieves a list of proxies from free-proxy-list.net
func FetchProxies() ([]Proxy, error) {
	resp, err := http.Get("https://free-proxy-list.net/")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch proxies: %v", err)
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var proxies []Proxy

	// Find the proxy table and iterate through rows
	doc.Find("table.table-striped tbody tr").Each(func(i int, s *goquery.Selection) {
		// Get all cells in the row
		cells := s.Find("td")
		if cells.Length() >= 8 {
			proxy := Proxy{
				IP:        strings.TrimSpace(cells.Eq(0).Text()),
				Port:      strings.TrimSpace(cells.Eq(1).Text()),
				Code:      strings.TrimSpace(cells.Eq(2).Text()),
				Country:   strings.TrimSpace(cells.Eq(3).Text()),
				Anonymity: strings.TrimSpace(cells.Eq(4).Text()),
				Google:    strings.TrimSpace(cells.Eq(5).Text()),
				HTTPS:     strings.TrimSpace(cells.Eq(6).Text()),
				LastCheck: strings.TrimSpace(cells.Eq(7).Text()),
			}
			
			proxies = append(proxies, proxy)
		}
	})

	if len(proxies) == 0 {
		return nil, fmt.Errorf("no proxies found")
	}

	return proxies, nil
}

// GetProxyURL returns the proxy URL in format http://ip:port
func (p *Proxy) GetProxyURL() string {
	return fmt.Sprintf("http://%s:%s", p.IP, p.Port)
}
