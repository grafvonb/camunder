package c87operatev1

func (r *ProcessInstanceSearchResponse) FilterChildrenOnly() *ProcessInstanceSearchResponse {
	return r.filterByParent(false)
}

func (r *ProcessInstanceSearchResponse) FilterParentsOnly() *ProcessInstanceSearchResponse {
	return r.filterByParent(true)
}

func (r *ProcessInstanceSearchResponse) filterByParent(hasParent bool) *ProcessInstanceSearchResponse {
	if r == nil || r.Items == nil {
		return r
	}
	for i := len(*r.Items) - 1; i >= 0; i-- {
		pk := (*r.Items)[i].ParentKey
		if (pk != nil) == hasParent {
			*r.Items = append((*r.Items)[:i], (*r.Items)[i+1:]...)
		}
	}
	nt := int64(len(*r.Items))
	r.Total = &nt
	return r
}
