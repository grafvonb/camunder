package cmd

import (
	"fmt"
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
			log.Error(fmt.Sprintf("%v", err))
			return
		}

		switch rn {
		case "process-instance", "pi":
			svc, err := processinstance.New(svcs.Config, svcs.HTTP.Client(), log)
			if err != nil {
				log.Error(fmt.Sprintf("creating process instance service: %v", err))
				return
			}
			_, err = svc.CancelProcessInstance(cmd.Context(), flagCancelKey)
			if err != nil {
				log.Error(fmt.Sprintf("cancelling process instance: %v", err))
				return
			}
		default:
			log.Error(fmt.Sprintf("unknown resource type: %s, supported: %s", rn, supportedResourcesForCancel))
		}
	},
}

func init() {
	rootCmd.AddCommand(cancelCmd)

	AddBackoffFlagsAndBindings(cancelCmd, viper.GetViper())

	cancelCmd.Flags().Int64VarP(&flagCancelKey, "key", "k", 0, "resource key (e.g. process instance) to cancel")
	_ = cancelCmd.MarkFlagRequired("key")
}
