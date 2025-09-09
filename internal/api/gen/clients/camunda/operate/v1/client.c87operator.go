package c87operatev1

func (r *ProcessInstanceSearchResponse) FilterChildrenOnly() *ProcessInstanceSearchResponse {
	if r == nil || r.Items == nil {
		return r
	}
	for i := len(*r.Items) - 1; i >= 0; i-- {
		if (*r.Items)[i].ParentKey == nil {
			*r.Items = append((*r.Items)[:i], (*r.Items)[i+1:]...)
		}
	}
	nt := int64(len(*r.Items))
	r.Total = &nt
	return r
}
