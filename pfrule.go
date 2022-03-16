package vbox

import (
	"fmt"
	"net"
)

// PFRule represents a port forwarding rule.
type PFRule struct {
	Proto     PFProto
	HostIP    net.IP // can be nil to match any host interface
	GuestIP   net.IP // can be nil if guest IP is leased from built-in DHCP
	HostPort  uint16
	GuestPort uint16
}

// PFProto represents the protocol of a port forwarding rule.
type PFProto string

const (
	// PFTCP when forwarding a TCP port.
	PFTCP = PFProto("tcp")
	// PFUDP when forwarding an UDP port.
	PFUDP = PFProto("udp")
)

// String returns a human-friendly representation of the port forwarding rule.
func (r PFRule) String() string {
	hostip, guestip := grab(r)
	return fmt.Sprintf("%s://%s:%d --> %s:%d",
		r.Proto, hostip, r.HostPort,
		guestip, r.GuestPort)
}

// Format returns the string needed as a command-line argument to VBoxManage.
func (r PFRule) Format() string {
	hostip, guestip := grab(r)
	return fmt.Sprintf("%s,%s,%d,%s,%d", r.Proto, hostip, r.HostPort, guestip, r.GuestPort)
}

func grab(r PFRule) (string, string) {
	hostip := ""
	if r.HostIP != nil {
		hostip = r.HostIP.String()
	}
	guestip := ""
	if r.GuestIP != nil {
		guestip = r.GuestIP.String()
	}
	return hostip, guestip
}
