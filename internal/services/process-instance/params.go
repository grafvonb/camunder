package processinstance

type GetCmdFilterOpts struct {
	BpmnProcessId string
	State         PIState
	Key           *int64
}
