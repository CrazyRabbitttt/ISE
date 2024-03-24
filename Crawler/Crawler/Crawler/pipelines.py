# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html


# useful for handling different item types with a single interface
from itemadapter import ItemAdapter
from Crawler.items import CrawlerItem
import requests
import json

class CrawlerPipeline(object):
    def __init__(self):
        self.counter = 199999    # 初始化计数器
        self.api_url = 'http://localhost:8080/api/addIndex?database=default'
        
    def process_item(self, item: CrawlerItem, spider):
        # 1. 对于这个索引原始的数据生成一个唯一的 key、填充到 item 中
        # 2. 将传过来的body生成http请求发送到后端，构建索引
        # todo： 需要定向的对于热搜词进行抓取、构建热搜词Trie树
        print("轮到我pipeline对于item进行处理和执行了...")
        self.counter += 1
        docKey = self.counter
        body = item.gen_http_body()
        # print("需要处理的URL:", body['attrs']['page_url'])
        body['key'] = docKey
        try:
            if body['terms'] == "" or body['terms'] == ' ':
                print("len of terms is:", len(body["terms"]))
                return item
        except Exception as e:
            print("error occur:", e)
        # json_body = json.dumps(body)

        headers = {'Content-Type': 'application/json'}
        try:
            if body['attrs']['title'] != " " and body['attrs']['title'] != "":
                print("send post method 2 server")
                response = requests.post(self.api_url, json=body, headers=headers)
                print("post body:", body)
            else:
                print("Crawl nothing, do not send http request")
        except Exception as e:
            print("Post method error:",e)
        if response.status_code == 200:
            print(f"Send index failed, status code:{response.status_code}, text:{response.text}")
        else:
            print(f"Send index failed, status code:{response.status_code}, text:{response.text}")
        # return item
    