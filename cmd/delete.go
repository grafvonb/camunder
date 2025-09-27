package cmd

import (
	"fmt"
	"strings"

	"github.com/grafvonb/camunder/internal/logging"
	"github.com/grafvonb/camunder/internal/services/common"
	"github.com/grafvonb/camunder/internal/services/processinstance"
	piapi "github.com/grafvonb/camunder/pkg/camunda/processinstance"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var supportedResourcesForDelete = common.ResourceTypes{
	"pi": "process-instance",
}

var (
	flagDeleteKey        int64
	flagDeleteWithCancel bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete [resource name] [key]",
	Short:   "Delete a resource of a given type by its key. " + supportedResourcesForDelete.PrettyString(),
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"d", "del", "remove", "rm"},
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
			var pidr piapi.ChangeStatus
			if flagDeleteWithCancel {
				pidr, err = svc.DeleteProcessInstanceWithCancel(cmd.Context(), flagDeleteKey)
			} else {
				pidr, err = svc.DeleteProcessInstance(cmd.Context(), flagDeleteKey)
			}
			if err != nil {
				log.Error(fmt.Sprintf("deleting process instance with key %d: %v", flagDeleteKey, err))
				return
			}
			log.Debug(pidr.String())
		default:
			log.Error(fmt.Sprintf("unknown resource type: %s, supported: %s", rn, supportedResourcesForDelete))
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	AddBackoffFlagsAndBindings(deleteCmd, viper.GetViper())

	deleteCmd.Flags().Int64VarP(&flagDeleteKey, "key", "k", 0, "resource key (e.g. process instance) to delete")
	_ = deleteCmd.MarkFlagRequired("key")

	deleteCmd.Flags().BoolVarP(&flagDeleteWithCancel, "cancel", "c", false, "tries to cancel the process instance before deleting it (if not in the state COMPLETED or CANCELED)")
}
