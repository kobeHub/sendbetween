package cfg

import (
	"fmt"

	"github.com/kataras/golog"
	"github.com/kobeHub/sendbetween/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetChatCfg() *ChatCfg {
	cfg := &ChatCfg{
		Address:    viper.GetString("address"),
		Dest:       viper.GetString("dest"),
		Rendezvous: viper.GetString("rendezvous"),
	}

	return cfg
}

func ParseBaseFlag(cmd *cobra.Command) {
	ipv4, err := utill.GetIp4ViaUDP()
	if err != nil {
		golog.Fatal("Get ip error:", err)
	}
	defaultAddr := fmt.Sprintf("/ip4/%s/tcp/0", ipv4)
	cmd.PersistentFlags().BoolP("debug", "v", false, "Enable debug mode, log more detail and generate same host ID")
	cmd.PersistentFlags().StringP("address", "a", defaultAddr, "Address that peer listen to, default random address")
	cmd.PersistentFlags().StringP("dest", "d", "", "The peer destination connect to, start client if not speficied")
	cmd.PersistentFlags().StringP("rendezvous", "r", "talkme", "Unique string to identify group of nodes.")

	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		golog.Fatalf("Err: %v\n", err)
	}
}
