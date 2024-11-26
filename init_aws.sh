#!/bin/bash

echo "Setting up Localstack..."

SIGNING_KEY_ID=$(awslocal kms create-key \
    --key-spec RSA_2048 \
    --key-usage SIGN_VERIFY \
    --description "LocalStack Signing KMS Key for PS256" \
    --query 'KeyMetadata.KeyId' --output text)

echo "KMS Signing Key ID: $SIGNING_KEY_ID"
awslocal kms get-public-key --key-id $SIGNING_KEY_ID --query 'PublicKey' --output text | \
  base64 --decode | \
  openssl rsa -pubin -inform DER -outform PEM

awslocal kms create-alias \
    --alias-name "alias/mockin/signing-key" \
    --target-key-id $SIGNING_KEY_ID

ENCRYPTION_KEY_ID=$(awslocal kms create-key \
    --key-spec RSA_2048 \
    --key-usage ENCRYPT_DECRYPT \
    --description "LocalStack KMS Key for RSA-OAEP encryption" \
    --query 'KeyMetadata.KeyId' --output text)

echo "KMS Encryption Key ID: $ENCRYPTION_KEY_ID"
awslocal kms get-public-key --key-id $ENCRYPTION_KEY_ID --query 'PublicKey' --output text | \
  base64 --decode | \
  openssl rsa -pubin -inform DER -outform PEM

awslocal kms create-alias \
    --alias-name "alias/mockin/encryption-key" \
    --target-key-id $ENCRYPTION_KEY_ID
