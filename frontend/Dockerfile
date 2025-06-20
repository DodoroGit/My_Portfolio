# 使用官方 Nginx 映像
FROM nginx:1.25

# 複製自訂的 nginx 設定
COPY nginx.conf /etc/nginx/nginx.conf

# 複製前端靜態檔案
COPY html/ /usr/share/nginx/html/

# 開放 port 80
EXPOSE 80

# 啟動 nginx
CMD ["nginx", "-g", "daemon off;"]
