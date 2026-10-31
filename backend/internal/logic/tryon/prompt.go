// Package tryon 换装提示词模板，针对不同模型优化
package tryon

// BuildPoseTryon 构建姿势引导换装专用提示词（针对 3D 人体白模/姿势素体）
func BuildPoseTryon(lang string) string {
	if lang == "zh" {
		return poseTryonZhPrompt
	}
	return poseTryonEnPrompt
}

// BuildPoseAvatarTryon 构建姿势骨架 + 服装 + 头像参考三合一生图专用提示词
func BuildPoseAvatarTryon(lang string) string {
	if lang == "zh" {
		return poseAvatarTryonZhPrompt
	}
	return poseAvatarTryonEnPrompt
}

// poseAvatarTryonEnPrompt 姿势引导换装+头像定制英文提示词
const poseAvatarTryonEnPrompt = `Task: Professional fashion editorial photograph generated from a 3D pose reference, a garment photo, and an avatar face reference.

Strict Requirements:
1. Pose & Angle Transfer: Faithfully follow the exact pose, body orientation, limb positioning, camera angle, and perspective shown in Image 1 (the wireframe/mannequin pose reference).
2. Facial Identity & Avatar Reproduction: Faithfully transfer the face, facial features, eyes, skin tone, facial expression, and identity of the person shown in Image 3 (the avatar / face reference) onto the model. Seamlessly integrate the head and neck with the body pose, maintaining realistic skin texture and anatomically accurate proportions.
3. Wardrobe: Dress the model in the exact garment and accessories shown in Image 2. Reproduce the colors, fabrics, drape, textures, neckline, and tailored silhouette faithfully.
4. Lighting & Environment: Professional studio editorial lighting, clean neutral backdrop, shallow depth of field, 8k resolution, ultra-detailed fashion magazine standard.

Output: ONLY the final dressed-up fashion photograph.`

// poseAvatarTryonZhPrompt 姿势引导换装+头像定制中文提示词
const poseAvatarTryonZhPrompt = `任务：基于人体姿势参考、服装设计图及人物头像参考，生成专业商业时尚写实人像大片。

严格要求：
1. 姿态与视角迁移：严格还原图 1（人体姿势参考素体）中的身体姿态、四肢摆放、躯干朝向、透视关系以及拍摄视角（正面/俯视/仰视/侧面/斜侧）。
2. 人物人脸与身份复刻：精准复刻图 3（头像/人脸参考图）中的人物五官特征、脸型轮廓、神态气质与肤色质感，将该头像自然无缝地融合到姿势模特的头部，头部朝向与光影同身体姿势完美契合。
3. 服装精准着装：让模特完整穿上图 2 中的服装，忠实还原服装的剪裁款式、面料质感、褶皱受力与颜色细节。
4. 商业画质：采用高级商业影棚布光，干净明亮的中性背景，8K 超高细节与写实质感。

输出：仅输出最终的高清着装合成人像照片。`

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

// BuildGPTSafety 构建 GPT-Image 专用的防封控安全提示词
func BuildGPTSafety(lang string) string {
	if lang == "zh" {
		return gptSafetyZhPrompt
	}
	return gptSafetyEnPrompt
}

// BuildGPT 构建 GPT-Image 的提示词（默认指向安全防封控模板）
func BuildGPT(lang string) string {
	return BuildGPTSafety(lang)
}

// BuildGemini 构建 Gemini（Nano Banana 2）图像模型的提示词（保留原提示词）
func BuildGemini(lang string) string {
	// Gemini 对中英文都支持良好，但换装任务英文指令更稳定
	return enPrompt
}

// BuildDefault 通用/其他模型（如豆包等）提示词（保留原提示词）
func BuildDefault(lang string) string {
	if lang == "zh" {
		return zhPrompt
	}
	return enPrompt
}

// gptSafetyEnPrompt GPT-Image-2 专用英文防封控安全提示词
const gptSafetyEnPrompt = `Task: Professional commercial fashion virtual try-on photograph.

Strict Requirements:
1. Model Identity: Maintain the adult female model from Image 1 - faithfully preserve her facial features, ethnic identity, eye shape, smile, skin tone, and hairstyle. Keep her identity consistent.
2. Wardrobe: Dress her in the exact clothing and accessories shown in Image 2. Faithfully reproduce the colors, fabric textures, patterns, and silhouette.
3. Safety & Fit: Professional high-end commercial lookbook standard. Fully clothed, modest neckline, standard tailored fit, upright standing posture, non-revealing.
4. Composition: Clean bright studio lighting, neutral minimal background, photorealistic 8k, high detail fashion editorial photograph.

Output: ONLY the final dressed-up photograph.`

// gptSafetyZhPrompt GPT-Image-2 专用中文防封控安全提示词
const gptSafetyZhPrompt = `任务：专业商业时尚虚拟试衣摄影。

严格要求：
1. 模特身份：保持第一张图片中的成年女性模特身份——忠实保留其面部特征、种族相貌、眼部轮廓、肤色与发型，保持人物相貌风格一致。
2. 服装穿搭：让人物完整穿上第二张图片中展示的服装与配饰，忠实还原服装的色彩、面料质感、图案与款式剪裁。
3. 安全与合规：采用高端商业型录摄影标准。穿着完整端庄，标准合身剪裁，端正站姿，避免过度露肤与暴露。
4. 画面效果：纯净明亮的商业影棚光照，简洁中性背景，逼真细节与质感，专业时尚画册品质。

输出：仅输出最终的完整着装照片。`

// enPrompt 英文提示词（主用，效果最稳定）
const enPrompt = `You are a professional fashion virtual try-on assistant.

Task: Dress the person in the first image wearing the garment in the second image, and output a single photorealistic full-body or half-body photo.

Strict requirements:
1. Keep the person's identity EXACTLY the same: same face, same facial features, same skin tone, same hairstyle, same body type and proportions. Do not change who they are.
2. Keep the person's original pose, body angle and orientation as much as possible; only adjust the arms/limbs if strictly needed to fit the garment naturally.
3. Put the exact garment from the second image onto the person: preserve its color, fabric texture, pattern, neckline, sleeve length, hem length and silhouette faithfully. Do not redesign or simplify the garment.
4. The garment must drape naturally on the body, follow the body contours, and cast realistic folds, shadows and highlights consistent with the original lighting.
5. Keep the original background, lighting and overall mood of the person's photo.
6. Photorealistic, high detail, professional fashion editorial quality, 8k, natural skin texture.

Output: ONLY the final dressed-up photograph, nothing else.`

// zhPrompt 中文提示词（备用）
const zhPrompt = `你是一名专业的时尚虚拟试衣助手。

任务：让第一张图片中的人物穿上第二张图片中的服装，输出一张逼真的全身或半身照片。

严格要求：
1. 人物身份必须完全一致：相同的脸、五官、肤色、发型、身材比例。不得改变人物是谁。
2. 尽量保持人物原有的姿势、身体角度和朝向；仅在必要时微调手臂/肢体以自然贴合服装。
3. 将第二张图中的服装原样穿到人物身上：忠实保留其颜色、面料质感、图案、领型、袖长、下摆长度和廓形。不得重新设计或简化服装。
4. 服装必须自然垂坠贴合身体曲线，并根据原始光照产生真实的褶皱、阴影和高光。
5. 保留人物照片原有的背景、光线和整体氛围。
6. 逼真、高细节、专业时尚大片质感、8k、自然皮肤纹理。

输出：仅输出最终的换装照片，不要其他内容。`
