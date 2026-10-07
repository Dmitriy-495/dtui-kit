package digits

import (
	"strings"
	"testing"
	"time"
)

func TestGlyphsConsistent(t *testing.T) {
	for r, g := range glyphs {
		w := len(g[0])
		for i, row := range g {
			if len(row) != w {
				t.Errorf("глиф %q: строка %d шириной %d, ожидалось %d", r, i, len(row), w)
			}
		}
		if r >= '0' && r <= '9' && w != DigitWidth {
			t.Errorf("цифра %q шириной %d, ожидалось %d", r, w, DigitWidth)
		}
	}
}

func TestRenderShape(t *testing.T) {
	out := Render("12:34")
	lines := strings.Split(out, "\n")
	if len(lines) != Height {
		t.Fatalf("строк %d, ожидалось %d", len(lines), Height)
	}
	for i, l := range lines {
		if len(l) != Width("12:34") {
			t.Errorf("строка %d шириной %d, Width = %d", i, len(l), Width("12:34"))
		}
	}
}

func TestRenderGolden(t *testing.T) {
	got := Render("1")
	if got != "     \n    |\n    |" {
		t.Errorf("Render(\"1\") = %q", got)
	}
	if Render("") != "\n\n" {
		t.Errorf("Render(\"\") = %q", Render(""))
	}
}

func TestRenderUnknownIsBlank(t *testing.T) {
	if Render("я") != Render(" ") {
		t.Error("неизвестный символ должен рисоваться как пробел")
	}
	if Width("1я") != DigitWidth+1+1 {
		t.Errorf("Width = %d", Width("1я"))
	}
}

func TestClockFixedWidthAndBlink(t *testing.T) {
	base := time.Date(2026, 10, 7, 0, 0, 0, 0, time.UTC)
	want := Width("00:00")
	for m := 0; m < 24*60; m++ {
		for _, ns := range []int{0, 499_999_999, 500_000_000, 999_999_999} {
			tm := base.Add(time.Duration(m)*time.Minute + time.Duration(ns))
			lines := strings.Split(Clock(tm), "\n")
			if len(lines) != Height {
				t.Fatalf("%s: строк %d", tm.Format("15:04"), len(lines))
			}
			for _, l := range lines {
				if len(l) != want {
					t.Fatalf("%s ns=%d: ширина %d, ожидалось %d", tm.Format("15:04"), ns, len(l), want)
				}
			}
		}
	}
}

func TestColonBlinks(t *testing.T) {
	on := time.Date(2026, 10, 7, 12, 34, 5, 100_000_000, time.UTC)
	off := time.Date(2026, 10, 7, 12, 34, 5, 600_000_000, time.UTC)
	if !ColonOn(on) || ColonOn(off) {
		t.Fatal("ColonOn: ожидалось on в первые полсекунды и off во вторые")
	}
	if Clock(on) == Clock(off) {
		t.Error("двоеточие не мигает")
	}
	if !strings.Contains(Clock(on), ".") || strings.Contains(Clock(off), ".") {
		t.Error("точки двоеточия должны быть только в фазе on")
	}
}
