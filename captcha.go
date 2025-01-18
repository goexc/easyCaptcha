package easyCaptcha

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "image"
    "image/color"
    "image/jpeg"
    "image/png"
    "math"
    "math/rand"
    "os"
    "path/filepath"
    "runtime"
    "time"

    "github.com/fogleman/gg"
)

// CaptchaConfig 包含生成验证码的配置参数
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
    LineWidth   float64     // 线条宽度
}

// Captcha 包含验证码图像数据
type Captcha struct {
    img  image.Image // 验证码图像
    text string      // 验证码文本
}

// ToPNG 将验证码导出为PNG格式
// savePath为空时，将图片保存在主文件所在目录下，文件名为验证码内容+时间戳
func (c *Captcha) ToPNG(savePath string) ([]byte, error) {
    var buf bytes.Buffer
    if err := png.Encode(&buf, c.img); err != nil {
        return nil, fmt.Errorf("PNG编码失败: %v", err)
    }
    
    data := buf.Bytes()
    
    if savePath == "" {
        _, filename, _, _ := runtime.Caller(0)
        dir := filepath.Dir(filename)
        timestamp := time.Now().UnixNano() / int64(time.Millisecond)
        savePath = filepath.Join(dir, fmt.Sprintf("%s_%d.png", c.text, timestamp))
    }
    
    if err := os.WriteFile(savePath, data, 0666); err != nil {
        return nil, fmt.Errorf("保存PNG文件失败: %v", err)
    }
    return data, nil
}

// ToJPG 将验证码导出为JPG格式
// savePath为空时，将图片保存在主文件所在目录下，文件名为验证码内容+时间戳
func (c *Captcha) ToJPG(savePath string) ([]byte, error) {
    var buf bytes.Buffer
    if err := jpeg.Encode(&buf, c.img, nil); err != nil {
        return nil, fmt.Errorf("JPG编码失败: %v", err)
    }
    
    data := buf.Bytes()
    
    if savePath == "" {
        _, filename, _, _ := runtime.Caller(0)
        dir := filepath.Dir(filename)
        timestamp := time.Now().UnixNano() / int64(time.Millisecond)
        savePath = filepath.Join(dir, fmt.Sprintf("%s_%d.jpg", c.text, timestamp))
    }
    
    if err := os.WriteFile(savePath, data, 0666); err != nil {
        return nil, fmt.Errorf("保存JPG文件失败: %v", err)
    }
    return data, nil
}

// ToString 将验证码导出为Base64编码的字符串
func (c *Captcha) ToString() (string, error) {
    pngData, err := c.ToPNG("")
    if err != nil {
        return "", fmt.Errorf("生成Base64字符串失败: %v", err)
    }
    return base64.StdEncoding.EncodeToString(pngData), nil
}

// getDefaultFontPath 返回默认字体文件的绝对路径
func getDefaultFontPath() string {
    _, filename, _, _ := runtime.Caller(0)
    return filepath.Join(filepath.Dir(filename), "font", "default.ttf")
}

// GenerateCaptcha 根据提供的配置生成验证码图像
// 返回一个Captcha实例
func GenerateCaptcha(config CaptchaConfig) (*Captcha, error) {
    // 设置默认值
    if config.FontPath == "" {
        config.FontPath = getDefaultFontPath()
    }
    if config.FontSize == 0 {
        config.FontSize = 36 // 默认字体大小
    }
    if config.CharColor == (color.Color)(nil) {
        config.CharColor = color.RGBA{R: 255, G: 0, B: 0, A: 255} // 默认红色
    }
    if config.BgColor == (color.Color)(nil) {
        config.BgColor = color.RGBA{R: 255, G: 255, B: 255, A: 255} // 默认白色
    }
    if config.NoiseCount == 0 {
        config.NoiseCount = 100 // 默认噪点数量
    }
    if config.CurveCount == 0 {
        config.CurveCount = 2 // 默认曲线数量
    }
    if config.LineWidth == 0 {
        config.LineWidth = 1 // 默认线条宽度为1像素
    }

    source := rand.NewSource(time.Now().UnixNano())
    r := rand.New(source)
    dc := gg.NewContext(config.Width, config.Height)

    // 设置全局线条宽度
    dc.SetLineWidth(config.LineWidth)

    // 加载背景图片（如果提供）
    if config.BgImagePath != "" {
        img, err := gg.LoadImage(config.BgImagePath)
        if err != nil {
            return nil, fmt.Errorf("加载背景图片失败: %v", err)
        }
        dc.DrawImage(img, 0, 0)
    } else {
        dc.SetColor(config.BgColor)
        dc.Clear()
    }

    // 计算字符间距
    charSpacing := float64(config.Width) / float64(len(config.Text)+1)
    for i, char := range config.Text {
        // 随机旋转角度（-30到30度）
        angle := r.Float64()*60 - 30

        // 随机字体大小（在给定大小的±5范围内）
        fontSize := config.FontSize + r.Float64()*10 - 5
        if err := dc.LoadFontFace(config.FontPath, fontSize); err != nil {
            return nil, fmt.Errorf("加载字体失败: %v", err)
        }

        // 设置随机字符颜色
        charColor := color.RGBA{
            R: uint8(r.Intn(256)),
            G: uint8(r.Intn(256)),
            B: uint8(r.Intn(256)),
            A: 255,
        }
        dc.SetColor(charColor)

        // 计算每个字符的位置
        x := charSpacing * float64(i+1)
        y := float64(config.Height) / 2

        dc.Push()
        dc.RotateAbout(gg.Radians(angle), x, y)
        dc.DrawStringAnchored(string(char), x, y, 0.5, 0.5)
        dc.Pop()
    }

    // 添加噪点
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

    // 绘制随机曲线
    maxTotalLength := float64(config.Width + config.Height)
    totalLength := 0.0

    for i := 0; i < config.CurveCount && totalLength < maxTotalLength; i++ {
        curveColor := color.RGBA{
            R: uint8(r.Intn(256)),
            G: uint8(r.Intn(256)),
            B: uint8(r.Intn(256)),
            A: 255,
        }
        dc.SetColor(curveColor)

        // 起始点
        startX := float64(r.Intn(config.Width))
        startY := float64(r.Intn(config.Height))
        dc.MoveTo(startX, startY)
        lastX, lastY := startX, startY

        // 绘制曲线段并检查长度
        remainingLength := maxTotalLength - totalLength
        segmentCount := 3 // 曲线段数量
        maxSegmentLength := remainingLength / float64(segmentCount)

        for j := 0; j < segmentCount; j++ {
            // 在最大段长度范围内生成控制点和终点
            controlX := lastX + (r.Float64()*2-1)*maxSegmentLength
            controlY := lastY + (r.Float64()*2-1)*maxSegmentLength
            endX := lastX + (r.Float64()*2-1)*maxSegmentLength
            endY := lastY + (r.Float64()*2-1)*maxSegmentLength

            // 将点限制在图片边界内
            controlX = math.Max(0, math.Min(float64(config.Width), controlX))
            controlY = math.Max(0, math.Min(float64(config.Height), controlY))
            endX = math.Max(0, math.Min(float64(config.Width), endX))
            endY = math.Max(0, math.Min(float64(config.Height), endY))

            dc.QuadraticTo(controlX, controlY, endX, endY)

            // 计算段长度（使用控制点近似）
            segmentLength := math.Sqrt(math.Pow(controlX-lastX, 2) + math.Pow(controlY-lastY, 2)) +
                math.Sqrt(math.Pow(endX-controlX, 2) + math.Pow(endY-controlY, 2))
            totalLength += segmentLength

            lastX, lastY = endX, endY

            // 如果超过最大长度则停止
            if totalLength >= maxTotalLength {
                break
            }
        }
        dc.Stroke()
    }

    return &Captcha{img: dc.Image(), text: config.Text}, nil
}
