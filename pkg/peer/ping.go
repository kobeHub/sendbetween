package peer

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kataras/golog"
	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Ping(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	node, err := libp2p.New(ctx,
		libp2p.ListenAddrStrings(viper.GetString("address")),
		libp2p.Ping(false),
	)
	if err != nil {
		golog.Error("Peer start failed: ", err)
		return err
	}
	golog.Info("Peer start listen to: ", node.Addrs())

	// configure customer ping protocol
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	// print node's peer info in multi format
	nodeInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}

	selfAddrs, err := peerstore.AddrInfoToP2pAddrs(&nodeInfo)
	if err != nil {
		golog.Error(err)
		return err
	}
	golog.Info("Slef libp2p address:", selfAddrs[0])

	peerAddrStr := viper.GetString("connect")
	if peerAddrStr == "" {
		golog.Info("Start to accept peer ping...")
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
		<-ch
	} else {
		peerAddr, err := multiaddr.NewMultiaddr(peerAddrStr)
		if err != nil {
			golog.Error(err)
			return err
		}
		peer, err := peerstore.AddrInfoFromP2pAddr(peerAddr)
		if err != nil {
			golog.Error(err)
			return err
		}
		golog.Info("Peer libp2p address:", peerAddr)

		golog.Info("Send 5 ping to peer...")
		ch := pingService.Ping(ctx, peer.ID)
		for i := 0; i < 5; i++ {
			res := <-ch
			golog.Warn("Got ping response", "RTT: ", res.RTT)
		}
	}

	// wait for exit singnal
	if err := node.Close(); err != nil {
		golog.Error(err)
		return err
	}
	return nil
}
