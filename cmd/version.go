package cmd

import (
	"github.com/spf13/cobra"
)

var (
	version = "dev" // set by ldflags
	commit  = "none"
	date    = "unknown"
)

var (
	flagAsJSON bool
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		if flagAsJSON {
			out := map[string]string{
				"version": version,
				"commit":  commit,
				"date":    date,
			}
			cmd.Println(ToJSONString(out))
			return
		}
		cmd.Printf("Camunder version %s, commit %s, built at %s\n", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVarP(&flagAsJSON, "json", "j", false, "output as json")
}
