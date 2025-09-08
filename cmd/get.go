package cmd

import (
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
	"github.com/spf13/viper"
)

const maxSearchSize int32 = 1000

var supportedResourcesForGet = common.ResourceTypes{
	"ct": "cluster-topology",
	"pd": "process-definition",
	"pi": "process-instance",
}

var (
	flagKey               int64
	flagBpmnProcessID     string
	flagProcessVersion    int32
	flagProcessVersionTag string
	flagState             string
	flagParentKey         int64
)

var searchFilterOpts processinstance.SearchFilterOpts

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [resource type]",
	Short: "List resources of a defined type. " + supportedResourcesForGet.PrettyString(),
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rn := strings.ToLower(args[0])
		searchFilterOpts = populateSearchFilterOpts()
		ko, err := cmd.Flags().GetBool(FlagKeyOnlyName)
		if err != nil {
			cmd.PrintErrf("error reading flag --%s: %v\n", FlagKeyOnlyName, err)
			return
		}
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
			if ko {
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
			if searchFilterOpts.Key != nil {
				pi, err := svc.GetProcessInstanceByKey(cmd.Context(), *searchFilterOpts.Key)
				if err != nil {
					cmd.PrintErrf("error fetching process instance by key %d: %v\n", *searchFilterOpts.Key, err)
					return
				}
				cmd.Println(ToJSONString(pi))
			} else {
				state, err := processinstance.PIStateFilterFromString(flagState)
				if err != nil {
					cmd.PrintErrf("error parsing state %q filter: %v\n", flagState, err)
					return
				}
				searchFilterOpts.State = state
				pisr, err := svc.SearchForProcessInstances(cmd.Context(), searchFilterOpts, maxSearchSize)
				if err != nil {
					cmd.PrintErrf("error fetching process instances: %v\n", err)
					return
				}
				if ko {
					keys := common.KeysFromItems(pisr.Items, func(it c87operatev1.ProcessInstanceItem) int64 {
						return *it.Key
					})
					for _, k := range keys {
						cmd.Println(k)
					}
					return
				}
				cmd.Println(ToJSONString(pisr))
			}

		default:
			cmd.PrintErrf("unknown resource type: %s\n", rn)
			cmd.Println(supportedResourcesForGet.PrettyString())
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	common.AddBackoffFlagsAndBindings(getCmd, viper.GetViper())

	fs := getCmd.Flags()
	fs.Int64VarP(&flagKey, "key", "k", 0, "resource key (e.g. process instance) to fetch")
	fs.StringVarP(&flagBpmnProcessID, "bpmn-process-id", "b", "", "BPMN process ID to filter process instances")
	fs.Int32VarP(&flagProcessVersion, "process-version", "v", 0, "process definition version")
	fs.StringVar(&flagProcessVersionTag, "process-version-tag", "", "process definition version tag")
	fs.Int64Var(&flagParentKey, "parent-key", 0, "parent process instance key")
	fs.StringVarP(&flagState, "state", "s", "all", "state to filter process instances: all, active, completed, canceled")
}

func populateSearchFilterOpts() processinstance.SearchFilterOpts {
	var opts processinstance.SearchFilterOpts
	if flagKey != 0 {
		opts.Key = &flagKey
	}
	if flagParentKey != 0 {
		opts.ParentKey = &flagParentKey
	}
	if flagBpmnProcessID != "" {
		opts.BpmnProcessId = &flagBpmnProcessID
	}
	if flagProcessVersion != 0 {
		opts.ProcessVersion = &flagProcessVersion
	}
	if flagProcessVersionTag != "" {
		opts.ProcessVersionTag = &flagProcessVersionTag
	}
	return opts
}
