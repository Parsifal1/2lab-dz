package signData

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/fullsailor/pkcs7"
)

// signData - имя папки с файлами, которые нужно добавить в архив
func SignData(data []byte) (signed []byte, err error) {
	var signedData *pkcs7.SignedData

	if signedData, err = pkcs7.NewSignedData(data); err != nil {
		return
	}

	var cert tls.Certificate

	if len(cert.Certificate) == 0 {
		return nil, fmt.Errorf("Не удалось загрузить сертификат")
	}

	rsaKey := cert.PrivateKey
	var rsaCert *x509.Certificate

	if rsaCert, err = x509.ParseCertificate(cert.Certificate[0]); err != nil {
		return
	}

	if err = signedData.AddSigner(rsaCert, rsaKey, pkcs7.SignerInfoConfig{}); err != nil {
		return
	}

	return signedData.Finish()
}
