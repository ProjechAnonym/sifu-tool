package controller

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"sifu-tool/cert"
	"sifu-tool/ent"
	"sifu-tool/models"

	"go.uber.org/zap"
)

func AddCertJob(form models.AcmeForm, key string, entClient *ent.Client, logger *zap.Logger) (error) {

	if form.Auto {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			logger.Error(fmt.Sprintf("生成私钥失败: [%s]", err.Error()))
			return fmt.Errorf("生成私钥失败")
		}
		privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			logger.Error(fmt.Sprintf("存储私钥字节失败: [%s]", err.Error()))
			return fmt.Errorf("存储私钥字节失败")
		}
		
		if _, err := cert.Encrypt(privateKeyBytes, key); err != nil {
			logger.Error(fmt.Sprintf("加密私钥失败: [%s]", err.Error()))
			return fmt.Errorf("加密私钥失败")
		}
		
		
		
	}
	return nil
}