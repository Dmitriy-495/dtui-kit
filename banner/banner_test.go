package banner

import (
	"strings"
	"testing"
)

func TestRender_ProducesNonEmptyMultilineOutput(t *testing.T) {
	out := Render("dt", Options{})
	if out == "" {
		t.Fatal("Render вернул пустую строку")
	}
	lines := strings.Split(out, "\n")
	if len(lines) < 2 {
		t.Errorf("ожидался многострочный ASCII-арт, получено %d строк(и)", len(lines))
	}
}

func TestRender_EmptyFontDefaultsToStandard(t *testing.T) {
	withEmpty := Render("A", Options{})
	withExplicit := Render("A", Options{Font: FontStandard})
	if withEmpty != withExplicit {
		t.Error("Render с пустым Font должен давать тот же результат, что и явный FontStandard")
	}
}

func TestRender_DifferentFontsProduceDifferentOutput(t *testing.T) {
	standard := Render("A", Options{Font: FontStandard})
	small := Render("A", Options{Font: FontSmall})
	if standard == small {
		t.Error("разные шрифты должны давать разный рендер одного и того же текста")
	}
}

func TestRender_LeftAlignIgnoresWidth(t *testing.T) {
	// AlignLeft должен давать тот же результат независимо от Width —
	// проверка через сравнение с "нулевым" вызовом надёжнее, чем
	// предположение о том, с чего начинается конкретная буква в
	// конкретном шрифте (некоторые буквы, например "A" в FontStandard,
	// сами по себе имеют треугольную форму с ведущими пробелами на
	// части строк — это свойство самого шрифта, не бага выравнивания).
	unaligned := Render("A", Options{})
	leftWithWidth := Render("A", Options{Align: AlignLeft, Width: 200})
	if unaligned != leftWithWidth {
		t.Error("AlignLeft с заданным Width должен давать тот же результат, что и рендер вообще без Width")
	}
}

func TestRender_CenterAlignAddsLeadingSpacesWhenWidthIsLarger(t *testing.T) {
	unaligned := Render("A", Options{Align: AlignLeft})
	centered := Render("A", Options{Align: AlignCenter, Width: 200})

	unalignedLines := strings.Split(unaligned, "\n")
	centeredLines := strings.Split(centered, "\n")

	if len(unalignedLines) != len(centeredLines) {
		t.Fatalf("центрирование не должно менять число строк: %d vs %d", len(unalignedLines), len(centeredLines))
	}

	// Хотя бы непустая строка должна получить отступ слева при
	// центрировании в широком поле.
	foundIndented := false
	for i, line := range centeredLines {
		if strings.TrimSpace(unalignedLines[i]) == "" {
			continue
		}
		if strings.HasPrefix(line, " ") {
			foundIndented = true
			break
		}
	}
	if !foundIndented {
		t.Error("при AlignCenter с большим Width ожидался отступ слева хотя бы на одной строке")
	}
}

func TestRender_RightAlignPushesTextToRightEdge(t *testing.T) {
	unaligned := Render("A", Options{Align: AlignLeft})
	right := Render("A", Options{Align: AlignRight, Width: 200})

	unalignedLines := strings.Split(unaligned, "\n")
	rightLines := strings.Split(right, "\n")

	centeredPad := 0
	rightPad := 0
	for i, line := range rightLines {
		if strings.TrimSpace(unalignedLines[i]) == "" {
			continue
		}
		rightPad = len(line) - len(strings.TrimLeft(line, " "))
		break
	}

	centered := Render("A", Options{Align: AlignCenter, Width: 200})
	centeredLines := strings.Split(centered, "\n")
	for i, line := range centeredLines {
		if strings.TrimSpace(unalignedLines[i]) == "" {
			continue
		}
		centeredPad = len(line) - len(strings.TrimLeft(line, " "))
		break
	}

	if rightPad <= centeredPad {
		t.Errorf("отступ при AlignRight (%d) должен быть больше, чем при AlignCenter (%d), в широком поле", rightPad, centeredPad)
	}
}

func TestRender_ZeroOrNegativeWidthFallsBackToUnaligned(t *testing.T) {
	unaligned := Render("A", Options{Align: AlignLeft})
	zeroWidth := Render("A", Options{Align: AlignCenter, Width: 0})
	negativeWidth := Render("A", Options{Align: AlignCenter, Width: -5})

	if unaligned != zeroWidth {
		t.Error("Width=0 с AlignCenter должен давать тот же результат, что и без выравнивания")
	}
	if unaligned != negativeWidth {
		t.Error("отрицательный Width должен давать тот же результат, что и без выравнивания")
	}
}

func TestRender_WidthSmallerThanContentDoesNotPanic(t *testing.T) {
	// Ширина меньше самого баннера — не должно быть паники или
	// отрицательного количества пробелов.
	_ = Render("Hello World", Options{Align: AlignCenter, Width: 1})
}

func TestFonts_AllReturnedFontsAreKnown(t *testing.T) {
	for _, f := range Fonts() {
		if !IsKnownFont(f) {
			t.Errorf("шрифт %q из Fonts() должен считаться известным (IsKnownFont)", f)
		}
	}
}

func TestIsKnownFont_RejectsArbitraryString(t *testing.T) {
	if IsKnownFont(Font("definitely-not-a-real-font-xyz")) {
		t.Error("произвольная строка не должна считаться известным шрифтом")
	}
}

func TestMustRender_PanicsOnUnknownFont(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("ожидалась паника при неизвестном шрифте")
		}
	}()
	MustRender("test", Options{Font: Font("totally-unknown-font")})
}

func TestMustRender_DoesNotPanicOnKnownFont(t *testing.T) {
	// Успех теста — отсутствие паники.
	_ = MustRender("test", Options{Font: FontStandard})
}

func TestMustRender_DoesNotPanicOnEmptyFont(t *testing.T) {
	// Пустой Font трактуется как FontStandard (известный), не как
	// "неизвестный шрифт" — не должен вызывать панику.
	_ = MustRender("test", Options{})
}

func TestFonts_AllKnownFontsRenderDigitsAndColon(t *testing.T) {
	for _, f := range Fonts() {
		if out := Render("07:06:34", Options{Font: f}); out == "" {
			t.Errorf("шрифт %q вернул пустой результат для цифр и двоеточия", f)
		}
	}
}

func TestFontBanner3_IsKnownAndRenders(t *testing.T) {
	if !IsKnownFont(FontBanner3) {
		t.Fatal("FontBanner3 должен входить в Fonts()")
	}
	if out := Render("dtrader", Options{Font: FontBanner3}); out == "" {
		t.Error("FontBanner3 вернул пустой результат")
	}
}
