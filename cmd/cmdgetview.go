package cmd

import (
	"fmt"
	"strings"

	c87operatev1 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v1"
	"github.com/spf13/cobra"
)

func listKeyOnlyProcessInstancesView(c *cobra.Command, resp *c87operatev1.ProcessInstanceSearchResponse) error {
	return renderListView(c, resp, func(r *c87operatev1.ProcessInstanceSearchResponse) *[]c87operatev1.ProcessInstanceItem {
		return r.Items
	}, keyOnlyProcessInstanceView)
}

func listProcessInstancesView(c *cobra.Command, resp *c87operatev1.ProcessInstanceSearchResponse) error {
	if flagOneLine {
		return renderListView(c, resp, func(r *c87operatev1.ProcessInstanceSearchResponse) *[]c87operatev1.ProcessInstanceItem {
			return r.Items
		}, oneLineProcessInstanceView)
	}
	return listJSONView(c, resp, func(r *c87operatev1.ProcessInstanceSearchResponse) *[]c87operatev1.ProcessInstanceItem {
		return r.Items
	})
}

func keyOnlyProcessInstanceView(c *cobra.Command, item *c87operatev1.ProcessInstanceItem) error {
	if item == nil {
		return nil
	}
	c.Println(valueOr(item.Key, int64(0)))
	return nil
}

func processInstanceView(c *cobra.Command, item *c87operatev1.ProcessInstanceItem) error {
	if flagOneLine {
		return oneLineProcessInstanceView(c, item)
	}
	if flagKeysOnly {
		return keyOnlyProcessInstanceView(c, item)
	}
	if item == nil {
		c.Println("{}")
		return nil
	}
	c.Println(ToJSONString(item))
	return nil
}

func oneLineProcessInstanceView(c *cobra.Command, item *c87operatev1.ProcessInstanceItem) error {
	if item == nil {
		return nil
	}

	key := valueOr(item.Key, int64(0))
	tenant := valueOr(item.TenantId, "")
	bpmnID := valueOr(item.BpmnProcessId, "")
	version := valueOr(item.ProcessVersion, int32(0))
	versionTag := valueOr(item.ProcessVersionTag, "")
	state := valueOr(item.State, "")
	start := valueOr(item.StartDate, "")
	end := valueOr(item.EndDate, "")
	parent := valueOr(item.ParentKey, int64(0))
	incident := valueOr(item.Incident, false)

	var pTag, eTag, vTag string
	if parent > 0 {
		pTag = fmt.Sprintf(" p:%d", parent)
	}
	if end != "" {
		eTag = fmt.Sprintf(" e:%s", end)
	}
	if versionTag != "" {
		vTag = "/" + versionTag
	}

	out := fmt.Sprintf(
		"%-16d %s %s v%d%s %s s:%s%s%s i:%t",
		key, tenant, bpmnID, version, vTag, state, start, eTag, pTag, incident,
	)
	c.Println(strings.TrimSpace(out))
	return nil
}

func listKeyOnlyProcessDefinitionsView(c *cobra.Command, resp *c87operatev1.ProcessDefinitionSearchResponse) error {
	return renderListView(c, resp, func(r *c87operatev1.ProcessDefinitionSearchResponse) *[]c87operatev1.ProcessDefinitionItem {
		return r.Items
	}, keyOnlyProcessDefinitionView)
}

func listProcessDefinitionsView(c *cobra.Command, resp *c87operatev1.ProcessDefinitionSearchResponse) error {
	if flagOneLine {
		return renderListView(c, resp, func(r *c87operatev1.ProcessDefinitionSearchResponse) *[]c87operatev1.ProcessDefinitionItem {
			return r.Items
		}, oneLineProcessDefinitionView)
	}
	return listJSONView(c, resp, func(r *c87operatev1.ProcessDefinitionSearchResponse) *[]c87operatev1.ProcessDefinitionItem {
		return r.Items
	})
}

func keyOnlyProcessDefinitionView(c *cobra.Command, item *c87operatev1.ProcessDefinitionItem) error {
	if item == nil {
		return nil
	}
	c.Println(valueOr(item.Key, int64(0)))
	return nil
}

func processDefinitionView(c *cobra.Command, item *c87operatev1.ProcessDefinitionItem) error {
	if flagOneLine {
		return oneLineProcessDefinitionView(c, item)
	}
	if flagKeysOnly {
		return keyOnlyProcessDefinitionView(c, item)
	}
	if item == nil {
		c.Println("{}")
		return nil
	}
	c.Println(ToJSONString(item))
	return nil
}

func oneLineProcessDefinitionView(c *cobra.Command, item *c87operatev1.ProcessDefinitionItem) error {
	if item == nil {
		return nil
	}

	key := valueOr(item.Key, int64(0))
	tenant := valueOr(item.TenantId, "")
	bpmnID := valueOr(item.BpmnProcessId, "")
	version := valueOr(item.Version, int32(0))
	versionTag := valueOr(item.VersionTag, "")

	vTag := ""
	if versionTag != "" {
		vTag = "/" + versionTag
	}

	out := fmt.Sprintf("%-16d %s %s v%d%s",
		key, tenant, bpmnID, version, vTag,
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

func printFound[T any](c *cobra.Command, items *[]T) {
	if items == nil {
		c.Println("found: 0")
		return
	}
	c.Println("found:", len(*items))
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
