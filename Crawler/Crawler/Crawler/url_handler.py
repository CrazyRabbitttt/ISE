# from urllib.parse import urlparse, urlunparse, parse_qs, urlencode

# def remove_query_params(url):
#     parsed_url = urlparse(url)
#     # 重建URL而不包含查询参数
#     return urlunparse(parsed_url._replace(query=""))


# url = "http://example.com/path/dwi.html"
# cleaned_url = remove_query_params(url)

# print(url)
# print(cleaned_url)

from datasketch import MinHash, MinHashLSH

# 定义一组示例字符串
texts = [
    "The quick brown fox jumps over the lazy dog",
    "The quick brown fox jumps over the lazy dog!",
    "A quick brown dog jumps over the lazy fox",
    "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
    "Lorem Ipsum is simply dummy text of the printing",
]

# 创建LSH对象，设置相似度阈值为0.8
lsh = MinHashLSH(threshold=0.8, num_perm=128)

# 创建一个字典来存储每个字符串的MinHash
minhashes = {}

for i, text in enumerate(texts):
    # 为每个字符串创建一个MinHash对象
    m = MinHash(num_perm=128)
    for word in text.split():
        m.update(word.encode('utf8'))
    # 将MinHash对象添加到LSH中
    lsh.insert(f"text{i}", m)
    minhashes[f"text{i}"] = m

# 检查每个字符串的相似项
unique_texts = []
for i, mh in minhashes.items():
    result = lsh.query(mh)
    # 只选择LSH查询结果中的第一个元素（如果存在的话）作为代表
    if i == result[0]:
        unique_texts.append(texts[int(i.replace("text", ""))])

# 打印去重后的字符串
print("Unique texts after LSH similarity check:")
for text in unique_texts:
    print(text)


