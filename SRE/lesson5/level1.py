import requests

# --- GET 请求 ---
url = "https://www.yfs365.com/"
response = requests.get(url)

print(f"状态码: {response.status_code}")
# 自动推测编码并解码内容
response.encoding = response.apparent_encoding  # 自动检测
print(f"响应内容 (Text): {response.text}") 
# 获取二进制内容 (如图片、文件)
# print(response.content) 