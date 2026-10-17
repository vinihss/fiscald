package sefaz

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"fiscald/internal/signer"
)

type SEFAZClient struct {
	httpClient *http.Client
	urlWS      string
}

func NewSEFAZClient(certA1 *signer.CertificadoA1, urlWS string, timeout time.Duration) *SEFAZClient {
	tlsCert := tls.Certificate{
		Certificate: [][]byte{certA1.Cert.Raw},
		PrivateKey:  certA1.PrivateKey,
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{tlsCert},
		InsecureSkipVerify: false,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	return &SEFAZClient{
		httpClient: client,
		urlWS:      urlWS,
	}
}

func (s *SEFAZClient) TransmitirNFCe(ctx context.Context, xmlAssinado []byte) ([]byte, error) {
	soapEnvelope := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?><soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope"><soap12:Header><nfeCabecMsg xmlns="http://www.portalfiscal.inf.br/nfe/wsdl/NFeAutorizacao4"><cUF>35</cUF><versaoDados>4.00</versaoDados></nfeCabecMsg></soap12:Header><soap12:Body><nfeDadosMsg xmlns="http://www.portalfiscal.inf.br/nfe/wsdl/NFeAutorizacao4">%s</nfeDadosMsg></soap12:Body></soap12:Envelope>`, string(xmlAssinado))

	req, err := http.NewRequestWithContext(ctx, "POST", s.urlWS, bytes.NewBufferString(soapEnvelope))
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisicao HTTP: %w", err)
	}

	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha no envio mTLS para SEFAZ: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sefaz retornou HTTP STATUS %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta do corpo SOAP: %w", err)
	}

	return body, nil
}
