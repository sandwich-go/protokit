# protokit

主要实现了对Google Protobuff文件的解析.

- 通过注释的方式扩展proto语义
- 将proto schema解析为ProfoFile，便于借助`Service`定义实现自定义的服务代码生成
- 为单个生成目标(Service,DB等)引入`ImportSet`概念，自动管理依赖

