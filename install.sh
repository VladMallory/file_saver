#!/bin/bash

REPO="VladMallory/file_saver"
DEST_DIR="/root/save-file"
FINAL_NAME="file_saver"

# Определяем архитектуру процессора (x86_64 или ARM)
ARCH=$(uname -m)
case $ARCH in
    x86_64) BINARY_NAME="saveFile-linux-amd64" ;;
    aarch64|arm64) BINARY_NAME="file_saver-linux-arm64" ;;
    *) echo "Ошибка: архитектура $ARCH не поддерживается"; exit 1 ;;
esac

echo "Ищем последнюю версию для $ARCH..."

# Получаем ссылку на скачивание из последнего релиза на GitHub
DOWNLOAD_URL=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep "browser_download_url.*$BINARY_NAME" | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Ошибка: не удалось найти релиз $BINARY_NAME в репозитории $REPO."
    echo "Убедись, что на GitHub создан Release и к нему прикреплен скомпилированный файл!"
    exit 1
fi

echo "Скачиваем: $DOWNLOAD_URL"
curl -L -o /tmp/$FINAL_NAME "$DOWNLOAD_URL"

echo "Даем права на выполнение..."
chmod +x /tmp/$FINAL_NAME

echo "Создаем папку и устанавливаем утилиту..."
# Проверяем, запущен ли скрипт от root (ID 0). Если нет — используем sudo.
if [ "$(id -u)" -eq 0 ]; then
    mkdir -p "$DEST_DIR"
    mv /tmp/$FINAL_NAME "$DEST_DIR/$FINAL_NAME"
else
    sudo mkdir -p "$DEST_DIR"
    sudo mv /tmp/$FINAL_NAME "$DEST_DIR/$FINAL_NAME"
fi

echo "---"
echo "Готово! Утилита успешно установлена."
echo "Запустить её можно командой: $DEST_DIR/$FINAL_NAME"
