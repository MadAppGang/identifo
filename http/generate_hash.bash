#!/bin/bash

#here is an example how to generate binary digest
echo -n $'test\nme' | openssl dgst -sha256 -hmac "secret"