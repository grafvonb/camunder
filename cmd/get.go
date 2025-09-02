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
		cfg := config.MustFrom(cmd.Context())

		/*
			token = viper.GetString("token")
			if token == "" {
				token = os.Getenv("CAMUNDA8_API_TOKEN")
				if token == "" {
					cmd.PrintErrln("Error: Bearer token must be provided via --token flag or CAMUNDA8_API_TOKEN environment variable")
					return
				}
			}
		*/

		httpClient := &http.Client{
			Timeout: cfg.HTTP.Timeout * time.Second,
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

	getCmd.Flags().String("token", "", "Bearer token for authentication")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
