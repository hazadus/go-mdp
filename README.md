# go-mdp

CLI tool для конвертирования файлов в формате Markdown в HTML и их просмотра в браузере.

## Интересные фичи

- ✅ Определение целевой ОС для выбора стандартного средства просмотра
- ✅ Использование golden-файлов для тестирования
- ✅ Создание временных файлов средствами ОС
- ✅ Запуск внешних программ
- ✅ Использование внешних модулей для конвертации Markdown и очистки HTML

## Сборка

```bash
make build
```

## Использование

```bash
# Конвертировать в HTML и открыть полученный файл
mdp -file ./README.md
# Только конвертировать в HTML
mdp -file ./README.md -s
```

В результате сборки будет создан файл `./bin/mdp`.

## Запуск тестов

```bash
make test
```
