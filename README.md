# Cubox Archiver

我很喜欢 [Cubox]((https://cubox.pro/)) 稍后读的多端收集功能。

我没有太多其他需求（如标注），只是想在读完后点击「归档」，并把归档内容作为一个知识库，以后可以通过标题和描述快速检索。

但 Cubox 免费版限制总数据条数200，归档的内容也会占用容量，开会员的话一年要 98。

所以我写了这个程序，用来把 Cubox 的归档内容自动同步到其他地方（Notion等），然后再把 Cubox 的归档内容删除了。

## 特性

- 抽象化归档器，支持多种归档方式
- 流式处理
- 自动去重
- 同时支持配置文件和命令行参数
- 自动发版

## 归档器

* [x] Notion：自动创建数据库，并同步数据，自动去重。
* [ ] CSV 文件（写了一半，去重逻辑还没写，建议用于测试）
* [ ] 数据库

## 使用

先去 [Release 页](https://github.com/aFlyBird0/cubox-archiver/releases)下载二进制文件，Linux 和 macOS 用 `chmod +x` 给可执行权限。

如果使用 Notion 作为归档器（虽然现在只支持Notion），请先创建一个 Notion 机器人，然后选定一个页面把机器人 Connection 进来。

1. 运行一次程序，传入 page id，这时候会提示你已经自动创建了一个新数据库
2. 把新数据库的 database id 写到配置文件里，再运行，就真正启动了

配置相关往下看：

### 用文件传入配置

```bash
./cubox-archiver from-file -f config.yaml
```

具体配置看 `config.example.yaml`。

### 用命令行参数一个个传入配置

```bash
./cubox-archiver from-flag --help
```

然后根据提示自己拼参数。对于参数的详细解释，可以看 `config.example.yaml`。

## 文档没写的

* 怎么扒 Cubox 的 Auth 和 Cookie（应该挺简单）
* Notion 申请机器人和 Connection 的方法（看官网吧）

## 后续计划

* [ ] 支持更多归档器
* [ ] 做成 GitHub Action，自动定期运行
