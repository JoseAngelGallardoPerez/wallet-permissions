package service

import (
	"github.com/Confialink/wallet-permissions/internal/srvdiscovery"
	"context"
	"net/http"

	pb "github.com/Confialink/wallet-settings/rpc/proto/settings"
)

type CardModule struct {
	IsEnabled bool
}

// CardModuleSettings returns new CardModule from settings service or err if can not get it
func CardModuleSettings() (*CardModule, error) {
	settings := CardModule{}
	client, err := getSystemSettingsClient()
	if err != nil {
		return &settings, err
	}

	setting, err := client.Get(context.Background(), &pb.Request{Path: "regional/modules/velmie_wallet_cards"})
	if err != nil {
		return &settings, err
	}

	settings.IsEnabled = setting.Setting.Value == "enable"

	return &settings, nil
}

// getSystemSettingsClient returns rpc client to settings service
func getSystemSettingsClient() (pb.SettingsHandler, error) {
	settingsUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameSettings)
	if err != nil {
		return nil, err
	}
	return pb.NewSettingsHandlerProtobufClient(settingsUrl.String(), http.DefaultClient), nil
}
