package peer

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	mrand "math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/kataras/golog"
	"github.com/kobeHub/sendbetween/pkg/cfg"
	sendio "github.com/kobeHub/sendbetween/pkg/io"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	peerlib "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func handleStream(s network.Stream) {
	golog.Debug("Got new stream")

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)
}

func readData(rw *bufio.ReadWriter) {
	scanner := bufio.NewScanner(rw)
	scanner.Split(sendio.ScanMsg)
	for scanner.Scan() {
		item := scanner.Text()
		if item != "\n" {
			fmt.Printf("\x1b[32m%s\x1b[0m\n", item)
		}
		fmt.Print(">> ")
		if err := scanner.Err(); err != nil {
			golog.Error("Input invalid: ", err)
		}
	}
}

func writeData(rw *bufio.ReadWriter) {
	strReader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(strReader)
	scanner.Split(sendio.ScanMsg)

	for scanner.Scan() {
		item := scanner.Text()
		if err := scanner.Err(); err != nil {
			golog.Error("Received invalid: ", err)
			continue
		}

		fmt.Printf(">> ")
		if _, err := rw.WriteString(fmt.Sprintf("%s%s", item, sendio.MSG_DELIM)); err != nil {
			golog.Error("Write string to buffer error: ", err)
			continue
		}
		if err := rw.Flush(); err != nil {
			golog.Error(err)
			continue
		}
	}
}

func StartPeer(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	config := cfg.GetChatCfg()

	var r io.Reader
	nodeAddrStr := config.Address
	if viper.GetBool("debug") {
		nodeAddrParts := strings.Split(nodeAddrStr, "/")
		if len(nodeAddrParts) == 0 {
			golog.Error("Debug mode must spefic listen address")
			return nil
		}
		sp, err := strconv.ParseInt(nodeAddrParts[len(nodeAddrParts)-1], 10, 64)
		if err != nil {
			golog.Error("Debug mode must spefic listen port")
			return nil
		}
		r = mrand.New(mrand.NewSource(sp))
	} else {
		r = rand.Reader
	}

	// Generate RSA pair
	privKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		golog.Error("Create private key failed")
		return nil
	}

	node, err := libp2p.New(ctx,
		libp2p.ListenAddrStrings(nodeAddrStr),
		libp2p.Identity(privKey),
	)
	if err != nil {
		golog.Error("Self node start failed: ", err)
		return nil
	}
	golog.Debug("Self node listen address: ", node.Addrs())
	golog.Info("Welcome to sendbetween!")

	destAddrStr := config.Dest
	// Only start peer client
	if destAddrStr == "" {
		nodeInfo := peerlib.AddrInfo{
			ID:    node.ID(),
			Addrs: node.Addrs(),
		}
		selfAddrs, err := peerlib.AddrInfoToP2pAddrs(&nodeInfo)
		if err != nil {
			golog.Error(err)
			return nil
		}
		golog.Infof("Run ./sendbetween -d %s", selfAddrs[0])
		golog.Info("Waiting for incoming connection\n\n")

		// set stream hander
		node.SetStreamHandler("/send/1.0", handleStream)
		fmt.Print(">> ")

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		golog.Debug("Received signal, shutdowning...")
	} else {
		fmt.Println("The node's multiaddress:")
		for _, la := range node.Addrs() {
			fmt.Printf(" - %v\n", la)
		}
		fmt.Println()

		// Extra peer address info and add to node peerstore
		peerMaddr, err := multiaddr.NewMultiaddr(destAddrStr)
		if err != nil {
			golog.Fatal(err)
			return nil
		}
		peerInfo, err := peerlib.AddrInfoFromP2pAddr(peerMaddr)
		if err != nil {
			golog.Fatal(err)
			return nil
		}
		node.Peerstore().AddAddrs(peerInfo.ID, peerInfo.Addrs, peerstore.PermanentAddrTTL)

		// start a new stream with destination
		s, err := node.NewStream(ctx, peerInfo.ID, "/send/1.0")
		if err != nil {
			golog.Fatal(err)
			return nil
		}

		// Create a buffered stream so that read and writes are not bolcked
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		go writeData(rw)
		go readData(rw)
		fmt.Print(">> ")

		select {}
	}

	if err := node.Close(); err != nil {
		golog.Error(err)
		return err
	}
	return nil
}

func StartPeerWithMdns(cmd *cobra.Command, args []string) {
	return
}
