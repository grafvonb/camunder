package cmd

import (
	"strings"

	c87operatev1 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v1"
	"github.com/grafvonb/camunder/internal/services/common"
	processinstance "github.com/grafvonb/camunder/internal/services/process-instance"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var supportedResourcesForDelete = common.ResourceTypes{
	"pi": "process-instance",
}

var (
	flagDeleteKey        string
	flagDeleteWithCancel bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [resource name] [key]",
	Short: "Delete a resource of a given type by its key. " + supportedResourcesForDelete.PrettyString(),
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
			svc, err := processinstance.New(svcs.Config, svcs.HTTP.Client(), svcs.Auth,
				processinstance.WithQuietEnabled(flagQuiet))
			if err != nil {
				cmd.PrintErrf("Error creating process instance service: %v\n", err)
				return
			}
			var pidr *c87operatev1.ProcessInstanceDeleteResponse
			if flagDeleteWithCancel {
				pidr, err = svc.DeleteProcessInstanceWithCancel(cmd.Context(), flagDeleteKey)
			} else {
				pidr, err = svc.DeleteProcessInstance(cmd.Context(), flagDeleteKey)
			}
			if err != nil {
				cmd.PrintErrf("Error deleting process instance with key %s: %v\n", flagDeleteKey, err)
				return
			}
			cmd.Println(ToJSONString(pidr))
		default:
			cmd.PrintErrf("Unknown resource type: %s\n", rn)
			cmd.Println(supportedResourcesForDelete.PrettyString())
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	AddBackoffFlagsAndBindings(deleteCmd, viper.GetViper())

	deleteCmd.Flags().StringVarP(&flagDeleteKey, "key", "k", "", "resource key (e.g. process instance) to delete")
	_ = deleteCmd.MarkFlagRequired("key")

	deleteCmd.Flags().BoolVarP(&flagDeleteWithCancel, "cancel", "c", false, "tries to cancel the process instance before deleting it (if not in the state COMPLETED or CANCELED)")
}
