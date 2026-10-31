# AGENTS.md · AI 姿势库协作智能体开发与提交规范

本文档面向所有参与 `ai-pose-gallery-zh` 仓库维护的 AI Agent（如 OpenCode、Claude Code、Cursor、GitHub Copilot 等）及协作开发者。任何代理或贡献者在修改、生成代码和提交 Git 记录时，必须严格执行本文档所规定的安全红线与工程规范。

---

## 1. 核心安全红线（绝对零泄露原则）

### 1.1 严禁提交任何真实密钥与凭据
- **绝对禁止**将任何真实的 API Key（包括但不限于 OpenAI `sk-...`、Google Gemini `AIza...`、火山方舟 `ARK...`、通义千问 `sk-...` 等）、私有 Token、密码、私钥证书提交到 Git 仓库或暂存区。
- 仓库内仅允许保留配置文件模板（如 `backend/manifest/config/config.example.yaml`），且模板内的密钥字段**必须为空或通用占位符**（如 `""`、`"sk-你的key"`）。
- 本地真实配置文件（`backend/manifest/config/config.yaml`、`backend/.env` 等）必须在 `.gitignore` 中被严格忽略。

### 1.2 严禁提交个人隐私与敏感信息
- **绝对禁止**在代码、注释、提交信息（Commit Message）、PR 描述或生成文档中包含个人真实姓名、私人邮箱、手机号码、私人局域网 IP、服务器公网 IP、个人绝对开发路径（如包含开发者个人用户名的路径）。
- Git 提交者信息应保持项目已有的一致风格，禁止在无授权情况下将私人联系方式作为作者写入共有历史。

### 1.3 外部输入安全与 SSRF 防护
- 后端服务内置针对客户端 Base URL 的安全过滤，禁止指向内网私有 IP 范围（如 `10.0.0.0/8`, `192.168.0.0/16`, `127.0.0.1` 等），防止 SSRF 漏洞。
- 绝不允许静默向未受信任的第三方域外发包含用户私有 API Key（如 `client_config.api_key`）、认证 Token 或用户照片的敏感载荷。

---

## 2. 架构与防御性编码规范

1. **前端原生零构建架构**：
   - 前端页面保持原生 HTML5 / CSS3 / Vanilla JS 架构，零打包构建工具、零第三方重量级框架依赖，保障极速加载与轻量体验。
   - 所有用户私有配置（API Key、Base URL 等）仅保存在浏览器本地 `localStorage`，单次请求临时封装在 `client_config` 传递，后端仅在内存中使用，绝不持久化至磁盘。
   - 本地生成结果采用浏览器 IndexedDB (`PoseStudioDB`) 安全存储，支持单条删除与本地清空。

2. **后端服务健壮性与兜底（Fail-Safe Defaults）**：
   - 后端服务（GoFrame）在启动时，若检测到宿主机未显式创建 `config.yaml`，必须具备代码层自动回退加载 `config.example.yaml` 的能力。
   - 服务端口在未配置时，必须硬兜底显式绑定 `:8000`，杜绝因未配置 address 而进入系统 `:0` 随机高位临时端口。
   - CORS 中间件需安全放行常见本地前端开发预览端口（`:8000`, `:8080`, `:5500`, `:5173`, `:3000` 等），同时阻断未授权的跨域恶意来源。

3. **图像安全校验防线**：
   - 针对前端或客户端提交的 Base64 与 Data URL 图片，必须执行严格的真实二进制文件头探测（Magic Bytes / `http.DetectContentType`）。
   - 仅放行合规且安全的 `image/jpeg`、`image/png`、`image/webp` 格式，自动归一化 `image/jpg`，并放行通用二进制流 `application/octet-stream` 进入底层真实 Magic Bytes 检测。
   - 严格限制图片最小/最大分辨率及总像素阈值，杜绝乘法整数溢出与解压炸弹风险。

4. **姿势生图与换装规范**：
   - 姿势白模素体生图时，必须传递 `keep_face: false`，自动匹配专属姿势引导提示词（`poseTryonEnPrompt` / `poseTryonZhPrompt`），确保真实模特在素体骨架上自然渲染；换装模式则保留原人脸一致性。

---

## 3. Git 提交自我审查 Checklist（提交前必查）

在任何 `git commit` 与 `git push` 执行之前，Agent 必须在后台严格执行以下核验步骤：

- [ ] **步骤 1：排查 Untracked 与暂存文件**
  运行 `git status`，确保没有临时文件、日志文件（`*.log`）、编译产物（`main`、`*.exe`）、真实配置文件（`config.yaml`）进入暂存区。
- [ ] **步骤 2：逐行审查 Git Diff**
  运行 `git diff --cached`，逐行检查每一处新增和修改的字符，确认：
  - 无真实 Key、Token、密码；
  - 无个人姓名、电话、私有私人邮箱、个人电脑私有路径；
  - 无调试时遗留的临时硬编码凭证或测试图片。
- [ ] **步骤 3：敏感字符串安全扫描**
  对修改文件执行正则扫描，确保无符合正则 `(sk-[a-zA-Z0-9]{20,}|AIza[a-zA-Z0-9_-]{20,})` 的高危私钥。
- [ ] **步骤 4：自动化测试与构建校验**
  运行后端单元测试 `go test ./...` 与 `go vet ./...`，确保代码编译通过、测试全绿。

---

## 4. Commit Message 提交规范

遵循 **Conventional Commits** 规范，提交信息简洁、清晰、明确表达变更动机与影响：

- `feat: ...`：新增功能或用户可见特性（如接入 AI 创作工坊）
- `fix: ...`：修复缺陷、修复漏洞或修复配置异常
- `docs: ...`：文档调整（README、AGENTS.md 等）
- `test: ...`：新增或重构测试用例
- `refactor: ...`：既不修复错误也不添加新功能的代码重构
- `chore: ...`：构建流程、依赖维护或辅助工具变动
