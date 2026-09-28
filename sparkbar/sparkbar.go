// Package sparkbar реализует простой вертикальный бар-чарт для
// отображения истории числовой метрики (0-100%) в терминале — общего
// назначения компонент для TUI проектов семейства dtrader (CPU/RAM/
// Disk, и в будущем любые другие метрики 0-100%, которые имеет смысл
// показывать как временной ряд). Каждый столбец — один замер, высота
// столбца кодируется через блочные символы Unicode "▁▂▃▄▅▆▇█" (восемь
// уровней заполнения одной ячейки текста), несколько столбцов рядом
// дают вид знакомого bar-чарта без необходимости в полноценной графике.
package sparkbar

import (
	"strings"

	"github.com/Dmitriy-495/dtui-kit/theme"
)

// barLevels — блочные символы от пустого до полного, в порядке
// возрастания. len(barLevels)-1 — максимальный индекс уровня.
var barLevels = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// History хранит скользящее окно последних значений одной метрики
// (0-100, проценты). Значения новее в конце среза, старее — в начале;
// при заполнении окна старые значения вытесняются, как в кольцевом
// буфере (реализовано через простой срез с обрезкой, а не настоящий
// кольцевой буфер — при maxLen в пределах нескольких десятков значений
// разница в производительности незаметна, а код проще).
type History struct {
	values []float64
	maxLen int
}

// NewHistory создаёт историю ёмкостью maxLen замеров.
func NewHistory(maxLen int) History {
	if maxLen < 1 {
		maxLen = 1
	}
	return History{values: make([]float64, 0, maxLen), maxLen: maxLen}
}

// Push добавляет новое значение, вытесняя самое старое, если окно уже
// заполнено. Значение обрезается в диапазон [0, 100] — метрики,
// отображаемые этим компонентом, все являются процентами,
// отрицательные или превышающие 100% значения были бы признаком
// ошибки в источнике данных, а не легитимным состоянием, которое
// стоит визуализировать как есть.
func (h *History) Push(value float64) {
	if value < 0 {
		value = 0
	}
	if value > 100 {
		value = 100
	}
	h.values = append(h.values, value)
	if len(h.values) > h.maxLen {
		h.values = h.values[len(h.values)-h.maxLen:]
	}
}

// Latest возвращает самое свежее значение и true, либо (0, false),
// если история ещё пуста (ни одного замера не было).
func (h History) Latest() (float64, bool) {
	if len(h.values) == 0 {
		return 0, false
	}
	return h.values[len(h.values)-1], true
}

// Render рендерит текущую историю h как строку из блочных символов, по
// одному символу на замер, самые старые слева, самые свежие справа —
// то же самое направление чтения, что и у любого обычного графика
// временного ряда. Пустая история рендерится как width символов
// минимального уровня, чтобы не ломать выравнивание в layout, пока
// данные ещё не начали поступать.
func Render(h History, width int) string {
	if width < 1 {
		width = 1
	}

	var b strings.Builder
	// Если история короче width, дополняем слева минимальными
	// столбиками — иначе график "прыгал" бы по ширине от нескольких
	// символов до width по мере накопления данных, вместо равномерного
	// заполнения слева направо.
	padding := width - len(h.values)
	for i := 0; i < padding; i++ {
		b.WriteRune(barLevels[0])
	}

	start := 0
	if len(h.values) > width {
		start = len(h.values) - width
	}
	for _, v := range h.values[start:] {
		idx := int(v / 100 * float64(len(barLevels)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(barLevels) {
			idx = len(barLevels) - 1
		}
		b.WriteRune(barLevels[idx])
	}

	return b.String()
}

// RenderStyled — то же самое, что Render, но красит результат по
// последнему значению истории: зелёный (OK), жёлтый (Warn, ближе к
// порогу) или красный (SOS, критично) — три те же пороговые зоны, что
// используются во всём проекте dtrader для индикации состояния (см.
// пакет theme). warnAt/sosAt — пороги в процентах, ПОСЛЕ которых
// включается соответствующий цвет (warnAt <= sosAt).
func RenderStyled(h History, width int, warnAt, sosAt float64) string {
	bar := Render(h, width)
	latest, ok := h.Latest()
	if !ok {
		return theme.MutedStyle.Render(bar)
	}
	switch {
	case latest >= sosAt:
		return theme.SOSStyle.Render(bar)
	case latest >= warnAt:
		return theme.WarnStyle.Render(bar)
	default:
		return theme.OKStyle.Render(bar)
	}
}
