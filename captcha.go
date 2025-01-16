package easyCaptcha

import (
	"bytes"
	"encoding/base64"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math/rand"
	"time"
)



// CaptchaConfig 包含生成验证码的配置参数
// CaptchaConfig holds the configuration for the captcha
type CaptchaConfig struct {
	Width       int         // 验证码图片的宽度
	Height      int         // 验证码图片的高度
	FontPath    string      // 字体文件的路径
	FontSize    float64     // 字体大小
	CharColor   color.Color // 字符颜色
	BgColor     color.Color // 背景颜色
	BgImagePath string      // 背景图片的路径
	Text        string      // 验证码中显示的文本
	NoiseCount  int         // 噪点数量
	CurveCount  int         // 曲线数量
}

// Captcha 包含验证码图像数据
// Captcha holds the captcha image data
type Captcha struct {
	img image.Image
}

// ToPNG 将验证码导出为PNG格式
// ToPNG exports the captcha as a PNG
func (c *Captcha) ToPNG() ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, c.img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ToJPG 将验证码导出为JPG格式
// ToJPG exports the captcha as a JPG
func (c *Captcha) ToJPG() ([]byte, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, c.img, nil); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ToString 将验证码导出为Base64编码的字符串
// ToString exports the captcha as a base64-encoded string
func (c *Captcha) ToString() (string, error) {
	pngData, err := c.ToPNG()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pngData), nil
}

// GenerateCaptcha 根据提供的配置生成验证码图像
// GenerateCaptcha creates a captcha image based on the provided configuration
// 返回一个Captcha实例
func GenerateCaptcha(config CaptchaConfig) (*Captcha, error) {
	// Set default values if not provided
	if config.FontPath == "" {
		config.FontPath = "./arial.ttf"
	}
	if config.FontSize == 0 {
		config.FontSize = 36
	}
	if config.CharColor == (color.Color)(nil) {
		config.CharColor = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Default to red
	}
	if config.BgColor == (color.Color)(nil) {
		config.BgColor = color.RGBA{R: 255, G: 255, B: 255, A: 255} // Default to white
	}
	if config.NoiseCount == 0 {
		config.NoiseCount = 100 // Default noise count
	}
	if config.CurveCount == 0 {
		config.CurveCount = 2 // Default curve count
	}

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	dc := gg.NewContext(config.Width, config.Height)

	// Load background image if provided
	if config.BgImagePath != "" {
		img, err := gg.LoadImage(config.BgImagePath)
		if err != nil {
			return nil, err
		}
		dc.DrawImage(img, 0, 0)
	} else {
		dc.SetColor(config.BgColor)
		dc.Clear()
	}

	charSpacing := float64(config.Width) / float64(len(config.Text)+1)
	for i, char := range config.Text {
		// Randomize rotation angle between -30 to 30 degrees
		angle := r.Float64()*60 - 30

		// Randomize font size within ±5 of the given font size
		fontSize := config.FontSize + r.Float64()*10 - 5
		if err := dc.LoadFontFace(config.FontPath, fontSize); err != nil {
			return nil, err
		}

		// Set random character color
		charColor := color.RGBA{
			R: uint8(r.Intn(256)),
			G: uint8(r.Intn(256)),
			B: uint8(r.Intn(256)),
			A: 255,
		}
		dc.SetColor(charColor)

		// Calculate position for each character
		x := charSpacing * float64(i+1)
		y := float64(config.Height) / 2

		dc.Push()
		dc.RotateAbout(gg.Radians(angle), x, y)
		dc.DrawStringAnchored(string(char), x, y, 0.5, 0.5)
		dc.Pop()
	}

	// Add noise
	for i := 0; i < config.NoiseCount; i++ {
		x := r.Intn(config.Width)
		y := r.Intn(config.Height)
		noiseColor := color.RGBA{
			R: uint8(r.Intn(256)),
			G: uint8(r.Intn(256)),
			B: uint8(r.Intn(256)),
			A: 255,
		}
		dc.SetColor(noiseColor)
		dc.SetPixel(x, y)
	}

	// Draw random curves
	for i := 0; i < config.CurveCount; i++ {
		curveColor := color.RGBA{
			R: uint8(r.Intn(256)),
			G: uint8(r.Intn(256)),
			B: uint8(r.Intn(256)),
			A: 255,
		}
		dc.SetColor(curveColor)
		dc.SetLineWidth(1 + r.Float64()*2)
		dc.MoveTo(float64(r.Intn(config.Width)), float64(r.Intn(config.Height)))
		for j := 0; j < 3; j++ {
			dc.QuadraticTo(
				float64(r.Intn(config.Width)), float64(r.Intn(config.Height)),
				float64(r.Intn(config.Width)), float64(r.Intn(config.Height)),
			)
		}
		dc.Stroke()
	}

	return &Captcha{img: dc.Image()}, nil
}
