package cmd

import (
	"strings"

	"github.com/grafvonb/camunder/internal/services/common"
	walkv87 "github.com/grafvonb/camunder/internal/services/walk/v87"
	"github.com/spf13/cobra"
)

var supportedResourcesForWalk = common.ResourceTypes{
	"pi": "process-instance",
}

var (
	flagStartKey int64
	flagWalkMode string
)

var validWalkModes = map[string]bool{
	"parent":   true,
	"children": true,
	"family":   true,
}

var walkCmd = &cobra.Command{
	Use:     "walk [resource type]",
	Short:   "Traverse (walk) the parent/child graph of resource type. " + supportedResourcesForWalk.PrettyString(),
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"w", "traverse"},
	Run: func(cmd *cobra.Command, args []string) {
		if !validWalkModes[flagWalkMode] {
			cmd.PrintErrf("invalid value for --walk: %q (must be parent, children, or family)", flagWalkMode)
			return
		}
		rn := strings.ToLower(args[0])
		svcs, err := NewFromContext(cmd.Context())
		if err != nil {
			cmd.PrintErrf("%v\n", err)
			return
		}
		switch rn {
		case "process-instance", "pi":
			svc, err := walkv87.New(svcs.Config, svcs.HTTP.Client(), svcs.Auth,
				walkv87.WithQuietEnabled(flagQuiet))
			if err != nil {
				cmd.PrintErrf("error creating walk service: %v\n", err)
				return
			}
			var path KeysPath
			var chain Chain
			switch flagWalkMode {
			case "parent":
				_, path, chain, err = svc.Ancestry(cmd.Context(), flagStartKey)
				if err != nil {
					return
				}
			case "children":
				path, _, chain, err = svc.Descendants(cmd.Context(), flagStartKey)
				if err != nil {
					return
				}
			case "family":
				path, _, chain, err = svc.Family(cmd.Context(), flagStartKey)
				if err != nil {
					return
				}
			default:
				return
			}
			if flagKeysOnly {
				cmd.Println(path.KeysOnly(chain))
				return
			}
			cmd.Println(path.StandardLine(chain))
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
	fs.StringVarP(&flagWalkMode, "mode", "m", "", "walk mode: parent, children, family")
	_ = walkCmd.MarkFlagRequired("mode")

	// view options
	fs.BoolVarP(&flagKeysOnly, "keys-only", "", false, "only print the keys of the resources")
}
