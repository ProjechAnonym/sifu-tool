package cert

import (
	"fmt"
	"sifu-tool/models"

	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
)

func resolverConfig(resolver string, config map[string]string) (any, error) {
	switch resolver {
		case "cloudflare":
			cfg := cloudflare.NewDefaultConfig()
			cfg.ZoneToken = config[models.CFTOKEN]
			return cfg, nil
	}
	return nil, fmt.Errorf("暂不支持%s", resolver)
}