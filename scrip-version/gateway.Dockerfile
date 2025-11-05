# 使用Nginx作为网关服务
FROM nginx:alpine

# 删除默认的nginx配置文件
RUN rm /etc/nginx/conf.d/default.conf

# 复制自定义的nginx配置文件
COPY nginx.conf /etc/nginx/nginx.conf

# 设置工作目录
WORKDIR /usr/share/nginx/html

# 复制前端文件到容器中
COPY . .

# 暴露80端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]