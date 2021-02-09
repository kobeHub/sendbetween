package mdns

import (
	"context"
	"time"

	"github.com/kataras/golog"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// Initialize mdns discovery service
func InitMDNS(ctx context.Context, host host.Host, rendezvous string) chan peer.AddrInfo {
	srv, err := discovery.NewMdnsService(ctx, host, time.Hour, rendezvous)
	if err != nil {
		golog.Fatal(err)
	}

	// register with service so that we get notified about peer discovery
	n := &discoveryNotifee{
		PeerChan: make(chan peer.AddrInfo),
	}
	srv.RegisterNotifee(n)
	return n.PeerChan
}
