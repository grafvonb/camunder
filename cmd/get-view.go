package cmd

import (
	"fmt"
	"strings"

	c87operatev1 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v1"
	"github.com/spf13/cobra"
)

func ListKeyOnlyProcessInstancesView(c *cobra.Command, resp *c87operatev1.ProcessInstanceSearchResponse) error {
	return renderListView(c, resp, func(r *c87operatev1.ProcessInstanceSearchResponse) *[]c87operatev1.ProcessInstanceItem {
		return r.Items
	}, KeyOnlyProcessInstanceView)
}

func ListProcessInstancesView(c *cobra.Command, resp *c87operatev1.ProcessInstanceSearchResponse) error {
	if flagOneLine {
		return renderListView(c, resp, func(r *c87operatev1.ProcessInstanceSearchResponse) *[]c87operatev1.ProcessInstanceItem {
			return r.Items
		}, OneLineProcessInstanceView)
	}
	return listJSONView(c, resp)
}

func KeyOnlyProcessInstanceView(c *cobra.Command, item *c87operatev1.ProcessInstanceItem) error {
	if item == nil {
		return nil
	}
	c.Println(valueOr(item.Key, int64(0)))
	return nil
}

func ProcessInstanceView(c *cobra.Command, item *c87operatev1.ProcessInstanceItem) error {
	if flagOneLine {
		return OneLineProcessInstanceView(c, item)
	}
	if flagKeysOnly {
		return KeyOnlyProcessInstanceView(c, item)
	}
	return listJSONView(c, item)
}

func OneLineProcessInstanceView(c *cobra.Command, item *c87operatev1.ProcessInstanceItem) error {
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
		"%-16d %s %s v%d%s %s s:%s%s%s",
		key, tenant, bpmnID, version, vTag, state, start, eTag, pTag,
	)
	c.Println(strings.TrimSpace(out))
	return nil
}

func ListKeyOnlyProcessDefinitionsView(c *cobra.Command, resp *c87operatev1.ProcessDefinitionSearchResponse) error {
	return renderListView(c, resp, func(r *c87operatev1.ProcessDefinitionSearchResponse) *[]c87operatev1.ProcessDefinitionItem {
		return r.Items
	}, KeyOnlyProcessDefinitionView)
}

func ListProcessDefinitionsView(c *cobra.Command, resp *c87operatev1.ProcessDefinitionSearchResponse) error {
	if flagOneLine {
		return renderListView(c, resp, func(r *c87operatev1.ProcessDefinitionSearchResponse) *[]c87operatev1.ProcessDefinitionItem {
			return r.Items
		}, OneLineProcessDefinitionView)
	}
	return listJSONView(c, resp)
}

func KeyOnlyProcessDefinitionView(c *cobra.Command, item *c87operatev1.ProcessDefinitionItem) error {
	if item == nil {
		return nil
	}
	c.Println(valueOr(item.Key, int64(0)))
	return nil
}

func ProcessDefinitionView(c *cobra.Command, item *c87operatev1.ProcessDefinitionItem) error {
	if flagOneLine {
		return OneLineProcessDefinitionView(c, item)
	}
	if flagKeysOnly {
		return KeyOnlyProcessDefinitionView(c, item)
	}
	return listJSONView(c, item)
}

func OneLineProcessDefinitionView(c *cobra.Command, item *c87operatev1.ProcessDefinitionItem) error {
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

func listJSONView[Resp any](c *cobra.Command, resp *Resp) error {
	if resp == nil {
		c.Println("{}")
		return nil
	}
	switch r := any(resp).(type) {
	case *c87operatev1.ProcessInstanceSearchResponse:
		if r.Items != nil {
			c.Println("found:", len(*r.Items))
		} else {
			c.Println("found: 0")
		}
	case *c87operatev1.ProcessDefinitionSearchResponse:
		if r.Items != nil {
			c.Println("found:", len(*r.Items))
		} else {
			c.Println("found: 0")
		}
	}
	c.Println(ToJSONString(resp))
	return nil
}

func renderListView[Resp any, Item any](c *cobra.Command, resp *Resp, itemsOf func(*Resp) *[]Item,
	render func(*cobra.Command, *Item) error) error {
	if resp == nil {
		c.Println("found: 0")
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

func valueOr[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}
	return def
}
