package theme

import (
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"charm.land/lipgloss/v2"
)

// TestColors_AreNotEmpty — простая защита от опечатки вида
// lipgloss.Color("") при будущем редактировании палитры: пустой цвет
// не паникует сам по себе, но рендерится некорректно и молча, что
// сложно заметить визуально сразу.
func TestColors_AreNotEmpty(t *testing.T) {
	colors := map[string]color.Color{
		"ColorBorder":  ColorBorder,
		"ColorOK":      ColorOK,
		"ColorWarn":    ColorWarn,
		"ColorSOS":     ColorSOS,
		"ColorData":    ColorData,
		"ColorMuted":   ColorMuted,
		"ColorNeutral": ColorNeutral,
	}
	for name, c := range colors {
		if c == nil {
			t.Errorf("%s пуст — вероятная опечатка в определении палитры", name)
		}
	}
}

// TestStyles_RenderWithoutPanicking — стили должны безопасно
// применяться к произвольному тексту, включая пустую строку.
func TestStyles_RenderWithoutPanicking(t *testing.T) {
	_ = MutedStyle.Render("тест")
	_ = DataStyle.Render("тест")
	_ = OKStyle.Render("тест")
	_ = WarnStyle.Render("тест")
	_ = SOSStyle.Render("тест")
	_ = BorderStyle.Render("тест")
	_ = MutedStyle.Render("")
}

// resetToDefault возвращает активную палитру к значениям по умолчанию
// после теста — Load меняет пакетный уровень (current и производные
// Color*/*Style переменные), так что без явного сброса один тест мог
// бы повлиять на порядок выполнения других (Go не гарантирует порядок
// тестов внутри пакета, но и не гарантирует его отсутствие — полагаться
// на оба варианта одинаково хрупко, поэтому сброс обязателен).
func resetToDefault(t *testing.T) {
	t.Helper()
	current = defaultPalette
	applyCurrent()
	t.Cleanup(func() {
		current = defaultPalette
		applyCurrent()
	})
}

func writeTempTheme(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "theme.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("не удалось записать временный theme.yaml: %v", err)
	}
	return path
}

func TestLoad_OverridesSpecifiedFieldsOnly(t *testing.T) {
	resetToDefault(t)
	path := writeTempTheme(t, "border: \"99\"\n")

	if err := Load(path); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if ColorBorder != lipgloss.Color("99") {
		t.Errorf("ColorBorder = %q, want %q (переопределено из файла)", ColorBorder, "99")
	}
	if ColorOK != lipgloss.Color(defaultPalette.OK) {
		t.Errorf("ColorOK = %q, want %q (не переопределялось, должен остаться дефолт)", ColorOK, defaultPalette.OK)
	}
}

func TestLoad_UpdatesDerivedStyles(t *testing.T) {
	resetToDefault(t)
	path := writeTempTheme(t, "ok: \"46\"\n")

	if err := Load(path); err != nil {
		t.Fatalf("Load: %v", err)
	}

	// OKStyle должен рендерить С НОВЫМ цветом, не со старым дефолтным —
	// проверяем через сравнение с "эталонным" стилем, построенным
	// напрямую с ожидаемым цветом, а не полагаемся на конкретный ANSI-
	// код в строке (хрупко к деталям lipgloss).
	want := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render("x")
	got := OKStyle.Render("x")
	if got != want {
		t.Errorf("OKStyle после Load = %q, want %q", got, want)
	}
}

func TestLoad_ReadsBannerFont(t *testing.T) {
	resetToDefault(t)
	path := writeTempTheme(t, "banner_font: \"slant\"\n")

	if err := Load(path); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if Current().BannerFont != "slant" {
		t.Errorf("Current().BannerFont = %q, want %q", Current().BannerFont, "slant")
	}
}

func TestLoad_ReturnsErrorForMissingFile(t *testing.T) {
	resetToDefault(t)
	if err := Load("/definitely/does/not/exist/theme.yaml"); err == nil {
		t.Error("ожидалась ошибка для отсутствующего файла")
	}
}

func TestLoad_ReturnsErrorForMalformedYAML(t *testing.T) {
	resetToDefault(t)
	path := writeTempTheme(t, "not: [valid: yaml: :::")
	if err := Load(path); err == nil {
		t.Error("ожидалась ошибка для некорректного YAML")
	}
}

func TestLoadIfExists_SilentlyKeepsDefaultsWhenFileMissing(t *testing.T) {
	resetToDefault(t)
	before := Current()

	if err := LoadIfExists("/definitely/does/not/exist/theme.yaml"); err != nil {
		t.Errorf("LoadIfExists для отсутствующего файла не должен возвращать ошибку, получено: %v", err)
	}

	if Current() != before {
		t.Error("палитра не должна была измениться при отсутствующем файле")
	}
}

func TestLoadIfExists_LoadsWhenFilePresent(t *testing.T) {
	resetToDefault(t)
	path := writeTempTheme(t, "border: \"55\"\n")

	if err := LoadIfExists(path); err != nil {
		t.Fatalf("LoadIfExists: %v", err)
	}
	if ColorBorder != lipgloss.Color("55") {
		t.Errorf("ColorBorder = %q, want %q", ColorBorder, "55")
	}
}

func TestLoadIfExists_ReturnsErrorForMalformedExistingFile(t *testing.T) {
	resetToDefault(t)
	path := writeTempTheme(t, "not: [valid: yaml: :::")
	if err := LoadIfExists(path); err == nil {
		t.Error("ожидалась ошибка для существующего, но некорректного файла (не должна маскироваться под 'файла нет')")
	}
}
