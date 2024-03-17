# jpstudy-back
用自己的脚手架写的一个日语学习后端

---
## 环境配置
golang 1.21.3
mysql 8.0
protobuf的可执行文件直接一并上传

### 编译
运行build路径下的build.bat

### 运行
注意查看build/config里的config.json配置（涉及数据库密码未上传）。配置正确运行goserver.exe即可。如运行至linux服务器则修改build.bat成linux可执行文件

---
## 其他说明
详情表的数据库内容如有需要再导出一份excel。数据源也是从git上某个整理好的项目，有空时将相关地址补上；
思路是想像百词斩一样学日语，查词时会把翻译结果写入轻词本，轻词本记熟后就可以翻看详情中的词典；
翻译使用第三方的翻译api。可配合https://github.com/Spirild/jpstudy-front 前端部分使用。