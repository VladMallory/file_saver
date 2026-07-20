# Установка

```bash
curl -s https://raw.githubusercontent.com/VladMallory/file_saver/main/install.sh | bash
```

# Настройка
В папке `save-file` нужно создать `.env` и там прописать переменные с своими параметрами
```.env
TELEGRAM_TOKEN=
TELEGRAM_CHAT_ID=
```

```bash
cd save-file
nano .env
```

- `TELEGRAM_TOKEN` - берется из @BotFather. Нужно будет создать бота, или если уже есть, взять токен из существующего
- `TELEGRAM_CHAT_ID` - берется из @Getmyid_bot. **Вроде бы не требует подписку и не показывает рекламу**

# Запуск 
```bash
make
```

Посмотреть логи
```bash
make log
```
