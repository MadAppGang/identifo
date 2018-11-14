
#!/bin/bash

ssh-keygen -t rsa -b 4096 -m PEM -f private.pem -C "identifo@madappgang.com" -N ""
mv private.pem.pub public.pem