package cmd

import (
	"fmt"
	"strings"

	"github.com/grafvonb/camunder/internal/logging"
	"github.com/grafvonb/camunder/internal/services/common"
	"github.com/grafvonb/camunder/internal/services/processinstance"
	piapi "github.com/grafvonb/camunder/pkg/camunda/processinstance"
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
		log := logging.FromContext(cmd.Context())
		if !validWalkModes[flagWalkMode] {
			log.Error(fmt.Sprintf("invalid value for --walk: %q (must be parent, children, or family)", flagWalkMode))
			return
		}
		rn := strings.ToLower(args[0])
		svcs, err := NewFromContext(cmd.Context())
		if err != nil {
			log.Error(fmt.Sprintf("%v", err))
			return
		}
		switch rn {
		case "process-instance", "pi":
			svc, err := processinstance.New(svcs.Config, svcs.HTTP.Client(), svcs.Auth, log)
			if err != nil {
				log.Error(fmt.Sprintf("error creating walk service: %v", err))
				return
			}
			walkerSvc, ok := svc.(piapi.Walker)
			if !ok {
				log.Error(fmt.Sprintf("walk command not supported by this API version %s\n", svcs.Config.APIs.Version))
				return
			}

			var path KeysPath
			var chain Chain
			switch flagWalkMode {
			case "parent":
				_, path, chain, err = walkerSvc.Ancestry(cmd.Context(), flagStartKey)
				if err != nil {
					return
				}
			case "children":
				path, _, chain, err = walkerSvc.Descendants(cmd.Context(), flagStartKey)
				if err != nil {
					return
				}
			case "family":
				path, _, chain, err = walkerSvc.Family(cmd.Context(), flagStartKey)
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
			log.Error(fmt.Sprintf("unknown resource type: %s, supported: %s", rn, supportedResourcesForWalk))
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
