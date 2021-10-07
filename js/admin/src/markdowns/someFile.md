&nbsp;

#Its title
####How to generate RS256 private key (widely supported by all framwroks)
```sh
ssh-keygen -t rsa -b 2048 -m PEM -f private.pem -C "identifo@madappgang.com" -N ""
rm private.pem.pub
openssl rsa -in private.pem -pubout -outform PEM -out public.pem
```
#Its title

&nbsp;

####How to generate EC secp256k1  private key (RECOMMENDED)
```sh
openssl ecparam -name prime256v1 -genkey -noout -out private_ec.pem
openssl pkcs8 -topk8 -nocrypt -inform PEM -outform PEM -in private_ec.pem -out private.pem
rm private_ec.pem
openssl ec -in private.pem -pubout -out public.pem
```
