package ui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	chartService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/chart"
	"github.com/pkg/browser"

	"os"
	"strings"
	"time"
)

func (s *Service) showMainMenu(ctx context.Context) {
	cell := fyne.NewSize(400, 100)

	importsBtn := widget.NewButton("Импортирование данных", func() { s.showImportView(ctx) })
	reportsBtn := widget.NewButton("Создание Технологической карты", func() { s.showReportsView(ctx) })
	chartsBtn := widget.NewButton("Создание графиков", func() { s.showChartsView(ctx) })
	exportBtn := widget.NewButton("Экспортирование данных", func() { s.showExportView(ctx) })
	guidebooksBtn := widget.NewButton("Редактирование справочников", func() { s.showGuidebookView(ctx) })

	// NewGridWrap принимает размер ячейки — и «упаковывает» каждый элемент в box этого размера
	grid := container.NewGridWrap(cell,
		importsBtn,
		reportsBtn,
		chartsBtn,
		exportBtn,
		guidebooksBtn,
	)

	s.window.SetContent(container.NewCenter(grid))
}

func (s *Service) showImportView(ctx context.Context) {
	back := widget.NewButton("◀ Домой", func() { s.showMainMenu(ctx) })
	content := s.buildImportContent(ctx)
	s.window.SetContent(container.NewBorder(back, nil, nil, nil, content))
}

func (s *Service) showReportsView(ctx context.Context) {
	back := widget.NewButton("◀ Домой", func() { s.showMainMenu(ctx) })
	// Отчёты: только блок 5
	reportBtn := widget.NewButton("Заполнить Шапку Отчета (Блок 5)", func() {
		s.showBlockFiveForm(ctx)
	})
	s.window.SetContent(container.NewBorder(back, nil, nil, nil,
		container.NewVBox(
			widget.NewLabel("Отчёты"),
			widget.NewSeparator(),
			reportBtn,
		),
	))
}

func (s *Service) showChartsView(ctx context.Context) {
	back := widget.NewButton("◀ Домой", func() { s.showMainMenu(ctx) })
	chartBtn1 := widget.NewButton("2. Интерактивный График Pзаб/Тзаб (Блок 1)", func() {
		// вставьте сюда тот же код из Run для chartBtn1
		blockOneData, err := s.memStorage.GetTableOneData()
		if err != nil {
			dialog.ShowError(fmt.Errorf("не удалось получить данные Блока 1: %w", err), s.window)
			return
		}
		if len(blockOneData) == 0 {
			dialog.ShowInformation("Нет данных", "Недостаточно данных для построения графика", s.window)
			return
		}
		htmlPath, err := s.chart.GenerateTableOneChart(blockOneData)
		if err != nil {
			dialog.ShowError(fmt.Errorf("ошибка генерации HTML графика: %w", err), s.window)
			return
		}
		s.serverMutex.Lock()
		s.chartHtmlToServe = htmlPath
		s.serverMutex.Unlock()
		if err := s.startLocalWebServer(); err != nil {
			s.serverMutex.Lock()
			s.chartHtmlToServe = ""
			s.serverMutex.Unlock()
			dialog.ShowError(fmt.Errorf("ошибка запуска веб-сервера для графика: %w", err), s.window)
			return
		}
		time.Sleep(150 * time.Millisecond)
		s.serverMutex.Lock()
		port := s.serverPort
		s.serverMutex.Unlock()
		if port == "" {
			dialog.ShowError(fmt.Errorf("не удалось получить порт веб-сервера"), s.window)
			return
		}
		url := fmt.Sprintf("http://127.0.0.1:%s/", port)
		s.zLog.Debugw("Opening chart via web server", "url", url, "serving", htmlPath) // Логгируем

		if err := browser.OpenURL(url); err != nil {
			dialog.ShowError(fmt.Errorf("не удалось открыть браузер: %w", err), s.window)
		}
	})

	chartBtn2 := widget.NewButton("1. Интерактивный График Ртр, Рзтр, Рлин (Блок 2)", func() {
		// 1) Создаём select
		unitSelect := widget.NewSelect([]string{"kgf/cm2", "bar", "atm"}, func(string) {})
		unitSelect.PlaceHolder = "Выберите единицу"

		// 2) Упаковываем в контейнер
		content := container.NewVBox(
			widget.NewLabel("Единица измерения:"),
			unitSelect,
		)

		// 3) Создаём диалог с «OK»/«Отмена»
		dlg := dialog.NewCustomConfirm(
			"Выбор единицы",
			"OK", "Отмена",
			content,
			func(ok bool) {
				if !ok {
					// Отмена
					return
				}
				unit := unitSelect.Selected
				if unit == "" {
					dialog.ShowInformation("Не выбрана единица",
						"Пожалуйста, выберите единицу измерения",
						s.window)
					return
				}

				// 4) Ваш существующий код построения графика
				blockTwoData, err := s.memStorage.GetTableTwoData()
				if err != nil {
					dialog.ShowError(fmt.Errorf("не удалось получить данные Блока 2: %w", err), s.window)
					return
				}
				if len(blockTwoData) == 0 {
					dialog.ShowInformation("Нет данных",
						"Недостаточно данных для построения графика",
						s.window)
					return
				}
				htmlPath, err := s.chart.GenerateTableTwoChart(blockTwoData, unit)
				if err != nil {
					dialog.ShowError(fmt.Errorf("ошибка генерации HTML графика: %w", err), s.window)
					return
				}
				s.serverMutex.Lock()
				s.chartHtmlToServe = htmlPath
				s.serverMutex.Unlock()
				if err := s.startLocalWebServer(); err != nil {
					s.serverMutex.Lock()
					s.chartHtmlToServe = ""
					s.serverMutex.Unlock()
					dialog.ShowError(fmt.Errorf("ошибка запуска веб-сервера: %w", err), s.window)
					return
				}
				// Небольшая задержка, чтобы сервер успел подняться
				time.Sleep(150 * time.Millisecond)
				s.serverMutex.Lock()
				port := s.serverPort
				s.serverMutex.Unlock()
				if port == "" {
					dialog.ShowError(fmt.Errorf("не удалось получить порт веб-сервера"), s.window)
					return
				}
				url := fmt.Sprintf("http://127.0.0.1:%s/", port)
				s.zLog.Debugw("Opening chart via web server", "url", url, "serving", htmlPath) // Логгируем

				if err := browser.OpenURL(url); err != nil {
					dialog.ShowError(fmt.Errorf("не удалось открыть браузер: %w", err), s.window)
				}
			},
			s.window,
		)

		// Опционально: сразу задаём размер диалога, чтобы Select не обрезался
		dlg.Resize(fyne.NewSize(240, 120))
		dlg.Show()
	})

	chartBtn3 := widget.NewButton("3 Интерактивный график Дебитов (Блок 3)", func() {
		blockThreeData, err := s.memStorage.GetTableThreeData()
		if err != nil {
			dialog.ShowError(fmt.Errorf("не удалось получить данные Блока 3: %w", err), s.window)
			return
		}
		if len(blockThreeData) == 0 {
			dialog.ShowInformation("Нет данных", "Недостаточно данных для построения графика", s.window)
			return
		}
		htmlPath, err := s.chart.GenerateTableThreeChart(blockThreeData)
		if err != nil {
			dialog.ShowError(fmt.Errorf("ошибка генерации HTML графика: %w", err), s.window)
			return
		}
		s.serverMutex.Lock()
		s.chartHtmlToServe = htmlPath
		s.serverMutex.Unlock()
		if err := s.startLocalWebServer(); err != nil {
			s.serverMutex.Lock()
			s.chartHtmlToServe = ""
			s.serverMutex.Unlock()
			dialog.ShowError(fmt.Errorf("ошибка запуска веб-сервера для графика: %w", err), s.window)
			return
		}
		time.Sleep(150 * time.Millisecond)
		s.serverMutex.Lock()
		port := s.serverPort
		s.serverMutex.Unlock()
		if port == "" {
			dialog.ShowError(fmt.Errorf("не удалось получить порт веб-сервера"), s.window)
			return
		}
		url := fmt.Sprintf("http://127.0.0.1:%s/", port)
		s.zLog.Debugw("Opening chart via web server", "url", url, "serving", htmlPath) // Логгируем

		if err := browser.OpenURL(url); err != nil {
			dialog.ShowError(fmt.Errorf("не удалось открыть браузер: %w", err), s.window)
		}
	})

	s.window.SetContent(container.NewBorder(back, nil, nil, nil,
		container.NewVBox(
			widget.NewLabel("Графики"),
			widget.NewSeparator(),
			chartBtn2,
			chartBtn1,
			chartBtn3,
		),
	))
}

func (s *Service) showExportView(ctx context.Context) {
	s.zLog.Debugw("Open export view")

	t1, _ := s.memStorage.GetTableOneData()
	t2, _ := s.memStorage.GetTableTwoData()
	t3, _ := s.memStorage.GetTableThreeData()
	t4, _ := s.memStorage.GetTableFourData()
	t5, _ := s.memStorage.GetTableFiveData()

	//Проверяем, что все блоки импортированы
	if (len(t1) * len(t2) * len(t3) * len(t4)) == 0 {
		dialog.ShowError(fmt.Errorf("сначала ипортируйте все 4 блока"), s.window)
		return
	}
	//Проверяем, что техническая карта составлена
	if t5 == (models.TableFive{}) {
		dialog.ShowError(fmt.Errorf("сначала заполните Тех. карту"), s.window)
		return
	}
	// Проверяем, есть ли элементы в директории
	entries, err := os.ReadDir(chartService.HtmlChartsDirectory)
	if err != nil {
		s.zLog.Errorw("Директория "+chartService.HtmlChartsDirectory+" не создалась", err)
		dialog.ShowError(fmt.Errorf("отсутствует директория с графиками"), s.window)
		return
	}
	//Проверяем, что все три графика нарисованы
	//TODO надо удалять старые иначе именно они будут отправляться
	if len(entries) != 3 {
		s.zLog.Errorw("Директория "+chartService.HtmlChartsDirectory+" содержит не правильное число графиков", err)
		dialog.ShowError(fmt.Errorf("директория "+chartService.HtmlChartsDirectory+" содержит не правильное число графиков, должно быть 3"), s.window)
		return
	}
	arch, err := s.archiver.Archive(t1, t2, t3, t4, t5)

	if err != nil {
		s.zLog.Errorw("Ошибка инициализации архива в буфер")
		dialog.ShowError(err, s.window)
		return
	}
	filenamebase := strings.TrimSpace(t5.FieldName + "_" + t5.ResearchType)
	filenamebase = strings.ReplaceAll(filenamebase, " ", "-")

	go func() {
		uploadCtx := context.Background()                          // Используем новый контекст для фоновой задачи
		s.showLoadingIndicator("Выгрузка архива: " + filenamebase) // Показываем индикатор

		// Вызываем Upload из exporter, передавая ему имя и буфер
		info, err := s.exporter.Upload(uploadCtx, filenamebase, arch)
		finalErr := err // Сохраняем ошибку выгрузки

		if finalErr == nil {
			// Если выгрузка УСПЕШНА, пытаемся сохранить метаданные
			s.zLog.Infow("Экспорт данных успешен, сохраняем метаданные", "data", info)

			// Создаем модель для сохранения в БД
			archiveInfoModel := models.ArchiveInfo{
				ObjectName: info.Key, // exporter.Upload должен вернуть имя объекта (info.Key?)
				BucketName: info.Bucket,
				Size:       info.Size,
				ETag:       info.ETag,
				UploadedAt: time.Now(), // Фиксируем время сохранения метаданных
			}

			// Вызываем метод сервиса БД
			errDb := s.db.SaveArchiveInfo(uploadCtx, archiveInfoModel)
			if errDb != nil {
				// Логируем ошибку сохранения метаданных, но не перезаписываем finalErr,
				// т.к. основная операция (выгрузка) прошла успешно.
				s.zLog.Errorw("Не удалось сохранить метаданные архива в БД", "object_name", info.Key, "error", errDb)
				// Можно показать пользователю предупреждение, что выгрузка прошла, но запись в БД не удалась.
				// Например, добавив к сообщению об успехе "(ошибка сохранения деталей)"
			}
		}

		// Скрываем индикатор и показываем результат (с учетом ошибки выгрузки)
		s.hideLoadingIndicator(finalErr) // Показываем только ошибку выгрузки, если она была

		if finalErr == nil {
			dialog.ShowInformation("Экспорт", "Данные экспортированы успешно!", s.window)
		}
		// Ошибку выгрузки покажет hideLoadingIndicator

	}() // Конец горутины

}

func (s *Service) showGuidebookView(ctx context.Context) {
	s.zLog.Debugw("Opening Guidebook Management view")

	backBtn := widget.NewButton("◀ Домой", func() { s.showMainMenu(ctx) })

	// Определяем типы справочников, которые можно редактировать
	guidebookTypes := []string{
		"Месторождение",         // -> oilfield
		"Продуктивный горизонт", // -> productive_horizon
		"Тип прибора",           // -> instrument_type
		"Вид исследования",      // -> research_type
	}
	guidebookTypeSelect := widget.NewSelect(guidebookTypes, nil)
	guidebookTypeSelect.PlaceHolder = "Выберите тип справочника"

	newValueEntry := widget.NewEntry()
	newValueEntry.PlaceHolder = "Введите новое значение"
	newValueEntry.Validator = validation.NewRegexp(`.+`, "Поле не может быть пустым")

	statusLabel := widget.NewLabel("")

	addBtn := widget.NewButton("Добавить значение", func() {
		statusLabel.SetText("")
		selectedType := guidebookTypeSelect.Selected
		newValue := strings.TrimSpace(newValueEntry.Text)

		if selectedType == "" || newValue == "" {
			dialog.ShowInformation("Ошибка", "Выберите тип справочника и введите непустое значение.", s.window)
			return
		}

		err := newValueEntry.Validate()
		if err != nil {
			statusLabel.SetText("Ошибка валидации: " + err.Error())
			return
		}

		err = s.addGuidebookEntry(ctx, selectedType, newValue)
		if err != nil {
			// Обрабатываем возможную ошибку уникальности (если уберем ON CONFLICT) или другую ошибку БД
			// В случае ON CONFLICT DO NOTHING ошибки не будет при дубликате
			s.zLog.Errorw("Failed to add guidebook entry", "type", selectedType, "value", newValue, "error", err)
			dialog.ShowError(fmt.Errorf("Не удалось добавить значение: %w", err), s.window)
			statusLabel.SetText("Ошибка при добавлении.")
		} else {
			s.zLog.Infow("Successfully added guidebook entry", "type", selectedType, "value", newValue)
			statusLabel.SetText(fmt.Sprintf("Значение '%s' добавлено в '%s'.", newValue, selectedType))
			newValueEntry.SetText("")
		}
	})

	content := container.NewVBox(
		widget.NewLabel("Управление справочниками"),
		widget.NewSeparator(),
		guidebookTypeSelect,
		newValueEntry,
		addBtn,
		statusLabel,
	)

	s.window.SetContent(container.NewBorder(backBtn, nil, nil, nil, content))
}
