import sys
import re
import scrapy
from urllib import parse
from Crawler.items import CrawlerItem
from gerapy_auto_extractor.extractors.datetime import extract_datetime
from ..html_extractor import MainContent

# 对于静态文件进行过滤操作，对于文件类别就不再进行爬取
FILE_WORDS = ['.gif','.png','.bmp','.jpeg','.jpg', '.svg',
              '.mp3','.wma','.flv','.mp4','.wmv','.ogg','.avi',
              '.doc','.docx','.xls','.xlsx','.ppt','.pptx','.txt','.pdf',
              '.zip','.exe','.tat','.ico','.css','.js','.swf','.apk','.m3u8','.ts']

# 还有就是如果包含很长一串数字的一般都是内容界面
# 不应该是各种文件名的后缀
def is_static_url(url):
    '''
    
    '''
    for w in FILE_WORDS:
        if w in url[-5:]:
            return True

    return False

class WebSpider(scrapy.Spider):
    name = 'webSpider'
    start_urls = ['https://tuijian.hao123.com/']

    # 抓取的规则
    rule_encode = "//meta/@charset"
    rule_keywords = "//meta[@name='keywords']/@content"
    rule_description = "//meta[@name='description']/@content"
    rule_lang = "//@lang"
    rule_url = "//@href"  # 简答地提取url的规则
    def parse(self, response):
        page_url = response.request.url
        print("-"*100)
        # print("开始爬取%s......" % page_url)
        if response.status == 200:
            # 获取内容
            encode = response.xpath(self.rule_encode).extract()
            keywords = response.xpath(self.rule_keywords).extract()
            description = response.xpath(self.rule_description).extract()
            lang = response.xpath(self.rule_lang).extract()
            # 这里代码检测有问题，实际没问题，只能说VS有点垃圾，继承关系都搞不懂
            # publish_time = extract_datetime(response.body.decode('utf-8'))

            urls = response.xpath(self.rule_url).extract()
            urls_cleaned = []
            for url in urls:        # 所有的外部链接
                # 如果说链接中的后缀是 static 类型的，不予查找
                if is_static_url(url) or "javascript:" in url.lower():
                    continue
                # 绝对链接不变，相对链接转换为绝对链接
                full_url = parse.urljoin(page_url, url)
                urls_cleaned.append(full_url)

            # 如果符合详情页规则，就下载该网页,提取其正文
            # print("该网页符合详情页规则.....")
            # print("提取[ %s ]携带的正文标题中......" % page_url)

            extractor = MainContent()
            title = extractor.extract(page_url, response.body)

            # 保存...
            detail_item = CrawlerItem()
            detail_item['page_url'] = page_url
            detail_item['keywords'] = keywords
            detail_item['description'] = description
            detail_item['title'] = title
            detail_item['urls'] = urls_cleaned
            # detail_item['publish_time'] = publish_time
            print("title:", title, " page_url:", page_url)

            yield detail_item

            for url in urls_cleaned:
                yield scrapy.Request(url=url, callback=self.parse)
        else:
            print("[ %s ]未爬取成功......" % page_url)
            return











