package cmd

import (
	"fmt"
	"log"

	"github.com/hauchu1196/wagmi/internal/proxy"
	"github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Lấy danh sách proxy từ free-proxy-list.net",
	Run: func(cmd *cobra.Command, args []string) {
		proxies, err := proxy.FetchProxies()
		if err != nil {
			log.Fatalf("❌ Lỗi khi lấy danh sách proxy: %v", err)
		}

		fmt.Printf("✅ Tìm thấy %d proxy:\n\n", len(proxies))
		for i, p := range proxies {
			fmt.Printf("%d. %s:%s (%s) - %s\n", i+1, p.IP, p.Port, p.Country, p.Anonymity)
		}
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
}
