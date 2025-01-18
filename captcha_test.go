package easyCaptcha

import (
    "os"
    "path/filepath"
    "runtime"
    "strings"
    "testing"
)

func TestPNGCaptcha(t *testing.T) {
    // 创建临时目录用于测试文件
    tempDir, err := os.MkdirTemp("", "captcha_test_png")
    if err != nil {
        t.Fatalf("创建临时目录失败: %v", err)
    }
    defer os.RemoveAll(tempDir) // 测试完成后清理

    config := CaptchaConfig{
        Width:      240,
        Height:     80,
        FontPath:   "./font/default.ttf",
        FontSize:   36,
        Text:       "TEST",
        NoiseCount: 100,
        CurveCount: 2,
        LineWidth:  1,
    }

    captchaInstance, err := GenerateCaptcha(config)
    if err != nil {
        t.Fatalf("生成验证码失败: %v", err)
    }

    // 测试空路径的PNG导出（应保存在项目目录下）
    pngData, err := captchaInstance.ToPNG("")
    if err != nil {
        t.Errorf("使用空路径导出PNG失败: %v", err)
    }
    if len(pngData) == 0 {
        t.Error("PNG数据为空")
    }

    // 验证PNG文件是否已在项目目录中创建
    _, thisFile, _, _ := runtime.Caller(0)
    projectDir := filepath.Dir(thisFile)
    files, err := os.ReadDir(projectDir)
    if err != nil {
        t.Errorf("读取项目目录失败: %v", err)
    }

    foundPNG := false
    for _, file := range files {
        if strings.HasPrefix(file.Name(), "TEST_") && strings.HasSuffix(file.Name(), ".png") {
            foundPNG = true
            // 清理测试文件
            os.Remove(filepath.Join(projectDir, file.Name()))
            break
        }
    }
    if !foundPNG {
        t.Error("PNG文件未在项目目录中创建")
    }

    // 测试自定义路径的PNG导出
    pngPath := filepath.Join(tempDir, "test.png")
    pngData2, err := captchaInstance.ToPNG(pngPath)
    if err != nil {
        t.Errorf("使用自定义路径导出PNG失败: %v", err)
    }
    if len(pngData2) == 0 {
        t.Error("自定义路径的PNG数据为空")
    }
    if _, err := os.Stat(pngPath); os.IsNotExist(err) {
        t.Error("PNG文件未在自定义路径创建")
    }
}

func TestJPGCaptcha(t *testing.T) {
    // 创建临时目录用于测试文件
    tempDir, err := os.MkdirTemp("", "captcha_test_jpg")
    if err != nil {
        t.Fatalf("创建临时目录失败: %v", err)
    }
    defer os.RemoveAll(tempDir) // 测试完成后清理

    config := CaptchaConfig{
        Width:      240,
        Height:     80,
        FontPath:   "./font/default.ttf",
        FontSize:   36,
        Text:       "TEST",
        NoiseCount: 100,
        CurveCount: 2,
        LineWidth:  1,
    }

    captchaInstance, err := GenerateCaptcha(config)
    if err != nil {
        t.Fatalf("生成验证码失败: %v", err)
    }

    // 测试空路径的JPG导出（应保存在项目目录下）
    jpgData, err := captchaInstance.ToJPG("")
    if err != nil {
        t.Errorf("使用空路径导出JPG失败: %v", err)
    }
    if len(jpgData) == 0 {
        t.Error("JPG数据为空")
    }

    // 验证JPG文件是否已在项目目录中创建
    _, thisFile, _, _ := runtime.Caller(0)
    projectDir := filepath.Dir(thisFile)
    files, err := os.ReadDir(projectDir)
    if err != nil {
        t.Errorf("读取项目目录失败: %v", err)
    }

    foundJPG := false
    for _, file := range files {
        if strings.HasPrefix(file.Name(), "TEST_") && strings.HasSuffix(file.Name(), ".jpg") {
            foundJPG = true
            // 清理测试文件
            os.Remove(filepath.Join(projectDir, file.Name()))
            break
        }
    }
    if !foundJPG {
        t.Error("JPG文件未在项目目录中创建")
    }

    // 测试自定义路径的JPG导出
    jpgPath := filepath.Join(tempDir, "test.jpg")
    jpgData2, err := captchaInstance.ToJPG(jpgPath)
    if err != nil {
        t.Errorf("使用自定义路径导出JPG失败: %v", err)
    }
    if len(jpgData2) == 0 {
        t.Error("自定义路径的JPG数据为空")
    }
    if _, err := os.Stat(jpgPath); os.IsNotExist(err) {
        t.Error("JPG文件未在自定义路径创建")
    }
}

func TestBase64Captcha(t *testing.T) {
    config := CaptchaConfig{
        Width:      240,
        Height:     80,
        FontPath:   "./font/default.ttf",
        FontSize:   36,
        Text:       "TEST",
        NoiseCount: 100,
        CurveCount: 2,
        LineWidth:  1,
    }

    captchaInstance, err := GenerateCaptcha(config)
    if err != nil {
        t.Fatalf("生成验证码失败: %v", err)
    }

    // 测试Base64导出
    base64String, err := captchaInstance.ToString()
    if err != nil {
        t.Errorf("导出Base64字符串失败: %v", err)
    }
    if len(base64String) == 0 {
        t.Error("Base64字符串为空")
    }
}
