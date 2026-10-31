# AI 姿势库 (ai-pose-gallery-zh) 接入 AI 换装/生图后端开发任务书 (AI Studio Integration Spec)

> **面向执行 Agent（Claude Code / OpenCode / Codex / Cursor）的标准化执行规范**  
> 目标仓库：`D:\code\code_ai\ai-pose-gallery-zh`  
> 资源源仓库：`D:\code\code_ai\nvzhuang`（服装库与换装后端）  
> 核心任务：将 `nvzhuang/backend` 的多模型图像生成/换装服务移植到 `ai-pose-gallery-zh`，并对其前端 (`index.html`) 进行交互打通，打造「以姿势为骨架、以服装为血肉」的 **AI 姿势生图与试穿工坊 (Pose AI Studio)**。

---

## 0. 架构定位与核心业务流

```
                    ┌──────────────────────────────────────────────┐
                    │     用户界面：ai-pose-gallery-zh (前端)       │
                    │  (288组姿势 / 1439张多视角白模参考图)          │
                    └──────────────────────┬───────────────────────┘
                                           │
                    [ 浏览姿势 / 角度切换 / 点击"以此姿势创作" ]
                                           │
                                           ▼
                    ┌──────────────────────────────────────────────┐
                    │          AI Studio 创作工坊抽屉 / 弹窗        │
                    │  ① 姿势参考底图 (选定姿势/角度，自动转 Base64)  │
                    │  ② 服装素材 (上传服装照 / 填入女装图库链接)   │
                    │  ③ 生图模型选择 (ComfyUI / Gemini / 豆包等)   │
                    │  ④ 设置中心 (浏览器本地保存个人 API Key / URL)  │
                    └──────────────────────┬───────────────────────┘
                                           │ POST /tryon
                                           ▼
                    ┌──────────────────────────────────────────────┐
                    │      移植后端：ai-pose-gallery-zh/backend     │
                    │           (GoFrame v2 高性能服务)             │
                    │  - 多模型适配 (ComfyUI / Gemini / GPT / 豆包) │
                    │  - 姿势引导生图/换装专属 Prompt 渲染          │
                    │  - 零持久化安全代理 (内存用完即弃，防泄露)    │
                    └──────────────────────┬───────────────────────┘
                                           │ 返回 Base64 图像
                                           ▼
                    ┌──────────────────────────────────────────────┐
                    │          生成结果展示 / 本地 IndexedDB 记录   │
                    └──────────────────────────────────────────────┘
```

---

## 1. 阶段一：后端服务迁移与姿势适配 (`backend/`)

### 1.1 文件完整复制清单
请将 `D:\code\code_ai\nvzhuang\backend\` 目录下的以下核心文件整体迁移至 `D:\code\code_ai\ai-pose-gallery-zh\backend\`：

```
D:\code\code_ai\ai-pose-gallery-zh\backend/
├── api/                   # API 路由定义 (tryon, hello)
├── internal/
│   ├── cmd/               # 服务启动入口与端口兜底绑定 (:8000)
│   ├── consts/            # 常量定义
│   ├── controller/        # 控制器实现
│   ├── imagevalidate/     # 图像安全检查与 Base64 校验
│   ├── logic/tryon/       # 核心生成逻辑 (Gemini, ComfyUI, GPT, 豆包, Prompt)
│   ├── middleware/        # CORS、限流、鉴权、BodyLimit 中间件
│   ├── packed/            # 静态资源打包
│   └── security/          # 安全并发与令牌桶
├── manifest/
│   ├── config/
│   │   ├── config.example.yaml # 配置模板
│   │   └── config.yaml         # 本地真实配置 (确保在 .gitignore)
│   └── docker/            # Docker 构建文件
├── workflows/             # ComfyUI 工作流 API JSON 文件
│   ├── comfyui_qwen_image_edit_2511_api.json
│   └── comfyui_h3_ref2va_still_api.json
├── go.mod                 # Go 模块定义
├── go.sum                 # Go 依赖校验
├── build.bat              # Windows 构建脚本
├── restart.bat            # 重启脚本
├── docker-compose.yml     # 容器编排
└── main.go                # 入口文件
```

### 1.2 专属姿势生图提示词优化 (`internal/logic/tryon/prompt.go`)
在原换装提示词基础上，新增针对「**人体姿势参考图 (Mannequin/Pose Reference)** + **服装 (Garment)**」的适配指令，替换或扩展 `prompt.go`：

```go
// poseTryonEnPrompt 姿势引导换装专用英文提示词（效果最佳）
const poseTryonEnPrompt = `Task: Professional fashion editorial photograph generated from a 3D pose reference and a garment photo.

Strict Requirements:
1. Pose & Angle Transfer: Faithfully follow the exact pose, body orientation, limb positioning, camera angle, and perspective shown in Image 1 (the wireframe/mannequin pose reference). 
2. Realistic Human Model: Transform the reference mannequin into an attractive, photorealistic adult fashion model with natural skin texture, realistic facial expression, modern styled hairstyle, and anatomically accurate proportions.
3. Wardrobe: Dress the model in the exact garment and accessories shown in Image 2. Reproduce the colors, fabrics, drape, textures, neckline, and tailored silhouette faithfully.
4. Lighting & Environment: Professional studio editorial lighting, clean neutral backdrop, shallow depth of field, 8k resolution, ultra-detailed fashion magazine standard.

Output: ONLY the final dressed-up fashion photograph.`

// poseTryonZhPrompt 姿势引导换装专用中文提示词
const poseTryonZhPrompt = `任务：基于人体姿势参考与服装设计图，生成专业商业时尚写实人像大片。

严格要求：
1. 姿态与视角迁移：严格还原图 1（人体姿势参考素体）中的身体姿态、四肢摆放、躯干朝向、透视关系以及拍摄视角（正面/俯视/仰视/侧面/斜侧）。
2. 真人模特渲染：将姿势素体渲染为相貌自然端庄、皮肤质感逼真、五官立体的成年时尚模特，发型与姿势风格自然契合。
3. 服装精准着装：让模特完整穿上图 2 中的服装，忠实还原服装的剪裁款式、面料质感、褶皱受力与颜色细节。
4. 商业画质：采用高级商业影棚布光，干净明亮的中性背景，8K 超高细节与写实质感。

输出：仅输出最终的高清着装合成人像照片。`
```

### 1.3 CORS 跨域支持微调 (`internal/middleware/cors.go`)
确保后端允许姿势图库的本地预览服务（如 `http://127.0.0.1:8000`、`http://localhost:8000`、`http://127.0.0.1:5500` 等），避免因跨域导致前端 Fetch 拦截。

---

## 2. 阶段二：前端 UI 与交互升级 (`index.html`)

在 `D:\code\code_ai\ai-pose-gallery-zh\index.html` 中进行轻量但高效的功能扩展，保持原有纯原生 JS、零依赖、无构建的高性能风格。

### 2.1 顶部导航与快捷入口
1. 在顶部 `.top-actions` 区域（原“使用指南”旁），增加两个操作按钮：
   - **`🎨 AI 创作工坊`**：直接唤起 AI Studio 抽屉。
   - **`⚙️ 设置`**：唤起个人模型设置中心（填入用户自己的 Gemini / OpenAI / 豆包 Key 与 BaseURL，只保存在本地 `localStorage`）。

```html
<button class="action-btn studio-btn" id="openStudioBtn" type="button">
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2v20M2 12h20"/></svg>
  <span>AI 创作工坊</span>
</button>
<button class="action-btn settings-btn" id="openSettingsBtn" type="button" aria-label="模型设置">
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
</button>
```

### 2.2 大图详情灯箱（Viewer）打通快捷创作
在大图详情窗（`#viewer`）的侧边栏下载按钮区（`.download-stack`），增加核心 CTA（Call to Action）按钮：

```html
<!-- 在 #downloadImage 和 #openImage 旁添加 -->
<button class="studio-cta-btn" id="viewerUsePoseBtn" type="button">
  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 3v18M3 12h18"/></svg>
  <span>以当前姿势生图 / 试衣</span>
</button>
```
*交互逻辑*：用户在查看任意姿势大图并切换角度（如“斜侧”）时，点击此按钮，会自动将**当前选中的姿势图与视角**作为输入底图，直接唤出 AI 创作抽屉。

### 2.3 创作工坊抽屉（AI Studio Drawer / Modal）
在页面末尾增加 `#studioDrawer`：
- **步骤 1：姿势底图（已就绪）**：展示当前选定姿势缩略图、姿势名称、编号与视角（如 `S形站姿 · 斜侧`）；支持点击更换。
- **步骤 2：搭配服装**：
  - 模式 A（本地上传）：支持拖拽或选择本地服装图片（裙子、西装、泳衣、日常私服等）。
  - 模式 B（女装图库预设）：提供 6-8 张来自 `nvzhuang` 项目的高频预设服装卡片，点击即可一键套用；或支持输入服装图 URL。
- **步骤 3：模型选择与配置**：
  - 下拉框选择：`ComfyUI (本地离线·免费)` / `Gemini 3.1 Flash Image (推荐)` / `豆包 Seedream` / `GPT-Image (兼容格式)`。
  - 尺寸选择：`1024x1536 (2:3 纵向，契合姿势)` / `1024x1024`。
- **步骤 4：一键生成**：
  - 点击「立即生成」，前端将姿势图通过 Canvas/Fetch 转为 Base64，与服装图 Base64 及配置拼接为 JSON 发往 `http://127.0.0.1:8000/tryon`。
  - 显示 Loading 进度环与生成耗时。
  - 结果直接展示在右侧/下方的生成画廊中，提供「下载高清图」与「复制 Midjourney / SD 英文提示词」。

### 2.4 设置中心模态框 (`#settingsModal`)
复用 `nvzhuang` 经过严格安全审计的设置架构：
- 字段包括：`后端服务地址`（默认 `http://127.0.0.1:8000`）、`协议 (Provider)`、`Base URL`、`API Key`、`真实模型名`。
- **安全保障**：所有设置存储在 `localStorage.getItem('ai_pose_studio_settings')`，单次请求临时封装在 `client_config` 传递，后端绝不向磁盘写 Key，杜绝代码和凭证泄漏。

---

## 3. 核心 JS 逻辑注入指引 (`index.html`)

在 `index.html` 底部 `<script>` 中加入以下关键模块：

### 3.1 姿势图像转 Base64 辅助函数
```javascript
async function urlToBase64(url) {
  const resp = await fetch(url);
  const blob = await resp.blob();
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => resolve(reader.result);
    reader.onerror = reject;
    reader.readAsDataURL(blob);
  });
}
```

### 3.2 触发生成的核心调用逻辑
```javascript
async function executeStudioGenerate(poseImgUrl, garmentBase64, options = {}) {
  const settings = JSON.parse(localStorage.getItem('ai_pose_studio_settings') || '{}');
  const backendUrl = (settings.backendUrl || 'http://127.0.0.1:8000').replace(/\/+$/, '');
  
  // 1. 姿势底图转 Base64
  const poseBase64 = await urlToBase64(poseImgUrl);
  
  // 2. 构造请求载荷
  const payload = {
    person_image: poseBase64,      // 姿势白模作为模特骨架
    garment_image: garmentBase64,  // 服装素材
    model: options.model || settings.modelAlias || 'gemini',
    size: options.size || '1024x1536',
    keep_face: false,              // 姿势是白模素体，无需保持原人脸
    lang: 'zh'
  };

  // 3. 附带前端自定义模型配置（如有）
  if (settings.apiKey && settings.provider) {
    payload.client_config = {
      provider: settings.provider,
      base_url: settings.baseUrl || '',
      api_key: settings.apiKey,
      model: settings.realModel || ''
    };
  }

  // 4. 调用后端
  const res = await fetch(`${backendUrl}/tryon`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  });

  const data = await res.json();
  if (data.code !== 0) {
    throw new Error(data.message || '生成失败');
  }
  
  return data.data; // { image: "base64...", mime: "image/png", prompt: "..." }
}
```

### 3.3 本地历史记录 (IndexedDB)
生成成功后，写入浏览器的 IndexedDB 数据库 `PoseStudioDB`，记录字段包括：
- `id`: 时间戳
- `poseId`: 姿势标识 (如 `full_standard_P001_front`)
- `resultImage`: 生成图 Base64
- `timestamp`: 时间
- 支持单条删除、高清保存、一键清空。

---

## 4. 关键合规与安全防线（必须遵守）

1. **绝对零泄露原则**：
   - 严禁在任何 `.git` 追踪的源码、示例配置或提交记录中硬编码真实的 `sk-...`、`ARK...` 或 `AIza...`。
   - `config.yaml` 必须包含在 `.gitignore` 中。
2. **SSRF 与安全边界防护**：
   - 后端已内置针对客户端 Base URL 的安全过滤（禁止指向内网私有 IP 范围 `10.0.0.0/8`, `192.168.0.0/16`, `127.0.0.1` 等），请完整保留该安全逻辑。
3. **内容安全处理**：
   - AI Studio 默认聚焦于 **「全身单人・标准系」** 与 **「上半身单人」**，确保生成的穿搭健康、合规、契合主流社区公约。

---

## 5. 验收测试清单（Checklist）

完成修改后，请按以下顺序执行联调核验：

- [ ] **步骤 1：后端编译运行测试**
  在 `ai-pose-gallery-zh/backend` 目录下执行：
  ```bash
  go vet ./...
  go test ./...
  go run main.go
  ```
  确认服务正常监听在 `http://127.0.0.1:8000`，输出日志无 panic。

- [ ] **步骤 2：静态前端跨域与连通性**
  在 `ai-pose-gallery-zh` 根目录启动静态服务：
  ```bash
  python -m http.server 8080
  ```
  浏览器打开 `http://127.0.0.1:8080`，打开控制台，确认姿势卡片正常渲染，无 JS 语法错误。

- [ ] **步骤 3：灯箱到 AI 创作工坊联调**
  点击任意姿势卡片进入灯箱，切换视角至「斜侧」，点击「以当前姿势生图 / 试衣」：
  1. 确认工坊弹窗正常唤起，姿势图与角度标签正确回显；
  2. 上传一张服装图或选中预设服装；
  3. 配置模型或使用本地 ComfyUI / Gemini 发起请求；
  4. 检查是否成功返回生成图，且生成的人体姿态与选中的白模姿势完全匹配。

- [ ] **步骤 4：Git Diff 安全排查**
  提交前执行 `git status` 与 `git diff`，确保无真实 API Key、无大体积测试生成图误入暂存区。
