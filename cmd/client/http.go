package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/corp4/gateway4/internal/client"
	"github.com/corp4/gateway4/internal/protocol"
	"github.com/spf13/cobra"
)

// Run the SSH server
var httpCmd = &cobra.Command{
	Use:   "http <address>",
	Short: "Expose an HTTP service to the internet",
	Long:  `Expose an HTTP service to the internet.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := client.NewClient(args[0])
		if err != nil {
			log.Errorf("failed to create client: %v", err)
			return
		}
		defer client.Close()

		for {
			err := protocol.ReadServer(client, client)
			if err != nil {
				log.Errorf("%v", err)
				return
			}
		}
	},
}

// Initialize the run command
func init() {
	rootCmd.AddCommand(httpCmd)
}
