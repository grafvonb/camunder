package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev" // set by ldflags
	commit  = "none"
	date    = "unknown"
)

var asJSON bool

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		if asJSON {
			out := map[string]string{
				"version": version,
				"commit":  commit,
				"date":    date,
			}
			b, err := json.MarshalIndent(out, "", "  ")
			if err != nil {
				cmd.PrintErrf("Error marshalling version info to JSON: %v\n", err)
				return
			}
			fmt.Println(string(b))
			return
		}
		fmt.Printf("Camunder version %s, commit %s, built at %s\n", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVarP(&asJSON, "json", "j", false, "output as json")
}
