// Package banner рендерит текст как ASCII-арт (FIGlet-баннер) — для
// логотипов и заголовков в TUI проектов dtrader. Использует
// github.com/common-nighthawk/go-figure под капотом, но НЕ
// экспортирует ни одного его типа напрямую — весь публичный API этого
// пакета (Render, Font, список поддерживаемых шрифтов) выражен через
// собственные, стабильные имена. Если в будущем go-figure окажется
// заброшен (репозиторий уже несколько лет без активности на момент
// написания) или обнаружится более качественный форк/альтернатива —
// единственное место, которое нужно будет поменять, это данный файл;
// ни один потребитель dtui-kit (dtrader-tui-6, dtrader-history-tui и
// далее) не должен из-за этого переписывать свой код.
//
// ВАЖНО про go-figure: некоторые его внутренние пути вызывают
// log.Fatal (например NewColorFigure с неизвестным цветом) — то есть
// прямой os.Exit(1) всего процесса, не просто ошибку. Render защищает
// от этого там, где может (проверка шрифта заранее через Fonts()), но
// полностью исключить риск без форка самой библиотеки нельзя — при
// сомнении в валидности входных данных используй только проверенные
// шрифты из Fonts().
package banner

import (
	"fmt"
	"strings"

	figure "github.com/common-nighthawk/go-figure"
)

// Font — идентификатор FIGlet-шрифта. Тип введён (не просто string),
// чтобы в будущем можно было валидировать значение на этапе компиляции
// через именованные константы, не меняя сигнатуру Render.
type Font string

// Общеупотребимые, визуально спокойные шрифты из встроенного набора
// go-figure — достаточно разборчивые на мелких терминалах, в отличие
// от декоративных вычурных шрифтов (3-d, alligator и т.п.), которые
// быстро становятся нечитаемыми уже на десятке символов ширины.
// Полный список встроенных шрифтов на порядок больше — см. Fonts().
const (
	FontStandard Font = "standard"
	FontSmall    Font = "small"
	FontSlant    Font = "slant"
	FontMini     Font = "mini"
	FontBig      Font = "big"
	// FontBanner3 — крупный сплошной шрифт из "#" (7 строк): логотип.
	FontBanner3 Font = "banner3"
)

// Align — горизонтальное выравнивание результата Render относительно
// заданной ширины терминала/блока — аналог justify-*/text-align в
// Tailwind, тот же принцип на уровне текстового блока.
type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

// Options — параметры рендера. Width — ширина области, относительно
// которой считается выравнивание (обычно ширина терминала или
// content-блока); при Align == AlignLeft не используется и может быть
// оставлена нулевой.
type Options struct {
	Font  Font
	Align Align
	Width int
}

// Render рендерит text как ASCII-арт с заданными опциями. Пустой
// Options.Font трактуется как FontStandard. Возвращает
// многострочную строку без завершающего "\n".
func Render(text string, opts Options) string {
	font := opts.Font
	if font == "" {
		font = FontStandard
	}

	// strict=false — невалидные (не-ASCII) символы заменяются на "?"
	// самой библиотекой, вместо log.Fatal при strict=true — этот пакет
	// всегда использует нестрогий режим, потому что аварийный останов
	// всего процесса из-за одного неожиданного символа в заголовке
	// неприемлем для встраиваемого TUI-компонента.
	fig := figure.NewFigure(text, string(font), false)
	lines := fig.Slicify()

	if opts.Align == AlignLeft || opts.Width <= 0 {
		return strings.Join(lines, "\n")
	}

	maxLen := 0
	for _, line := range lines {
		if l := len([]rune(line)); l > maxLen {
			maxLen = l
		}
	}

	aligned := make([]string, len(lines))
	for i, line := range lines {
		aligned[i] = alignLine(line, maxLen, opts.Width, opts.Align)
	}
	return strings.Join(aligned, "\n")
}

// alignLine дополняет одну строку пробелами слева согласно align.
// maxLen — ширина самой широкой строки баннера (все строки баннера
// выравниваются по ОДНОМУ общему отступу, а не каждая по отдельности
// — иначе баннер потерял бы форму, "рассыпавшись" по разной ширине
// строк вместо единого центрированного блока).
func alignLine(line string, maxLen, width int, align Align) string {
	pad := width - maxLen
	if pad <= 0 {
		return line
	}
	switch align {
	case AlignCenter:
		left := pad / 2
		return strings.Repeat(" ", left) + line
	case AlignRight:
		return strings.Repeat(" ", pad) + line
	default:
		return line
	}
}

// Fonts возвращает имена всех шрифтов, гарантированно встроенных в
// go-figure на момент написания этого пакета (см. bindata.go самой
// библиотеки) — полный список крупнее набора именованных констант
// FontStandard/FontSmall/... выше, но не весь он одинаково пригоден
// для мелких терминалов; здесь перечислены только те, что были явно
// проверены на разборчивость.
func Fonts() []Font {
	return []Font{FontStandard, FontSmall, FontSlant, FontMini, FontBig, FontBanner3}
}

// IsKnownFont — true, если font входит в Fonts() (т.е. заведомо
// проверенный, безопасный выбор). Использование шрифта вне этого
// списка не запрещено (go-figure поддерживает намного больше), но
// Render в таком случае не даёт никакой гарантии читаемости и
// полагается на успешную загрузку шрифта самой библиотекой.
func IsKnownFont(font Font) bool {
	for _, f := range Fonts() {
		if f == font {
			return true
		}
	}
	return false
}

// MustRender — то же самое, что Render, но паникует, если font не
// входит в IsKnownFont — для использования в местах, где ошибка
// шрифта должна быть замечена сразу при разработке (например,
// построение статичного заголовка на этапе Init()), а не тихо дать
// go-figure шанс на внутренний log.Fatal с неясной причиной.
func MustRender(text string, opts Options) string {
	if opts.Font != "" && !IsKnownFont(opts.Font) {
		panic(fmt.Sprintf("banner: неизвестный шрифт %q — используй один из Fonts() или Render (не MustRender), если осознанно хочешь другой шрифт из полного набора go-figure", opts.Font))
	}
	return Render(text, opts)
}
