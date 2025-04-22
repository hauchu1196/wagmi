package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

type SetupRPCRequest struct {
	Chain         string `json:"chain"`
	UseMailDomain bool   `json:"use_mail_domain"`
}

type SetupRPCResponse struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	RPC      string `json:"rpc"`
	WSS      string `json:"wss"`
}

func setupRPCHandler(w http.ResponseWriter, r *http.Request) {
	var req SetupRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := setupRPC(req.Chain, req.UseMailDomain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := SetupRPCResponse{
		Email:    result.Email,
		Password: result.Password,
		RPC:      result.RPC,
		WSS:      result.WSS,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Panic recovered: %v\n", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		if port == "" {
			port = "8080"
		}
		if port[0] != ':' {
			port = ":" + port
		}

		r := mux.NewRouter()
		r.Use(recoveryMiddleware)
		r.HandleFunc("/setup-rpc", setupRPCHandler).Methods("POST")

		fmt.Printf("API server running on http://localhost%s\n", port)
		http.ListenAndServe(port, r)
	},
}

func init() {
	apiCmd.Flags().StringP("port", "p", "8080", "Port to run the API server on")
	rootCmd.AddCommand(apiCmd)
}
