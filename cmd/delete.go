package cmd

import (
	"strings"

	"github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v87"
	"github.com/grafvonb/camunder/internal/services/common"
	v88 "github.com/grafvonb/camunder/internal/services/processinstance/v87"
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
		rn := strings.ToLower(args[0])
		svcs, err := NewFromContext(cmd.Context())
		if err != nil {
			cmd.PrintErrf("%v\n", err)
			return
		}
		switch rn {
		case "process-instance", "pi":
			svc, err := v88.New(svcs.Config, svcs.HTTP.Client(), svcs.Auth,
				v88.WithQuietEnabled(flagQuiet))
			if err != nil {
				cmd.PrintErrf("error creating process instance service: %v\n", err)
				return
			}
			var pidr *v87.ChangeStatus
			if flagDeleteWithCancel {
				pidr, err = svc.DeleteProcessInstanceWithCancel(cmd.Context(), flagDeleteKey)
			} else {
				pidr, err = svc.DeleteProcessInstance(cmd.Context(), flagDeleteKey)
			}
			if err != nil {
				cmd.PrintErrf("error deleting process instance with key %d: %v\n", flagDeleteKey, err)
				return
			}
			cmd.Println(ToJSONString(pidr))
		default:
			cmd.PrintErrf("unknown resource type: %s\n", rn)
			cmd.Println(supportedResourcesForDelete.PrettyString())
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
