package c87operatev1

func (r *ProcessInstanceSearchResponse) FilterByHavingIncidents(has bool) *ProcessInstanceSearchResponse {
	return r.filterByIncident(has)
}

func (r *ProcessInstanceSearchResponse) FilterChildrenOnly() *ProcessInstanceSearchResponse {
	// children have a parent
	return r.filterByParent(true)
}

func (r *ProcessInstanceSearchResponse) FilterParentsOnly() *ProcessInstanceSearchResponse {
	// parents have no parent
	return r.filterByParent(false)
}

func (r *ProcessInstanceSearchResponse) filterByParent(hasParent bool) *ProcessInstanceSearchResponse {
	return r.filterByBool(func(pi *ProcessInstanceItem) bool {
		return pi.ParentKey != nil
	}, hasParent)
}

func (r *ProcessInstanceSearchResponse) filterByIncident(hasIncident bool) *ProcessInstanceSearchResponse {
	return r.filterByBool(func(pi *ProcessInstanceItem) bool {
		return pi.Incident != nil && *pi.Incident
	}, hasIncident)
}

func (r *ProcessInstanceSearchResponse) filterByBool(predicate func(*ProcessInstanceItem) bool, want bool) *ProcessInstanceSearchResponse {
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
