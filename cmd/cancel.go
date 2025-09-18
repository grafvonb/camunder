package cmd

import (
	"strings"

	"github.com/grafvonb/camunder/internal/logging"
	"github.com/grafvonb/camunder/internal/services/common"
	"github.com/grafvonb/camunder/internal/services/processinstance"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var supportedResourcesForCancel = common.ResourceTypes{
	"pi": "process-instance",
}

var (
	flagCancelKey int64
)

// cancelCmd represents the cancel command
var cancelCmd = &cobra.Command{
	Use:     "cancel [resource name] [key]",
	Short:   "Cancel a resource of a given type by its key. " + supportedResourcesForCancel.PrettyString(),
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"c", "cn", "stop", "abort"},
	Run: func(cmd *cobra.Command, args []string) {
		log := logging.FromContext(cmd.Context())
		rn := strings.ToLower(args[0])
		svcs, err := NewFromContext(cmd.Context())
		if err != nil {
			cmd.PrintErrf("%v\n", err)
			return
		}

		switch rn {
		case "process-instance", "pi":
			svc, err := processinstance.New(svcs.Config, svcs.HTTP.Client(), svcs.Auth, log, flagQuiet)
			if err != nil {
				cmd.PrintErrf("error creating process instance service: %v\n", err)
				return
			}
			_, err = svc.CancelProcessInstance(cmd.Context(), flagCancelKey)
			if err != nil {
				cmd.PrintErrf("error cancelling process instance: %v\n", err)
				return
			}
		default:
			cmd.PrintErrf("unknown resource type: %s\n", rn)
			cmd.Println(supportedResourcesForGet.PrettyString())
		}
	},
}

func init() {
	rootCmd.AddCommand(cancelCmd)

	AddBackoffFlagsAndBindings(cancelCmd, viper.GetViper())

	cancelCmd.Flags().Int64VarP(&flagCancelKey, "key", "k", 0, "resource key (e.g. process instance) to cancel")
	_ = cancelCmd.MarkFlagRequired("key")
}
