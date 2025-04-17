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

type SetupRPCResult struct {
	Email    string
	Password string
	RPC      string
	WSS      string
}

func setupRPC(chain string) (*SetupRPCResult, error) {
	chainID := getChainID(chain)
	if chainID == 0 {
		return nil, fmt.Errorf("unsupported chain: %s", chain)
	}

	// Step 1: Create temporary email
	mailClient, err := tempmail.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary email: %v", err)
	}
	fmt.Println("📩 Email:", mailClient.Email)
	fmt.Println("🔐 Password:", mailClient.Password)

	// Step 2: Register BlockPI account
	client := blockpi.NewClient()
	if err := client.Register(mailClient.Email, "dd6064eab545bb5adb3d514c8398f1d0"); err != nil {
		return nil, fmt.Errorf("registration failed: %v", err)
	}
	fmt.Println("✅ Đăng ký thành công.")

	// Step 3: Wait for confirmation email
	fmt.Println("⏳ Đang đợi email xác nhận...")
	html, err := mailClient.WaitForConfirmationEmail(60 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to get confirmation email: %v", err)
	}

	code, err := blockpi.ExtractConfirmationCode(html)
	if err != nil {
		return nil, fmt.Errorf("failed to extract confirmation code: %v", err)
	}
	fmt.Println("✅ Mã xác nhận:", code)

	// Step 4: Send confirmation code
	token, err := client.EmailConfirm(code)
	if err != nil {
		return nil, fmt.Errorf("email confirmation failed: %v", err)
	}
	fmt.Println("✅ Email đã được xác minh.")

	// Step 5: Login to get JWT
	client.Token = token
	fmt.Println("✅ Đăng nhập thành công.")

	// Step 6: First confirmation
	if err := client.FirstConfirm(); err != nil {
		return nil, fmt.Errorf("first confirmation failed: %v", err)
	}
	fmt.Println("✅ Tài khoản đã xác nhận lần đầu.")

	// Step 7: Generate API Key
	rpc, wss, err := client.GenerateApiKey(chainID, "wagmi-"+chain)
	if err != nil {
		return nil, fmt.Errorf("failed to generate API Key: %v", err)
	}
	fmt.Println("✅ API Key đã tạo:")
	fmt.Println("🔌 RPC:", rpc)
	fmt.Println("🔌 WSS:", wss)

	return &SetupRPCResult{
		Email:    mailClient.Email,
		Password: mailClient.Password,
		RPC:      rpc,
		WSS:      wss,
	}, nil
}

var setuprpcCmd = &cobra.Command{
	Use:   "setup-rpc",
	Short: "Tạo tài khoản BlockPI, xác minh email và lấy RPC/WSS",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := setupRPC(chain)
		if err != nil {
			log.Fatalf("❌ %v", err)
		}
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
