package cert

import (
	"crypto"
	"fmt"
	"sifu-tool/models"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/registration"
	"go.uber.org/zap"
)

type acmeUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u acmeUser) GetEmail() string {
	return u.Email
}
func (u acmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u acmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}
func GetCert(email, cipher string, cfg map[string]string, domains []string, key crypto.PrivateKey, logger *zap.Logger) ( error) {
    user := acmeUser{
		Email: email,
		key:   key,
	}

	userCfg := lego.NewConfig(&user)
	userCfg.Certificate.KeyType = certcrypto.KeyType(certCrypto(cipher))
	
	client, err := lego.NewClient(userCfg)
	if err != nil {
		logger.Error(fmt.Sprintf("创建lego客户端失败: [%s]", err.Error()))
		return fmt.Errorf("创建lego客户端失败")
	}
	switch cfg[models.RESOLVER] {
		case "cloudflare":
			resolverCfg := cloudflare.NewDefaultConfig()
			resolverCfg.ZoneToken = cfg[models.CFTOKEN]
			cfResolver, err := cloudflare.NewDNSProviderConfig(resolverCfg)
			if err != nil {
				logger.Error(fmt.Sprintf(`配置"%s"失败: [%s]`, cfg[models.RESOLVER], err.Error()))
				return fmt.Errorf(`配置"%s"失败`, cfg[models.RESOLVER])
			}
			if err := client.Challenge.SetDNS01Provider(cfResolver); err != nil {
				logger.Error(fmt.Sprintf(`配置"%s"TXT失败: [%s]`, cfg[models.RESOLVER], err.Error()))
				return fmt.Errorf(`配置"%s"TXT失败`, cfg[models.RESOLVER])
			}
		
	}
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		logger.Error(fmt.Sprintf("注册let's encrypt账户失败: [%s]", err.Error()))
		return fmt.Errorf("注册let's encrypt账户失败")
	}
	fmt.Println(reg)
	
	request := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		logger.Error(fmt.Sprintf("获取证书失败: [%s]", err.Error()))
		return fmt.Errorf("获取证书失败")
	}
	fmt.Println(*certificates)
	return nil
}