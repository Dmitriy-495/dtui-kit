// Package theme — фирменная цветовая палитра и базовые стили для TUI
// проектов семейства dtrader (dtrader-tui-6, dtrader-history-tui и
// далее). Значения по умолчанию портированы 1:1 из dtrader-tui-6/
// internal/tui/symbol.go (см. раздел 11 CHECKPOINT.md проекта
// dtrader-6, откуда эти значения взяты изначально), не выдуманы
// заново.
//
// Вынесена в отдельный переиспользуемый пакет по решению автора
// проекта dtrader: единый визуальный язык на все консольные
// интерфейсы семейства, вместо копирования одних и тех же констант в
// каждый TUI-репозиторий по отдельности.
//
// Палитра переопределяема через внешний theme.yaml (см. Load) — по
// решению автора: конфиг должен быть внешним файлом, читаемым при
// старте приложения, не программным API, чтобы менять оформление без
// пересборки бинарника.
package theme

import (
	"fmt"
	"image/color"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/goccy/go-yaml"
)

// Palette — вся конфигурируемая часть темы: цвета (как строки ANSI
// 256-цветного кода — то же, что принимает lipgloss.Color напрямую)
// и шрифт баннера по умолчанию для этого TUI (см. пакет banner). Поля
// экспортированы для декодирования YAML через теги — само значение
// используется через Current()/примененные стили ниже, не напрямую.
type Palette struct {
	Border string `yaml:"border"`
	OK     string `yaml:"ok"`
	Warn   string `yaml:"warn"`
	SOS    string `yaml:"sos"`
	Data   string `yaml:"data"`
	Muted  string `yaml:"muted"`

	// BannerFont — имя шрифта из banner.Fonts() по умолчанию для этого
	// TUI (см. пакет banner) — например "standard", "small". Пустая
	// строка означает "использовать значение по умолчанию самого
	// пакета banner", не приводит к ошибке при Load.
	BannerFont string `yaml:"banner_font"`
}

// defaultPalette — встроенные значения, идентичные тем, что были
// захардкожены до введения конфигурируемости (см. историю этого
// файла) — Load(), если файл не найден, оставляет ИМЕННО ЭТИ значения
// активными, не пустую/нулевую палитру.
var defaultPalette = Palette{
	Border:     "214",
	OK:         "82",
	Warn:       "226",
	SOS:        "196",
	Data:       "214",
	Muted:      "239",
	BannerFont: "standard",
}

// current — активная палитра. Пакетный уровень (не завязана на
// конкретный экземпляр TUI-приложения) — та же договорённость, что и
// раньше для отдельных переменных ColorBorder и т.д.: один процесс TUI
// имеет ровно одну активную тему на всё время работы.
var current = defaultPalette

// Экспортированные стили и цвета, вычисленные из current — сохраняют
// прежние имена (ColorBorder, OKStyle и т.д.), чтобы код, уже
// написанный против theme.ColorBorder/theme.OKStyle (см. пакет
// sparkbar и потребителей dtui-kit), не требовал изменений после
// введения конфигурируемости. Пересчитываются в applyCurrent(),
// вызываемой из init() и из каждого успешного Load().
var (
	ColorBorder  color.Color
	ColorOK      color.Color
	ColorWarn    color.Color
	ColorSOS     color.Color
	ColorData    color.Color
	ColorMuted   color.Color
	ColorNeutral color.Color

	MutedStyle  lipgloss.Style
	DataStyle   lipgloss.Style
	OKStyle     lipgloss.Style
	WarnStyle   lipgloss.Style
	SOSStyle    lipgloss.Style
	BorderStyle lipgloss.Style
)

func init() {
	applyCurrent()
}

// applyCurrent пересчитывает все экспортированные Color*/*Style
// переменные из current — единственное место, где они присваиваются,
// чтобы не рассинхронизировать структуру Palette и производные стили
// при добавлении новых полей в будущем.
func applyCurrent() {
	ColorBorder = lipgloss.Color(current.Border)
	ColorOK = lipgloss.Color(current.OK)
	ColorWarn = lipgloss.Color(current.Warn)
	ColorSOS = lipgloss.Color(current.SOS)
	ColorData = lipgloss.Color(current.Data)
	ColorMuted = lipgloss.Color(current.Muted)
	ColorNeutral = ColorMuted // нейтральное состояние — своего цвета в исходной палитре нет

	MutedStyle = lipgloss.NewStyle().Foreground(ColorMuted)
	DataStyle = lipgloss.NewStyle().Foreground(ColorData)
	OKStyle = lipgloss.NewStyle().Foreground(ColorOK)
	WarnStyle = lipgloss.NewStyle().Foreground(ColorWarn)
	SOSStyle = lipgloss.NewStyle().Foreground(ColorSOS).Bold(true)

	BorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder)
}

// Current возвращает копию активной палитры — например, чтобы
// прочитать BannerFont и передать его в banner.Render.
func Current() Palette {
	return current
}

// Load читает theme.yaml по указанному path и делает его активной
// палитрой (пересчитывая все Color*/*Style переменные). Поля,
// отсутствующие в файле (пустая строка после декодирования),
// сохраняют значение из defaultPalette — YAML может переопределять
// только часть полей, не требуется указывать все сразу.
//
// Должен вызываться как можно раньше в main(), до того как остальной
// код успеет прочитать theme.ColorBorder/theme.OKStyle и т.п. — эти
// переменные меняются по месту (пакетный уровень), а не через
// внедрение зависимости, так что код, уже сохранивший значение стиля
// в свою переменную ДО вызова Load, не увидит новую тему (тот же
// принцип, что для любых пакетных var в Go — не специфика этого пакета).
func Load(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("theme: не удалось прочитать %q: %w", path, err)
	}

	loaded := defaultPalette // копия — незаполненные YAML-поля унаследуют дефолты
	if err := yaml.Unmarshal(raw, &loaded); err != nil {
		return fmt.Errorf("theme: не удалось разобрать %q: %w", path, err)
	}

	current = loaded
	applyCurrent()
	return nil
}

// LoadIfExists — то же самое, что Load, но отсутствие файла по path
// не считается ошибкой (остаётся активной палитра по умолчанию) —
// для случая, когда theme.yaml опционален и большинство запусков его
// не имеют. Ошибки чтения по ДРУГИМ причинам (например, есть файл, но
// недостаточно прав) по-прежнему возвращаются — молчаливое
// игнорирование оправдано только для "файла нет вообще", не для
// прочих сбоев.
func LoadIfExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return Load(path)
}
