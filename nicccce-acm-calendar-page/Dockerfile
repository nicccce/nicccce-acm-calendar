# 使用Node.js作为构建环境
FROM node:18-alpine as build

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json
COPY package*.json ./

# 安装依赖
RUN npm ci --only=production

# 复制源代码
COPY . .

# 构建React应用
RUN npm run build

# 使用Nginx作为生产环境
FROM nginx:alpine

# 删除默认的nginx配置文件
RUN rm /etc/nginx/conf.d/default.conf

# 创建一个空的default.conf文件以避免启动警告
RUN touch /etc/nginx/conf.d/default.conf

# 复制自定义的nginx配置文件
COPY nginx.conf /etc/nginx/nginx.conf

# 从构建阶段复制构建产物到nginx目录
COPY --from=build /app/build /usr/share/nginx/html

# 暴露80端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]