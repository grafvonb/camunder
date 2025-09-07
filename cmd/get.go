package cmd

import (
	"strconv"
	"strings"

	c87operatev1 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v1"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth"
	"github.com/grafvonb/camunder/internal/services/cluster"
	"github.com/grafvonb/camunder/internal/services/common"
	"github.com/grafvonb/camunder/internal/services/httpc"
	processdefinition "github.com/grafvonb/camunder/internal/services/process-definition"
	processinstance "github.com/grafvonb/camunder/internal/services/process-instance"
	"github.com/spf13/cobra"
)

var supportedResourcesForGet = common.ResourceTypes{
	"ct": "cluster-topology",
	"pd": "process-definition",
	"pi": "process-instance",
}

var (
	bpmnProcessId string
	state         string
	quick         bool
	getPIKey      string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [resource type]",
	Short: "List resources of a defined type. " + supportedResourcesForGet.PrettyString(),
	Args:  cobra.ExactArgs(1),
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
		case "cluster-topology", "ct":
			svc, err := cluster.New(cfg, httpSvc.Client(), authSvc, isQuiet)
			if err != nil {
				cmd.PrintErrf("error creating cluster service: %v\n", err)
				return
			}
			topology, err := svc.GetClusterTopology(cmd.Context())
			if err != nil {
				cmd.PrintErrf("error fetching topology: %v\n", err)
				return
			}
			cmd.Println(ToJSONString(topology))
		case "process-definition", "pd":
			svc, err := processdefinition.New(cfg, httpSvc.Client(), authSvc, isQuiet)
			if err != nil {
				cmd.PrintErrf("error creating process definition service: %v\n", err)
				return
			}
			pdsr, err := svc.SearchForProcessDefinitions(cmd.Context())
			if err != nil {
				cmd.PrintErrf("error fetching process definitions: %v\n", err)
				return
			}
			if quick {
				keys := common.KeysFromItems(pdsr.Items, func(it c87operatev1.ProcessDefinitionItem) int64 {
					return *it.Key
				})
				for _, k := range keys {
					cmd.Println(k)
				}
				return
			}
			cmd.Println(ToJSONString(pdsr))
		case "process-instance", "pi":
			svc, err := processinstance.New(cfg, httpSvc.Client(), authSvc, isQuiet)
			if err != nil {
				cmd.PrintErrf("error creating process instance service: %v\n", err)
				return
			}
			if getPIKey == "" {
				if bpmnProcessId == "" {
					cmd.PrintErrln("please provide a process ID to filter process instances using the --bpmn-process-id flag")
					return
				}
				piStateFilter, err := processinstance.PIStateFilterFromString(state)
				if err != nil {
					cmd.PrintErrf("error parsing state filter: %v\n", err)
					return
				}
				pisr, err := svc.SearchForProcessInstances(cmd.Context(), bpmnProcessId, piStateFilter)
				if err != nil {
					cmd.PrintErrf("error fetching process instances: %v\n", err)
					return
				}
				if quick {
					keys := common.KeysFromItems(pisr.Items, func(it c87operatev1.ProcessInstanceItem) int64 {
						return *it.Key
					})
					for _, k := range keys {
						cmd.Println(k)
					}
					return
				}
				cmd.Println(ToJSONString(pisr))
			} else {
				piKey, err := strconv.ParseInt(getPIKey, 10, 64)
				if err != nil {
					cmd.PrintErrf("error parsing key as int64: %v\n", err)
					return
				}
				pi, err := svc.GetProcessInstanceByKey(cmd.Context(), piKey)
				if err != nil {
					cmd.PrintErrf("error fetching process instance by key %d: %v\n", piKey, err)
					return
				}
				cmd.Println(ToJSONString(pi))
			}
		default:
			cmd.PrintErrf("unknown resource type: %s\n", rn)
			cmd.Println(supportedResourcesForGet.PrettyString())
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&bpmnProcessId, "bpmn-process-id", "b", "", "BPMN process ID to filter process instances")
	getCmd.Flags().StringVarP(&state, "state", "s", "all", "state to filter process instances: all, active, completed, canceled")
	getCmd.Flags().BoolVarP(&quick, "quick", "q", false, "quick output (only keys for process definitions/instances)")
	getCmd.Flags().StringVarP(&getPIKey, "key", "k", "", "process instance key to fetch (if provided, --bpmn-process-id and --state are ignored)")
}
