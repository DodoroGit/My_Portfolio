FROM nginx:1.25

# 建立 certbot 驗證目錄
RUN mkdir -p /var/www/certbot

# 複製 nginx 設定
COPY nginx/nginx.conf /etc/nginx/nginx.conf
COPY nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf

# 靜態檔案
COPY html/ /usr/share/nginx/html/

# 開放 HTTP/HTTPS 端口
EXPOSE 80
EXPOSE 443

# 提供 volume 讓 certbot 掛載 SSL 憑證
VOLUME [ "/etc/letsencrypt", "/var/www/certbot" ]

# 啟動 nginx
CMD ["nginx", "-g", "daemon off;"]
