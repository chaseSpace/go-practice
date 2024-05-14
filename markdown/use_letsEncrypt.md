## Lets Encrypt免费证书

本文档记录在Centos上使用certbot来完成证书的申请和自动续期。


```shell
# 1. 安装certbot工具（参考网上）

# 2. 通过以下命令生成 DNS TXT 记录
sudo certbot certonly --manual --preferred-challenges dns -d example.com -d www.example.com
```
使用DNS txt验证域名所有权的优势是**支持颁发通配符证书**。

下面是一个实际的申请案例【域名：gamii.me】（Yes/No的输入直接参照即可）：
```shell
[centos@ip-172-31-42-96 ~]$ sudo certbot certonly --manual --preferred-challenges dns -d gamii.me
Saving debug log to /var/log/letsencrypt/letsencrypt.log
Plugins selected: Authenticator manual, Installer None
Enter email address (used for urgent renewal and security notices)
 (Enter 'c' to cancel): random2035@qq.com     # 输入你的联系邮箱 用于续期和安全通知
Starting new HTTPS connection (1): acme-v02.api.letsencrypt.org

- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
Please read the Terms of Service at
https://letsencrypt.org/documents/LE-SA-v1.3-September-21-2022.pdf. You must
agree in order to register with the ACME server. Do you agree?
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
(Y)es/(N)o: Y    # 同意服务条款

- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
Would you be willing, once your first certificate is successfully issued, to
share your email address with the Electronic Frontier Foundation, a founding
partner of the Let's Encrypt project and the non-profit organization that
develops Certbot? We'd like to send you email about our work encrypting the web,
EFF news, campaigns, and ways to support digital freedom.
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
(Y)es/(N)o: N   # 拒绝垃圾邮件
Account registered.
Requesting a certificate for gamii.me
Performing the following challenges:
dns-01 challenge for gamii.me

- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
Please deploy a DNS TXT record under the name
_acme-challenge.gamii.me with the following value:

 # 输出这一项时 需要去你域名托管商那里配置一条DNS的TXT记录
 # 类型：TXT
 # 名称：_acme-challenge
 # 值：下列字符串 （每次申请都会重新生成，续期则不需要）
 # TTL: 1h
vLnVoe0oKW67bDIQ4nBNtAvJC1OzLAujRUWNbcaXXNg        

Before continuing, verify the record is deployed.
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
# 这一步就会去读取你域名的TXT记录，若读取到并匹配无误，则颁发证书到本地
Press Enter to Continue
Waiting for verification...
Resetting dropped connection: acme-v02.api.letsencrypt.org
Cleaning up challenges

IMPORTANT NOTES:
 - Congratulations! Your certificate and chain have been saved at:
   /etc/letsencrypt/live/gamii.me/fullchain.pem
   Your key file has been saved at:
   /etc/letsencrypt/live/gamii.me/privkey.pem
   Your certificate will expire on 2023-08-03. To obtain a new or
   tweaked version of this certificate in the future, simply run
   certbot again. To non-interactively renew *all* of your
   certificates, run "certbot renew"
 - If you like Certbot, please consider supporting our work by:

   Donating to ISRG / Let's Encrypt:   https://letsencrypt.org/donate
   Donating to EFF:                    https://eff.org/donate-le
```

- 证书全链文件（可导出证书以及公钥）：`/etc/letsencrypt/live/gamii.me/fullchain.pem`
- 证书私钥：`/etc/letsencrypt/live/gamii.me/privkey.pem`

>fullchain.pem相当于是证书，私钥则是单独一个文件，一般不会合并到一个文件中。

下面是导出步骤：
```shell
# fullchain.pem 导出公钥
openssl x509 -in /etc/letsencrypt/live/gamii.me/fullchain.pem -pubkey -noout > public_key.pem

# fullchain.pem 导出证书
openssl x509 -in /etc/letsencrypt/live/gamii.me/fullchain.pem -out cert.crt

# privkey.pem 转为key后缀
openssl rsa -in /etc/letsencrypt/live/gamii.me/privkey.pem -outform PEM -out private_key.pem
```

还可以通过证书文件查看证书颁发机构和有效期：
```shell
[centos@ip-172-31-42-96 ~]$ openssl x509 -in cert.crt -text -noout
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            04:68:a4:19:3b:24:f1:87:45:a8:bc:b6:74:e1:39:8b:bf:ba
    Signature Algorithm: sha256WithRSAEncryption
        Issuer: C=US, O=Let's Encrypt, CN=R3
        Validity
            Not Before: May  5 01:48:45 2023 GMT   // 颁发时间
            Not After : Aug  3 01:48:44 2023 GMT   // 有效截止时间
        Subject: CN=gamii.me
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (2048 bit)
                Modulus:
                    00:be:24:37:85:51:ee:da:14:70:e4:45:7f:fb:ef:
                    1e:b9:68:57:15:75:47:41:fb:da:b9:de:a8:97:67:
                    82:0a:d8:e1:55:99:15:8c:7b:23:1c:6c:30:7e:7e:
                    14:52:23:a2:df:24:aa:9d:f2:db:b6:62:ef:d7:5e:
                    45:51:75:b4:01:62:62:29:6a:b2:6b:28:00:48:51:
                    d6:4b:4d:37:75:45:69:29:ee:e6:75:3a:c2:be:da:
                    c7:d8:d6:09:10:ab:7c:70:11:d4:45:8d:7e:69:95:
                    ff:f9:0d:0b:9d:6f:bc:0b:a8:5f:bc:f5:b9:ff:36:
                    36:a8:ec:ef:35:3f:2c:45:fb:66:cd:97:4b:c6:6f:
                    10:7c:6f:09:57:e8:de:d5:ea:e2:bf:74:8b:48:96:
                    3e:90:a3:f4:b4:be:75:42:ea:22:02:a5:ef:dc:c7:
                    66:40:05:f5:8e:98:8b:9c:4d:98:53:18:32:87:be:
                    a9:da:a8:12:f0:25:e6:a8:d9:18:54:6e:be:dd:ed:
                    39:7f:1e:4b:ad:c9:41:5a:da:53:1a:3d:ab:b6:24:
                    8a:78:73:aa:fb:cd:37:8a:b9:5f:23:34:da:91:1a:
                    f7:e3:42:f9:70:9b:98:19:87:8c:ab:56:9e:48:c7:
                    81:90:9f:b6:2d:d4:44:52:0e:cd:9d:5e:a9:8e:a9:
                    f1:d9
                Exponent: 65537 (0x10001)
                。。。省略部分
```


### 配置到Nginx
注意其中 root的路径设置为你原本的路径。
```shell
   server {
        listen       80 default_server;
        listen       [::]:80 default_server;
        server_name  gamii.me;
        return 301 https://$server_name$request_uri;
    }

    server {
        listen       443 ssl http2 default_server;
        listen       [::]:443 ssl http2 default_server;
        server_name  gamii.me www.gamii.me;

        ssl_protocols TLSv1 TLSv1.1 TLSv1.2 TLSv1.3;
        # 这是一个通用的密钥套件配置列表，能够兼容多数浏览器
        ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;

        ssl_certificate /etc/letsencrypt/live/gamii.me/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/gamii.me/privkey.pem;
        ssl_session_cache shared:SSL:1m;
        ssl_session_timeout  10m;
#         ssl_ciphers PROFILE=SYSTEM;
        ssl_prefer_server_ciphers on;

        # Load configuration files for the default server block.
        include /etc/nginx/default.d/*.conf;

        root   /home/centos/www

        location / {

        }

        error_page 404 /404.html;
            location = /40x.html {
        }

        error_page 500 502 503 504 /50x.html;
            location = /50x.html {
        }
    }
```

## 自动续期

### 1. 先修改初次申请证书时生成的conf文件
```shell
vi /etc/letsencrypt/renewal/gamii.me.conf

# 将 authenticator = manual 改为 authenticator = nginx
```

### 2. 测试手动更新
```shell
sudo certbot renew  --force-renew  # 刚颁发的证书进行续期 需要添加--force-renew选项，默认是过期前60天内才能续期
```
成功续期的输出如下：
```shell
[centos@ip-172-31-42-96 ~]$ sudo certbot renew  --force-renew
Saving debug log to /var/log/letsencrypt/letsencrypt.log

- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
Processing /etc/letsencrypt/renewal/gamii.me.conf
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
Plugins selected: Authenticator nginx, Installer None
Starting new HTTPS connection (1): acme-v02.api.letsencrypt.org
Renewing an existing certificate for gamii.me and *.gamii.me

- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
new certificate deployed without reload, fullchain is
/etc/letsencrypt/live/gamii.me/fullchain.pem
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
Congratulations, all renewals succeeded: 
  /etc/letsencrypt/live/gamii.me/fullchain.pem (success)
```

### 2. 设定cron定时任务
```shell
# 每周一0点执行（sleep一个5min内的时间）
echo "0 0 * * 1 root python -c 'import random; import time; time.sleep(random.random() * 300)' && certbot renew  --force-renew" | sudo tee -a /var/spool/cron/root > /dev/null

[centos@ip-172-31-42-96 ~]$ sudo crontab -l  # 查看任务
0 0 * * 1 root python -c 'import random; import time; time.sleep(random.random() * 300)' && certbot renew  --force-renew
```

[关于LetsEncrypt免费证书的频率限制](https://letsencrypt.org/zh-cn/docs/rate-limits/)

#### 查看执行情况
首先查看是否执行：`sudo cat /var/log/cron`

然后查看执行结果：`sudo cat /var/log/letsencrypt/letsencrypt.log`