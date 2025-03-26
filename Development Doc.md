# Fake_BiLiBiLi 开发文档
## 引言
本设计文档旨在描述fake_bilibili后端系统的整体架构和详细设计，为后续开发工作提供指导和参考。

## 系统概述
fake_bilibili后端系统是一个基于Web的视频播放直播平台，为用户提供视频播放、上次、收藏，发表评论，直播等功能。系统采用Websocket，后端使用Go语言开发。
技术栈为：Gin、Grom、

## 详细设计
### 用户模块
微信快捷登录、注册、账号密码登录、发送邮箱验证码、忘记密码
* 微信登录：[POST] 

## 相关逻辑
### 用户
#### 注册Register
1. 判断数据库是否已存在注册邮箱，已存在返回“邮箱已注册
2. 与Redis中保存该的用户对应邮箱验证码对比，判断邮箱验证码是否正确/过期，不正确返回“验证码错误/过期”
3. 密码利用密码盐加密
4. 从本地四个头像随机算一个保存
5. 保存本地数据库表
6. 生成token
> 密码盐加密逻辑：
> 1. 从密码盐字符串随机选出6位作为密码盐salt
> 2. md5(salt + password + salt) 整体进行Md5加密得到加密密码
> 3. 将加密密码和密码盐salt保存到数据库
> 
## 杂项知识
1. 跨域配置中的预检请求是什么？
> 预检请求（Preflight Request）是浏览器在发送某些 非简单跨域请求 之前，自动发起的一次 HTTP OPTIONS 请求。目的是向服务器确认是否允许实际的跨域请求。
## （暂时）参考文章
* [Go系列：结构体标签](https://juejin.cn/post/7005465902804123679#heading-17)
* [Go 基础系列：17. 详解 20 个占位符](https://zhuanlan.zhihu.com/p/415843240)
* [处理中文字符的rune类型](https://www.cnblogs.com/cheyunhua/p/16007219.html)
* [令牌桶]
* * [Go 基于令牌桶实现的官方限流器实际使用](https://blog.csdn.net/ic_xcc/article/details/120418426)
* * [Golang 标准库限流器 time/rate 实现剖析](https://www.cyhone.com/articles/analisys-of-golang-rate/)
## 疑点
### 数据库表
在定义部分数据库表的时候，外键的引用还要再另外创建一个“多余”的表结构。
比如在创建视频相关评论数据库表的时候，不直接引用`./infrastructure/model/contribution/video/video.go`的表结构来进行外键绑定
而是又在`./infrastructure/model/contribution/video/comments.go`中创建一个
`VideoInfo`表结构。说是为了解决“依赖循环”，这是什么意思？
> Q：可不可以不另外创建表结构来解决“依赖循环”的问题？感觉这样好冗余。
> 
> A：可以，不过这种创建中间表结构是最简单解决“依赖循环”的方式。
> 其他方式有：使用接口、依赖注入、拆分模块...
> >循环依赖：如果 Video 表也引用了 Comment 表，或者通过其他间接方式引用了 Comment 表，就会形成一个循环依赖。这种情况下，数据库迁移或代码编译时可能会失败，因为无法确定哪个表应该先被创建。

```go
type Comment struct {
	gorm.Model
	...
	
	UserInfo  user.User `json:"user_info" gorm:"foreignKey:Uid"`
	VideoInfo VideoInfo `json:"video_info" gorm:"foreignKey:VideoID"`
}

// VideoInfo 临时加一个video模型解决依赖循环
type VideoInfo struct {
	gorm.Model
	Uid   uint           `json:"uid" gorm:"uid"`
	Title string         `json:"title" gorm:"title"`
	Video datatypes.JSON `json:"video" gorm:"video"`
	Cover datatypes.JSON `json:"cover" gorm:"cover"`
}

```