// Package tryon 换装与姿势生图提示词模板，针对不同模型与组合优化
package tryon

// BuildPosePrompt 根据是否有服装和头像构建对应的高强度姿势保留提示词
func BuildPosePrompt(lang string, hasGarment, hasAvatar bool) string {
	if hasGarment && hasAvatar {
		return BuildPoseAvatarGarmentTryon(lang)
	} else if hasGarment {
		return BuildPoseGarmentTryon(lang)
	} else if hasAvatar {
		return BuildPoseAvatarTryon(lang)
	}
	return BuildPoseOnlyTryon(lang)
}

// BuildPoseTryon 构建姿势引导换装专用提示词（姿势素体+服装）
func BuildPoseTryon(lang string) string {
	return BuildPoseGarmentTryon(lang)
}

// BuildPoseGarmentTryon 构建姿势骨架 + 服装素材组合生图提示词
func BuildPoseGarmentTryon(lang string) string {
	if lang == "zh" {
		return poseGarmentZhPrompt
	}
	return poseGarmentEnPrompt
}

// BuildPoseAvatarTryon 构建姿势骨架 + 人物头像组合生图提示词（无服装）
func BuildPoseAvatarTryon(lang string) string {
	if lang == "zh" {
		return poseAvatarZhPrompt
	}
	return poseAvatarEnPrompt
}

// BuildPoseAvatarGarmentTryon 构建姿势骨架 + 服装素材 + 头像参考三合一生图提示词
func BuildPoseAvatarGarmentTryon(lang string) string {
	if lang == "zh" {
		return poseAvatarGarmentZhPrompt
	}
	return poseAvatarGarmentEnPrompt
}

// BuildPoseOnlyTryon 构建纯姿势骨架生图提示词（无头像、无服装）
func BuildPoseOnlyTryon(lang string) string {
	if lang == "zh" {
		return poseOnlyZhPrompt
	}
	return poseOnlyEnPrompt
}

// ==================== 1. 姿势 + 服装 + 头像三合一 ====================

const poseAvatarGarmentEnPrompt = `Task: Professional commercial fashion editorial photograph generated from a 3D pose reference (Image 1), a garment photo (Image 2), and an avatar face reference (Image 3).

CRITICAL & MANDATORY REQUIREMENTS:
1. ABSOLUTE POSE & ANGLE PRESERVATION (Pose & Angle Transfer - HIGHEST MANDATORY PRIORITY): The generated model MUST 100% strictly duplicate and match the EXACT 3D human pose, skeletal orientation, bone alignment, limb angles, spine curvature, torso direction, hand and finger gestures, leg placement, and camera perspective (eye-level/high-angle/low-angle/side/diagonal) shown in Image 1 (the wireframe/mannequin pose reference). ABSOLUTELY FORBIDDEN: DO NOT alter, relax, mirror, flip, rotate, or re-orient any limb, joint, or the camera angle. Treat Image 1 as an unchangeable 3D spatial wireframe and pose ground truth.
2. FACIAL IDENTITY & AVATAR REPRODUCTION: Faithfully transfer the face, eyes, nose, mouth, facial contours, expression, skin tone, and identity of the person shown in Image 3 (the avatar / face reference) onto the model. Seamlessly integrate the head and neck onto the body in Image 1, matching the exact head angle and gaze direction.
3. WARDROBE FIDELITY: Dress the model in the exact garment shown in Image 2. Reproduce the colors, fabrics, drape, textures, neckline, sleeves, and tailored silhouette faithfully, fitted naturally over Image 1's pose.
4. COMMERCIAL EDITORIAL QUALITY: Professional high-end studio lighting, clean neutral backdrop, shallow depth of field, photorealistic 8k, natural human skin texture, ultra-detailed fashion magazine standard.

Output: ONLY the final photograph.`

const poseAvatarGarmentZhPrompt = `任务：基于人体姿势参考（图 1）、服装设计图（图 2）及人物头像参考（图 3），生成专业商业时尚写实人像大片。

严格核心要求：
1. 【最高绝对优先级核心指令：人体姿势与骨架必须 100% 绝对保留】生成的真人模特必须严格 100% 复刻并完全贴合图 1（人体姿势参考素体）中的人体三维骨骼姿态、四肢摆放角度、关节点位置、躯干朝向、头部倾角、脊柱弧度、手部与手指姿态、双腿姿势以及镜头的拍摄透视视角（正面/俯视/仰视/侧视/斜角等）。绝对严禁擅自修改肢体动作、擅自微调四肢角度、擅自旋转或镜像身体、擅自改变摄像机拍摄视角！图 1 为不可变更的绝对几何骨架底稿，生成模特的身体线条与构图视角必须与图 1 严丝合缝、完全吻合！
2. 人物面部与身份复刻：精准复刻图 3（头像/人脸参考图）中的人物五官特征、脸型轮廓、神态气质与肤色质感，自然无缝地融合到姿势模特的头部，头部朝向与视线方向同姿势完全契合。
3. 服装精准着装：让模特完整穿上图 2 中的服装，忠实还原服装的款式剪裁、面料质感、褶皱受力与颜色细节。
4. 商业画质：采用高级商业影棚布光，干净明亮的中性背景，8K 超高细节与自然逼真的皮肤毛孔质感。

输出：仅输出最终的高清合成人像照片。`

// ==================== 2. 姿势 + 服装 ====================

const poseGarmentEnPrompt = `Task: Professional commercial fashion editorial photograph generated from a 3D pose reference (Image 1) and a garment photo (Image 2).

CRITICAL & MANDATORY REQUIREMENTS:
1. ABSOLUTE POSE & ANGLE PRESERVATION (Pose & Angle Transfer - HIGHEST MANDATORY PRIORITY): The generated model MUST 100% strictly duplicate and match the EXACT 3D human pose, skeletal orientation, bone alignment, limb angles, spine curvature, torso direction, hand and finger gestures, leg placement, and camera perspective (eye-level/high-angle/low-angle/side/diagonal) shown in Image 1 (the wireframe/mannequin pose reference). ABSOLUTELY FORBIDDEN: DO NOT alter, relax, mirror, flip, rotate, or re-orient any limb, joint, or the camera angle. Treat Image 1 as an unchangeable 3D spatial wireframe and pose ground truth.
2. REALISTIC HUMAN MODEL: Transform the wireframe/mannequin in Image 1 into an attractive, photorealistic adult fashion model with natural skin texture, lifelike facial expression, stylish hairstyle, and anatomically accurate proportions.
3. WARDROBE FIDELITY: Dress the model in the exact garment shown in Image 2. Reproduce the colors, fabrics, drape, textures, neckline, sleeves, and tailored silhouette faithfully, fitted naturally over Image 1's pose.
4. COMMERCIAL EDITORIAL QUALITY: Professional high-end studio lighting, clean neutral backdrop, shallow depth of field, photorealistic 8k, ultra-detailed fashion magazine standard.

Output: ONLY the final photograph.`

const poseGarmentZhPrompt = `任务：基于人体姿势参考（图 1）与服装设计图（图 2），生成专业商业时尚写实人像大片。

严格核心要求：
1. 【最高绝对优先级核心指令：人体姿势与骨架必须 100% 绝对保留】生成的真人模特必须严格 100% 复刻并完全贴合图 1（人体姿势参考素体）中的人体三维骨骼姿态、四肢摆放角度、关节点位置、躯干朝向、头部倾角、脊柱弧度、手部与手指姿态、双腿姿势以及镜头的拍摄透视视角（正面/俯视/仰视/侧视/斜角等）。绝对严禁擅自修改肢体动作、擅自微调四肢角度、擅自旋转或镜像身体、擅自改变摄像机拍摄视角！图 1 为不可变更的绝对几何骨架底稿，生成模特的身体线条与构图视角必须与图 1 严丝合缝、完全吻合！
2. 真人模特渲染：将姿势素体渲染为相貌端庄自然、皮肤质感逼真、五官立体的成年时尚模特，发型与姿势风格自然契合。
3. 服装精准着装：让模特完整穿上图 2 中的服装，忠实还原服装的款式剪裁、面料质感、褶皱受力与颜色细节。
4. 商业画质：采用高级商业影棚布光，干净明亮的中性背景，8K 超高细节与自然逼真的皮肤毛孔质感。

输出：仅输出最终的高清合成人像照片。`

// ==================== 3. 姿势 + 头像（无服装） ====================

const poseAvatarEnPrompt = `Task: Professional commercial fashion editorial photograph generated from a 3D pose reference (Image 1) and an avatar face reference (Image 2).

CRITICAL & MANDATORY REQUIREMENTS:
1. ABSOLUTE POSE & ANGLE PRESERVATION (Pose & Angle Transfer - HIGHEST MANDATORY PRIORITY): The generated model MUST 100% strictly duplicate and match the EXACT 3D human pose, skeletal orientation, bone alignment, limb angles, spine curvature, torso direction, hand and finger gestures, leg placement, and camera perspective (eye-level/high-angle/low-angle/side/diagonal) shown in Image 1 (the wireframe/mannequin pose reference). ABSOLUTELY FORBIDDEN: DO NOT alter, relax, mirror, flip, rotate, or re-orient any limb, joint, or the camera angle. Treat Image 1 as an unchangeable 3D spatial wireframe and pose ground truth.
2. FACIAL IDENTITY & AVATAR REPRODUCTION: Faithfully transfer the face, eyes, nose, mouth, facial contours, expression, skin tone, and identity of the person shown in Image 2 (the avatar / face reference) onto the model. Seamlessly integrate the head and neck onto the body in Image 1, matching the exact head angle and gaze direction.
3. TASTEFUL MODERN WARDROBE: Dress the model in an elegant, stylish, well-tailored modern outfit that perfectly complements the pose and body contour, with natural fabric drape and realistic lighting folds.
4. COMMERCIAL EDITORIAL QUALITY: Professional high-end studio lighting, clean neutral backdrop, shallow depth of field, photorealistic 8k, natural human skin texture, ultra-detailed fashion magazine standard.

Output: ONLY the final photograph.`

const poseAvatarZhPrompt = `任务：基于人体姿势参考（图 1）与人物头像参考（图 2），生成专业商业时尚写实人像大片。

严格核心要求：
1. 【最高绝对优先级核心指令：人体姿势与骨架必须 100% 绝对保留】生成的真人模特必须严格 100% 复刻并完全贴合图 1（人体姿势参考素体）中的人体三维骨骼姿态、四肢摆放角度、关节点位置、躯干朝向、头部倾角、脊柱弧度、手部与手指姿态、双腿姿势以及镜头的拍摄透视视角（正面/俯视/仰视/侧视/斜角等）。绝对严禁擅自修改肢体动作、擅自微调四肢角度、擅自旋转或镜像身体、擅自改变摄像机拍摄视角！图 1 为不可变更的绝对几何骨架底稿，生成模特的身体线条与构图视角必须与图 1 严丝合缝、完全吻合！
2. 人物面部与身份复刻：精准复刻图 2（头像/人脸参考图）中的人物五官特征、脸型轮廓、神态气质与肤色质感，自然无缝地融合到姿势模特的头部，头部朝向与视线方向同姿势完全契合。
3. 时尚着装搭配：为模特搭配一套剪裁得体、优雅大方的现代时尚服饰，面料垂坠自然，贴合姿态线条。
4. 商业画质：采用高级商业影棚布光，干净明亮的中性背景，8K 超高细节与自然逼真的皮肤毛孔质感。

输出：仅输出最终的高清合成人像照片。`

// ==================== 4. 纯姿势生图（无头像、无服装） ====================

const poseOnlyEnPrompt = `Task: Professional commercial fashion editorial photograph generated from a 3D pose reference (Image 1).

CRITICAL & MANDATORY REQUIREMENTS:
1. ABSOLUTE POSE & ANGLE PRESERVATION (Pose & Angle Transfer - HIGHEST MANDATORY PRIORITY): The generated model MUST 100% strictly duplicate and match the EXACT 3D human pose, skeletal orientation, bone alignment, limb angles, spine curvature, torso direction, hand and finger gestures, leg placement, and camera perspective (eye-level/high-angle/low-angle/side/diagonal) shown in Image 1 (the wireframe/mannequin pose reference). ABSOLUTELY FORBIDDEN: DO NOT alter, relax, mirror, flip, rotate, or re-orient any limb, joint, or the camera angle. Treat Image 1 as an unchangeable 3D spatial wireframe and pose ground truth.
2. REALISTIC HUMAN MODEL: Transform the wireframe/mannequin in Image 1 into an attractive, photorealistic adult fashion model with natural human skin texture, lifelike facial expression, stylish hairstyle, and anatomically accurate proportions.
3. TASTEFUL MODERN WARDROBE: Dress the model in an elegant, stylish, well-tailored modern outfit that perfectly complements the pose and body contour, with natural fabric drape and realistic lighting folds.
4. COMMERCIAL EDITORIAL QUALITY: Professional high-end studio lighting, clean neutral backdrop, shallow depth of field, photorealistic 8k, ultra-detailed fashion magazine standard.

Output: ONLY the final photograph.`

const poseOnlyZhPrompt = `任务：基于人体姿势参考（图 1），生成专业商业时尚写实人像大片。

严格核心要求：
1. 【最高绝对优先级核心指令：人体姿势与骨架必须 100% 绝对保留】生成的真人模特必须严格 100% 复刻并完全贴合图 1（人体姿势参考素体）中的人体三维骨骼姿态、四肢摆放角度、关节点位置、躯干朝向、头部倾角、脊柱弧度、手部与手指姿态、双腿姿势以及镜头的拍摄透视视角（正面/俯视/仰视/侧视/斜角等）。绝对严禁擅自修改肢体动作、擅自微调四肢角度、擅自旋转或镜像身体、擅自改变摄像机拍摄视角！图 1 为不可变更的绝对几何骨架底稿，生成模特的身体线条与构图视角必须与图 1 严丝合缝、完全吻合！
2. 真人模特渲染：将姿势素体渲染为相貌端庄自然、皮肤质感逼真、五官立体的成年时尚模特，发型与姿势风格自然契合。
3. 时尚着装搭配：为模特搭配一套剪裁得体、优雅大方的现代时尚服饰，面料垂坠自然，贴合姿态线条。
4. 商业画质：采用高级商业影棚布光，干净明亮的中性背景，8K 超高细节与自然逼真的皮肤毛孔质感。

输出：仅输出最终的高清合成人像照片。`

// ==================== 保留换装与通用提示词 ====================

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
	return enPrompt
}

// BuildDefault 通用/其他模型（如豆包等）提示词（保留原提示词）
func BuildDefault(lang string) string {
	if lang == "zh" {
		return zhPrompt
	}
	return enPrompt
}

const gptSafetyEnPrompt = `Task: Professional commercial fashion virtual try-on photograph.

Strict Requirements:
1. Model Identity: Maintain the adult female model from Image 1 - faithfully preserve her facial features, ethnic identity, eye shape, smile, skin tone, and hairstyle. Keep her identity consistent.
2. Wardrobe: Dress her in the exact clothing and accessories shown in Image 2. Faithfully reproduce the colors, fabric textures, patterns, and silhouette.
3. Safety & Fit: Professional high-end commercial lookbook standard. Fully clothed, modest neckline, standard tailored fit, upright standing posture, non-revealing.
4. Composition: Clean bright studio lighting, neutral minimal background, photorealistic 8k, high detail fashion editorial photograph.

Output: ONLY the final dressed-up photograph.`

const gptSafetyZhPrompt = `任务：专业商业时尚虚拟试衣摄影。

严格要求：
1. 模特身份：保持第一张图片中的成年女性模特身份——忠实保留其面部特征、种族相貌、眼部轮廓、肤色与发型，保持人物相貌风格一致。
2. 服装穿搭：让人物完整穿上第二张图片中展示的服装与配饰，忠实还原服装的色彩、面料质感、图案与款式剪裁。
3. 安全与合规：采用高端商业型录摄影标准。穿着完整端庄，标准合身剪裁，端正站姿，避免过度露肤与暴露。
4. 画面效果：纯净明亮的商业影棚光照，简洁中性背景，逼真细节与质感，专业时尚画册品质。

输出：仅输出最终的完整着装照片。`

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
