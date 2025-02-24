package platform

// var net_platform NetworkPlatform

// message for get address request
// request for all the known address
// maybe we could handle nodes directly???
type GetAddress struct {
	AddrFrom   string
	message_id uint64
}

type GetNodes struct {
	AddrFrom string
	Address  []string
}

// send node object message
type NodesMessage struct {
	AddrFrom string
	Nodes    []NetworkNode // array to make is generic
}

type NodesInfo struct {
	AddrFrom string
	Nodes    []string // array to make is generic
}

// payload to send when we will be sending the data
type ConnectionRequest struct {
	AddrFrom  string
	ConnectId uint64
}

type ConnectionReply struct {
	AddrFrom string
	Node     NetworkNode // piggyback network Node
	IsReply  bool        // is this the reply to the reply or just the reply
}

type RequestMessage struct {
	AddrFrom string
	Id       uint64
}

type AckMessage struct {
	AddrFrom string
	Id       uint64
}

// a generic struct to request a information
// this struct can be used, so no parameter needs to be provided to the request
type GetInformation struct {
	AddrFrom string
}
