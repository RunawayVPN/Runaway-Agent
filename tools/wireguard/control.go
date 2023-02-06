package wireguard

import (
	dsnet "github.com/RunawayVPN/dsnet/cmd/cli"
	dsnet_utils "github.com/RunawayVPN/dsnet/utils"
)

// Up starts the wireguard interface if it is not already running
func Up() error {
	// Check report for error
	_, err := dsnet.GenerateReport()
	if err == nil {
		return nil
	}
	config, err := dsnet.MustLoadConfigFile()
	if err != nil {
		return err
	}
	server := dsnet.GetServer(config)
	if e := server.Up(); e != nil {
		return e
	}
	if e := dsnet_utils.ShellOut(config.PostUp, "PostUp"); e != nil {
		return e
	}
	return nil
}

func Down() error {
	config, err := dsnet.MustLoadConfigFile()
	if err != nil {
		return err
	}
	server := dsnet.GetServer(config)
	if e := server.DeleteLink(); e != nil {
		return e
	}
	if e := dsnet_utils.ShellOut(config.PostDown, "PostDown"); e != nil {
		return e
	}
	return nil
}
