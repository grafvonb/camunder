package cmd

import (
	"fmt"
	"strings"

	c87operatev1 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v1"
	"github.com/spf13/cobra"
)

func ListItemsView(c *cobra.Command, resp *c87operatev1.ProcessInstanceSearchResponse) error {
	c.Println("found:", len(*resp.Items))
	if flagOneLine {
		return listOnelineItemsView(c, resp)
	}
	return listJsonItemsView(c, resp)
}

func ListKeyOnlyItemsView(c *cobra.Command, resp *c87operatev1.ProcessInstanceSearchResponse) error {
	c.Println("found:", len(*resp.Items))
	for i := range *resp.Items {
		if err := KeyOnlyItemView(c, &(*resp.Items)[i]); err != nil {
			return err
		}
	}
	return nil
}

func listOnelineItemsView(c *cobra.Command, resp *c87operatev1.ProcessInstanceSearchResponse) error {
	for i := range *resp.Items {
		if err := OneLineItemView(c, &(*resp.Items)[i]); err != nil {
			return err
		}
	}
	return nil
}

func listJsonItemsView(c *cobra.Command, resp *c87operatev1.ProcessInstanceSearchResponse) error {
	c.Println(ToJSONString(resp))
	return nil
}

func OneLineItemView(c *cobra.Command, item *c87operatev1.ProcessInstanceItem) error {
	p := valueOr(item.ParentKey, int64(0))
	ps := ""
	if p > 0 {
		ps = fmt.Sprintf(" p:%d", p)
	}
	eds := valueOr(item.EndDate, "")
	if eds != "" {
		eds = fmt.Sprintf(" e:%s", eds)
	}
	pvt := valueOr(item.ProcessVersionTag, "")
	if pvt != "" {
		pvt = fmt.Sprintf("/%s", pvt)
	}
	out := fmt.Sprintf("%-16d %s %s v%d%s %s s:%s%s%s",
		*item.Key,
		*item.TenantId,
		*item.BpmnProcessId,
		*item.ProcessVersion,
		pvt,
		*item.State,
		*item.StartDate,
		eds,
		ps,
	)
	c.Println(strings.TrimSpace(out))
	return nil
}

func KeyOnlyItemView(c *cobra.Command, item *c87operatev1.ProcessInstanceItem) error {
	c.Println(item.Key)
	return nil
}

func valueOr[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}
	return def
}
