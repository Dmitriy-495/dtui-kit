# dtui-kit

Общая библиотека визуальных компонентов для консольных интерфейсов
(TUI) проектов семейства `dtrader`: фирменная палитра и переиспользуемые
элементы отображения (графики, баннеры и далее по мере необходимости),
чтобы не копировать одни и те же цвета и компоненты в каждый
TUI-репозиторий отдельно (`dtrader-tui-6`, `dtrader-history-tui` и
последующие).

Текущая версия: `v0.3.0`.

## Философия

Кит — это клей, а не фреймворк. Он строится поверх
[`bubbletea`](https://github.com/charmbracelet/bubbletea),
[`bubbles`](https://github.com/charmbracelet/bubbles) и
[`lipgloss`](https://github.com/charmbracelet/lipgloss) и опирается на
них открыто: то, что они уже умеют (раскладка, таблицы, рамки,
выравнивание), здесь не переписывается. В библиотеку попадает только то,
что уже реально используется и понравилось, без компонентов «на вырост».

Четыре правила и их обоснование — в [`PHILOSOPHY.md`](PHILOSOPHY.md).
Как компоненты связаны со слоем раскладки (`dtui`) и общими свойствами
стилей — в [`ARCHITECTURE.md`](ARCHITECTURE.md).

## Требования

- Go **1.25** или новее.
- Charm **v2**: `charm.land/lipgloss/v2` (и `charm.land/bubbletea/v2`,
  `charm.land/bubbles/v2` в приложении-потребителе).
  Экспортируемые цвета `theme` имеют тип `image/color.Color`
  (в `v0.1.0` это был `lipgloss.Color`).

## Пакеты

### `theme`

Фирменная цветовая палитра и базовые lipgloss-стили: оранжевый для
рамок и акцентов, зелёный/жёлтый/красный для статусов OK/Warning/SOS.
Переопределяется внешним `theme.yaml` (см. `theme.example.yaml` в корне
репозитория): читается при старте приложения, без пересборки бинарника.

```go
import "github.com/Dmitriy-495/dtui-kit/theme"

func main() {
    // опционально: переопределить палитру из файла рядом с бинарником
    if err := theme.LoadIfExists("theme.yaml"); err != nil {
        log.Fatal(err)
    }

    fmt.Println(theme.OKStyle.Render("работает"))
    style := lipgloss.NewStyle().Foreground(theme.ColorBorder)
}
```

### `sparkbar`

Вертикальный бар-чарт истории числовой метрики (0–100%) на блочных
символах Unicode (`▁▂▃▄▅▆▇█`): для CPU/RAM/Disk и любых других
процентных метрик, которые удобно показывать как временной ряд прямо в
терминале.

```go
import "github.com/Dmitriy-495/dtui-kit/sparkbar"

hist := sparkbar.NewHistory(40) // окно на последние 40 замеров
hist.Push(42.5)
fmt.Println(sparkbar.RenderStyled(hist, 40, 70, 90)) // warnAt=70%, sosAt=90%
```

Свой, а не `ntcharts`: для однострочной истории метрики он не даёт
явного выигрыша и тяжелее. Если понадобятся многострочные графики или
оси, брать `ntcharts`, а не писать своё.

### `banner`

ASCII-арт (FIGlet) заголовки и логотипы. Обёртка над
[`common-nighthawk/go-figure`](https://github.com/common-nighthawk/go-figure)
со своим стабильным API: наружу не торчит ни одного типа `go-figure`,
и при необходимости реализация под капотом заменяется без изменений в
проектах-потребителях. Есть собственное горизонтальное выравнивание
всего блока (`AlignLeft`/`AlignCenter`/`AlignRight`): все строки
сдвигаются на один общий отступ, поэтому форма фигуры не искажается.

```go
import "github.com/Dmitriy-495/dtui-kit/banner"

fmt.Println(banner.Render("dtrader", banner.Options{
    Font:  banner.FontBanner3,
    Align: banner.AlignCenter,
    Width: 80,
}))
```

Проверенные шрифты (`banner.Fonts()`): `standard`, `small`, `slant`,
`mini`, `big`, `banner3`. `FontBanner3` — крупный сплошной шрифт из `#`
высотой 7 строк, подходит для логотипа.

Что важно знать:

- `MustRender` паникует на шрифте вне `Fonts()`. `Render` не паникует,
  но `go-figure` на непроверенном имени может вызвать `log.Fatal`
  (прямой выход из программы), поэтому используй шрифты из `Fonts()`.
- Формат `.tlf` не поддерживается, а часть внешних `.flf` роняет парсер.
- Кегль в консоли — это число строк шрифта; размер самих символов
  задаёт терминал.

## Установка

```bash
go get github.com/Dmitriy-495/dtui-kit
```

## Тесты

```bash
go vet ./... && go test ./... -race
```

## Версии

| Тег | Что изменилось |
|---|---|
| `v0.1.0` | `theme`, `sparkbar`, `banner`, `align` |
| `v0.2.0` | переезд на Charm v2 (ломающее изменение: цвета `theme` стали `color.Color`); пакет `align` удалён как дубль `lipgloss` |
| `v0.3.0` | `banner.FontBanner3` |

Изменение публичного API `theme` и `sparkbar` (у них уже есть
потребители) делается обратно совместимо либо с явным планом миграции и
новым минорным тегом.
