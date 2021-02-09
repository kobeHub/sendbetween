package cfg

import (
	"github.com/kataras/golog"
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
	cmd.PersistentFlags().BoolP("debug", "v", false, "Enable debug mode, log more detail and generate same host ID")
	cmd.PersistentFlags().StringP("address", "a", "/ip4/127.0.0.1/tcp/0", "Address that peer listen to, default random address")
	cmd.PersistentFlags().StringP("dest", "d", "", "The peer destination connect to, start client if not speficied")
	cmd.PersistentFlags().StringP("rendezvous", "r", "talkme", "Unique string to identify group of nodes.")

	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		golog.Fatalf("Err: %v\n", err)
	}
}
