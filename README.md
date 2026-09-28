# dtui-kit

Общая библиотека визуальных компонентов для консольных интерфейсов
(TUI) проектов семейства `dtrader` — фирменная палитра и переиспользуемые
элементы отображения (графики, баннеры, выравнивание, и далее по мере
необходимости), чтобы не копировать одни и те же цвета/компоненты в
каждый TUI-репозиторий отдельно (`dtrader-tui-6`, `dtrader-history-tui`
и последующие).

## Философия

Минимализм: в библиотеку выносится только то, что реально
дублировалось между как минимум двумя проектами, а не то, что
теоретически может пригодиться. Никакого фреймворка поверх
[`bubbletea`](https://github.com/charmbracelet/bubbletea)/[`lipgloss`](https://github.com/charmbracelet/lipgloss)
— `dtui-kit` не заменяет их, а дополняет несколькими готовыми,
уже выверенными деталями.

## Пакеты

### `theme`

Фирменная цветовая палитра и базовые lipgloss-стили — оранжевый для
рамок/акцентов, зелёный/жёлтый/красный для статусов OK/Warning/SOS.
Переопределяема через внешний `theme.yaml` (см. `theme.example.yaml`
в корне репозитория) — читается при старте приложения, без
пересборки бинарника.

```go
import "github.com/Dmitriy-495/dtui-kit/theme"

func main() {
    // опционально — переопределить палитру из файла рядом с бинарником
    if err := theme.LoadIfExists("theme.yaml"); err != nil {
        log.Fatal(err)
    }

    fmt.Println(theme.OKStyle.Render("работает"))
    style := lipgloss.NewStyle().Foreground(theme.ColorBorder)
}
```

### `sparkbar`

Вертикальный бар-чарт истории числовой метрики (0-100%) на блочных
символах Unicode (`▁▂▃▄▅▆▇█`) — для CPU/RAM/Disk и любых других
процентных метрик, которые имеет смысл показывать как временной ряд
прямо в терминале.

```go
import "github.com/Dmitriy-495/dtui-kit/sparkbar"

hist := sparkbar.NewHistory(40) // окно на последние 40 замеров
hist.Push(42.5)
fmt.Println(sparkbar.RenderStyled(hist, 40, 70, 90)) // warnAt=70%, sosAt=90%
```

### `banner`

ASCII-арт (FIGlet) заголовки и логотипы — обёртка над
[`common-nighthawk/go-figure`](https://github.com/common-nighthawk/go-figure)
со своим стабильным API (см. комментарий в коде — при необходимости
реализация под капотом заменяется без изменений в проектах-потребителях)
и собственным горизонтальным выравниванием (аналог `text-align` в
Tailwind, которого нет во внешней библиотеке).

```go
import "github.com/Dmitriy-495/dtui-kit/banner"

fmt.Println(banner.Render("dtrader", banner.Options{
    Font:  banner.FontSlant,
    Align: banner.AlignCenter,
    Width: 80,
}))
```

### `align`

Горизонтальное выравнивание произвольного многострочного текста
относительно заданной ширины — аналог `justify-*`/`text-align` в
Tailwind. В отличие от выравнивания в `banner` (весь баннер как одна
фигура с общим отступом), здесь каждая строка выравнивается
независимо — подходит для обычного текста, логов, абзацев.

```go
import "github.com/Dmitriy-495/dtui-kit/align"

fmt.Println(align.Block("привет\nмир", 20, align.Center))
```

## Установка

```bash
go get github.com/Dmitriy-495/dtui-kit
```

## Тесты

```bash
go test ./... -race
```

