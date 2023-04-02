## Openssl使用


### 生成私钥和公钥
```shell
$ openssl
OpenSSL> genrsa -out rsa_private_key.pem 1024  # 可调整长度为2048，生成的pem文件是pkcs#1的私钥文件
OpenSSL> openssl pkcs8 -topk8 -inform PEM -in rsa_private_key.pem -outform PEM -nocrypt > rsa_private_key_pkcs8.pem  # 如需要 可转为pkcs#8格式 
OpenSSL> rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem # 生成对应公钥
```

生成两个文件 `rsa_private_key.pem`, `rsa_public_key.pem`，都是PKCS1格式（带头尾），如下：
```shell
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCqx/yNundmEO646WHZ/BbxmUJt
2eLTq4DtrGfPRioEYGuCinlsdR7MtRYwY/lGe1qQpI75hi7RKlfpzuZ3FNmRe0tb
HQ9i5hczwtLA6lYixc8D8MEZZQ+Ch7PI7lr2zy7hUcaDEpMIlStTJ7IZworHiER6
59Qc3+opxLRmTMJX9QIDAQAB
-----END PUBLIC KEY-----
```

去掉头尾就是PKCS8格式：
```shell
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCqx/yNundmEO646WHZ/BbxmUJt
2eLTq4DtrGfPRioEYGuCinlsdR7MtRYwY/lGe1qQpI75hi7RKlfpzuZ3FNmRe0tb
HQ9i5hczwtLA6lYixc8D8MEZZQ+Ch7PI7lr2zy7hUcaDEpMIlStTJ7IZworHiER6
59Qc3+opxLRmTMJX9QIDAQAB
```