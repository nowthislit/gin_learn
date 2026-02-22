package main

import (
	"time"
)

// seedData 初始化测试数据
func seedData() {
	// 清空数据
	DB.Exec("DELETE FROM article_tags")
	DB.Exec("DELETE FROM user_follows")
	DB.Exec("DELETE FROM comments")
	DB.Exec("DELETE FROM articles")
	DB.Exec("DELETE FROM tags")
	DB.Exec("DELETE FROM categories")
	DB.Exec("DELETE FROM users")

	now := time.Now()

	// ========== 1. 创建用户 ==========
	users := []User{
		{
			Username:  "zhangsan",
			Email:     "zhangsan@example.com",
			Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Avatar:    "https://api.dicebear.com/7.x/avataaars/svg?seed=zhangsan",
			Bio:       "热爱技术，喜欢分享的全栈开发者",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Username:  "lisi",
			Email:     "lisi@example.com",
			Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
			Avatar:    "https://api.dicebear.com/7.x/avataaars/svg?seed=lisi",
			Bio:       "专注于后端开发和系统架构",
			CreatedAt: now.Add(-24 * time.Hour),
			UpdatedAt: now.Add(-24 * time.Hour),
		},
		{
			Username:  "wangwu",
			Email:     "wangwu@example.com",
			Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
			Avatar:    "https://api.dicebear.com/7.x/avataaars/svg?seed=wangwu",
			Bio:       "前端工程师，热爱 React 和 Vue",
			CreatedAt: now.Add(-48 * time.Hour),
			UpdatedAt: now.Add(-48 * time.Hour),
		},
	}
	for i := range users {
		DB.Create(&users[i])
	}

	// ========== 2. 创建分类（树形结构）==========
	// 父分类
	parentCategories := []Category{
		{Name: "技术", Slug: "tech", Level: 1, SortOrder: 1, CreatedAt: now},
		{Name: "生活", Slug: "life", Level: 1, SortOrder: 2, CreatedAt: now},
		{Name: "随笔", Slug: "essay", Level: 1, SortOrder: 3, CreatedAt: now},
	}
	for i := range parentCategories {
		DB.Create(&parentCategories[i])
	}

	// 子分类
	childCategories := []Category{
		{Name: "后端开发", Slug: "backend", ParentID: &parentCategories[0].ID, Level: 2, SortOrder: 1, CreatedAt: now},
		{Name: "前端开发", Slug: "frontend", ParentID: &parentCategories[0].ID, Level: 2, SortOrder: 2, CreatedAt: now},
		{Name: "数据库", Slug: "database", ParentID: &parentCategories[0].ID, Level: 2, SortOrder: 3, CreatedAt: now},
		{Name: "DevOps", Slug: "devops", ParentID: &parentCategories[0].ID, Level: 2, SortOrder: 4, CreatedAt: now},
		{Name: "旅行", Slug: "travel", ParentID: &parentCategories[1].ID, Level: 2, SortOrder: 1, CreatedAt: now},
		{Name: "美食", Slug: "food", ParentID: &parentCategories[1].ID, Level: 2, SortOrder: 2, CreatedAt: now},
	}
	for i := range childCategories {
		DB.Create(&childCategories[i])
	}

	// ========== 3. 创建标签 ==========
	tags := []Tag{
		{Name: "Go", CreatedAt: now},
		{Name: "Gin", CreatedAt: now},
		{Name: "GORM", CreatedAt: now},
		{Name: "MySQL", CreatedAt: now},
		{Name: "Redis", CreatedAt: now},
		{Name: "Docker", CreatedAt: now},
		{Name: "Kubernetes", CreatedAt: now},
		{Name: "React", CreatedAt: now},
		{Name: "Vue", CreatedAt: now},
		{Name: "JavaScript", CreatedAt: now},
		{Name: "微服务", CreatedAt: now},
		{Name: "性能优化", CreatedAt: now},
		{Name: "架构设计", CreatedAt: now},
		{Name: "面试题", CreatedAt: now},
	}
	for i := range tags {
		DB.Create(&tags[i])
	}

	// ========== 4. 创建文章 ==========
	articles := []Article{
		{
			Title:        "Gin框架入门指南",
			Summary:      "从零开始学习Gin框架，构建高性能的Web应用",
			Content:      "## 什么是Gin？\n\nGin是一个用Go语言编写的Web框架，它具有高性能、易用性强的特点。\n\n### 主要特性\n\n1. 快速的路由性能\n2. 中间件支持\n3. JSON验证\n4. 错误管理\n\n## 快速开始\n\n```go\npackage main\n\nimport \"github.com/gin-gonic/gin\"\n\nfunc main() {\n    r := gin.Default()\n    r.GET(\"/ping\", func(c *gin.Context) {\n        c.JSON(200, gin.H{\n            \"message\": \"pong\",\n        })\n    })\n    r.Run()\n}\n```\n\n## 总结\n\nGin框架是Go语言Web开发的首选框架之一。",
			Status:       1,
			ViewCount:    1250,
			LikeCount:    89,
			CommentCount: 15,
			UserID:       users[0].ID,
			CategoryID:   childCategories[1].ID,
			CreatedAt:    now.Add(-72 * time.Hour),
			UpdatedAt:    now.Add(-24 * time.Hour),
		},
		{
			Title:        "GORM最佳实践",
			Summary:      "深入理解GORM的使用技巧和性能优化方法",
			Content:      "## GORM简介\n\nGORM是Go语言最受欢迎的ORM库。\n\n### 核心概念\n\n- **Model定义**：使用struct标签定义数据库表\n- **关联关系**：一对一、一对多、多对多\n- **回调钩子**：BeforeCreate、AfterUpdate等\n\n## 性能优化\n\n1. 使用预加载避免N+1问题\n2. 合理使用索引\n3. 批量操作\n\n## 总结\n\n掌握这些技巧能让你的应用性能提升数倍。",
			Status:       1,
			ViewCount:    980,
			LikeCount:    76,
			CommentCount: 12,
			UserID:       users[0].ID,
			CategoryID:   childCategories[2].ID,
			CreatedAt:    now.Add(-48 * time.Hour),
			UpdatedAt:    now.Add(-12 * time.Hour),
		},
		{
			Title:        "微服务架构设计原则",
			Summary:      "从零构建微服务系统的核心原则和实践",
			Content:      "## 什么是微服务？\n\n微服务是一种架构风格，将应用程序构建为一组小型服务。\n\n### 设计原则\n\n1. **单一职责**：每个服务只做一件事\n2. **独立部署**：服务可以独立部署和扩展\n3. **去中心化治理**：技术栈可以多样化\n\n### 服务间通信\n\n- 同步通信：REST、gRPC\n- 异步通信：消息队列\n\n## 总结\n\n微服务不是银弹，需要权衡利弊。",
			Status:       1,
			ViewCount:    2100,
			LikeCount:    156,
			CommentCount: 28,
			UserID:       users[1].ID,
			CategoryID:   childCategories[0].ID,
			CreatedAt:    now.Add(-96 * time.Hour),
			UpdatedAt:    now.Add(-48 * time.Hour),
		},
		{
			Title:        "Docker容器化部署实战",
			Summary:      "手把手教你使用Docker部署Go应用",
			Content:      "## Docker基础\n\nDocker是一个开源的应用容器引擎。\n\n### Dockerfile示例\n\n```dockerfile\nFROM golang:1.21-alpine AS builder\nWORKDIR /app\nCOPY . .\nRUN go build -o main .\n\nFROM alpine:latest\nRUN apk --no-cache add ca-certificates\nWORKDIR /root/\nCOPY --from=builder /app/main .\nCMD [\"./main\"]\n```\n\n## 多阶段构建\n\n使用多阶段构建可以大幅减小镜像体积。\n\n## 总结\n\n容器化是现代应用部署的标准做法。",
			Status:       1,
			ViewCount:    750,
			LikeCount:    45,
			CommentCount: 8,
			UserID:       users[1].ID,
			CategoryID:   childCategories[3].ID,
			CreatedAt:    now.Add(-120 * time.Hour),
			UpdatedAt:    now.Add(-96 * time.Hour),
		},
		{
			Title:        "React Hooks深入解析",
			Summary:      "理解useState、useEffect等Hooks的工作原理",
			Content:      "## Hooks简介\n\nHooks让我们可以在函数组件中使用状态和其他React特性。\n\n### useState\n\n```javascript\nconst [count, setCount] = useState(0);\n```\n\n### useEffect\n\n处理副作用，相当于componentDidMount、componentDidUpdate和componentWillUnmount的组合。\n\n## 自定义Hooks\n\n可以提取组件逻辑到可重用的函数中。\n\n## 总结\n\nHooks改变了我们编写React组件的方式。",
			Status:       1,
			ViewCount:    1850,
			LikeCount:    134,
			CommentCount: 22,
			UserID:       users[2].ID,
			CategoryID:   childCategories[1].ID,
			CreatedAt:    now.Add(-144 * time.Hour),
			UpdatedAt:    now.Add(-72 * time.Hour),
		},
		{
			Title:        "2024年东京之旅",
			Summary:      "记录我在东京一周的旅行见闻和美食推荐",
			Content:      "## 行程概览\n\n这次东京之行为期一周，游览了浅草、涩谷、新宿等地。\n\n### 美食推荐\n\n1. **一兰拉面**：必吃的豚骨拉面\n2. **筑地市场**：新鲜的海鲜刺身\n3. **代官山**：文艺咖啡馆\n\n### 购物攻略\n\n- 涩谷：潮流品牌聚集地\n- 银座：高端百货商店\n\n## 总结\n\n东京是一座值得多次探索的城市。",
			Status:       1,
			ViewCount:    3200,
			LikeCount:    267,
			CommentCount: 45,
			UserID:       users[2].ID,
			CategoryID:   childCategories[4].ID,
			CreatedAt:    now.Add(-168 * time.Hour),
			UpdatedAt:    now.Add(-120 * time.Hour),
		},
		{
			Title:        "（草稿）Kubernetes入门教程",
			Summary:      "",
			Content:      "## 准备工作\n\nTODO: 完善内容\n\n...",
			Status:       0,
			ViewCount:    0,
			LikeCount:    0,
			CommentCount: 0,
			UserID:       users[0].ID,
			CategoryID:   childCategories[3].ID,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}
	for i := range articles {
		DB.Create(&articles[i])
	}

	// ========== 5. 关联文章和标签 ==========
	// 文章1：Gin框架 - 标签：Go, Gin, 后端
	DB.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[0].ID, tags[0].ID)
	DB.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[0].ID, tags[1].ID)
	DB.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[0].ID, tags[10].ID)

	// 文章2：GORM - 标签：Go, GORM, MySQL
	DB.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[1].ID, tags[0].ID)
	DB.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[1].ID, tags[2].ID)
	DB.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[1].ID, tags[3].ID)

	// 文章3：微服务 - 标签：架构设计, 微服务
	DB.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[2].ID, tags[12].ID)
	DB.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[2].ID, tags[10].ID)

	// 文章4：Docker - 标签：Docker, DevOps
	DB.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[3].ID, tags[5].ID)

	// 文章5：React - 标签：React, JavaScript, 前端
	DB.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[4].ID, tags[7].ID)
	DB.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[4].ID, tags[9].ID)

	// 文章6：东京之旅 - 标签：旅行
	// 无需添加标签

	// ========== 6. 创建评论 ==========
	comments := []Comment{
		// 文章1的评论
		{
			Content:   "写得太好了，对Gin框架的讲解很清晰！",
			UserID:    users[1].ID,
			ArticleID: articles[0].ID,
			CreatedAt: now.Add(-24 * time.Hour),
		},
		{
			Content:   "收藏了，学习中。",
			UserID:    users[2].ID,
			ArticleID: articles[0].ID,
			CreatedAt: now.Add(-20 * time.Hour),
		},
		// 文章2的评论
		{
			Content:   "GORM的预加载确实很重要，之前踩过N+1的坑。",
			UserID:    users[2].ID,
			ArticleID: articles[1].ID,
			CreatedAt: now.Add(-18 * time.Hour),
		},
		// 文章3的评论（有回复）
		{
			Content:   "微服务拆分的粒度怎么把握呢？",
			UserID:    users[0].ID,
			ArticleID: articles[2].ID,
			CreatedAt: now.Add(-48 * time.Hour),
		},
		{
			Content:   "这确实是个难题，建议先从单体开始，遇到瓶颈再拆分。",
			UserID:    users[1].ID,
			ArticleID: articles[2].ID,
			ParentID:  func() *uint { id := uint(4); return &id }(),
			CreatedAt: now.Add(-46 * time.Hour),
		},
		{
			Content:   "同意楼上，过早优化是万恶之源。",
			UserID:    users[2].ID,
			ArticleID: articles[2].ID,
			ParentID:  func() *uint { id := uint(4); return &id }(),
			CreatedAt: now.Add(-44 * time.Hour),
		},
		// 文章6的评论
		{
			Content:   "一兰拉面确实好吃！",
			UserID:    users[0].ID,
			ArticleID: articles[5].ID,
			CreatedAt: now.Add(-100 * time.Hour),
		},
		{
			Content:   "下次去东京要试试你推荐的这些。",
			UserID:    users[1].ID,
			ArticleID: articles[5].ID,
			CreatedAt: now.Add(-96 * time.Hour),
		},
	}
	for i := range comments {
		DB.Create(&comments[i])
	}

	// ========== 7. 创建用户关注关系 ==========
	// zhangsan 关注 lisi 和 wangwu
	DB.Exec("INSERT INTO user_follows (user_id, follower_id) VALUES (?, ?)", users[1].ID, users[0].ID)
	DB.Exec("INSERT INTO user_follows (user_id, follower_id) VALUES (?, ?)", users[2].ID, users[0].ID)

	// lisi 关注 wangwu
	DB.Exec("INSERT INTO user_follows (user_id, follower_id) VALUES (?, ?)", users[2].ID, users[1].ID)
}
