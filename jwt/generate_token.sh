#!/bin/bash

cd "$(dirname "$0")"
#list all supported curves
#openssl ecparam -list_curves

#we are using secp256k1 curve
#Generate an EC private key, of size 256, and output it to a file named private.pem
openssl ecparam -name prime256v1 -genkey -noout -out private_ec.pem

#Generate pkcs8 instead of SEC 1
openssl pkcs8 -topk8 -nocrypt -inform PEM -outform PEM -in private_ec.pem -out private.pem

#remove 
rm private_ec.pem

#Extract the public key from the key pair, which can be used in a certificate:
openssl ec -in private.pem -pubout -out public.pem

