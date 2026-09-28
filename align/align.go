// Package align — горизонтальное выравнивание многострочных текстовых
// блоков относительно заданной ширины — тот же принцип, что
// justify-*/text-align в Tailwind, применительно к тексту в терминале.
//
// lipgloss (используемый bubbletea-экосистемой) уже умеет выравнивание
// через Style.Align()/AlignHorizontal() — этот пакет не конкурирует с
// lipgloss и не переизобретает его целиком. Он нужен для двух
// случаев, которые lipgloss не закрывает единообразно: (1)
// выравнивание уже готового многострочного текста (например,
// результата banner.Render), где отдельный lipgloss.Style ещё нужно
// завести и правильно настроить под конкретный блок; (2)
// переиспользуемая, одна и та же функция во всех TUI-проектах dtrader
// без необходимости каждый раз собирать
// lipgloss.NewStyle().Width(w).Align(...) заново.
package align

import "strings"

// Horizontal — горизонтальное выравнивание, аналог text-align.
type Horizontal int

const (
	Left Horizontal = iota
	Center
	Right
)

// Block выравнивает многострочный текст text по горизонтали
// относительно width — каждая строка дополняется пробелами независимо
// от длины остальных строк, в отличие от banner.Render с Align, где
// весь баннер выравнивается по ОДНОМУ общему отступу (там это
// осознанное решение — баннер как единая фигура; здесь, наоборот,
// каждая строка текста выравнивается сама по себе, как и ожидается от
// обычного многострочного текстового блока, где строки разной длины —
// например абзац или лог).
func Block(text string, width int, h Horizontal) string {
	if width <= 0 || h == Left {
		return text
	}
	lines := strings.Split(text, "\n")
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = Line(line, width, h)
	}
	return strings.Join(out, "\n")
}

// Line выравнивает одну строку line относительно width. Если line
// длиннее width — возвращается без изменений (обрезка — не задача
// этого пакета, это выбор вызывающего кода, если он вообще нужен).
func Line(line string, width int, h Horizontal) string {
	lineLen := len([]rune(line))
	pad := width - lineLen
	if pad <= 0 || h == Left {
		return line
	}
	switch h {
	case Center:
		left := pad / 2
		right := pad - left
		return strings.Repeat(" ", left) + line + strings.Repeat(" ", right)
	case Right:
		return strings.Repeat(" ", pad) + line
	default:
		return line
	}
}
