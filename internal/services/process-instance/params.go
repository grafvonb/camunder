package processinstance

type SearchFilterOpts struct {
	Key               *int64
	BpmnProcessId     *string
	ProcessVersion    *int32
	ProcessVersionTag *string
	State             PIState
	ParentKey         *int64
}
