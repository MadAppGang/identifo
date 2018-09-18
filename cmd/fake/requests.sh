#!/bin/bash

curl -X POST \
  http://127.0.0.1:8080/auth/login \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'Digest: SHA-256=OPI53rqqR7HkmThF5/DG/+Sd4iM9ckot30l/eg5lggA=' \
  -H 'X-Identifo-ClientID: 12345' \
  -d '{"username": "test@madappgang.com","password": "secret","scope":  ["offline", "chat"]}'