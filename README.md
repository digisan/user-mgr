## user-mgr

Lib for user sign-up, sign-in, sign-out, and reset password

go upgrade to 1.22.1
version => v0.9.2

### *Run command below first to create both private and public RSA keys*

> `openssl genrsa -out cert/id_rsa 4096`

> `openssl rsa -in cert/id_rsa -pubout -out cert/id_rsa.pub`