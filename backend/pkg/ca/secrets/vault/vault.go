package vault

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"

	"enroller/pkg/ca/secrets"

	"github.com/hashicorp/vault/api"
)

type vaultSecrets struct {
	client   *api.Client
	roleID   string
	secretID string
	CA       string
}

func NewVaultSecrets(address string, roleID string, secretID string) (*vaultSecrets, error) {
	conf := api.DefaultConfig()
	conf.Address = strings.ReplaceAll(conf.Address, "https://127.0.0.1:8200", address)
	tlsConf := &api.TLSConfig{Insecure: true}
	conf.ConfigureTLS(tlsConf)
	client, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}

	err = Login(client, roleID, secretID)
	if err != nil {
		return nil, err
	}
	return &vaultSecrets{client: client, roleID: roleID, secretID: secretID}, nil
}

func Login(client *api.Client, roleID string, secretID string) error {
	loginPath := "auth/approle/login"
	options := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}
	resp, err := client.Logical().Write(loginPath, options)
	if err != nil {
		return err
	}
	client.SetToken(resp.Auth.ClientToken)
	return nil
}

func (vs *vaultSecrets) GetCAs() (secrets.CAs, error) {
	resp, err := vs.client.Sys().ListMounts()
	if err != nil {
		return secrets.CAs{}, err
	}
	var CAs []secrets.CA
	for mount, mountOutput := range resp {
		if mountOutput.Type == "pki" {
			CAs = append(CAs, secrets.CA{Name: strings.TrimSuffix(mount, "/")})
		}
	}
	return secrets.CAs{CAs: CAs}, nil
}

func (vs *vaultSecrets) GetCAInfo(CA string) (secrets.CAInfo, error) {
	caPath := CA + "/cert/ca"
	resp, err := vs.client.Logical().Read(caPath)
	if err != nil {
		return secrets.CAInfo{}, err
	}
	pemBlock, _ := pem.Decode([]byte(resp.Data["certificate"].(string)))
	if pemBlock == nil {
		return secrets.CAInfo{}, errors.New("Cannot find the next PEM formatted block")
	}
	if pemBlock.Type != "CERTIFICATE" || len(pemBlock.Headers) != 0 {
		return secrets.CAInfo{}, errors.New("Unmatched type of headers")
	}
	caCert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return secrets.CAInfo{}, err
	}
	key := caCert.PublicKeyAlgorithm.String()
	var keyBits int
	switch key {
	case "RSA":
		keyBits = caCert.PublicKey.(*rsa.PublicKey).N.BitLen()
	case "ECDSA":
		keyBits = caCert.PublicKey.(*ecdsa.PublicKey).Params().BitSize
	}

	CAInfo := secrets.CAInfo{
		C:       strings.Join(caCert.Subject.Country, " "),
		L:       strings.Join(caCert.Subject.Locality, " "),
		O:       strings.Join(caCert.Subject.Organization, " "),
		ST:      strings.Join(caCert.Subject.Province, " "),
		CN:      caCert.Subject.CommonName,
		KeyType: key,
		KeyBits: keyBits,
	}

	return CAInfo, nil

}

func (vs *vaultSecrets) DeleteCA(CA string) error {
	deletePath := CA + "/root"
	_, err := vs.client.Logical().Delete(deletePath)
	if err != nil {
		return err
	}
	return nil
}
