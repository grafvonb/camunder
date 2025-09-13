package c87operate

func (r *ResultsProcessInstance) FilterByHavingIncidents(has bool) *ResultsProcessInstance {
	return r.filterByIncident(has)
}

func (r *ResultsProcessInstance) FilterChildrenOnly() *ResultsProcessInstance {
	// children have a parent
	return r.filterByParent(true)
}

func (r *ResultsProcessInstance) FilterParentsOnly() *ResultsProcessInstance {
	// parents have no parent
	return r.filterByParent(false)
}

func (r *ResultsProcessInstance) filterByParent(hasParent bool) *ResultsProcessInstance {
	return r.filterByBool(func(pi *ProcessInstance) bool {
		return pi.ParentKey != nil
	}, hasParent)
}

func (r *ResultsProcessInstance) filterByIncident(hasIncident bool) *ResultsProcessInstance {
	return r.filterByBool(func(pi *ProcessInstance) bool {
		return pi.Incident != nil && *pi.Incident
	}, hasIncident)
}

func (r *ResultsProcessInstance) filterByBool(predicate func(*ProcessInstance) bool, want bool) *ResultsProcessInstance {
	if r == nil || r.Items == nil {
		return r
	}
	items := *r.Items
	for i := len(items) - 1; i >= 0; i-- {
		if predicate(&items[i]) != want {
			items = append(items[:i], items[i+1:]...)
		}
	}
	*r.Items = items
	nt := int64(len(items))
	r.Total = &nt
	return r
}
