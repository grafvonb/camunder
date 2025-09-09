package cmd

import (
	"strings"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth"
	"github.com/grafvonb/camunder/internal/services/common"
	"github.com/grafvonb/camunder/internal/services/httpc"
	processinstance "github.com/grafvonb/camunder/internal/services/process-instance"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var supportedResourcesForCancel = common.ResourceTypes{
	"pi": "process-instance",
}

var (
	flagCancelKey string
)

// cancelCmd represents the cancel command
var cancelCmd = &cobra.Command{
	Use:   "cancel [resource name] [key]",
	Short: "Cancel a resource of a given type by its key. " + supportedResourcesForCancel.PrettyString(),
	Run: func(cmd *cobra.Command, args []string) {
		rn := strings.ToLower(args[0])
		cfg, err := config.FromContext(cmd.Context())
		if err != nil {
			cmd.PrintErrf("error retrieving config from context: %v\n", err)
			return
		}
		httpSvc, err := httpc.FromContext(cmd.Context())
		if err != nil {
			cmd.PrintErrf("error creating http service: %v\n", err)
			return
		}
		authSvc, err := auth.FromContext(cmd.Context())
		if err != nil {
			cmd.PrintErrf("error retrieving auth service: %v\n", err)
			return
		}
		switch rn {
		case "process-instance", "pi":
			svc, err := processinstance.New(cfg, httpSvc.Client(), authSvc, flagQuiet)
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

	cancelCmd.Flags().StringVarP(&flagCancelKey, "key", "k", "", "resource key (e.g. process instance) to cancel")
	_ = cancelCmd.MarkFlagRequired("key")
}
