package cmd

import (
	"strings"

	"github.com/grafvonb/camunder/internal/services/common"
	"github.com/grafvonb/camunder/internal/services/walk"
	"github.com/spf13/cobra"
)

var supportedResourcesForWalk = common.ResourceTypes{
	"pi": "process-instance",
}

var (
	flagStartKey int64
)

var walkCmd = &cobra.Command{
	Use:   "walk",
	Short: "Traverse (walk) the parent/child graph process instances.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rn := strings.ToLower(args[0])
		svcs, err := NewFromContext(cmd.Context())
		if err != nil {
			cmd.PrintErrf("%v\n", err)
			return
		}
		switch rn {
		case "process-instance", "pi":
			svc, err := walk.New(svcs.Config, svcs.HTTP.Client(), svcs.Auth,
				walk.WithQuietEnabled(flagQuiet))
			if err != nil {
				cmd.PrintErrf("error creating walk service: %v\n", err)
				return
			}
			_, path, chain, err := svc.Ancestry(cmd.Context(), flagStartKey)
			if err != nil {
				return
			}
			cmd.Println(printKeys(path, chain))
		default:
			cmd.PrintErrf("unknown resource type: %s\n", rn)
			cmd.Println(supportedResourcesForWalk.PrettyString())
		}
	},
}

func init() {
	rootCmd.AddCommand(walkCmd)

	fs := walkCmd.Flags()
	fs.Int64VarP(&flagStartKey, "start-key", "w", 0, "start walking from this process instance key")
	_ = walkCmd.MarkFlagRequired("start-key")
}
