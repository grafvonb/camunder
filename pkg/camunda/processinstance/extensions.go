package processinstance

func (r ProcessInstances) FilterByHavingIncidents(has bool) ProcessInstances {
	return r.filterByIncident(has)
}

func (r ProcessInstances) FilterChildrenOnly() ProcessInstances {
	// children have a parent
	return r.filterByParent(true)
}

func (r ProcessInstances) FilterParentsOnly() ProcessInstances {
	// parents have no parent
	return r.filterByParent(false)
}

func (r ProcessInstances) filterByParent(hasParent bool) ProcessInstances {
	return r.filterByBool(func(pi *ProcessInstance) bool {
		return pi.ParentKey > 0
	}, hasParent)
}

func (r ProcessInstances) filterByIncident(hasIncident bool) ProcessInstances {
	return r.filterByBool(func(pi *ProcessInstance) bool {
		return pi.Incident
	}, hasIncident)
}

func (r ProcessInstances) filterByBool(predicate func(*ProcessInstance) bool, want bool) ProcessInstances {
	if r.Items == nil {
		return r
	}
	items := r.Items
	for i := len(items) - 1; i >= 0; i-- {
		if predicate(&items[i]) != want {
			items = append(items[:i], items[i+1:]...)
		}
	}
	r.Total = int32(len(r.Items))
	return r
}
