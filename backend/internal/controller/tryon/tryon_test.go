package tryon

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"backend/api/tryon/v1"
	"backend/internal/middleware"
	"backend/internal/security"
)

func createTestPNG(w, h int) string {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func TestControllerTryonImageValidation(t *testing.T) {
	// 启动测试服务
	s := g.Server("test-ctrl-tryon")
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.RequestID)
		group.Bind(NewV1())
	})
	s.SetPort(0)
	_ = s.Start()
	defer s.Shutdown()

	port := strings.Split(s.GetListenedAddress(), ":")[1]
	client := &http.Client{}

	// 1. 无效 Base64
	badReq := v1.TryonReq{
		PersonImage:  "!!not-valid-base64!!",
		GarmentImage: createTestPNG(300, 300),
	}
	badBody, _ := json.Marshal(badReq)
	req1, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:"+port+"/tryon", bytes.NewReader(badBody))
	req1.Header.Set("Content-Type", "application/json")
	resp1, err := client.Do(req1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp1.Body.Close()
	if resp1.StatusCode != 422 {
		t.Fatalf("expected 422, got %d", resp1.StatusCode)
	}
	b1, _ := io.ReadAll(resp1.Body)
	var res1 struct {
		Code int `json:"code"`
		Data struct {
			ErrorCode string `json:"error_code"`
		} `json:"data"`
	}
	_ = json.Unmarshal(b1, &res1)
	if res1.Code != 42201 || res1.Data.ErrorCode != "INVALID_IMAGE_BASE64" {
		t.Fatalf("unexpected res1: %+v", res1)
	}

	// 2. 文本图片（不支持类型）
	textB64 := base64.StdEncoding.EncodeToString([]byte("plain text not image"))
	textReq := v1.TryonReq{
		PersonImage:  textB64,
		GarmentImage: createTestPNG(300, 300),
	}
	textBody, _ := json.Marshal(textReq)
	req2, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:"+port+"/tryon", bytes.NewReader(textBody))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 422 {
		t.Fatalf("expected 422, got %d", resp2.StatusCode)
	}
	b2, _ := io.ReadAll(resp2.Body)
	var res2 struct {
		Code int `json:"code"`
		Data struct {
			ErrorCode string `json:"error_code"`
		} `json:"data"`
	}
	_ = json.Unmarshal(b2, &res2)
	if res2.Code != 42202 || res2.Data.ErrorCode != "UNSUPPORTED_IMAGE_TYPE" {
		t.Fatalf("unexpected res2: %+v", res2)
	}

	// 3. 尺寸过小
	smallB64 := createTestPNG(100, 100)
	smallReq := v1.TryonReq{
		PersonImage:  smallB64,
		GarmentImage: createTestPNG(300, 300),
	}
	smallBody, _ := json.Marshal(smallReq)
	req3, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:"+port+"/tryon", bytes.NewReader(smallBody))
	req3.Header.Set("Content-Type", "application/json")
	resp3, err := client.Do(req3)
	if err != nil {
		t.Fatal(err)
	}
	defer resp3.Body.Close()
	if resp3.StatusCode != 422 {
		t.Fatalf("expected 422, got %d", resp3.StatusCode)
	}
	b3, _ := io.ReadAll(resp3.Body)
	var res3 struct {
		Code int `json:"code"`
		Data struct {
			ErrorCode string `json:"error_code"`
		} `json:"data"`
	}
	_ = json.Unmarshal(b3, &res3)
	if res3.Code != 42204 || res3.Data.ErrorCode != "IMAGE_DIMENSIONS_INVALID" {
		t.Fatalf("unexpected res3: %+v", res3)
	}

	// 4. 可选头像：头像非合法图片时报错
	badAvatarReq := v1.TryonReq{
		PersonImage:  createTestPNG(300, 300),
		GarmentImage: createTestPNG(300, 300),
		AvatarImage:  "data:image/png;base64,invalid-avatar-b64",
	}
	badAvatarBody, _ := json.Marshal(badAvatarReq)
	req4, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:"+port+"/tryon", bytes.NewReader(badAvatarBody))
	req4.Header.Set("Content-Type", "application/json")
	resp4, err := client.Do(req4)
	if err != nil {
		t.Fatal(err)
	}
	defer resp4.Body.Close()
	if resp4.StatusCode != 422 {
		t.Fatalf("expected 422, got %d", resp4.StatusCode)
	}

	// 5. 可选头像：头像低于放宽后的下限（20x20 < 32x32）报错 AVATAR_DIMENSIONS_INVALID (42206)
	tinyAvatarReq := v1.TryonReq{
		PersonImage:  createTestPNG(300, 300),
		GarmentImage: createTestPNG(300, 300),
		AvatarImage:  createTestPNG(20, 20),
	}
	tinyAvatarBody, _ := json.Marshal(tinyAvatarReq)
	req5, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:"+port+"/tryon", bytes.NewReader(tinyAvatarBody))
	req5.Header.Set("Content-Type", "application/json")
	resp5, err := client.Do(req5)
	if err != nil {
		t.Fatal(err)
	}
	defer resp5.Body.Close()
	if resp5.StatusCode != 422 {
		t.Fatalf("expected 422, got %d", resp5.StatusCode)
	}
	b5, _ := io.ReadAll(resp5.Body)
	var res5 struct {
		Code int `json:"code"`
		Data struct {
			ErrorCode string `json:"error_code"`
		} `json:"data"`
	}
	_ = json.Unmarshal(b5, &res5)
	if res5.Code != 42206 || res5.Data.ErrorCode != "AVATAR_DIMENSIONS_INVALID" {
		t.Fatalf("unexpected res5: %+v", res5)
	}

	// 6. 可选头像：头像在放宽范围（如 64x64）通过图片校验（不再报 422 尺寸错误）
	smallAvatarReq := v1.TryonReq{
		PersonImage:  createTestPNG(300, 300),
		GarmentImage: createTestPNG(300, 300),
		AvatarImage:  createTestPNG(64, 64),
	}
	smallAvatarBody, _ := json.Marshal(smallAvatarReq)
	req6, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:"+port+"/tryon", bytes.NewReader(smallAvatarBody))
	req6.Header.Set("Content-Type", "application/json")
	resp6, err := client.Do(req6)
	if err != nil {
		t.Fatal(err)
	}
	defer resp6.Body.Close()
	// 如果由于未配置模型或者缺少外部 API 报错，状态码不会是 422（图片校验错误），说明成功通过了图片安全尺寸校验
	if resp6.StatusCode == 422 {
		b6, _ := io.ReadAll(resp6.Body)
		t.Fatalf("expected 64x64 avatar to pass image validation, got 422: %s", string(b6))
	}

	// 7. 可选服装：省略服装（纯姿势生图），成功通过图片校验（状态码不应为 422）
	poseOnlyReq := v1.TryonReq{
		PersonImage: createTestPNG(300, 300),
	}
	poseOnlyBody, _ := json.Marshal(poseOnlyReq)
	req7, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:"+port+"/tryon", bytes.NewReader(poseOnlyBody))
	req7.Header.Set("Content-Type", "application/json")
	resp7, err := client.Do(req7)
	if err != nil {
		t.Fatal(err)
	}
	defer resp7.Body.Close()
	if resp7.StatusCode == 422 {
		b7, _ := io.ReadAll(resp7.Body)
		t.Fatalf("expected pose-only request to pass image validation, got 422: %s", string(b7))
	}

	// 8. 可选服装：省略服装但提供头像（姿势+头像生图），成功通过图片校验（状态码不应为 422）
	poseAvatarReq := v1.TryonReq{
		PersonImage: createTestPNG(300, 300),
		AvatarImage: createTestPNG(64, 64),
	}
	poseAvatarBody, _ := json.Marshal(poseAvatarReq)
	req8, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:"+port+"/tryon", bytes.NewReader(poseAvatarBody))
	req8.Header.Set("Content-Type", "application/json")
	resp8, err := client.Do(req8)
	if err != nil {
		t.Fatal(err)
	}
	defer resp8.Body.Close()
	if resp8.StatusCode == 422 {
		b8, _ := io.ReadAll(resp8.Body)
		t.Fatalf("expected pose+avatar request to pass image validation, got 422: %s", string(b8))
	}

	// 9. 提供了服装但为非法 base64，应触发 422 校验失败
	badGarmentReq := v1.TryonReq{
		PersonImage:  createTestPNG(300, 300),
		GarmentImage: "!!not-valid-base64!!",
	}
	badGarmentBody, _ := json.Marshal(badGarmentReq)
	req9, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:"+port+"/tryon", bytes.NewReader(badGarmentBody))
	req9.Header.Set("Content-Type", "application/json")
	resp9, err := client.Do(req9)
	if err != nil {
		t.Fatal(err)
	}
	defer resp9.Body.Close()
	if resp9.StatusCode != 422 {
		t.Fatalf("expected 422 for bad garment, got %d", resp9.StatusCode)
	}
}

func TestControllerTryonQuotaAndConcurrency(t *testing.T) {
	secService := security.GetService()
	// 设置配额为 1 进行测试
	secService.DailyQuotaTracker = security.NewDailyQuotaTracker(1)
	secService.ConcurrencyManager = security.NewConcurrencyManager(1, 1)

	s := g.Server("test-ctrl-quota")
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.RequestID)
		group.Bind(NewV1())
	})
	s.SetPort(0)
	_ = s.Start()
	defer s.Shutdown()

	port := strings.Split(s.GetListenedAddress(), ":")[1]
	client := &http.Client{}

	validPerson := createTestPNG(300, 300)
	validGarment := createTestPNG(300, 300)

	// 先直接耗尽配额
	secService.DailyQuotaTracker.CheckAndDeduct("anonymous")

	// 发起请求，应当直接触发 DAILY_QUOTA_EXCEEDED (42902)
	tryReq := v1.TryonReq{
		PersonImage:  validPerson,
		GarmentImage: validGarment,
	}
	reqBody, _ := json.Marshal(tryReq)
	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:"+port+"/tryon", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", resp.StatusCode)
	}
	b, _ := io.ReadAll(resp.Body)
	var res struct {
		Code int `json:"code"`
		Data struct {
			ErrorCode string `json:"error_code"`
		} `json:"data"`
	}
	_ = json.Unmarshal(b, &res)
	if res.Code != 42902 || res.Data.ErrorCode != "DAILY_QUOTA_EXCEEDED" {
		t.Fatalf("unexpected quota res: %+v", res)
	}
}

func init() {
	// 确保测试环境未设置任何破坏性环境变量
	os.Unsetenv("TRYON_API_KEYS")
}
