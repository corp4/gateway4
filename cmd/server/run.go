package main

import (
	"github.com/spf13/cobra"
)

var (
	// Host to listen on
	host string

	// API port to listen on
	apiPort int

	// HTTP port to listen on
	httpPort int

	// HTTPS port to listen on
	httpsPort int

	// Path to TLS certificate
	certPath string

	// Path to TLS key
	keyPath string

	// SSH port to listen on
	sshPort int

	// Path to SSH key
	sshKeyPath string
)

// Run the SSH server
var runCmd = &cobra.Command{
	Use:        "run",
	Short:      "Run the server",
	Long:       `Initialize the server and start listening for incoming connections.`,
	SuggestFor: []string{"serve", "start"},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Initialize the run command
func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&host, "host", "H", "0.0.0.0", "Host to listen on")
	runCmd.Flags().IntVarP(&apiPort, "api-port", "a", 8080, "API port to listen on")

	runCmd.Flags().IntVarP(&httpPort, "http-port", "p", 80, "HTTP port to listen on")
	runCmd.Flags().IntVarP(&httpsPort, "https-port", "P", 443, "HTTPS port to listen on")
	runCmd.Flags().StringVarP(&certPath, "cert", "c", "", "Path to TLS certificate")
	runCmd.Flags().StringVarP(&keyPath, "key", "k", "", "Path to TLS key")
	runCmd.MarkFlagRequired("cert")
	runCmd.MarkFlagRequired("key")

	runCmd.Flags().IntVarP(&sshPort, "ssh-port", "s", 22, "SSH port to listen on")
	runCmd.Flags().StringVarP(&sshKeyPath, "ssh-key", "K", "", "Path to SSH key")
	runCmd.MarkFlagRequired("ssh-key")
}
