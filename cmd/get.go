package cmd

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/cluster"
	"github.com/grafvonb/camunder/internal/services/common"
	processdefinition "github.com/grafvonb/camunder/internal/services/process-definition"
	processinstance "github.com/grafvonb/camunder/internal/services/process-instance"
	"github.com/spf13/cobra"
)

var supportedResourcesForGet = common.ResourceTypes{
	"ct": "cluster-topology",
	"pd": "process-definition",
	"pi": "process-instance",
}

var bpmnProcessId string

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [resource type]",
	Short: "List resources of a defined type. " + supportedResourcesForGet.PrettyString(),
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rn := strings.ToLower(args[0])
		cfg, err := config.FromContext(cmd.Context())
		if err != nil {
			cmd.PrintErrf("Error retrieving config from context: %v\n", err)
			return
		}
		timeout, err := time.ParseDuration(cfg.HTTP.Timeout)
		if err != nil {
			cmd.PrintErrf("Error parsing '%s' as timeout duration: %v\n", cfg.HTTP.Timeout, err)
			return
		}
		httpClient := &http.Client{
			Timeout: timeout,
		}
		switch rn {
		case "cluster-topology", "ct":
			svc, err := cluster.New(cfg, httpClient)
			if err != nil {
				cmd.PrintErrf("Error creating cluster service: %v\n", err)
				return
			}
			topology, err := svc.GetClusterTopology(cmd.Context())
			if err != nil {
				cmd.PrintErrf("Error fetching topology: %v\n", err)
				return
			}
			b, _ := json.MarshalIndent(topology, "", "  ")
			cmd.Println(string(b))
		case "process-definition", "pd":
			svc, err := processdefinition.New(cfg, httpClient)
			if err != nil {
				cmd.PrintErrf("Error creating process definition service: %v\n", err)
				return
			}
			pds, err := svc.SearchForProcessDefinitions(cmd.Context())
			if err != nil {
				cmd.PrintErrf("Error fetching process definitions: %v\n", err)
				return
			}
			b, _ := json.MarshalIndent(pds, "", "  ")
			cmd.Println(string(b))
		case "process-instance", "pi":
			if bpmnProcessId == "" {
				cmd.PrintErrln("Please provide a process ID to filter process instances using the --bpmn-process-id flag.")
				return
			}
			svc, err := processinstance.New(cfg, httpClient, isQuiet)
			if err != nil {
				cmd.PrintErrf("Error creating process instance service: %v\n", err)
				return
			}
			pis, err := svc.SearchForProcessInstances(cmd.Context(), bpmnProcessId)
			if err != nil {
				cmd.PrintErrf("Error fetching process instances: %v\n", err)
				return
			}
			b, _ := json.MarshalIndent(pis, "", "  ")
			cmd.Println(string(b))
		default:
			cmd.PrintErrf("Unknown resource type: %s\n", rn)
			cmd.Println(supportedResourcesForGet.PrettyString())
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&bpmnProcessId, "bpmn-process-id", "b", "", "BPMN process ID to filter process instances")
}
