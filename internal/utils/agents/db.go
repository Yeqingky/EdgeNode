package agents

type DB interface {
	Init() error
	InsertAgentIP(ipId int64, ip string, agentCode string) error
	ListAgentIPs(offset int64, size int64) (agentIPs []*AgentIP, err error)
}
