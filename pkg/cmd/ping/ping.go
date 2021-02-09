package peer

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kobeHub/sendbetween/pkg/peer"
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "ping",
		Short: "Connect to a peer of speficied p2p address",
		Long:  "Connect to a peer of speficied p2p address",
		RunE:  peer.Ping,
	}

	c.Flags().StringP("connect", "t", "", "The peer address")
	viper.BindPFlag("connect", c.Flags().Lookup("connect"))
	return c
}
