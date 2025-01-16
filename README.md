# EasyCaptcha

EasyCaptcha 是一个用于生成可定制图形验证码的 Go 包。它允许您配置验证码的各个方面，例如图像大小、字体、颜色、噪点和曲线。

## 功能
- **可定制的图像大小**：设置验证码图像的宽度和高度。
- **字体配置**：指定验证码文本的字体路径和大小。每个字符可以有略微不同的字体大小。
- **随机字符颜色**：字符以随机颜色绘制。
- **背景选项**：使用纯色或图像作为背景。
- **噪点和曲线**：为验证码添加随机噪点和曲线以增加复杂性。

## 安装

使用 EasyCaptcha 需要在系统上安装 Go。

```bash
# 克隆仓库
git clone git@github.com:goexc/easyCaptcha.git

# 进入项目目录
cd easyCaptcha

# 安装依赖
go mod tidy
```

## 使用方法

以下是如何使用 EasyCaptcha 的基本示例：

```go
package easyCaptcha

import (
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
}
```

## 配置选项

- `Width` 和 `Height`：验证码图像的尺寸。
- `FontPath`：字体文件的路径。
- `FontSize`：验证码文本的基本字体大小。
- `Text`：验证码中显示的字符。
- `NoiseCount`：添加的随机噪点数量（默认 100）。
- `CurveCount`：绘制的随机曲线数量（默认 2）。
- `BgColor`：背景颜色，如果不使用图像。
- `BgImagePath`：背景图像文件的路径。
- `CharColor`：字符颜色。

## 许可证

此项目根据 MIT 许可证授权 - 有关详细信息，请参阅 LICENSE 文件。
