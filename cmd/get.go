package cmd

import (
	"strings"

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

// filter options
var (
	flagKey               int64
	flagBpmnProcessID     string
	flagProcessVersion    int32
	flagProcessVersionTag string
	flagState             string
	flagParentKey         int64
)

// command options
var (
	flagParentsOnly       bool
	flagChildrenOnly      bool
	flagWithOrphanParents bool
)

// view options
var (
	flagKeysOnly bool
	flagOneLine  bool
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
			svc, err := cluster.New(cfg, httpSvc.Client(), authSvc, flagQuiet)
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
			searchFilterOpts := populatePDSearchFilterOpts()
			svc, err := processdefinition.New(cfg, httpSvc.Client(), authSvc, flagQuiet)
			if err != nil {
				cmd.PrintErrf("error creating process definition service: %v\n", err)
				return
			}
			if searchFilterOpts.Key != nil {
				pd, err := svc.GetProcessDefinitionByKey(cmd.Context(), *searchFilterOpts.Key)
				if err != nil {
					cmd.PrintErrf("error fetching process definition by key %d: %v\n", *searchFilterOpts.Key, err)
					return
				}
				err = ProcessDefinitionView(cmd, pd)
			} else {
				pdsr, err := svc.SearchForProcessDefinitions(cmd.Context(), searchFilterOpts, maxSearchSize)
				if err != nil {
					cmd.PrintErrf("error fetching process definitions: %v\n", err)
					return
				}
				if flagKeysOnly {
					err = ListKeyOnlyProcessDefinitionsView(cmd, pdsr)
					if err != nil {
						cmd.PrintErrf("error rendering keys-only view: %v\n", err)
					}
					return
				}
				err = ListProcessDefinitionsView(cmd, pdsr)
				if err != nil {
					cmd.PrintErrf("error rendering items view: %v\n", err)
				}
			}

		case "process-instance", "pi":
			searchFilterOpts := populatePISearchFilterOpts()
			svc, err := processinstance.New(cfg, httpSvc.Client(), authSvc, flagQuiet)
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
				err = ProcessInstanceView(cmd, pi)
				if err != nil {
					cmd.PrintErrf("error rendering key-only view: %v\n", err)
					return
				}
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
				if flagChildrenOnly && flagParentsOnly {
					cmd.PrintErrf("using both --children-only and --parents-only filters returns always no results\n")
					return
				}
				if flagChildrenOnly {
					pisr = pisr.FilterChildrenOnly()
				}
				if flagParentsOnly {
					pisr = pisr.FilterParentsOnly()
				}
				if flagWithOrphanParents {
					pisr.Items, err = svc.FilterProcessInstanceWithOrphanParent(cmd.Context(), pisr.Items)
					if err != nil {
						cmd.PrintErrf("error filtering process instances with orphan parents: %v\n", err)
						return
					}
					nt := int64(len(*pisr.Items))
					pisr.Total = &nt
				}
				if flagKeysOnly {
					err = ListKeyOnlyProcessInstancesView(cmd, pisr)
					if err != nil {
						cmd.PrintErrf("error rendering keys-only view: %v\n", err)
					}
					return
				}
				err = ListProcessInstancesView(cmd, pisr)
				if err != nil {
					cmd.PrintErrf("error rendering items view: %v\n", err)
				}
			}

		default:
			cmd.PrintErrf("unknown resource type: %s\n", rn)
			cmd.Println(supportedResourcesForGet.PrettyString())
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	AddBackoffFlagsAndBindings(getCmd, viper.GetViper())

	fs := getCmd.Flags()
	fs.Int64VarP(&flagKey, "key", "k", 0, "resource key (e.g. process instance) to fetch")
	fs.StringVarP(&flagBpmnProcessID, "bpmn-process-id", "b", "", "BPMN process ID to filter process instances")
	fs.Int32VarP(&flagProcessVersion, "process-version", "v", 0, "process definition version")
	fs.StringVar(&flagProcessVersionTag, "process-version-tag", "", "process definition version tag")
	fs.Int64Var(&flagParentKey, "parent-key", 0, "parent process instance key")
	fs.StringVarP(&flagState, "state", "s", "all", "state to filter process instances: all, active, completed, canceled")

	// command options
	fs.BoolVar(&flagParentsOnly, "parents-only", false, "show only parent process instances, meaning instances with no parent key set")
	fs.BoolVar(&flagChildrenOnly, "children-only", false, "show only child process instances, meaning instances that have a parent key set")
	fs.BoolVar(&flagWithOrphanParents, "with-orphan-parents", false, "show only child instances whose parent does not exist (return 404 on get by key)")

	// view options
	fs.BoolVar(&flagKeysOnly, "keys-only", false, "show only keys in output")
	fs.BoolVar(&flagOneLine, "one-line", false, "output one line per item")
}

func populatePISearchFilterOpts() processinstance.SearchFilterOpts {
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

func populatePDSearchFilterOpts() processdefinition.SearchFilterOpts {
	var opts processdefinition.SearchFilterOpts
	if flagKey != 0 {
		opts.Key = &flagKey
	}
	if flagBpmnProcessID != "" {
		opts.BpmnProcessId = &flagBpmnProcessID
	}
	if flagProcessVersion != 0 {
		opts.Version = &flagProcessVersion
	}
	if flagProcessVersionTag != "" {
		opts.VersionTag = &flagProcessVersionTag
	}
	return opts
}
