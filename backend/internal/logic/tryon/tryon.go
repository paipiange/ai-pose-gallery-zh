package tryon

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	"backend/internal/imagevalidate"
)

// DoTryon 执行换装：选择模型 → 构建提示词 → 调用生成 → 返回结果
func DoTryon(ctx context.Context, personB64, garmentB64, model, size, lang string, avatarB64 ...string) (resultB64, usedModel, usedPrompt string, err error) {
	return DoTryonWithClientConfig(ctx, personB64, garmentB64, model, size, lang, nil, avatarB64...)
}

// DoTryonWithClientConfig 支持请求级临时模型配置。
// clientConfig 不会进入全局配置、缓存、数据库或日志。
func DoTryonWithClientConfig(ctx context.Context, personB64, garmentB64, model, size, lang string, clientConfig *ModelConfig, avatarB64 ...string) (resultB64, usedModel, usedPrompt string, err error) {
	// 参数校验：姿势底图必选
	if personB64 == "" {
		return "", "", "", gerror.New("姿势底图不能为空")
	}
	if err := decodeB64(personB64); err != nil {
		return "", "", "", gerror.New("姿势底图 base64 格式错误")
	}

	// 搭配服装可选
	hasGarment := strings.TrimSpace(garmentB64) != ""
	if hasGarment {
		if err := decodeB64(garmentB64); err != nil {
			return "", "", "", gerror.New("服装图 base64 格式错误")
		}
	}

	// 人物头像可选
	hasAvatar := len(avatarB64) > 0 && strings.TrimSpace(avatarB64[0]) != ""
	if hasAvatar {
		if err := decodeB64(avatarB64[0]); err != nil {
			return "", "", "", gerror.New("头像参考图 base64 格式错误")
		}
	}

	// 默认模型：Gemini（Nano Banana 2）
	if model == "" {
		model = "gemini"
	}

	var gen ImageGenerator
	if clientConfig != nil {
		gen, err = NewGeneratorFromClientConfig(clientConfig)
	} else {
		gen, err = NewGenerator(ctx, model)
	}
	if err != nil {
		return "", "", "", err
	}

	// 构建高强度姿势保留提示词（全面支持 纯姿势、姿势+头像、姿势+服装、三合一）
	promptStr := BuildPosePrompt(lang, hasGarment, hasAvatar)

	// 调用生成
	var result string
	if hasAvatar {
		result, err = gen.Generate(ctx, personB64, garmentB64, promptStr, size, avatarB64[0])
	} else {
		result, err = gen.Generate(ctx, personB64, garmentB64, promptStr, size)
	}
	if err != nil {
		// 不记录原始错误：网络错误可能包含用户自定义 Base URL，供应商响应也可能含敏感信息。
		return "", "", "", gerror.Wrap(err, "生图生成失败，请稍后重试")
	}

	return result, gen.Name(), promptStr, nil
}

// DoTryonWithValidatedImages 使用已完成安全校验的图片结构执行换装与姿势生图
func DoTryonWithValidatedImages(ctx context.Context, personImg, garmentImg *imagevalidate.ValidatedImage, model, size, lang string, keepFace bool, clientConfig *ModelConfig, avatarImg ...*imagevalidate.ValidatedImage) (resultB64, usedModel, usedPrompt string, err error) {
	if personImg == nil || len(personImg.Data) == 0 {
		return "", "", "", gerror.New("姿势底图不能为空")
	}

	personB64 := base64.StdEncoding.EncodeToString(personImg.Data)

	var garmentB64 string
	if garmentImg != nil && len(garmentImg.Data) > 0 {
		garmentB64 = base64.StdEncoding.EncodeToString(garmentImg.Data)
	}

	var avatarB64 string
	if len(avatarImg) > 0 && avatarImg[0] != nil && len(avatarImg[0].Data) > 0 {
		avatarB64 = base64.StdEncoding.EncodeToString(avatarImg[0].Data)
	}

	hasGarment := garmentB64 != ""
	hasAvatar := avatarB64 != ""

	if model == "" {
		model = "gemini"
	}

	var gen ImageGenerator
	if clientConfig != nil {
		gen, err = NewGeneratorFromClientConfig(clientConfig)
	} else {
		gen, err = NewGenerator(ctx, model)
	}
	if err != nil {
		return "", "", "", err
	}

	var promptStr string
	if !keepFace {
		// 姿势白模生图模式（默认）：强制 100% 绝对保留姿势与镜头视角
		promptStr = BuildPosePrompt(lang, hasGarment, hasAvatar)
	} else if hasGarment && hasAvatar {
		promptStr = BuildPoseAvatarGarmentTryon(lang)
	} else if hasAvatar {
		promptStr = BuildPoseAvatarTryon(lang)
	} else if hasGarment {
		isGemini := (clientConfig != nil && strings.EqualFold(clientConfig.Provider, "gemini")) || gen.Name() == "gemini"
		if isGPTImageModel(model, clientConfig, gen.Name()) {
			promptStr = BuildGPTSafety(lang)
		} else if isGemini {
			promptStr = BuildGemini(lang)
		} else {
			promptStr = BuildDefault(lang)
		}
	} else {
		promptStr = BuildPoseOnlyTryon(lang)
	}

	var result string
	if hasAvatar {
		result, err = gen.Generate(ctx, personB64, garmentB64, promptStr, size, avatarB64)
	} else {
		result, err = gen.Generate(ctx, personB64, garmentB64, promptStr, size)
	}
	if err != nil {
		return "", "", "", gerror.Wrap(err, "生图生成失败，请稍后重试")
	}

	return result, gen.Name(), promptStr, nil
}

// isGPTImageModel 判断当前请求是否使用 GPT-Image 系列模型
func isGPTImageModel(model string, clientConfig *ModelConfig, genName string) bool {
	check := func(s string) bool {
		lower := strings.ToLower(s)
		return strings.Contains(lower, "gpt-image") || strings.Contains(lower, "gpt_image")
	}
	if check(genName) || check(model) {
		return true
	}
	if clientConfig != nil && (check(clientConfig.Model) || check(clientConfig.Name)) {
		return true
	}
	return false
}
