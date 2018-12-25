package checkSzp

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/fullsailor/pkcs7"
)

func CheckSzp(szpLocation string, hash string) (*pkcs7.PKCS7, error) {
	var szp []byte
	var sign *pkcs7.PKCS7
	var err error

	if szp, err = ioutil.ReadFile(szpLocation); err != nil {
		return nil, err
	}

	if sign, err = pkcs7.Parse(szp); err != nil {
		return nil, err
	}

	err = sign.Verify()
	if err != nil {
		return nil, err
	}

	signer := sign.GetOnlySigner()
	if signer == nil {
		return nil, errors.New("Unable to obtain a single signer")
	}

	if hash != "UNDEF" {
		if hash != fmt.Sprintf("%x", sha1.Sum(signer.Raw)) {
			fmt.Println(fmt.Sprintf("%x", sha1.Sum(signer.Raw)))
			return nil, errors.New("ERROR: Certificate hash is corrupted")
		}
	}

	crt, err := tls.LoadX509KeyPair("./my.cer", "./my.key")
	if err != nil {
		return nil, err
	}

	parsedCrt, err := x509.ParseCertificate(crt.Certificate[0])
	if err != nil {
		return nil, err
	}

	if bytes.Compare(parsedCrt.Raw, signer.Raw) != 0 {
		return nil, errors.New("Certificates don't match")
	}
	return sign, nil
}
