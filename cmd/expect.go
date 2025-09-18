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

var supportedResourcesForExpect = common.ResourceTypes{
	"pi": "process-instance",
}

var (
	flagExpectKey int64
)

// expectCmd represents the cancel command
var expectCmd = &cobra.Command{
	Use:     "expect [resource name] [key]",
	Short:   "Expect a resource of a given type to change (e.g. its state) by its key. " + supportedResourcesForExpect.PrettyString(),
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"e", "exp", "await"},
	Run: func(cmd *cobra.Command, args []string) {
		log := logging.FromContext(cmd.Context())
		if err := requireAnyFlag(cmd, "state"); err != nil {
			log.Error(err.Error())
			return
		}
		rn := strings.ToLower(args[0])
		svcs, err := NewFromContext(cmd.Context())
		if err != nil {
			log.Error(fmt.Sprintf("Error initializing service from context: %v", err))
			return
		}

		switch rn {
		case "process-instance", "pi":
			svc, err := processinstance.New(svcs.Config, svcs.HTTP.Client(), svcs.Auth, log, flagQuiet)
			if err != nil {
				log.Error(fmt.Sprintf("error creating process instance service: %v", err))
				return
			}

			state, err := piapi.ParseState(flagState)
			if err == nil && state != piapi.StateAll {
				log.Info(fmt.Sprintf("waiting for process instance %d to reach state %q", flagExpectKey, state))
				err = svc.WaitForProcessInstanceState(cmd.Context(), flagExpectKey, state)
				if err != nil {
					log.Error(fmt.Sprintf("error waiting for a process instance to reach a %q state: %v", state, err))
					return
				}
			}
		default:
			log.Error(fmt.Sprintf("unknown resource type %q, supported: %s", rn, supportedResourcesForGet))
		}
	},
}

func init() {
	rootCmd.AddCommand(expectCmd)

	AddBackoffFlagsAndBindings(expectCmd, viper.GetViper())

	expectCmd.Flags().Int64VarP(&flagExpectKey, "key", "k", 0, "resource key (e.g. process instance)")
	_ = expectCmd.MarkFlagRequired("key")

	expectCmd.Flags().StringVarP(&flagState, "state", "s", "", "state of a process instance: active, completed, canceled or absent")
}
