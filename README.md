# protokit

主要实现了对Google Protobuff文件的解析.

- 通过注释的方式扩展proto语义
- 将proto schema解析为ProfoFile，便于借助`Service`定义实现自定义的服务代码生成
- 为单个生成目标(Service,DB等)引入`ImportSet`概念，自动管理依赖

通用 message option 定义在 `protokit/misc.proto`：

```proto
import "protokit/misc.proto";

message ServerOnly {
  option (protokit.message_opts).no_csharp = true;
}
```

`no_csharp` 默认 `false`。解析后通过
`message.GetCommonMessageOption().GetNoCsharp()` 读取，由生成器消费该标记；
该 option 本身不会改变原生 protoc 的 C# 输出。
