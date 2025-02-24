package platform

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

/*
Initiate TCP Connection, creates connection and returns the connection
*/
func intiateTCPConnection(node *NetworkNode) (*net.Conn, error) {
	tcp_con, err := net.Dial("tcp", node.Socket.String())
	if err != nil {
		fmt.Println("Failed to initiate tcp connection with the host : ", node)
		return nil, err
	}
	return &tcp_con, err
}

/*
Connect to the node, to a specific address
Returns:

	error if any has occured
*/
func (net_platform *NetworkPlatform) ConnectToNode(address string) error {
	rand_num, err := rand.Prime(rand.Reader, 64)
	payload := ConnectionRequest{
		AddrFrom:  net_platform.Self_node.Socket.String(),
		ConnectId: rand_num.Uint64(),
	}
	fmt.Printf("%s\n", payload.AddrFrom)
	// connect to the network
	data := GobEncode(payload)
	data = append(CommandStringToBytes("connect"), data...)

	err = sendDataToAddress(address, data, net_platform)
	if err != nil {
		return err
	}
	SendTableToNode(net_platform, address)

	return nil
}

/*
Respond to connect request from a node
*/
func HandleConnectionInitiation(request []byte, net_platform *NetworkPlatform) error {
	var payload ConnectionRequest
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	send_payload := ConnectionReply{
		AddrFrom: net_platform.GetNodeAddress(),
		Node:     *net_platform.Self_node,
		IsReply:  false,
	}

	err := sendDataToAddress(payload.AddrFrom, append(CommandStringToBytes("connection_reply"), GobEncode(send_payload)...), net_platform)
	if err != nil {
		log.Printf("Connection initiation error: %s\n", err)
		return err
	}
	SendTableToNode(net_platform, payload.AddrFrom)
	time.Sleep(time.Second * 2)

	dispatch_pending_call(payload.AddrFrom)
	return nil
}

/*
Final step of connection reply:
  - If the reply is received from the node, finally connection is establised
*/
func HandleConnectionReply(request []byte, net_platform *NetworkPlatform) error {
	var payload ConnectionReply
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	net_platform.AddNode(payload.Node)
	SendTableToNode(net_platform, payload.AddrFrom)
	if !payload.IsReply {
		log.Println("Reply of a reply")
		// then a reply is recieved, reply with the self node information
		send_payload := ConnectionReply{
			AddrFrom: net_platform.GetNodeAddress(),
			Node:     *net_platform.Self_node,
			IsReply:  true,
		}
		err := sendDataToAddress(payload.AddrFrom, append(CommandStringToBytes("connection_reply"), GobEncode(send_payload)...), net_platform)

		if err != nil {
			return err
		}
	}

	if net_platform.config.ShareNodes {
		SendConnectedNodesAddress(payload.AddrFrom)
	}

	return nil
}

func SendConnectedNodesAddress(addr string) {
	var nodes_address []string
	for _, node := range network_platform.Connected_nodes {
		nodes_address = append(nodes_address, node.GetAddressString())
	}
	nodes := NodesInfo{
		network_platform.GetNodeAddress(),
		nodes_address,
	}

	sendDataToAddress(addr, append(CommandStringToBytes("nodes_to_connect"), GobEncode(nodes)...), network_platform)
}

func handleNodesInfoList(request []byte) error {
	var payload NodesInfo
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return err
	}

	for _, addr := range payload.Nodes {
		if network_platform.get_node_from_string(addr) == -1 {
			network_platform.ConnectToNode(addr)
		}
	}

	return nil
}
