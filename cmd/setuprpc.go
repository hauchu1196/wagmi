package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/hauchu1196/wagmi/internal/blockpi"
	"github.com/hauchu1196/wagmi/internal/tempmail"
	"github.com/spf13/cobra"
)

var chain string

var setuprpcCmd = &cobra.Command{
	Use:   "setup-rpc",
	Short: "Tạo tài khoản BlockPI, xác minh email và lấy RPC/WSS",
	Run: func(cmd *cobra.Command, args []string) {
		chainID := getChainID(chain)
		if chainID == 0 {
			log.Fatalf("❌ Chain không được hỗ trợ: %s", chain)
		}

		// Step 1: Tạo email tạm thời
		mailClient, err := tempmail.NewClient()
		if err != nil {
			log.Fatalf("❌ Lỗi tạo email tạm thời: %v", err)
		}

		email := mailClient.Email
		password := mailClient.Password
		fmt.Println("📩 Email:", email)
		fmt.Println("🔐 Password:", password)

		// Step 2: Đăng ký tài khoản BlockPI
		client := blockpi.NewClient()
		// raw password: 123456Aa@
		// hashed password: dd6064eab545bb5adb3d514c8398f1d0
		if err := client.Register(email, "dd6064eab545bb5adb3d514c8398f1d0"); err != nil {
			log.Fatalf("❌ Đăng ký thất bại: %v", err)
		}
		fmt.Println("✅ Đăng ký thành công.")

		// Step 3: Đợi email xác nhận và trích mã code
		fmt.Println("⏳ Đang đợi email xác nhận...")
		html, err := mailClient.WaitForConfirmationEmail(60 * time.Second)
		if err != nil {
			log.Fatalf("❌ Không lấy được email xác nhận: %v", err)
		}

		code, err := blockpi.ExtractConfirmationCode(html)
		if err != nil {
			log.Fatalf("❌ Không tìm được mã xác nhận: %v", err)
		}
		fmt.Println("✅ Mã xác nhận:", code)

		// Step 4: Gửi mã xác nhận
		token, err := client.EmailConfirm(code)
		if err != nil {
			log.Fatalf("❌ Xác nhận email thất bại: %v", err)
		}
		fmt.Println("✅ Email đã được xác minh.")

		// Step 5: Đăng nhập để lấy JWT
		client.Token = token
		fmt.Println("✅ Đăng nhập thành công.")

		// Step 6: Xác nhận lần đầu
		if err := client.FirstConfirm(); err != nil {
			log.Fatalf("❌ FirstConfirm thất bại: %v", err)
		}
		fmt.Println("✅ Tài khoản đã xác nhận lần đầu.")

		// Step 7: Sinh API Key cho chain
		rpc, wss, err := client.GenerateApiKey(chainID, "wagmi-"+chain)
		if err != nil {
			log.Fatalf("❌ Tạo API Key thất bại: %v", err)
		}

		// ✅ Output
		fmt.Println("✅ API Key đã tạo:")
		fmt.Println("🔌 RPC:", rpc)
		fmt.Println("🔌 WSS:", wss)
	},
}

func init() {
	setuprpcCmd.Flags().StringVarP(&chain, "chain", "c", "base", "Tên chain cần lấy RPC (vd: base, base-sepolia, ethereum, ethereum-sepolia)")
	rootCmd.AddCommand(setuprpcCmd)
}

// --- chain name to chain ID mapping ---

func getChainID(name string) int {
	switch name {
	case "base":
		return 2030
	case "base-sepolia":
		return 2041
	case "ethereum":
		return 1006
	case "ethereum-sepolia":
		return 1011
	default:
		return 0
	}
}
