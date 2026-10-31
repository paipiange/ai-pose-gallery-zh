# 拍片哥 AI 姿势库与工坊 · Docker 一键部署指南

本项目已提供完整的 **前端（Nginx 高性能托管与反向代理）+ 后端（GoFrame 姿势生成与 AI 换装引擎）** 全栈 Docker 容器化编排支持。只需一行命令即可在任何安装有 Docker 的机器上快速拉起整套服务。

---

## 一、架构概览

```text
 浏览器客户端 (Browser)
      │
      ▼
┌────────────────────────────────────────────────────────┐
│  拍片哥前端容器 (ai-pose-frontend · Nginx)              │
│  - 监听端口: 8080 (可由 FRONTEND_PORT 自定义)            │
│  - 静态资源: 姿势图库、动作分类、高清骨架素材           │
│  - 内置代理: /tryon, /hello 自动反向代理至后端，零跨域  │
└───────────────────────────┬────────────────────────────┘
                            │ (Docker 网络内部高速通信)
                            ▼
┌────────────────────────────────────────────────────────┐
│  拍片哥后端容器 (ai-pose-backend · GoFrame v2)          │
│  - 监听端口: 8000 (可由 BACKEND_PORT 自定义)            │
│  - 模型引擎: Gemini / 豆包 Seedream / OpenAI / ComfyUI  │
│  - 健康检查: GET /hello (自动健康探测保活)              │
└────────────────────────────────────────────────────────┘
```

---

## 二、5 秒极速上手

### 1. 复制环境变量文件
在项目根目录下运行：
```bash
cp .env.example .env
```

### 2. （可选）配置大模型 API Key
编辑 `.env` 文件，填入你使用的服务商 Key（按需填写即可，也可以留空，直接在前端页面的「设置中心」临时输入）：

```ini
# 推荐优先使用 Google Gemini 或 豆包 Seedream
GEMINI_API_KEY=AIzaSy...
ARK_API_KEY=xxxxxx-xxxx-xxxx...
OPENAI_API_KEY=sk-...
```

> **安全提示**：所有 API Key 均严格受 `.gitignore` 保护，绝不会被误提交至 Git 仓库。

### 3. 一键构建并启动
在项目根目录运行：
```bash
docker compose up -d --build
```

### 4. 立即使用
- **Web 前端界面**：[http://localhost:8080](http://localhost:8080)
- **后端 API 与文档**：[http://localhost:8000/swagger](http://localhost:8000/swagger)
- **健康探测地址**：[http://localhost:8000/hello](http://localhost:8000/hello)

---

## 三、宿主机 ComfyUI 本地离线生图联动

若你的宿主机拥有独立显卡并运行了 ComfyUI（如采用 Qwen-Image-Edit-2511 4步极速生图或 MiniMax H3 工作流）：

1. 容器已内置 `host.docker.internal` 主机网关映射；
2. 启动 ComfyUI 时请添加 `--enable-cors-header=*` 参数；
3. 后端容器即可免去任何云端 API 费用，全本地显卡离线生成！

---

## 四、常用维护与管理指令

### 查看容器运行状态
```bash
docker compose ps
```

### 查看实时运行日志
```bash
# 查看前后端所有服务日志
docker compose logs -f

# 仅查看后端生图推理日志
docker compose logs -f backend

# 仅查看前端 Nginx 访问日志
docker compose logs -f frontend
```

### 重启服务
```bash
docker compose restart
```

### 停止并移除容器
```bash
docker compose down
```

### 代码更新后重新构建并生效
```bash
git pull
docker compose up -d --build
```

---

## 五、自定义端口与高级配置

若服务器上的默认端口（8080 或 8000）已被占用，只需在 `.env` 中修改对应的端口变量：

```ini
FRONTEND_PORT=3000   # 将前端访问端口改为 3000
BACKEND_PORT=9000    # 将后端接口端口改为 9000
```
然后执行 `docker compose up -d` 即可生效。
