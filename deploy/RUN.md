# Clipper 部署运行文档

## 部署包内容

```
Clipper/
├── clipper.exe          # 主程序可执行文件
├── config.yaml          # 配置文件
├── web/dist/            # 前端静态文件
├── data/                # 数据库目录 (首次运行自动创建)
├── uploads/             # 图片上传目录 (首次运行自动创建)
└── deploy/              # 部署文档目录
    └── RUN.md           # 本文档
```

## 快速运行

### Windows

双击 `clipper.exe` 或在命令行运行：

```cmd
clipper.exe
```

访问 `http://localhost:8080`

### Linux/Mac

```bash
# 如需在 Linux/Mac 运行，需重新编译
go build -o clipper ./cmd/server
./clipper
```

## 配置说明

编辑 `config.yaml` 文件：

```yaml
server:
  http_port: 8080              # 服务端口
  session_secret: "your-secret-key"  # Session密钥 (生产环境请修改)

database:
  path: "./data/clipper.db"    # 数据库路径

storage:
  upload_dir: "./uploads"      # 图片存储目录

log:
  level: "info"
  format: "text"
```

**生产环境建议：**
- 修改 `session_secret` 为随机字符串
- 更换 HTTP 端口为其他端口 (如 80, 443)
- 配置 HTTPS (需要反向代理如 Nginx)

## 首次使用流程

1. 启动程序后访问 `http://localhost:8080`
2. 点击 "Create one" 注册账号
3. 输入用户名、密码 (密码至少6位)
4. 输入设备名称 (如 "Windows PC", "iPhone")
5. 登录后即可使用

## 功能说明

### 文本剪贴
- 左侧 "Text" 标签页输入文本内容
- 点击 "Create & Sync" 创建并同步
- 右侧列表显示所有剪贴内容
- 点击 "Copy" 复制到本地剪贴板

### 图片剪贴
- 左侧切换到 "Image" 标签页
- 拖拽或点击上传图片 (支持 JPEG/PNG/GIF/WebP, 最大20MB)
- 缩略图显示在列表中
- 点击缩略图打开预览弹窗
- 预览弹窗支持复制、下载、关闭

### 设备同步
- 多设备登录同一账号
- WebSocket 实时同步
- 设备上线/下线通知显示设备名称

### 分页浏览
- 支持 5/10/20 条每页
- 列表区域独立滚动

## 多设备访问

同一局域网内，其他设备可通过 IP 访问：

```text
http://<服务器IP>:8080
```

例如：`http://192.168.1.100:8080`

## 生产环境部署

### 使用 Nginx 反向代理

```nginx
server {
    listen 80;
    server_name clipper.example.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /api/v1/devices/ws {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }
}
```

### 使用 systemd 服务 (Linux)

创建 `/etc/systemd/system/clipper.service`：

```ini
[Unit]
Description=Clipper Clipboard Sync Service
After=network.target

[Service]
Type=simple
User=clipper
WorkingDirectory=/opt/clipper
ExecStart=/opt/clipper/clipper
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启用服务：

```bash
sudo systemctl enable clipper
sudo systemctl start clipper
```

## 数据备份

### 备份数据库

```bash
cp ./data/clipper.db ./backup/clipper_$(date +%Y%m%d).db
```

### 备份图片

```bash
tar -czf ./backup/uploads_$(date +%Y%m%d).tar.gz ./uploads/
```

## 常见问题

### Q: 端口被占用
修改 `config.yaml` 中的 `http_port` 为其他端口。

### Q: 无法跨设备同步
确保所有设备能访问服务器 IP，检查防火墙是否放行端口。

### Q: 图片上传失败
检查图片格式 (仅支持 JPEG/PNG/GIF/WebP) 和大小 (最大20MB)。

### Q: 登录后刷新页面状态丢失
正常现象，页面会自动重新验证登录状态。

### Q: iOS Safari 复制失败
iOS 有安全限制，点击 "Copy" 后需手动粘贴。

## 技术支持

如有问题，请查看日志输出或联系开发者。

## 文件清单

| 文件/目录 | 说明 | 必须 |
|----------|------|------|
| clipper.exe | 主程序 | ✓ |
| config.yaml | 配置文件 | ✓ |
| web/dist/ | 前端文件 | ✓ |
| data/ | 数据库目录 | 自动创建 |
| uploads/ | 图片目录 | 自动创建 |

## 系统要求

- Windows 10+ / Linux / macOS
- 无需额外依赖 (静态编译)
- 内存: 最低 64MB
- 磁盘: 根据使用量决定 (图片存储)

## 安全建议

1. 修改默认 Session 密钥
2. 使用 HTTPS (通过 Nginx 等反向代理)
3. 定期备份数据库和图片
4. 不要在公网直接暴露服务 (使用 VPN 或局域网)