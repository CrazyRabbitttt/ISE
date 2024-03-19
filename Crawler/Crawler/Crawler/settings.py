from fake_useragent import UserAgent
import time
import sys

BOT_NAME = "Crawler"

SPIDER_MODULES = ["Crawler.spiders"]
NEWSPIDER_MODULE = "Crawler.spiders"

ROBOTSTXT_OBEY = True

REQUEST_FINGERPRINTER_IMPLEMENTATION = "2.7"
TWISTED_REACTOR = "twisted.internet.asyncioreactor.AsyncioSelectorReactor"
FEED_EXPORT_ENCODING = "utf-8"


# user-agent : random user agent
USER_AGENT = UserAgent().random
# 设置 pipeline, 对于item怎么进行处理
ITEM_PIPELINES = {
    'Crawler.pipelines.CrawlerPipeline': 200
}
# 增加并发--Scrapy 下载器将执行的最大并发（即同时）请求数。
CONCURRENT_REQUESTS = 100
# 禁用 cookie
COOKIES_ENABLED = False
# 降低日志级别
LOG_LEVEL = 'INFO'
# the depth of crawler
DEPTH_LIMIT = 10
# log
LOG_FILE = "all.log"
