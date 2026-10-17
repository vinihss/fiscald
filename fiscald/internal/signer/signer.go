package signer

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/pkcs12"
)

type CertificadoA1 struct {
	PrivateKey *rsa.PrivateKey
	Cert       *x509.Certificate
}

func LoadCertificatePFX(pfxData []byte, password string) (*CertificadoA1, error) {
	privateKey, cert, err := pkcs12.Decode(pfxData, password)
	if err != nil {
		return nil, fmt.Errorf("falha ao decodificar arquivo PFX: %w", err)
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("chave privada do certificado nao e do tipo RSA")
	}

	return &CertificadoA1{
		PrivateKey: rsaKey,
		Cert:       cert,
	}, nil
}

func SignXML(xmlData []byte, certA1 *CertificadoA1) ([]byte, error) {
	hasher := sha1.New()
	hasher.Write(xmlData)
	digest := hasher.Sum(nil)
	digestBase64 := base64.StdEncoding.EncodeToString(digest)

	signedInfoXML := fmt.Sprintf(`<SignedInfo xmlns="http://www.w3.org/2000/09/xmldsig#"><CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/><SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"/><Reference URI="#NFe"><Transforms><Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/><Transform Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/></Transforms><DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/><DigestValue>%s</DigestValue></Reference></SignedInfo>`, digestBase64)

	siHasher := sha1.New()
	siHasher.Write([]byte(signedInfoXML))
	siDigest := siHasher.Sum(nil)

	signatureBytes, err := rsa.SignPKCS1v15(rand.Reader, certA1.PrivateKey, crypto.SHA1, siDigest)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar assinatura RSA: %w", err)
	}
	signatureBase64 := base64.StdEncoding.EncodeToString(signatureBytes)

	certBase64 := base64.StdEncoding.EncodeToString(certA1.Cert.Raw)

	signatureNode := fmt.Sprintf(`<Signature xmlns="http://www.w3.org/2000/09/xmldsig#">%s<SignatureValue>%s</SignatureValue><KeyInfo><X509Data><X509Certificate>%s</X509Certificate></X509Data></KeyInfo></Signature>`, signedInfoXML, signatureBase64, certBase64)

	signedXML := append(xmlData, []byte(signatureNode)...)
	return signedXML, nil
}
