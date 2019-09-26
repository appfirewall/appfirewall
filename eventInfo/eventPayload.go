package eventInfo

import (
	"net"
	"time"

	"github.com/appfirewall/appfirewall/protocol"
	"github.com/p-/socket-connect-bpf/as"
)

// EventPayload contains information about the connection
// TODO refactor connectionInfo/EventInfo/EventPayload
type EventPayload struct {
	KernelTime    string
	GoTime        time.Time
	AddressFamily string
	Pid           uint32
	ProcessPath   string
	ProcessArgs   string // TODO change to array
	User          string
	UserID        uint32
	Comm          string
	Host          string
	DestIP        net.IP
	DestPort      uint16
	ASInfo        as.ASInfo
}

func (e *EventPayload) ToAFConnectionInfo() *protocol.AFConnectionInfo {
	connectionInfo := &protocol.AFConnectionInfo{
		DstIp:       e.DestIP.String(),
		DstHost:     e.Host,
		DstPort:     uint32(e.DestPort),
		DstAsNumber: e.ASInfo.AsNumber,
		DstAsName:   e.ASInfo.Name,
		UserId:      e.UserID,
		UserName:    e.User,
		ProcessId:   e.Pid,
		ProcessPath: e.ProcessPath,
		ProcessArgs: []string{e.ProcessArgs}, // TODO pass array
	}
	return connectionInfo
}
