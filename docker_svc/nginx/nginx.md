## Docker Nginx

```shell

mkdir /root/nginx
mkdir /root/nginx/conf.d
mkdir /root/nginx/html

# vi /root/nginx/conf.d/default.conf

chmod -R 755 /root/nginx/html

docker run --name nginx-admin -p 80:80 \
  -v /root/nginx/conf.d:/etc/nginx/conf.d \
  -v /root/nginx/html:/usr/share/nginx/html -d nginx
```