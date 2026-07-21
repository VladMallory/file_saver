package installer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ─── Стили (Gruvbox) ───────────────────────────────────────────────────────
// Палитра Gruvbox: https://github.com/morhetz/gruvbox
// Каждый стиль отвечает за свой визуальный слой.
// SRP: стили не смешивают цвета, рамки и отступы — каждый за одно.

var (
	// Gruvbox dark palette: https://github.com/morhetz/gruvbox
	// Оставляем только используемые цвета — unused ловит неиспользуемые переменные.
	gbg0 = lipgloss.Color("#282828") // dark0_hard — фон кнопок.
	gbg4 = lipgloss.Color("#7c6f64") // dark4 — цвет разделителя.
	gfg4 = lipgloss.Color("#a89984") // gray — приглушённый текст.

	gBlue         = lipgloss.Color("#458588")
	gBrightAqua   = lipgloss.Color("#8ec07c")
	gBrightBlue   = lipgloss.Color("#83a598")
	gBrightGreen  = lipgloss.Color("#b8bb26")
	gBrightRed    = lipgloss.Color("#fb4934")
	gBrightYellow = lipgloss.Color("#fabd2f")

	// accent — основной акцентный цвет (bright aqua — мягкий зелёный).
	accent = gBrightAqua

	labelStyle = lipgloss.NewStyle().
			Foreground(gBrightBlue).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(gfg4)

	hintStyle = lipgloss.NewStyle().
			Foreground(gBrightYellow).
			Italic(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(gBrightRed).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(gBrightGreen).
			Bold(true)

	// buttonActive — яркая кнопка, когда она выбрана (тёмный текст на акцентном фоне).
	buttonActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(gbg0).
			Background(accent).
			Padding(0, 2)

	// buttonInactive — тусклая кнопка, когда не выбрана.
	buttonInactive = lipgloss.NewStyle().
			Foreground(gfg4).
			Padding(0, 2)

	// linkStyle — подчёркнутый текст для кликабельных ссылок.
	linkStyle = lipgloss.NewStyle().
			Foreground(accent).
			Underline(true)

	// inputBox — рамка вокруг поля ввода (синий акцент).
	inputBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gBlue).
			Padding(0, 1).
			Width(62)

	// separator — горизонтальная линия-разделитель.
	separator = lipgloss.NewStyle().
			Foreground(gbg4).
			Render(strings.Repeat("─", 50))

	// progressDot — закрашенный шаг прогресс-бара.
	progressDot = lipgloss.NewStyle().
			Foreground(accent).
			Render("●")

	// progressEmpty — незакрашенный шаг прогресс-бара.
	progressEmpty = lipgloss.NewStyle().
			Foreground(gfg4).
			Render("○")

	// dialogBox — внешняя рамка всего диалога.
	dialogBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gBlue).
			Padding(1, 2).
			Width(68)
)

// ─── Константы шагов ────────────────────────────────────────────────────────

const (
	stepToken = iota
	stepChatID
	stepCronAsk
	stepCronTime
	stepCronInterval
	stepPaths
	stepConfirm
	stepDone
)

const totalSteps = 7 // без учёта stepDone

// ─── CronSettings ────────────────────────────────────────────────────────────

type CronSettings struct {
	Enabled  bool
	Time     string // "HH:MM"
	Interval string // "daily", "weekly", "monthly"
}

// ─── Модель ─────────────────────────────────────────────────────────────────

type setupModel struct {
	step int

	tokenInput  textinput.Model
	chatIDInput textinput.Model

	timeInput     textinput.Model
	intervalInput textinput.Model
	pathInput     textinput.Model

	token        string
	chatID       string
	cronSettings CronSettings

	showError bool
	errMsg    string

	confirmFocus int
	cronYesNo    int
	pathFocus    int
	skipFocus    int // 0 — инпут, 1 — кнопка «Пропустить»

	paths       []string
	showTemplates bool
	templateCursor int

	width  int
	height int
}

// NewSetupModel — конструктор. Все инпуты создаются здесь и сохраняются в модель.
func NewSetupModel() setupModel {
	ti := textinput.New()
	ti.Placeholder = "Вставьте сюда токен бота"
	ti.CharLimit = 100
	ti.Width = 55
	ti.Focus()
	ti.Prompt = "▸ "

	ci := textinput.New()
	ci.Placeholder = "Например: 873925520"
	ci.CharLimit = 30
	ci.Width = 55
	ci.Prompt = "▸ "

	cronTime := textinput.New()
	cronTime.Placeholder = "HH:MM (например: 02:00)"
	cronTime.CharLimit = 5
	cronTime.Width = 20
	cronTime.Prompt = "▸ "

	cronInterval := textinput.New()
	cronInterval.Placeholder = "daily / weekly / monthly"
	cronInterval.CharLimit = 10
	cronInterval.Width = 20
	cronInterval.Prompt = "▸ "

	pathInput := textinput.New()
	pathInput.Placeholder = "Например: /panel/*"
	pathInput.CharLimit = 200
	pathInput.Width = 55
	pathInput.Prompt = "▸ "

	return setupModel{
		step:          stepToken,
		tokenInput:    ti,
		chatIDInput:   ci,
		timeInput:     cronTime,
		intervalInput: cronInterval,
		pathInput:     pathInput,
		cronSettings:  CronSettings{Enabled: false, Time: "02:00", Interval: "daily"},
		confirmFocus:  0,
		cronYesNo:     1,
		pathFocus:     0,
		skipFocus:     0,
	}
}

// ─── tea.Model ──────────────────────────────────────────────────────────────

func (m setupModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m setupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.step == stepPaths && m.showTemplates {
				m.showTemplates = false
				return m, nil
			}
			if m.step == stepToken {
				return m, tea.Quit
			}
			m.step--
			m.showError = false

			return m, m.updateFocusState()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m.handleStep(msg)
}

// View — главный метод рендеринга.
func (m setupModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Загрузка..."
	}

	showProgress := m.step != stepDone
	content := m.viewContent()

	return m.renderDialog(content, showProgress)
}

// viewContent — возвращает содержимое для текущего шага.
func (m setupModel) viewContent() string {
	switch m.step {
	case stepToken:
		return m.viewTokenStep()
	case stepChatID:
		return m.viewChatIDStep()
	case stepCronAsk:
		return m.viewCronAskStep()
	case stepCronTime:
		return m.viewCronTimeStep()
	case stepCronInterval:
		return m.viewCronIntervalStep()
	case stepPaths:
		return m.viewPathsStep()
	case stepConfirm:
		return m.viewConfirmStep()
	case stepDone:
		return m.viewDoneStep()
	default:
		return ""
	}
}

// ─── Фокус ──────────────────────────────────────────────────────────────────

func (m *setupModel) updateFocusState() tea.Cmd {
	m.tokenInput.Blur()
	m.chatIDInput.Blur()
	m.timeInput.Blur()
	m.intervalInput.Blur()
	m.pathInput.Blur()

	switch m.step {
	case stepToken:
		return m.tokenInput.Focus()
	case stepChatID:
		return m.chatIDInput.Focus()
	case stepCronTime:
		return m.timeInput.Focus()
	case stepCronInterval:
		return m.intervalInput.Focus()
	case stepPaths:
		return m.pathInput.Focus()
	}

	return nil
}

// ─── Делегирование Update по шагам ──────────────────────────────────────────

func (m setupModel) handleStep(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case stepToken:
		return m.updateStepToken(msg)
	case stepChatID:
		return m.updateStepChatID(msg)
	case stepCronAsk:
		return m.updateStepCronAsk(msg)
	case stepCronTime:
		return m.updateStepCronTime(msg)
	case stepCronInterval:
		return m.updateStepCronInterval(msg)
	case stepPaths:
		return m.updateStepPaths(msg)
	case stepConfirm:
		return m.updateStepConfirm(msg)
	case stepDone:
		return m, tea.Quit
	}

	return m, nil
}

// ─── Шаг 1: Токен ───────────────────────────────────────────────────────────

func (m setupModel) updateStepToken(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "left", "right":
			m.skipFocus = 1 - m.skipFocus
			return m, nil
		case "enter":
			if m.skipFocus == 1 {
				m.showError = false
				m.step = stepChatID
				m.skipFocus = 0
				return m, m.updateFocusState()
			}

			m.token = strings.TrimSpace(m.tokenInput.Value())
			if m.token == "" {
				m.showError = true
				m.errMsg = "Токен не может быть пустым"
				return m, nil
			}
			m.showError = false
			m.step = stepChatID
			m.skipFocus = 0

			return m, m.updateFocusState()
		}
	}

	var cmd tea.Cmd
	m.tokenInput, cmd = m.tokenInput.Update(msg)

	return m, cmd
}

func (m setupModel) viewTokenStep() string {
	b := new(strings.Builder)

	b.WriteString(labelStyle.Render("Токен Telegram-бота"))
	b.WriteString("\n\n")
	b.WriteString(inputBox.Render(m.tokenInput.View()))

	if m.showError {
		b.WriteString("\n\n")
		b.WriteString(errorStyle.Render("✗ " + m.errMsg))
	}

	b.WriteString("\n\n")
	b.WriteString(hintStyle.Render("Где взять?"))
	b.WriteString("  ")
	b.WriteString(
		hyperlink("https://t.me/BotFather", linkStyle.Render("@BotFather → Открыть в Telegram")),
	)

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Создайте бота в @BotFather и скопируйте токен"))

	b.WriteString("\n\n")
	inputBtn := buttonInactive
	skipBtn := buttonInactive
	if m.skipFocus == 0 {
		inputBtn = buttonActive
	} else {
		skipBtn = buttonActive
	}
	b.WriteString(inputBtn.Render("✎ Ввести"))
	b.WriteString("  ")
	b.WriteString(skipBtn.Render("Пропустить"))

	return b.String()
}

// ─── Шаг 2: Chat ID ─────────────────────────────────────────────────────────

func (m setupModel) updateStepChatID(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "left", "right":
			m.skipFocus = 1 - m.skipFocus
			return m, nil
		case "enter":
			if m.skipFocus == 1 {
				m.showError = false
				m.step = stepCronAsk
				m.skipFocus = 0
				return m, m.updateFocusState()
			}

			m.chatID = strings.TrimSpace(m.chatIDInput.Value())
			if m.chatID == "" {
				m.showError = true
				m.errMsg = "Chat ID не может быть пустым"
				return m, nil
			}
			m.showError = false
			m.step = stepCronAsk
			m.skipFocus = 0

			return m, m.updateFocusState()
		}
	}

	var cmd tea.Cmd
	m.chatIDInput, cmd = m.chatIDInput.Update(msg)

	return m, cmd
}

func (m setupModel) viewChatIDStep() string {
	b := new(strings.Builder)

	b.WriteString(labelStyle.Render("Chat ID получателя"))
	b.WriteString("\n\n")
	b.WriteString(inputBox.Render(m.chatIDInput.View()))

	if m.showError {
		b.WriteString("\n\n")
		b.WriteString(errorStyle.Render("✗ " + m.errMsg))
	}

	b.WriteString("\n\n")
	b.WriteString(hintStyle.Render("Где взять?"))
	b.WriteString("  ")
	b.WriteString(
		hyperlink(
			"https://t.me/userinfobot",
			linkStyle.Render("@userinfobot → Открыть в Telegram"),
		),
	)

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Напишите @userinfobot — он пришлёт ваш Chat ID"))

	b.WriteString("\n\n")
	inputBtn := buttonInactive
	skipBtn := buttonInactive
	if m.skipFocus == 0 {
		inputBtn = buttonActive
	} else {
		skipBtn = buttonActive
	}
	b.WriteString(inputBtn.Render("✎ Ввести"))
	b.WriteString("  ")
	b.WriteString(skipBtn.Render("Пропустить"))

	return b.String()
}

// ─── Шаг 3: Cron — да/нет ───────────────────────────────────────────────────

func (m setupModel) updateStepCronAsk(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "left", "right", "h", "l":
			m.cronYesNo = 1 - m.cronYesNo
		case "enter":
			if m.cronYesNo == 0 {
				m.step = stepCronTime

				return m, m.updateFocusState()
			}
			m.cronSettings.Enabled = false
			m.step = stepPaths

			return m, m.updateFocusState()
		}
	}

	return m, nil
}

func (m setupModel) viewCronAskStep() string {
	b := new(strings.Builder)

	b.WriteString(labelStyle.Render("Автоматические бекапы по расписанию?"))
	b.WriteString("\n\n")

	yesBtn := buttonInactive
	noBtn := buttonInactive
	if m.cronYesNo == 0 {
		yesBtn = buttonActive
	} else {
		noBtn = buttonActive
	}
	b.WriteString(yesBtn.Render("Да, настроить"))
	b.WriteString("  ")
	b.WriteString(noBtn.Render("Нет, пропустить"))

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Cron будет запускать бекап автоматически"))

	return b.String()
}

// ─── Шаг 4: Время Cron ─────────────────────────────────────────────────────

func (m setupModel) updateStepCronTime(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		m.cronSettings.Time = strings.TrimSpace(m.timeInput.Value())
		if m.cronSettings.Time == "" {
			m.cronSettings.Time = "02:00"
		}
		m.step = stepCronInterval

		return m, m.updateFocusState()
	}

	var cmd tea.Cmd
	m.timeInput, cmd = m.timeInput.Update(msg)

	return m, cmd
}

func (m setupModel) viewCronTimeStep() string {
	b := new(strings.Builder)

	b.WriteString(labelStyle.Render("Время бекапа"))
	b.WriteString("\n\n")
	b.WriteString(inputBox.Render(m.timeInput.View()))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Формат HH:MM. Например 02:00 — в два часа ночи"))

	return b.String()
}

// ─── Шаг 5: Интервал Cron ──────────────────────────────────────────────────

func (m setupModel) updateStepCronInterval(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		m.cronSettings.Enabled = true
		m.cronSettings.Interval = strings.TrimSpace(m.intervalInput.Value())
		if m.cronSettings.Interval == "" {
			m.cronSettings.Interval = "daily"
		}
		m.step = stepPaths

		return m, m.updateFocusState()
	}

	var cmd tea.Cmd
	m.intervalInput, cmd = m.intervalInput.Update(msg)

	return m, cmd
}

func (m setupModel) viewCronIntervalStep() string {
	b := new(strings.Builder)

	b.WriteString(labelStyle.Render("Периодичность"))
	b.WriteString("\n\n")
	b.WriteString(inputBox.Render(m.intervalInput.View()))
	b.WriteString("\n\n")
	b.WriteString(
		helpStyle.Render("daily — каждый день, weekly — раз в неделю, monthly — раз в месяц"),
	)

	return b.String()
}

// ─── Шаг 6: Пути для бекапа ─────────────────────────────────────────────────

func (m setupModel) updateStepPaths(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.showTemplates {
		return m.updateTemplateSelector(msg)
	}

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "left":
			m.pathFocus = (m.pathFocus - 1 + 3) % 3

			return m, nil

		case "right":
			m.pathFocus = (m.pathFocus + 1) % 3

			return m, nil

		case "enter":
			if m.pathFocus == 2 {
				m.showError = false
				m.step = stepConfirm

				return m, nil
			}

			if m.pathFocus == 1 {
				m.showTemplates = true
				m.templateCursor = 0
				m.showError = false

				return m, nil
			}

			path := strings.TrimSpace(m.pathInput.Value())
			if path == "" {
				m.showError = true
				m.errMsg = "Путь не может быть пустым"

				return m, nil
			}
			m.paths = append(m.paths, path)
			m.pathInput.SetValue("")
			m.showError = false

			return m, m.pathInput.Focus()
		}
	}

	var cmd tea.Cmd
	m.pathInput, cmd = m.pathInput.Update(msg)

	return m, cmd
}

type pathTemplate struct {
	Name  string
	Paths []string
}

var pathTemplates = []pathTemplate{
	{Name: "remnawave", Paths: []string{"/opt/remnawave/*", "/root/remnawave_backups/*"}},
}

func (m setupModel) updateTemplateSelector(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if m.templateCursor > 0 {
				m.templateCursor--
			}
		case "down", "j":
			if m.templateCursor < len(pathTemplates) {
				m.templateCursor++
			}
		case "enter":
			if m.templateCursor < len(pathTemplates) {
				t := pathTemplates[m.templateCursor]
				for _, p := range t.Paths {
					if !contains(m.paths, p) {
						m.paths = append(m.paths, p)
					}
				}
			}
			m.showTemplates = false
		}
	}

	return m, nil
}

func (m setupModel) viewPathsStep() string {
	if m.showTemplates {
		return m.viewTemplateSelector()
	}

	b := new(strings.Builder)

	b.WriteString(labelStyle.Render("Пути для бекапа"))
	b.WriteString("\n\n")
	b.WriteString(inputBox.Render(m.pathInput.View()))

	if len(m.paths) > 0 {
		b.WriteString("\n\n")
		b.WriteString(hintStyle.Render("Уже добавлены:"))
		for _, p := range m.paths {
			b.WriteString("\n  ")
			b.WriteString(helpStyle.Render("• " + p))
		}
	}

	if m.showError {
		b.WriteString("\n\n")
		b.WriteString(errorStyle.Render("✗ " + m.errMsg))
	}

	b.WriteString("\n\n")
	addBtn := buttonInactive
	tmplBtn := buttonInactive
	doneBtn := buttonInactive
	switch m.pathFocus {
	case 0:
		addBtn = buttonActive
	case 1:
		tmplBtn = buttonActive
	case 2:
		doneBtn = buttonActive
	}
	b.WriteString(addBtn.Render("+ Добавить"))
	b.WriteString("  ")
	b.WriteString(tmplBtn.Render("Шаблоны"))
	b.WriteString("  ")
	b.WriteString(doneBtn.Render("Готово"))

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Добавьте все нужные пути, затем нажмите «Готово»"))

	return b.String()
}

func (m setupModel) viewTemplateSelector() string {
	b := new(strings.Builder)

	b.WriteString(labelStyle.Render("Выберите шаблон"))
	b.WriteString("\n\n")

	var listItems []string
	for i, t := range pathTemplates {
		cursor := "  "
		style := helpStyle
		if i == m.templateCursor {
			cursor = "▸ "
			style = labelStyle
		}
		listItems = append(listItems, style.Render(cursor+t.Name))
	}

	listBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(gfg4).
		Padding(0, 1).
		Width(40).
		Render(strings.Join(listItems, "\n"))

	b.WriteString(listBox)

	b.WriteString("\n\n")

	// Кнопка «Назад» — всегда последняя в списке
	cursor := "  "
	btnStyle := buttonInactive
	if m.templateCursor == len(pathTemplates) {
		cursor = "▸ "
		btnStyle = buttonActive
	}
	b.WriteString(btnStyle.Render(cursor + "← Назад"))

	return b.String()
}

// ─── Шаг 7: Подтверждение ──────────────────────────────────────────────────

func (m setupModel) updateStepConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "left", "right", "h", "l":
			m.confirmFocus = 1 - m.confirmFocus
		case "enter":
			if m.confirmFocus == 0 {
				m.save() // внутри устанавливает step = stepDone или показывает ошибку

				return m, nil
			}

			return m, tea.Quit
		}
	}

	return m, nil
}

func (m setupModel) viewConfirmStep() string {
	b := new(strings.Builder)

	b.WriteString(labelStyle.Render("Проверьте данные"))
	b.WriteString("\n\n")

	// Карточка с введёнными данными — используем рамку для визуального выделения
	card := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(gBlue).
		Padding(0, 1)

	lines := ""
	lines += fmt.Sprintf("  Token:   %s\n", m.token)
	lines += fmt.Sprintf("  Chat ID: %s\n", m.chatID)
	if m.cronSettings.Enabled {
		lines += fmt.Sprintf("  Бекап:   %s (%s)\n", m.cronSettings.Time, m.cronSettings.Interval)
	} else {
		lines += "  Бекап:   не настроен\n"
	}
	b.WriteString(card.Render(strings.TrimRight(lines, "\n")))

	b.WriteString("\n\n")

	yesBtn := buttonInactive
	noBtn := buttonInactive
	if m.confirmFocus == 0 {
		yesBtn = buttonActive
	} else {
		noBtn = buttonActive
	}
	b.WriteString(yesBtn.Render("Да, сохранить"))
	b.WriteString("  ")
	b.WriteString(noBtn.Render("Нет, выйти"))

	return b.String()
}

// ─── Шаг 8: Готово ─────────────────────────────────────────────────────────

func (m setupModel) viewDoneStep() string {
	b := new(strings.Builder)

	b.WriteString(successStyle.Render("Настройки сохранены"))
	b.WriteString("\n\n")
	b.WriteString("• Файл .env обновлён\n")

	if m.cronSettings.Enabled {
		fmt.Fprintf(b, "• Cron: %s (%s)\n", m.cronSettings.Time, m.cronSettings.Interval)
	}

	fmt.Fprintf(b, "• Путей для бекапа: %d\n", len(m.paths))

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Запустить бекап: ./saveFile run"))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Нажмите любую клавишу для выхода..."))

	return b.String()
}

// ─── Сохранение ─────────────────────────────────────────────────────────────

func (m *setupModel) save() {
	if err := writePathsFile(m.paths); err != nil {
		m.errMsg = fmt.Sprintf("Ошибка записи path.txt: %v", err)
		m.showError = true

		return
	}

	// Если пользователь пропустил шаг — сохраняем старое значение из .env
	token := m.token
	chatID := m.chatID
	if token == "" {
		token = readEnvValue("TELEGRAM_TOKEN")
	}
	if chatID == "" {
		chatID = readEnvValue("TELEGRAM_CHAT_ID")
	}

	if err := writeEnvFile(token, chatID); err != nil {
		m.errMsg = fmt.Sprintf("Ошибка записи .env: %v", err)
		m.showError = true

		return
	}

	if m.cronSettings.Enabled {
		if err := installCronJob(m.cronSettings); err != nil {
			m.errMsg = fmt.Sprintf("Ошибка установки cron: %v", err)
			m.showError = true

			return
		}
	} else {
		// При повторной установке могли остаться старые записи — чистим.
		_ = uninstallCronJob()
	}

	m.step = stepDone
}

// ─── hyperlink — OSC-8 гиперссылка ────────────────────────────────────────
// Современные терминалы (iTerm2, Terminal.app, kitty, Windows Terminal)
// отображают такой текст как кликабельную ссылку. Формат: \x1b]8;;URL\x07TEXT\x1b]8;;\x07.
func hyperlink(url, text string) string {
	return fmt.Sprintf("\x1b]8;;%s\x07%s\x1b]8;;\x07", url, text)
}

// ─── renderDialog — центрирование и рамка ───────────────────────────────────
// Этот метод собирает все визуальные слои:
// 1. progress bar (шаг X из Y)
// 2. separator
// 3. content (то что нарисовала конкретная view-функция)
// 4. footer с подсказками по клавишам
// Всё это оборачивается в lipgloss.Place для центрирования по вертикали и горизонтали.

func (m setupModel) renderDialog(content string, showProgress bool) string {
	// Собираем внутренности диалога
	var parts []string

	if showProgress {
		parts = append(parts, m.progressBar())
		parts = append(parts, separator)
	}

	parts = append(parts, content)
	parts = append(parts, "")
	parts = append(parts, m.footer())

	dialogContent := lipgloss.JoinVertical(lipgloss.Left, parts...)
	dialogContent = dialogBox.Render(dialogContent)

	// Центрируем диалог на всём экране
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialogContent,
	)
}

// ─── progressBar — визуальный индикатор шагов ──────────────────────────────

func (m setupModel) progressBar() string {
	current := m.step // 0-based
	if current >= totalSteps {
		current = totalSteps - 1
	}

	var dots []string
	for i := range totalSteps {
		if i <= current {
			dots = append(dots, progressDot)
		} else {
			dots = append(dots, progressEmpty)
		}
	}

	bar := strings.Join(dots, " ")

	stepLabel := fmt.Sprintf("шаг %d из %d", current+1, totalSteps)

	return helpStyle.Render(bar + "   " + stepLabel)
}

// contains — проверяет, есть ли строка в срезе.
func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// ─── footer — подсказки по управлению ──────────────────────────────────────

func (m setupModel) footer() string {
	parts := []string{"enter — подтвердить", "esc — назад", "ctrl+c — выход"}

	if m.step == stepPaths && m.showTemplates {
		parts = []string{
			"↑ ↓ — выбрать",
			"enter — подтвердить",
			"esc — назад",
			"ctrl+c — выход",
		}

		return helpStyle.Render(strings.Join(parts, "  │  "))
	}

	switch m.step {
	case stepToken, stepChatID:
		parts = []string{
			"← → — переключить",
			"enter — подтвердить",
			"esc — назад",
			"ctrl+c — выход",
		}
	case stepPaths:
		parts = []string{
			"enter — добавить",
			"← → — переключить",
			"esc — назад",
			"ctrl+c — выход",
		}
	case stepCronAsk, stepConfirm:
		parts = []string{
			"← → / h l — выбрать",
			"enter — подтвердить",
			"esc — назад",
			"ctrl+c — выход",
		}
	}

	return helpStyle.Render(strings.Join(parts, "  │  "))
}
