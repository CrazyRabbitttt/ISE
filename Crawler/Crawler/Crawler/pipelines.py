# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html


# useful for handling different item types with a single interface
from itemadapter import ItemAdapter
from Crawler.items import CrawlerItem
from Crawler.lru_cache import LRUCache
import requests
import json
from urllib.parse import urlparse, urlunparse, parse_qs, urlencode

def remove_query_params(url):
    parsed_url = urlparse(url)
    # 重建URL而不包含查询参数
    return urlunparse(parsed_url._replace(query=""))

class CrawlerPipeline(object):
    def __init__(self):
        self.counter = 199999    # 初始化计数器
        self.api_url = 'http://localhost:8080/api/addIndex?database=default'
        self.lrucache = LRUCache(capacity=1000) #初始化LRUcache，用于对于构建的URL的去重
        
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
                return 
        except Exception as e:
            print("error occur:", e)
        # URL去重
        print("准备开始更新LRU")
        if self.lrucache.get(remove_query_params(body['attrs']['page_url'])) != "-1":    #这个URL之前存在过
            self.lrucache.put(remove_query_params(body['attrs']['page_url']), "1")
            print("url has crawled, cancel send to server:", remove_query_params(body['attrs']['page_url']))
            return
        print("walawala.....")
        try:
            self.lrucache.put(remove_query_params(body['attrs']['page_url']), "1")
        except Exception as e:
            print("error:", e)
        headers = {'Content-Type': 'application/json'}
        print("running here")
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
    