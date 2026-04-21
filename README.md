# Clipper

跨平台剪贴板同步应用，支持文本和图片的多设备实时同步。

## 功能特性

- **用户认证**：注册、登录、登出（Session Cookie + SameSite=Lax）
- **内容类型**：文本、图片（支持缩略图）
- **实时同步**：WebSocket 实现多设备剪贴板同步
- **设备命名**：自定义设备名称，在线/离线通知显示设备名
- **图片支持**：上传、预览、下载、复制到剪贴板（20MB限制）
- **快捷粘贴**：Ctrl+V 全局粘贴检测，自动识别文本/图片并切换标签页
- **跨平台访问**：PC、平板、手机浏览器访问
- **iOS 兼容**：复制功能支持 iOS Safari
- **分页显示**：剪贴内容分页浏览，支持 5/10/20 条切换
- **响应式布局**：左右双栏布局，移动端自动切换单栏
- **用户名校验**：用户名重复时提示，移除邮箱字段

## 技术栈

### 后端
- Go 1.21+
- Gin (HTTP框架)
- Gorilla WebSocket (实时通信)
- Gorilla Sessions (Session认证)
- SQLite (数据存储)
- Viper (配置管理)
- bcrypt (密码加密)
- imaging (图片缩略图生成)

### 前端
- Vue 3 + TypeScript
- Pinia (状态管理)
- Vue Router (路由守卫 + 认证检查)
- Element Plus (UI组件)
- Axios (HTTP客户端, withCredentials)
- Vite (构建工具)

## 项目结构

```
Clipper/
├── cmd/server/main.go          # 程序入口 + CORS中间件 + 静态文件路由
├── config.yaml                 # 配置文件
├── clipper.exe                 # 编译后的可执行文件
├── internal/
│   ├── config/config.go        # 配置加载
│   ├── model/                  # 数据模型
│   │   ├── user.go             # 用户模型 (无email字段)
│   │   ├── clip.go             # 剪贴内容模型 (ClipType: text/image + ClipMetaImage)
│   │   └── device.go           # 设备模型
│   ├── handler/                # HTTP handlers
│   │   ├── auth.go             # 认证接口 (register/login/logout/me)
│   │   ├── clip.go             # 剪贴内容接口 (list/create/upload/delete)
│   │   └── device.go           # 设备管理 & WebSocket升级
│   ├── service/                # 业务逻辑
│   │   ├── auth.go             # 认证服务 (bcrypt密码)
│   │   ├── clip.go             # 剪贴服务 (含UploadImage)
│   │   └── device.go           # 设备服务
│   ├── repository/             # 数据访问层
│   │   ├── db.go               # 数据库初始化 + 自动迁移
│   │   ├── user.go
│   │   ├── clip.go
│   │   └── device.go
│   ├── ws/                     # WebSocket
│   │   ├── hub.go              # 连接中心 (用户隔离广播)
│   │   ├── client.go           # ReadPump/WritePump
│   │   ├── message.go          # 消息类型定义
│   ├── middleware/auth.go      # Session 认证 (SameSite=Lax)
│   ├── pkg/crypto/password.go  # bcrypt密码加密
│   └── pkg/storage/storage.go  # 图片存储 + 缩略图生成
├── uploads/                    # 图片存储目录
│   └── {user_id}/              # 用户隔离目录
│       ├── {clip_id}_orig.{ext}    # 原图
│       └── {clip_id}_thumb.jpg    # 缩略图 (300px)
├── web/                        # 前端
│   ├── dist/                   # 构建输出
│   │   └── favicon.svg         # 浏览器标签图标
│   ├── public/                 # 静态资源
│   │   └── favicon.svg         # Clipper SVG图标
│   └── src/
│       ├── views/              # 页面组件
│       │   ├── Login.vue       # 登录页 (设备名称输入)
│       │   ├── Register.vue    # 注册页 (密码确认校验 + 用户名重复提示)
│       │   ├── Dashboard.vue   # 主页 (左右双栏 + 图片上传 + 预览弹窗)
│       │   └── Admin.vue       # 管理页
│       ├── stores/             # Pinia stores
│       │   ├── auth.ts         # 认证状态 (withCredentials)
│       │   ├── clip.ts         # 剪贴内容状态 (含uploadImage/getThumbUrl)
│       │   └── device.ts       # 设备状态
│       ├── composables/        # 组合式函数
│       │   └── useWebSocket.ts # WebSocket (自动重连 + 设备名称传递)
│       ├── types/              # TypeScript 类型
│       │   └── index.ts        # User/Clip/ClipMetaImage类型
│       └── router/index.ts     # 路由守卫 (认证检查)
└── data/                       # 数据库文件目录
```

## 实现状态

**已完成:**
- [x] 后端完整实现 (Repository/Service/Handler/WebSocket/Middleware)
- [x] 前端认证流程 (注册/登录/登出 + 状态持久化)
- [x] 文本剪贴内容创建、列表、删除
- [x] WebSocket 实时同步 (用户隔离广播)
- [x] Session Cookie 认证 (SameSite=Lax)
- [x] iOS Safari 复制兼容 (execCommand)
- [x] 注册页密码确认校验 (红色提示)
- [x] 密码框显示/隐藏切换 (眼睛图标)
- [x] 页面排版优化 (渐变背景、圆角卡片)
- [x] Dashboard 双栏布局 (左创建/右列表)
- [x] Tab切换动画 (水平滑动 + 弹性曲线)
- [x] 分页功能 (5/10/20条切换)
- [x] 右侧列表独立滚动 (固定最大高度)
- [x] 图片剪贴上传、预览、下载功能
- [x] 图片缩略图生成 (300px宽度, JPEG 85%质量)
- [x] 用户隔离图片存储 (`uploads/{user_id}/`)
- [x] 20MB图片大小限制, 格式校验 (JPEG/PNG/GIF/WebP)
- [x] 设备命名 (登录页 + Dashboard可编辑)
- [x] 设备在线/离线通知显示设备名称
- [x] 用户名重复校验 (toast提示 + 输入框变红)
- [x] 移除邮箱字段 (仅用户名注册)
- [x] 移动端响应式布局
- [x] 图片预览弹窗 (复制/下载/关闭按钮)
- [x] Clip卡片布局修复 (无压缩, 列表容器滚动)
- [x] 浏览器标签图标 (Clipper SVG favicon)
- [x] 全局 Ctrl+V 粘贴检测 (自动识别文本/图片, 自动切换标签页)
- [x] Session错误修复 (无效旧cookie不阻止登录/登出)

**待实现:**
- [ ] 剪贴内容搜索/筛选
- [ ] 设备管理页面 (在线设备列表)
- [ ] 剪贴内容加密存储

## 快速开始

### 生产模式（推荐）

```bash
# 构建前端
cd web && npm run build

# 运行后端 (前端已内嵌)
go run ./cmd/server/main.go

# 或编译后运行
go build -o clipper.exe ./cmd/server
./clipper.exe
```

访问 `http://localhost:8080`

### 开发模式

```bash
# 启动后端 (端口 8080)
go run ./cmd/server/main.go

# 启动前端开发服务器 (端口 3000)
cd web && npm run dev
```

注意：开发模式存在跨域问题，推荐直接访问 8080 端口。

## 配置说明

`config.yaml` 配置文件：

```yaml
server:
  http_port: 8080
  session_secret: "clipper-secret-change-in-production"

database:
  path: "./data/clipper.db"

storage:
  upload_dir: "./uploads"

log:
  level: "info"
  format: "text"
```

## API 接口

### 认证接口
- `POST /api/v1/auth/register` - 用户注册 (仅用户名+密码)
- `POST /api/v1/auth/login` - 用户登录（设置 Session Cookie）
- `POST /api/v1/auth/logout` - 用户登出
- `GET /api/v1/auth/me` - 获取当前用户信息

### 剪贴内容接口
- `GET /api/v1/clips` - 获取用户所有剪贴内容
- `POST /api/v1/clips?device_id=xxx` - 创建文本剪贴内容
- `POST /api/v1/clips/upload?device_id=xxx` - 上传图片 (multipart/form-data)
- `GET /api/v1/clips/:id` - 获取单个剪贴内容
- `DELETE /api/v1/clips/:id?device_id=xxx` - 删除剪贴内容 (同时删除图片文件)

### 设备接口
- `GET /api/v1/devices` - 获取用户在线设备
- `GET /api/v1/devices/ws?device_id=xxx&device_name=xxx` - WebSocket 连接（需认证）

### 静态文件
- `GET /uploads/{user_id}/{clip_id}_orig.{ext}` - 获取原图
- `GET /uploads/{user_id}/{clip_id}_thumb.jpg` - 获取缩略图
- `GET /favicon.svg` - 浏览器标签图标

## WebSocket 消息类型

| 类型 | 说明 |
|------|------|
| `clip_created` | 新剪贴内容创建通知 (含clip元信息) |
| `clip_deleted` | 剪贴内容删除通知 |
| `device_online` | 设备上线通知 (显示device_name) |
| `device_offline` | 设备离线通知 (显示device_name) |

## 数据库结构

```sql
-- 用户表 (无email字段)
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME,
    updated_at DATETIME
);

-- 剪贴内容表
CREATE TABLE clips (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    device_id TEXT NOT NULL,
    type TEXT NOT NULL,        -- 'text' 或 'image'
    content TEXT NOT NULL,     -- 文本内容或图片路径
    meta TEXT,                 -- JSON: {width, height, size, format, thumb_path}
    created_at DATETIME,
    updated_at DATETIME
);

-- 设备表
CREATE TABLE devices (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    ip TEXT,
    last_seen DATETIME,
    is_online BOOLEAN
);
```

## 已修复的问题

1. **登录状态刷新丢失** - Router guard 现在等待认证检查完成
2. **登录/登出不跳转** - 重置认证检查标志后正确跳转
3. **iOS 复制失败** - 使用 execCommand + textarea 兼容方案
4. **跨域 Cookie 不携带** - CORS 动态返回请求 Origin + SameSite=Lax
5. **Tab 切换穿帮** - 使用 overflow:hidden + 水平滑动动画
6. **WebSocket显示Unknown Device** - 现在传递device_name参数
7. **Clip卡片内容被裁剪** - 列表容器滚动, 卡片不压缩
8. **预览弹窗白色底板** - 全局样式覆盖dialog背景为透明
9. **浏览器无favicon** - 添加SVG图标到public/, 服务端提供路由

## License

MIT