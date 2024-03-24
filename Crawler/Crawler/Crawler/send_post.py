# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html


# useful for handling different item types with a single interface
# from Crawler.items import CrawlerItem
import requests
import json

def sendPost():
    api_url = 'http://localhost:8080/api/addIndex?database=default'
    data = {}
    data['key'] = 115656
    data['terms'] = "极客时间"
    if 'attrs' not in data:
        data['attrs'] = {}
    data['attrs']['title'] = "This is the title of test page_url"
    data['attrs']['number'] = 123456
    data['attrs']['page_url'] = "https://www.runoob.com/w3cnote/scrapy-detail.html"

    # json_data = json.dumps(data)
    headers = {'Content-Type': 'application/json'}
    try:
        response = requests.post(api_url, json=data, headers= headers)
    except Exception as e:
        print("发生了异常的现象:", e)

    print("服务端返回的响应数据:", response.text)
    if response.status_code == 200:
        print("请求发送成功")
    else:
        print("请求发送失败")

if __name__ == "__main__":
    sendPost()

# import requests

# # 定义要发送的JSON数据
# data = {
#     "key": 12399,
#     "terms": "大马哈鱼",
#     "attrs": {
#         "attribute1": "value1",
#         "attribute2": "value2",
#         "title": "这是一个关于大～马哈鱼的标题."
#     }
# }

# # 发送POST请求
# url = "http://localhost:8080/api/addIndex"
# response = requests.post(url, json=data)

# # 检查响应状态码
# if response.status_code == 200:
#     print("POST请求成功")
#     # 解析响应的JSON数据为IndexDoc结构
#     print("response:", response.text)
#     # print("Key:", index_doc["key"])
#     # print("Text:", index_doc["text"])
#     # print("Attrs:", index_doc["attrs"])
# else:
#     print("POST请求失败，状态码:", response.status_code)
