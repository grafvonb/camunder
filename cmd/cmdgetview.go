package cmd

import (
	"fmt"
	"strings"

	"github.com/grafvonb/camunder/pkg/camunda/processdefinition"
	"github.com/grafvonb/camunder/pkg/camunda/processinstance"
	"github.com/spf13/cobra"
)

func listKeyOnlyProcessInstancesView(c *cobra.Command, resp processinstance.ProcessInstances) error {
	return renderListViewV(c, resp, func(r processinstance.ProcessInstances) []processinstance.ProcessInstance {
		return r.Items
	}, keyOnlyProcessInstanceView)
}

func listProcessInstancesView(c *cobra.Command, resp processinstance.ProcessInstances) error {
	if flagOneLine {
		return renderListViewV(c, resp, func(r processinstance.ProcessInstances) []processinstance.ProcessInstance {
			return r.Items
		}, oneLineProcessInstanceView)
	}
	return listJSONViewV(c, resp, func(r processinstance.ProcessInstances) []processinstance.ProcessInstance {
		return r.Items
	})
}

func keyOnlyProcessInstanceView(c *cobra.Command, item processinstance.ProcessInstance) error {
	c.Println(item.Key)
	return nil
}

func processInstanceView(c *cobra.Command, item processinstance.ProcessInstance) error {
	if flagOneLine {
		return oneLineProcessInstanceView(c, item)
	}
	if flagKeysOnly {
		return keyOnlyProcessInstanceView(c, item)
	}
	c.Println(ToJSONString(item))
	return nil
}

func oneLineProcessInstanceView(c *cobra.Command, item processinstance.ProcessInstance) error {
	var pTag, eTag, vTag string
	if item.ParentKey > 0 {
		pTag = fmt.Sprintf(" p:%d", item.ParentKey)
	} else {
		pTag = " p:<root>"
	}
	if item.EndDate != "" {
		eTag = fmt.Sprintf(" e:%s", item.EndDate)
	}
	if item.ProcessVersionTag != "" {
		vTag = "/" + item.ProcessVersionTag
	}

	out := fmt.Sprintf(
		"%-16d %s %s v%d%s %s s:%s%s%s i:%t",
		item.Key, item.TenantId, item.BpmnProcessId, item.ProcessVersion, vTag, item.State, item.StartDate, eTag, pTag, item.Incident,
	)
	c.Println(strings.TrimSpace(out))
	return nil
}

func listKeyOnlyProcessDefinitionsView(c *cobra.Command, resp processdefinition.ProcessDefinitions) error {
	return renderListViewV(c, resp, func(r processdefinition.ProcessDefinitions) []processdefinition.ProcessDefinition {
		return r.Items
	}, keyOnlyProcessDefinitionView)
}

func listProcessDefinitionsView(c *cobra.Command, resp processdefinition.ProcessDefinitions) error {
	if flagOneLine {
		return renderListViewV(c, resp, func(r processdefinition.ProcessDefinitions) []processdefinition.ProcessDefinition {
			return r.Items
		}, oneLineProcessDefinitionView)
	}
	return listJSONViewV(c, resp, func(r processdefinition.ProcessDefinitions) []processdefinition.ProcessDefinition {
		return r.Items
	})
}

func keyOnlyProcessDefinitionView(c *cobra.Command, item processdefinition.ProcessDefinition) error {
	c.Println(item.Key)
	return nil
}

func processDefinitionView(c *cobra.Command, item processdefinition.ProcessDefinition) error {
	if flagOneLine {
		return oneLineProcessDefinitionView(c, item)
	}
	if flagKeysOnly {
		return keyOnlyProcessDefinitionView(c, item)
	}
	c.Println(ToJSONString(item))
	return nil
}

func oneLineProcessDefinitionView(c *cobra.Command, item processdefinition.ProcessDefinition) error {
	vTag := ""
	if item.VersionTag != "" {
		vTag = "/" + item.VersionTag
	}
	out := fmt.Sprintf("%-16d %s %s v%s%s",
		item.Key, item.TenantId, item.BpmnProcessId, version, vTag,
	)
	c.Println(strings.TrimSpace(out))
	return nil
}

func listJSONView[Resp any, Item any](c *cobra.Command, resp *Resp, itemsOf func(*Resp) *[]Item) error {
	if resp == nil {
		c.Println("{}")
		return nil
	}
	printFound(c, itemsOf(resp))
	c.Println(ToJSONString(resp))
	return nil
}

func listJSONViewV[Resp any, Item any](c *cobra.Command, resp Resp, itemsOf func(Resp) []Item) error {
	items := itemsOf(resp)
	printFoundV(c, items)
	c.Println(ToJSONString(resp))
	return nil
}

func renderListView[Resp any, Item any](c *cobra.Command, resp *Resp, itemsOf func(*Resp) *[]Item,
	render func(*cobra.Command, *Item) error) error {
	if resp == nil {
		return nil
	}
	itemsPtr := itemsOf(resp)
	if itemsPtr == nil {
		c.Println("found: 0")
		return nil
	}
	items := *itemsPtr
	c.Println("found:", len(items))
	for i := range items {
		if err := render(c, &items[i]); err != nil {
			return err
		}
	}
	return nil
}

func renderListViewV[Resp any, Item any](c *cobra.Command, resp Resp, itemsOf func(Resp) []Item,
	render func(*cobra.Command, Item) error) error {
	items := itemsOf(resp)
	if len(items) == 0 {
		c.Println("found: 0")
		return nil
	}
	c.Println("found:", len(items))
	for _, it := range items {
		if err := render(c, it); err != nil {
			return err
		}
	}
	return nil
}

func printFound[T any](c *cobra.Command, items *[]T) {
	if items == nil {
		c.Println("found: 0")
		return
	}
	c.Println("found:", len(*items))
}

func printFoundV[T any](c *cobra.Command, items []T) {
	c.Println("found:", len(items))
}

func printFilter(c *cobra.Command) {
	var filters []string
	if flagParentKey != 0 {
		filters = append(filters, fmt.Sprintf("parent-key=%d", flagParentKey))
	}
	if flagState != "" && flagState != "all" {
		filters = append(filters, fmt.Sprintf("state=%s", flagState))
	}
	if flagParentsOnly {
		filters = append(filters, "parents-only=true")
	}
	if flagChildrenOnly {
		filters = append(filters, "children-only=true")
	}
	if flagOrphanParentsOnly {
		filters = append(filters, "orphan-parents-only=true")
	}
	if flagIncidentsOnly {
		filters = append(filters, "incidents-only=true")
	}
	if flagNoIncidentsOnly {
		filters = append(filters, "no-incidents-only=true")
	}
	if len(filters) > 0 {
		c.Println("filter: " + strings.Join(filters, ", "))
	}
}

func valueOr[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}
	return def
}
