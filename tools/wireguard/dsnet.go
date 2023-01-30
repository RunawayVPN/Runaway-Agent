package wireguard

import (
	"strings"
	"time"

	_ "github.com/RunawayVPN/dsnet"
	"github.com/spf13/viper"
)

func init() {
	// Environment variable handling.
	viper.AutomaticEnv()
	viper.SetEnvPrefix("DSNET")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("config_file", "/etc/dsnetconfig.json")
	viper.SetDefault("fallback_wg_bing", "wireguard-go")
	viper.SetDefault("listen_port", 51820)
	viper.SetDefault("interface_name", "dsnet")

	// if last handshake (different from keepalive, see https://www.wireguard.com/protocol/)
	viper.SetDefault("peer_timeout", 3*time.Minute)

	// when is a peer considered gone forever? (could remove)
	// viper.SetDefault("peer_expiry", 28*time.Hour*24)
}
