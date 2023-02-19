package main

import (
	"github.com/spf13/cobra"
)

// Run the SSH server
var httpCmd = &cobra.Command{
	Use:   "http <address>",
	Short: "Expose an HTTP service to the internet",
	Long:  `Expose an HTTP service to the internet.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Initialize the run command
func init() {
	rootCmd.AddCommand(httpCmd)
}
