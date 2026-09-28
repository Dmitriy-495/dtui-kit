package sparkbar

import (
	"strings"
	"testing"
)

func TestHistory_PushAndLatest(t *testing.T) {
	h := NewHistory(3)
	if _, ok := h.Latest(); ok {
		t.Error("Latest() пустой истории должен вернуть ok=false")
	}

	h.Push(10)
	h.Push(20)
	latest, ok := h.Latest()
	if !ok || latest != 20 {
		t.Errorf("Latest() = (%f, %v), want (20, true)", latest, ok)
	}
}

func TestHistory_EvictsOldestWhenFull(t *testing.T) {
	h := NewHistory(3)
	h.Push(1)
	h.Push(2)
	h.Push(3)
	h.Push(4) // вытесняет 1

	if len(h.values) != 3 {
		t.Fatalf("len(values) = %d, want 3 (maxLen)", len(h.values))
	}
	if h.values[0] != 2 {
		t.Errorf("values[0] = %f, want 2 (самое старое (1) должно быть вытеснено)", h.values[0])
	}
	latest, _ := h.Latest()
	if latest != 4 {
		t.Errorf("Latest() = %f, want 4", latest)
	}
}

func TestHistory_ClampsOutOfRangeValues(t *testing.T) {
	h := NewHistory(5)
	h.Push(-10)
	h.Push(150)

	if h.values[0] != 0 {
		t.Errorf("отрицательное значение должно обрезаться до 0, получено %f", h.values[0])
	}
	if h.values[1] != 100 {
		t.Errorf("значение больше 100 должно обрезаться до 100, получено %f", h.values[1])
	}
}

func TestHistory_ZeroOrNegativeMaxLenDefaultsToOne(t *testing.T) {
	h := NewHistory(0)
	if h.maxLen != 1 {
		t.Errorf("maxLen при вызове с 0 = %d, want 1 (защита от вырожденного случая)", h.maxLen)
	}
	h.Push(10)
	h.Push(20)
	if len(h.values) != 1 {
		t.Errorf("len(values) = %d, want 1", len(h.values))
	}
}

func TestRender_EmptyHistoryPadsWithMinimalLevel(t *testing.T) {
	h := NewHistory(10)
	bar := Render(h, 5)
	runes := []rune(bar)
	if len(runes) != 5 {
		t.Fatalf("len(bar) = %d, want 5", len(runes))
	}
	for _, r := range runes {
		if r != barLevels[0] {
			t.Errorf("пустая история должна рендериться минимальным уровнем %q, получено %q", string(barLevels[0]), string(r))
		}
	}
}

func TestRender_FullValueRendersMaxLevel(t *testing.T) {
	h := NewHistory(1)
	h.Push(100)
	bar := Render(h, 1)
	if []rune(bar)[0] != barLevels[len(barLevels)-1] {
		t.Errorf("значение 100 должно рендериться максимальным уровнем %q, получено %q", string(barLevels[len(barLevels)-1]), bar)
	}
}

func TestRender_ZeroValueRendersMinLevel(t *testing.T) {
	h := NewHistory(1)
	h.Push(0)
	bar := Render(h, 1)
	if []rune(bar)[0] != barLevels[0] {
		t.Errorf("значение 0 должно рендериться минимальным уровнем %q, получено %q", string(barLevels[0]), bar)
	}
}

func TestRender_ShowsOnlyLastWidthValuesWhenHistoryLonger(t *testing.T) {
	h := NewHistory(10)
	for i := 1; i <= 10; i++ {
		h.Push(float64(i) * 10) // 10, 20, ..., 100
	}
	bar := Render(h, 3)
	if len([]rune(bar)) != 3 {
		t.Fatalf("len(bar) = %d, want 3 (ограничено width)", len([]rune(bar)))
	}
	// Последние три значения — 80, 90, 100 — должны дать три РАЗНЫХ
	// (возрастающих) уровня, не все одинаковый максимальный.
	runes := []rune(bar)
	if runes[0] == runes[2] {
		t.Errorf("ожидалось различие между первым и последним символом при разных значениях, получено одинаковые: %q", bar)
	}
}

func TestRender_WidthLessThanOneDefaultsToOne(t *testing.T) {
	h := NewHistory(5)
	h.Push(50)
	bar := Render(h, 0)
	if len([]rune(bar)) != 1 {
		t.Errorf("Sparkbar с width=0 должен вернуть строку длины 1, получено длину %d", len([]rune(bar)))
	}
}

func TestRenderStyled_EmptyHistoryUsesMutedStyle(t *testing.T) {
	h := NewHistory(5)
	out := RenderStyled(h, 3, 70, 90)
	// Не проверяем точный ANSI-код (хрупко к деталям lipgloss), только
	// то, что результат не паникует и содержит сам бар.
	if !strings.Contains(out, string(barLevels[0])) {
		t.Errorf("пустая история должна содержать минимальный уровень в выводе, получено: %q", out)
	}
}

func TestRenderStyled_DoesNotPanicAtThresholdBoundaries(t *testing.T) {
	for _, v := range []float64{0, 69.9, 70, 89.9, 90, 100} {
		h := NewHistory(1)
		h.Push(v)
		_ = RenderStyled(h, 5, 70, 90) // не должно паниковать ни на одной границе
	}
}
