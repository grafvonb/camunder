/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/grafvonb/camunder/internal/config"
	processinstance "github.com/grafvonb/camunder/internal/services/process-instance"
	"github.com/spf13/cobra"
)

var piKey string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [resource name] [key]",
	Short: "Delete a resource of a given type e.g. process instance by its key.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		rn := strings.ToLower(args[0])
		key := strings.ToLower(args[1])
		cfg, err := config.FromContext(cmd.Context())
		if err != nil {
			cmd.PrintErrf("Error retrieving config from context: %v\n", err)
			return
		}
		timeout, err := time.ParseDuration(cfg.HTTP.Timeout)
		if err != nil {
			cmd.PrintErrf("Error parsing '%s' as timeout duration: %v\n", cfg.HTTP.Timeout, err)
			return
		}
		httpClient := &http.Client{
			Timeout: timeout,
		}
		switch rn {
		case "process-instance", "pi":
			svc, err := processinstance.New(cfg.OperateAPI.BaseURL, httpClient, cfg.OperateAPI.Token)
			if err != nil {
				cmd.PrintErrf("Error creating process instance service: %v\n", err)
				return
			}
			pds, err := svc.DeleteProcessInstance(cmd.Context(), key)
			if err != nil {
				cmd.PrintErrf("Error deleting process instance with key %s: %v\n", key, err)
				return
			}
			b, _ := json.MarshalIndent(pds, "", "  ")
			cmd.Println(string(b))
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVar(&piKey, "pikey", "", "The process instance key to delete.")
}
