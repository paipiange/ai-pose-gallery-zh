package v1

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTryonResIncludesImageMIME(t *testing.T) {
	encoded, err := json.Marshal(TryonRes{Image: "image-data", MIME: "image/jpeg", Model: "gemini"})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if !strings.Contains(string(encoded), `"mime":"image/jpeg"`) {
		t.Fatalf("TryonRes JSON missing mime: %s", encoded)
	}
}

func TestTryonKeepFaceFalse(t *testing.T) {
	var req TryonReq
	err := json.Unmarshal([]byte(`{"keep_face": false}`), &req)
	if err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if req.KeepFace == nil || *req.KeepFace != false {
		t.Fatalf("expected req.KeepFace to be false, got %v", req.KeepFace)
	}

	var reqDefault TryonReq
	err = json.Unmarshal([]byte(`{}`), &reqDefault)
	if err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if reqDefault.KeepFace != nil {
		t.Fatalf("expected omitted reqDefault.KeepFace to be nil, got %v", reqDefault.KeepFace)
	}
}
