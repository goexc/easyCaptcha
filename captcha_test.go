package easyCaptcha

import (
	"fmt"
	"testing"
)

func TestCaptchaGeneration(t *testing.T) {
	config := CaptchaConfig{
		Width:      240,
		Height:     80,
		FontPath:   "./font/monaco.ttf",
		FontSize:   36,
		Text:       "TEST",
		NoiseCount: 100,
		CurveCount: 2,
		BgImagePath: "./bg.jpg",
	}

	captchaInstance, err := GenerateCaptcha(config)
	if err != nil {
		t.Fatalf("Failed to generate captcha: %v", err)
	}

	// Test PNG export
	pngData, err := captchaInstance.ToPNG()
	if err != nil {
		t.Errorf("Failed to export PNG: %v", err)
	}
	if len(pngData) == 0 {
		t.Error("PNG data is empty")
	}

	// Test JPG export
	jpgData, err := captchaInstance.ToJPG()
	if err != nil {
		t.Errorf("Failed to export JPG: %v", err)
	}
	if len(jpgData) == 0 {
		t.Error("JPG data is empty")
	}

	// Test Base64 export
	base64String, err := captchaInstance.ToString()
	if err != nil {
		t.Errorf("Failed to export Base64 string: %v", err)
	}
	if len(base64String) == 0 {
		t.Error("Base64 string is empty")
	}
	fmt.Println(base64String)
}
