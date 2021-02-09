package sendbetween

import (
	"github.com/kataras/golog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kobeHub/sendbetween/pkg/cfg"
	pingcli "github.com/kobeHub/sendbetween/pkg/cmd/ping"
	"github.com/kobeHub/sendbetween/pkg/peer"
)

func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "sendbetween",
		Short:            "send",
		Long:             `An easy cli tool to send text and files between different PCs across OS`,
		TraverseChildren: true,
		PreRun:           Init,
		RunE:             peer.StartPeer,
	}
	rootCmd.AddCommand(pingcli.NewCommand())
	cfg.ParseBaseFlag(rootCmd)
	return rootCmd
}

// Initialize command
func Init(cmd *cobra.Command, args []string) {
	if viper.GetBool("debug") {
		golog.SetLevel("debug")
	} else {
		golog.SetLevel("info")
	}
	golog.Debug("sendbetween initialize...")
}
