/*
Package cmd

Copyright © 2026 Adam Boczek <adam@boczek.com>
*/
package cmd

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/cluster"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [resource type]",
	Short: "List resources of a defined type e.g. cluster-topology, process-instances etc.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rn := strings.ToLower(args[0])
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
		case "cluster-topology", "ct":
			svc, err := cluster.New(cfg.API.BaseURL, httpClient, cfg.API.Token)
			if err != nil {
				cmd.PrintErrf("Error creating cluster service: %v\n", err)
				return
			}
			topology, err := svc.GetClusterTopology(cmd.Context())
			if err != nil {
				cmd.PrintErrf("Error fetching topology: %v\n", err)
				return
			}
			b, _ := json.MarshalIndent(topology, "", "  ")
			cmd.Println(string(b))
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
