#!/bin/bash

set -o nounset \
    -o errexit \
    -o verbose \
    -o xtrace

# Generate CA key
openssl req -new -x509 -keyout soltesandbox-ca.key -out soltesandbox-ca.crt -days 365 -subj '/CN=ca.solte.sandbox.io/OU=TEST/O=SANDBOX/L=Hesinki/S=Hesinki/C=FI' -passin pass:confluent -passout pass:confluent

# Kafkacat
openssl genrsa -des3 -passout "pass:confluent" -out kafkacat.client.key 1024
openssl req -passin "pass:confluent" -passout "pass:confluent" -key kafkacat.client.key -new -out kafkacat.client.req -subj '/CN=ca.solte.sandbox.io/OU=TEST/O=SANDBOX/L=Hesinki/S=Hesinki/C=FI'
openssl x509 -req -CA soltesandbox-ca.crt -CAkey soltesandbox-ca.key -in kafkacat.client.req -out kafkacat-ca1-signed.pem -days 9999 -CAcreateserial -passin "pass:confluent"

for i in broker producer consumer
do
	echo $i
	# Create keystores
	keytool -genkey -noprompt \
				 -alias $i \
				 -dname "CN=$i.test.confluent.io, OU=TEST, O=CONFLUENT, L=PaloAlto, S=Ca, C=US" \
				 -keystore kafka.$i.keystore.jks \
				 -keyalg RSA \
				 -storepass confluent \
				 -keypass confluent

	# Create CSR, sign the key and import back into keystore
	keytool -keystore kafka.$i.keystore.jks -alias $i -certreq -file $i.csr -storepass confluent -keypass confluent

	openssl x509 -req -CA soltesandbox-ca.crt -CAkey soltesandbox-ca.key -in $i.csr -out $i-ca-signed.crt -days 9999 -CAcreateserial -passin pass:confluent

	keytool -keystore kafka.$i.keystore.jks -alias CARoot -import -file soltesandbox-ca.crt -storepass confluent -keypass confluent

	keytool -keystore kafka.$i.keystore.jks -alias $i -import -file $i-ca-signed.crt -storepass confluent -keypass confluent

	# Create truststore and import the CA cert.
	keytool -keystore kafka.$i.truststore.jks -alias CARoot -import -file soltesandbox-ca.crt -storepass confluent -keypass confluent

  echo "confluent" > ${i}_sslkey_creds
  echo "confluent" > ${i}_keystore_creds
  echo "confluent" > ${i}_truststore_creds
done

openssl genrsa -out go-client.key 2048

# Generate a CSR using the private key
openssl req -new -key go-client.key -out go-client.csr -subj "/CN=localhost/OU=TEST/O=SANDBOX/L=Hesinki/C=FI"

# Sign CSR with the CA cert, generate the client certificate
openssl x509 -req -in go-client.csr -CA soltesandbox-ca.crt -CAkey soltesandbox-ca.key -CAcreateserial -out go-client.crt -days 365 -passin pass:confluent

# Remove the CSR, we won't need it anymore
rm go-client.csr


find ./ -type f ! -name 'create_certs.sh' -exec mv {} ../certs/ \;