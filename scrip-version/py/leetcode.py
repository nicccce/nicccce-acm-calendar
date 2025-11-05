#!/usr/bin/env python3

# -*- coding: utf-8 -*-
"""
LeetCode比赛信息抓取脚本 - 优化保存路径版本
确保文件保存到与脚本同级的json文件夹
"""

import json
import requests
import os
import sys
import time
import logging
from pathlib import Path

# 配置常量
API_URL = "https://algcontest.rainng.com/contests"
OUTPUT_FILENAME = "leetcode-recent.json"

def setup_logging():
    """设置日志记录"""
    # 创建日志记录器
    logger = logging.getLogger('leetcode_contests')
    logger.setLevel(logging.INFO)
    
    # 创建控制台处理器
    console_handler = logging.StreamHandler()
    console_handler.setLevel(logging.INFO)
    
    # 创建格式化器
    formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
    console_handler.setFormatter(formatter)
    
    # 添加处理器
    logger.addHandler(console_handler)
    return logger

logger = setup_logging()

def get_output_path():
    """确定文件保存路径：同级目录的json文件夹"""
    try:
        # 获取当前脚本所在目录
        script_dir = Path(__file__).parent.resolve()
        logger.info(f"脚本所在目录: {script_dir}")
        
        # 创建json目录（如果不存在）
        json_dir = script_dir.parent / "json"
        json_dir.mkdir(parents=True, exist_ok=True)
        logger.info(f"JSON目录: {json_dir}")
        
        # 返回完整文件路径
        output_path = json_dir / OUTPUT_FILENAME
        logger.info(f"输出文件路径: {output_path}")
        return output_path
    except Exception as e:
        logger.error(f"确定输出路径时出错: {str(e)}")
        # 回退方案：保存到当前工作目录
        return Path.cwd() / OUTPUT_FILENAME

def fetch_leetcode_contests():
    """从API获取LeetCode比赛信息"""
    try:
        logger.info(f"开始请求API: {API_URL}")
        response = requests.get(API_URL, timeout=15)
        response.raise_for_status()
        
        contests_data = response.json()
        leetcode_contests = [
            contest for contest in contests_data 
            if contest.get("oj") == "LeetCode"
        ]
        
        logger.info(f"成功获取到 {len(leetcode_contests)} 场LeetCode比赛")
        return leetcode_contests
    except Exception as e:
        logger.error(f"获取比赛数据失败: {str(e)}")
        return []

def transform_contest_data(contests):
    """转换比赛数据格式"""
    transformed = []
    for contest in contests:
        try:
            # 计算持续时间（秒）
            duration = contest["endTimeStamp"] - contest["startTimeStamp"]
            
            # 创建标准格式字典
            transformed.append({
                "start_time": contest["startTimeStamp"],
                "contest_url": contest["link"],
                "name": contest["name"],
                "durationSeconds": duration,
                "end_time": contest["endTimeStamp"],
                "platform": "LeetCode"
            })
        except KeyError as e:
            logger.warning(f"比赛数据字段缺失: {e}, 跳过此比赛: {contest.get('name', '未知比赛')}")
        except Exception as e:
            logger.warning(f"处理比赛数据时出错: {e}, 跳过此比赛")
    return transformed

def save_contests_to_file(contests, filepath):
    """保存比赛数据到文件"""
    try:
        with open(filepath, 'w', encoding='utf-8') as f:
            json.dump(contests, f, indent=2, ensure_ascii=False)
        
        logger.info(f"比赛数据已保存到: {filepath}")
        
        # 返回成功状态
        return True, filepath
    except Exception as e:
        logger.error(f"保存文件失败: {str(e)}")
        return False, str(e)

def main():
    """主函数"""
    logger.info("===== 开始获取LeetCode比赛数据 =====")
    
    # 获取输出文件路径
    output_path = get_output_path()
    
    # 获取并处理比赛数据
    raw_contests = fetch_leetcode_contests()
    if not raw_contests:
        logger.error("没有获取到LeetCode比赛数据，退出程序")
        sys.exit(1)
    
    formatted_contests = transform_contest_data(raw_contests)
    
    # 保存数据
    success, result = save_contests_to_file(formatted_contests, output_path)
    
    if success:
        logger.info("===== 操作成功完成 =====")
        print(f"SUCCESS: 数据已保存到 {result}")
        sys.exit(0)
    else:
        logger.error("===== 操作失败 =====")
        print(f"ERROR: 文件保存失败 - {result}")
        sys.exit(1)

# if __name__ == "__main__":
#     main()
main()

# 致谢：https://www.rainng.com/alg-contest-info/