# 使用官方Python镜像作为基础
FROM python:3.9-slim

# 安装系统依赖（BeautifulSoup需要的解析库）
RUN apt-get update && \
    apt-get install -y --no-install-recommends gcc python3-dev && \
    rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /app

# 复制当前目录所有文件到容器的/app目录
COPY . .

# 安装Python依赖
RUN pip install --no-cache-dir \
    flask \
    beautifulsoup4 \
    requests

# 给脚本添加执行权限
RUN chmod +x ./py/RunForContests

# 暴露Flask应用的端口
EXPOSE 5000

# 启动命令（同时运行两个进程）
CMD bash -c "./py/RunForContests & python app.py"