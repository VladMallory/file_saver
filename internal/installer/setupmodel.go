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
	// Gruvbox dark palette
	gbg0  = lipgloss.Color("#282828") // dark0_hard — фон
	gbg4  = lipgloss.Color("#7c6f64") // dark4 — тусклый
	gfg0  = lipgloss.Color("#ebdbb2") // light1 — основной текст
	gfg4  = lipgloss.Color("#a89984") // gray — приглушённый

	gRed    = lipgloss.Color("#cc241d")
	gGreen  = lipgloss.Color("#98971a")
	gYellow = lipgloss.Color("#d79921")
	gBlue   = lipgloss.Color("#458588")
	gPurple = lipgloss.Color("#b16286")
	gAqua   = lipgloss.Color("#689d6a")
	gOrange = lipgloss.Color("#d65d0e")

	gBrightRed    = lipgloss.Color("#fb4934")
	gBrightGreen  = lipgloss.Color("#b8bb26")
	gBrightYellow = lipgloss.Color("#fabd2f")
	gBrightBlue   = lipgloss.Color("#83a598")
	gBrightPurple = lipgloss.Color("#d3869b")
	gBrightAqua   = lipgloss.Color("#8ec07c")
	gBrightOrange = lipgloss.Color("#fe8019")

	// accent — основной акцентный цвет (bright aqua — мягкий зелёный)
	accent = gBrightAqua

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accent)

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

	// buttonActive — яркая кнопка, когда она выбрана (тёмный текст на акцентном фоне)
	buttonActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(gbg0).
			Background(accent).
			Padding(0, 2)

	// buttonInactive — тусклая кнопка, когда не выбрана
	buttonInactive = lipgloss.NewStyle().
			Foreground(gfg4).
			Padding(0, 2)

	// linkStyle — подчёркнутый текст для кликабельных ссылок
	linkStyle = lipgloss.NewStyle().
			Foreground(accent).
			Underline(true)

	// inputBox — рамка вокруг поля ввода (синий акцент)
	inputBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gBlue).
			Padding(0, 1).
			Width(62)

	// separator — горизонтальная линия-разделитель (тёмный серый)
	separator = lipgloss.NewStyle().
			Foreground(gbg4).
			Render(strings.Repeat("─", 50))

	// progressDot — закрашенный шаг прогресс-бара
	progressDot = lipgloss.NewStyle().
			Foreground(accent).
			Render("●")

	// progressEmpty — незакрашенный шаг прогресс-бара
	progressEmpty = lipgloss.NewStyle().
			Foreground(gfg4).
			Render("○")

	// dialogBox — внешняя рамка всего диалога
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
	stepConfirm
	stepDone
)

const totalSteps = 6 // без учёта stepDone

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

	token        string
	chatID       string
	cronSettings CronSettings

	showError bool
	errMsg    string

	confirmFocus int
	cronYesNo    int

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

	return setupModel{
		step:           stepToken,
		tokenInput:     ti,
		chatIDInput:    ci,
		timeInput:      cronTime,
		intervalInput:  cronInterval,
		cronSettings:   CronSettings{Enabled: false, Time: "02:00", Interval: "daily"},
		confirmFocus:   0,
		cronYesNo:      1,
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
// Все шаги собираются через m.renderDialog, который центрирует контент на экране.
func (m setupModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Загрузка..."
	}

	var content string
	switch m.step {
	case stepToken:
		content = m.viewTokenStep()
	case stepChatID:
		content = m.viewChatIDStep()
	case stepCronAsk:
		content = m.viewCronAskStep()
	case stepCronTime:
		content = m.viewCronTimeStep()
	case stepCronInterval:
		content = m.viewCronIntervalStep()
	case stepConfirm:
		content = m.viewConfirmStep()
	case stepDone:
		content = m.viewDoneStep()
	}

	// stepDone не участвует в progress bar (это финальный экран)
	showProgress := m.step != stepDone
	return m.renderDialog(content, showProgress)
}

// ─── Фокус ──────────────────────────────────────────────────────────────────

func (m *setupModel) updateFocusState() tea.Cmd {
	m.tokenInput.Blur()
	m.chatIDInput.Blur()
	m.timeInput.Blur()
	m.intervalInput.Blur()

	switch m.step {
	case stepToken:
		return m.tokenInput.Focus()
	case stepChatID:
		return m.chatIDInput.Focus()
	case stepCronTime:
		return m.timeInput.Focus()
	case stepCronInterval:
		return m.intervalInput.Focus()
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
	case stepConfirm:
		return m.updateStepConfirm(msg)
	case stepDone:
		return m, tea.Quit
	}
	return m, nil
}

// ─── Шаг 1: Токен ───────────────────────────────────────────────────────────

func (m setupModel) updateStepToken(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		m.token = strings.TrimSpace(m.tokenInput.Value())
		if m.token == "" {
			m.showError = true
			m.errMsg = "Токен не может быть пустым"
			return m, nil
		}
		m.showError = false
		m.step = stepChatID
		return m, m.updateFocusState()
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
	// Сначала стиль (linkStyle), потом гиперссылка — чтобы lipgloss правильно измерил ширину.
	b.WriteString(hyperlink("https://t.me/BotFather", linkStyle.Render("@BotFather → Открыть в Telegram")))

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Создайте бота в @BotFather и скопируйте токен"))

	return b.String()
}

// ─── Шаг 2: Chat ID ─────────────────────────────────────────────────────────

func (m setupModel) updateStepChatID(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		m.chatID = strings.TrimSpace(m.chatIDInput.Value())
		if m.chatID == "" {
			m.showError = true
			m.errMsg = "Chat ID не может быть пустым"
			return m, nil
		}
		m.showError = false
		m.step = stepCronAsk
		return m, m.updateFocusState()
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
	b.WriteString(hyperlink("https://t.me/userinfobot", linkStyle.Render("@userinfobot → Открыть в Telegram")))

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Напишите @userinfobot — он пришлёт ваш Chat ID"))

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
			m.step = stepConfirm
			return m, nil
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
		m.step = stepConfirm
		return m, nil
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
	b.WriteString(helpStyle.Render("daily — каждый день, weekly — раз в неделю, monthly — раз в месяц"))

	return b.String()
}

// ─── Шаг 6: Подтверждение ──────────────────────────────────────────────────

func (m setupModel) updateStepConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "left", "right", "h", "l":
			m.confirmFocus = 1 - m.confirmFocus
		case "enter":
			if m.confirmFocus == 0 {
				m.save()
				m.step = stepDone
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

// ─── Шаг 7: Готово ─────────────────────────────────────────────────────────

func (m setupModel) viewDoneStep() string {
	b := new(strings.Builder)

	b.WriteString(successStyle.Render("Настройки сохранены"))
	b.WriteString("\n\n")
	b.WriteString("• Файл .env обновлён\n")

	if m.cronSettings.Enabled {
		b.WriteString(fmt.Sprintf("• Cron: %s (%s)\n", m.cronSettings.Time, m.cronSettings.Interval))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Запустить бекап: ./saveFile run"))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Нажмите любую клавишу для выхода..."))

	return b.String()
}

// ─── Сохранение ─────────────────────────────────────────────────────────────

func (m setupModel) save() {
	writeEnvFile(m.token, m.chatID)

	if m.cronSettings.Enabled {
		installCronJob(m.cronSettings)
	}
}

// ─── hyperlink — OSC-8 гиперссылка ────────────────────────────────────────
// Современные терминалы (iTerm2, Terminal.app, kitty, Windows Terminal)
// отображают такой текст как кликабельную ссылку.
// \x1b]8;;URL\x07TEXT\x1b]8;;\x07
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
	for i := 0; i < totalSteps; i++ {
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

// ─── footer — подсказки по управлению ──────────────────────────────────────

func (m setupModel) footer() string {
	parts := []string{"enter — подтвердить", "esc — назад", "ctrl+c — выход"}

	if m.step == stepCronAsk || m.step == stepConfirm {
		parts = []string{"← → / h l — выбрать", "enter — подтвердить", "esc — назад", "ctrl+c — выход"}
	}

	return helpStyle.Render(strings.Join(parts, "  │  "))
}
