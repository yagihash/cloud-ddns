# cloud-ddns
DDNS based on Cloud DNS

```
git clone git@github.com:yagihashoo/cloud-ddns.git
cp your_sa_file.json secret/sa.json

cp _env .env
vim .env

cp config/default.yaml.sample config/default.yaml
vim config/default.yaml

docker-compose up -d --build
```
