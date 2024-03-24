# Defipublish_time=Nones for your scraped items
#
# See documentation in:
# https://docs.scrapy.org/en/latest/topics/items.html

import scrapy

class CrawlerItem(scrapy.Item):

    page_url = scrapy.Field() # 当前文档的URL
    title = scrapy.Field()      # 文档的标题
    content = scrapy.Field()    # 文档中的具体的内容
    # publish_time = scrapy.Field() # 文档的发布时间
    description = scrapy.Field() # 文档的描述信息
    keywords = scrapy.Field() # 文档的关键词
    urls = scrapy.Field() # 文档中包含的外部链接

    def gen_http_body(self):
        print("generate http body")
        data = {}
        # generated_terms = self['title'] + self['description'] + self['keywords']
        # term_list = self['title'] + self['description'] + self['keywords']
        # print("===title:", self['title'])
        # print("===description:", self['description'])
        # print("===page_url", self['page_url'])
        # print("===keywords", self['keywords'])
        # print("===urls", self['urls'])
        generated_term_list = self['description'] + self['keywords']
        generated_term_list.append(self['title'])
        generated_term = ' '.join(generated_term_list)
        
        print("generated term is:", generated_term, " page url is:", self['page_url'])
        # data['key'] =    ==> 生成的文档的唯一key值由 pipeline 中的函数来维护
        data['terms'] = generated_term
        # 填充用于排序、结果展示的 attrs
        if 'attrs' not in data:
            data['attrs'] = {}
        
        # data['attrs']['urls'] = self['urls']
        data['attrs']['title'] = self['title']
        data['attrs']['keywords'] = self['keywords']
        data['attrs']['description'] = self['description']
        # data['attrs']['publish_time'] = self['publish_time']
        data['attrs']['page_url'] = self['page_url']
        return data
        

        




