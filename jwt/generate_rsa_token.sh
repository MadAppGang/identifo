
#!/bin/bash

ssh-keygen -t rsa -b 2048 -m PEM -f private.pem -C "identifo@madappgang.com" -N ""
rm private.pem.pub
openssl rsa -in private.pem -pubout -outform PEM -out public.pem