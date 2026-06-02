package client

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type PartnerInfo struct {
	PartnerID     string
	Issuer        string
	PrivateKeyB64 string // PKCS#8 DER, base64-encoded (không phải PEM)
	Expire        time.Duration
}

func (p PartnerInfo) IsValid() bool {
	return p.PartnerID != "" && p.Issuer != "" && p.PrivateKeyB64 != "" && p.Expire > 0
}

func GenerateToken(partner PartnerInfo) (string, error) {
	if !partner.IsValid() {
		return "", errors.New("partner info invalid")
	}

	privateKey, err := parseRSAPrivateKeyPKCS8Base64(partner.PrivateKeyB64)
	if err != nil {
		return "", fmt.Errorf("generate token error: %w", err)
	}

	now := time.Now()
	exp := now.Add(partner.Expire)

	claims := jwt.MapClaims{
		"sub":       "partner",
		"iss":       partner.Issuer,
		"iat":       now.Unix(),
		"exp":       exp.Unix(),
		"partnerId": partner.PartnerID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token error: %w", err)
	}

	return signed, nil
}

func parseRSAPrivateKeyPKCS8Base64(b64 string) (*rsa.PrivateKey, error) {
	der, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("decode base64 private key error: %w", err)
	}
	keyAny, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, fmt.Errorf("parse PKCS#8 private key error: %w", err)
	}
	rsaPriv, ok := keyAny.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not RSA")
	}
	return rsaPriv, nil
}
