package oidc

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/go-jose/go-jose/v4"
	"github.com/luikyv/go-oidc/pkg/goidc"
)

func SignFunc(kmsClient *kms.Client, kmsSigKeyAlias string) goidc.SignFunc {
	return func(ctx context.Context, claims map[string]any, opts goidc.SignatureOptions) (string, error) {
		headerJSON, _ := json.Marshal(map[string]any{
			"alg": opts.Algorithm,
			"typ": opts.JWTType,
			"kid": kmsSigKeyAlias,
		})
		claimsJSON, _ := json.Marshal(claims)

		headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
		claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

		message := headerB64 + "." + claimsB64
		signOutput, err := kmsClient.Sign(ctx, &kms.SignInput{
			KeyId:            &kmsSigKeyAlias,
			Message:          []byte(message),
			MessageType:      types.MessageTypeRaw,
			SigningAlgorithm: signingAlg(opts.Algorithm),
		})
		if err != nil {
			return "", fmt.Errorf("failed to sign using KMS: %w", err)
		}

		signatureB64 := base64.RawURLEncoding.EncodeToString(signOutput.Signature)
		return headerB64 + "." + claimsB64 + "." + signatureB64, nil
	}
}

func signingAlg(alg jose.SignatureAlgorithm) types.SigningAlgorithmSpec {
	switch alg {
	case jose.ES256:
		return types.SigningAlgorithmSpecEcdsaSha256
	case jose.PS256:
		return types.SigningAlgorithmSpecRsassaPssSha256
	default:
		return types.SigningAlgorithmSpecRsassaPssSha256
	}
}

func DecryptFunc(kmsClient *kms.Client, kmsEncKeyAlias string) goidc.DecryptFunc {
	return func(ctx context.Context, jwe string, opts goidc.DecryptionOptions) (string, error) {
		jweParts := strings.Split(jwe, ".")

		encodedEncryptedKey := jweParts[1]
		encodedIV := jweParts[2]
		encodedCiphertext := jweParts[3]
		encodedTag := jweParts[4]

		encryptedKey, _ := base64.RawURLEncoding.DecodeString(encodedEncryptedKey)
		decryptOutput, err := kmsClient.Decrypt(ctx, &kms.DecryptInput{
			CiphertextBlob:      encryptedKey,
			KeyId:               &kmsEncKeyAlias,
			EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha1,
		})
		if err != nil {
			return "", fmt.Errorf("failed to decrypt CEK using KMS: %w", err)
		}

		iv, _ := base64.RawURLEncoding.DecodeString(encodedIV)
		ciphertext, _ := base64.RawURLEncoding.DecodeString(encodedCiphertext)
		authTag, _ := base64.RawURLEncoding.DecodeString(encodedTag)

		cek := decryptOutput.Plaintext
		fullCiphertext := append(ciphertext, authTag...)
		block, err := aes.NewCipher(cek)
		if err != nil {
			return "", fmt.Errorf("failed to create AES cipher block: %w", err)
		}

		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			return "", fmt.Errorf("failed to create AES-GCM: %w", err)
		}

		plaintext, err := aesgcm.Open(nil, iv, fullCiphertext, nil)
		if err != nil {
			return "", fmt.Errorf("failed to decrypt payload: %w", err)
		}

		return string(plaintext), nil
	}
}

func JWKSFunc(kmsClient *kms.Client, kmsSigKeyAlias, kmsEncKeyAlias string) goidc.JWKSFunc {
	return func(ctx context.Context) (jose.JSONWebKeySet, error) {
		signingJWK, err := fetchPublicKeyAsJWK(ctx, kmsClient, kmsSigKeyAlias, string(jose.PS256), string(goidc.KeyUsageSignature))
		if err != nil {
			return jose.JSONWebKeySet{}, fmt.Errorf("failed to fetch the signing key: %w", err)
		}

		encryptionJWK, err := fetchPublicKeyAsJWK(ctx, kmsClient, kmsEncKeyAlias, string(jose.RSA_OAEP), string(goidc.KeyUsageEncryption))
		if err != nil {
			return jose.JSONWebKeySet{}, fmt.Errorf("failed to fetch the encryption key: %w", err)
		}

		return jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{*signingJWK, *encryptionJWK},
		}, nil
	}
}

// fetchPublicKeyAsJWK retrieves the public key for a given alias and converts it to a JWK
func fetchPublicKeyAsJWK(
	ctx context.Context,
	kmsClient *kms.Client,
	alias, algorithm, keyUse string,
) (
	*jose.JSONWebKey,
	error,
) {

	kmsOutput, err := kmsClient.GetPublicKey(ctx, &kms.GetPublicKeyInput{
		KeyId: &alias,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	pub, err := x509.ParsePKIXPublicKey(kmsOutput.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return &jose.JSONWebKey{
		Key:       pub,
		KeyID:     alias,
		Algorithm: algorithm,
		Use:       keyUse,
	}, nil
}
