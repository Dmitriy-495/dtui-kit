package align

import "testing"

func TestLine_LeftReturnsUnchanged(t *testing.T) {
	if got := Line("hi", 10, Left); got != "hi" {
		t.Errorf("Line с Left = %q, want %q (без изменений)", got, "hi")
	}
}

func TestLine_CenterAddsEqualOrNearEqualPadding(t *testing.T) {
	got := Line("hi", 10, Center)
	want := "    hi    " // pad=8, left=4, right=4
	if got != want {
		t.Errorf("Line с Center = %q, want %q", got, want)
	}
	if len([]rune(got)) != 10 {
		t.Errorf("длина результата = %d, want 10", len([]rune(got)))
	}
}

func TestLine_CenterOddPaddingSplitsLeftShorter(t *testing.T) {
	// pad=7 (width=9, len=2): left=3, right=4 — левая сторона короче
	// при нечётном отступе, это осознанный выбор реализации (left =
	// pad/2 округляется вниз), не требование извне, но должен быть
	// стабильным и задокументированным поведением.
	got := Line("hi", 9, Center)
	if len([]rune(got)) != 9 {
		t.Fatalf("длина результата = %d, want 9", len([]rune(got)))
	}
	leadingSpaces := 0
	for _, r := range got {
		if r != ' ' {
			break
		}
		leadingSpaces++
	}
	if leadingSpaces != 3 {
		t.Errorf("leadingSpaces = %d, want 3", leadingSpaces)
	}
}

func TestLine_RightPushesToRightEdge(t *testing.T) {
	got := Line("hi", 10, Right)
	want := "        hi"
	if got != want {
		t.Errorf("Line с Right = %q, want %q", got, want)
	}
}

func TestLine_WidthSmallerThanContentReturnsUnchanged(t *testing.T) {
	if got := Line("hello world", 3, Center); got != "hello world" {
		t.Errorf("Line с width меньше содержимого = %q, want без изменений", got)
	}
}

func TestLine_ExactWidthReturnsUnchanged(t *testing.T) {
	if got := Line("hi", 2, Center); got != "hi" {
		t.Errorf("Line с width==len(line) = %q, want без изменений", got)
	}
}

func TestBlock_AppliesPerLineIndependently(t *testing.T) {
	text := "a\nbb\nccc"
	got := Block(text, 5, Right)
	want := "    a\n   bb\n  ccc"
	if got != want {
		t.Errorf("Block с Right = %q, want %q", got, want)
	}
}

func TestBlock_ZeroOrNegativeWidthReturnsUnchanged(t *testing.T) {
	text := "hello\nworld"
	if got := Block(text, 0, Center); got != text {
		t.Errorf("Block с width=0 = %q, want без изменений", got)
	}
	if got := Block(text, -5, Center); got != text {
		t.Errorf("Block с отрицательным width = %q, want без изменений", got)
	}
}

func TestBlock_LeftReturnsUnchanged(t *testing.T) {
	text := "hello\nworld"
	if got := Block(text, 20, Left); got != text {
		t.Errorf("Block с Left = %q, want без изменений", got)
	}
}

func TestBlock_SingleLineWorksLikeLine(t *testing.T) {
	got := Block("hi", 10, Center)
	want := Line("hi", 10, Center)
	if got != want {
		t.Errorf("Block с одной строкой = %q, want %q (как прямой вызов Line)", got, want)
	}
}
