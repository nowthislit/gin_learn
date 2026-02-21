# Gin 框架渐进式学习路线

本项目提供完整的 Gin 框架学习路径，从基础到实战，分四个阶段渐进式学习。

## 📚 学习路线图

| 阶段 | 名称 | 难度 | 预计时间 | 状态 |
|------|------|------|----------|------|
| Phase 1 | 基础入门 | ⭐⭐ | 3-5天 | ✅ 已完成 |
| Phase 2 | 中间件与验证 | ⭐⭐⭐ | 5-7天 | ✅ 已完成 |
| Phase 3 | 数据库集成 | ⭐⭐⭐ | 5-7天 | ⏳ 未开始 |
| Phase 4 | 进阶实战 | ⭐⭐⭐⭐ | 7-10天 | ⏳ 未开始 |

## 📂 项目结构

```
gin-learn/
├── README.md                    # 本文件 - 学习路线图
├── phase1_basic/               # 第一阶段：基础入门
│   ├── go.mod
│   ├── go.sum
│   ├── main.go                 # 完整示例代码
│   └── README.md               # 详细知识点
├── phase2_middleware/          # 第二阶段：中间件与验证
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── README.md
├── phase3_database/            # 第三阶段：数据库集成
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── test.db                 # 自动生成的SQLite数据库
│   └── README.md
└── phase4_advanced/            # 第四阶段：进阶实战
    ├── go.mod
    ├── go.sum
    ├── main.go                 # 项目入口
    ├── config.yaml             # 配置文件示例
    ├── config/                 # 配置管理
    ├── internal/               # 内部代码
    │   ├── api/               # API层
    │   ├── service/           # 业务逻辑层
    │   ├── repository/        # 数据访问层
    │   ├── model/             # 数据模型
    │   └── middleware/        # 中间件
    ├── pkg/                    # 公共包
    │   └── logger/            # 日志工具
    └── README.md
```

## 🎯 学习目标概览

### Phase 1: 基础入门
- ✅ Gin引擎的创建和配置
- ✅ HTTP方法路由（GET/POST/PUT/DELETE）
- ✅ 路由参数（:param, *path）
- ✅ 路由分组（Group）
- ✅ Query参数、Form参数、JSON绑定
- ✅ 多种响应格式（JSON/XML/String/HTML）

### Phase 2: 中间件与验证
- ✅ 中间件的概念和执行流程
- ✅ 内置中间件（Logger、Recovery）
- ✅ 自定义中间件开发
- ✅ 全局/分组/单路由中间件
- ✅ 请求数据验证（binding tags）
- ✅ 自定义验证规则
- ✅ 全局错误处理

### Phase 3: 数据库集成
- ✅ GORM安装和配置
- ✅ 数据模型定义
- ✅ CRUD操作
- ✅ 关联关系（BelongsTo/HasOne/HasMany/ManyToMany）
- ✅ 事务处理
- ✅ 高级查询和分页
- ✅ 数据库连接池

### Phase 4: 进阶实战
- ✅ 项目架构设计（分层架构）
- ✅ 配置管理（Viper）
- ✅ 日志系统（Zap）
- ✅ 完整的CRUD API
- ✅ 事务管理
- ✅ 复杂查询和搜索
- ✅ 生产级项目结构

## 🚀 开始学习

### 1. 选择阶段

根据你的当前水平选择合适的阶段：

- **零基础**: 从 Phase 1 开始
- **有基础**: 可以从 Phase 2 或 Phase 3 开始
- **想实战**: 直接跳到 Phase 4 看完整项目

### 2. 进入阶段目录

每个阶段都是独立的Go模块，可以单独运行：

```bash
# 进入第一阶段
cd phase1_basic

# 查看该阶段说明
cat README.md

# 运行示例
go run main.go
```

### 3. 学习路径建议

**路径A: 按部就班（推荐）**
```
Phase 1 → Phase 2 → Phase 3 → Phase 4
   ↓         ↓         ↓         ↓
 基础      进阶      数据      实战
```

**路径B: 按需学习**
- 需要了解中间件 → Phase 2
- 需要数据库操作 → Phase 3
- 想看完整项目 → Phase 4

**路径C: 快速上手**
- 先快速浏览 Phase 1-3 的知识点
- 重点学习 Phase 4 的实战项目

## 📝 学习检查清单

### Phase 1 检查清单 ✅ (已完成 - 2024-02-20)
- [x] 成功运行 Hello World
- [x] 理解 gin.Default() 和 gin.New() 的区别
- [x] 能创建 GET/POST/PUT/DELETE 路由
- [x] 能使用路由参数和路由分组
- [x] 能处理 Query/Form/JSON 请求
- [x] 能返回 JSON 和 String 响应
- [x] 完成练习：用户管理系统（评分：92/100）
  - 包含：注册、登录、CRUD、分页、软删除
  - 接口测试：Bruno API 测试文件已提供

### Phase 2 检查清单 ✅ (已完成 - 2024-02-21)
- [x] 理解中间件的执行流程
- [x] 能创建自定义中间件
- [x] 能在不同层级使用中间件
- [x] 掌握 JWT 认证实现
- [x] 实现请求限流中间件
- [x] 实现日志记录中间件
- [x] 实现错误恢复中间件
- [x] 实现 RBAC 权限控制
- [x] 完成练习：JWT 认证与验证系统（评分：100/100）
  - 包含：注册、登录、JWT Token、限流、权限控制
  - 接口测试：Bruno API 测试文件已提供

### Phase 3 检查清单
- [ ] 成功连接数据库
- [ ] 能定义数据模型
- [ ] 能完成基本的CRUD操作
- [ ] 理解各种关联关系
- [ ] 能正确使用事务
- [ ] 能实现分页和搜索

### Phase 4 检查清单
- [ ] 理解项目分层架构
- [ ] 能配置和使用日志系统
- [ ] 能读取配置文件
- [ ] 能完成完整的API开发
- [ ] 理解仓库/服务/控制器的职责
- [ ] 能部署到生产环境

## 🛠️ 常用命令

```bash
# 运行某个阶段
cd phase1_basic && go run main.go
cd phase2_middleware && go run main.go
cd phase3_database && go run main.go
cd phase4_advanced && go run main.go

# 编译
go build -o app .

# 运行测试
go test ./...

# 查看依赖
go mod graph

# 更新依赖
go get -u ./...
```

## 📖 参考资料

### 官方文档
- [Gin 官方文档](https://gin-gonic.com/zh-cn/docs/)
- [GORM 官方文档](https://gorm.io/zh_CN/docs/)
- [Go 官方文档](https://go.dev/doc/)

### 推荐书籍
- 《Go语言实战》
- 《Go Web编程》

### 推荐项目
- [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin)
- [go-zero](https://github.com/zeromicro/go-zero)
- [gin-boilerplate](https://github.com/Massad/gin-boilerplate)

## 🤝 贡献

如果你发现错误或有改进建议，欢迎提交 Issue 或 PR。

## 📄 License

MIT License - 详见 LICENSE 文件

---

**开始学习**: [点击进入 Phase 1](phase1_basic/README.md)
