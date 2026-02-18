# 红岩网校作业管理系统

## 1. 已实现功能清单

### 用户模块
- 用户注册（支持用户名、密码、昵称、部门，密码 bcrypt 加密）
- 用户登录（返回 Access Token 和 Refresh Token，双 token JWT 认证）
- 刷新 Token（使用 Refresh Token 获取新的双 token）
- 获取当前用户信息（需认证）
- 注销账号（软删除，需密码二次确认）

### 作业模块
- 发布作业（仅老登，需指定标题、描述、部门、截止时间、补交策略）
- 作业列表（支持按部门筛选、分页查询，返回部门标签）
- 作业详情（关联发布者信息，包含部门标签；小登可看到自己的提交记录）
- 修改作业（仅同部门老登，支持部分字段更新，含并发控制）
- 删除作业（仅同部门老登）

### 提交模块
- 提交作业（仅小登，自动判断是否迟交，检查作业存在性、部门匹配、截止时间）
- 我的提交列表（小登查看自己的所有提交及评语）
- 查看部门提交（老登查看本部门学员的提交列表）
- 批改作业（老登添加分数、评语，可标记优秀）
- 标记/取消优秀作业（老登）
- 优秀作业列表（所有用户可查看，支持部门筛选、分页）

### 其他
- 部门枚举与标签（响应中同时返回 department 和 department_label）
- 统一响应格式 `{code, message, data}`
- 统一错误码定义
- JWT 认证中间件
- 权限中间件（区分老登/小登）
- GORM 数据库操作，支持软删除
- 并发控制（更新作业时使用乐观锁）

## 2. 进阶功能说明

**暂未实现**以下进阶功能：
- 邮箱绑定与邮件通知（作业发布、截止提醒、批改通知）
- AI 作业评价
- 前端页面对接
- 部署（Docker、Linux）
- 考核系统（区分考核类型、多人阅卷）

后续可根据需求逐步扩展。

## 3. 项目简介

本项目是为红岩网校内部开发的作业管理系统，旨在实现「老登」（管理员/讲师）发布作业、「小登」（学员）提交作业的完整流程。系统支持多部门管理（后端、前端、SRE、产品、视觉设计、Android、iOS），提供作业的增删改查、提交批改、优秀作业展示等功能。采用 Go 语言编写，基于 Gin 框架和 GORM 数据库 ORM，使用 JWT 进行认证，代码结构清晰，遵循分层架构。

## 4. 技术栈说明

| 组件       | 技术              | 说明                         |
|------------|-------------------|------------------------------|
| 后端语言   | Go 1.21           | 核心开发语言                 |
| Web 框架   | Gin               | HTTP 路由和中间件            |
| 数据库     | MySQL 8.0         | 数据存储                     |
| ORM        | GORM              | 数据库操作                   |
| 认证       | JWT (双 token)    | 用户认证与授权               |
| 密码加密   | bcrypt            | 用户密码哈希                 |
| 配置文件   | 环境变量 / 硬编码 | 简单配置（可扩展）           |
| 项目结构   | 分层架构          | controller/service/repository/models |

## 5. 项目结构说明
```text
homework-system/
├── cmd/ # 程序入口
│ └── main.go
├── configs/ # 配置加载与数据库初始化
│ └── config.go
├── internal/ # 内部业务代码
│ ├── controller/ # HTTP 处理器（参数解析、响应返回）
│ │ ├── base.go
│ │ ├── user.go
│ │ ├── homework.go
│ │ └── submission.go
│ ├── service/ # 业务逻辑层
│ │ ├── service.go
│ │ ├── user.go
│ │ ├── homework.go
│ │ └── submission.go
│ ├── repository/ # 数据访问层
│ │ ├── errors.go
│ │ ├── user.go
│ │ ├── assignment.go
│ │ └── submission.go
│ └── models/ # 数据模型定义及辅助方法
│ ├── user.go
│ ├── assignment.go
│ ├── submission.go
│ ├── department.go
│ └── roles.go
├── middleware/ # Gin 中间件
│ └── auth.go # JWT 认证
├── pkg/ # 公共工具包
│ ├── jwt/ # JWT 生成与验证
│ ├── response/ # 统一响应格式
│ └── errcode/ # 错误码定义
├── router/ # 路由注册
│ └── router.go
├── docs/ # 文档
│ └── HomeworkSystem.postman_collection.json # Postman 集合
├── go.mod
├── go.sum
└── README.md
```
## 6. 本地运行指南

### 环境要求
- Go 1.21 或更高版本
- MySQL 8.0
- Git

### 步骤

1. **克隆仓库**
   git clone https://github.com/yourusername/homework-system.git
   cd homework-system

2. **配置数据库**
    - 创建数据库，例如 `homework`
    - 修改 `configs/config.go` 中的数据库连接信息（用户名、密码、地址、库名）

3. **安装依赖**
   go mod tidy

4. **运行项目**
   go run cmd/main.go
   服务默认启动在 `http://localhost:8080`。

5. **测试 API**
    - 导入 Postman 集合（见下一节）
    - 设置环境变量 `base_url` 为 `http://localhost:8080/api`
    - 依次调用接口测试