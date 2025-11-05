from flask import Flask, send_from_directory
import os

app = Flask(__name__)

# 获取当前目录作为静态文件根目录
BASE_DIR = os.path.dirname(os.path.abspath(__file__))

@app.route('/<path:filename>', methods=['GET'])
def serve_static(filename):
    """
    映射所有静态文件到/static/接口
    例如: /static/index.html, /static/css/style.css等
    """
    return send_from_directory(BASE_DIR, filename)

@app.route('/', methods=['GET'])
def serve_index():
    """
    特别映射根路径到index.html
    """
    return send_from_directory(BASE_DIR, 'index.html')

if __name__ == '__main__':
    # 启动服务器，监听所有网络接口
    app.run(host='0.0.0.0', port=5000, debug=True)