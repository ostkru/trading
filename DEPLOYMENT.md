# Инструкция по развертыванию PortalData Trading API

## 🚀 Обзор

Данная инструкция описывает процесс развертывания PortalData Trading API с использованием nginx reverse proxy для продакшена.

## 📋 Предварительные требования

### Системные требования
- Ubuntu 20.04+ или CentOS 8+
- Go 1.19+
- MySQL 8.0+
- Nginx 1.18+
- PHP 7.4+ (для тестирования)

### Сетевые требования
- Домен: `api.portaldata.ru`
- SSL сертификат для домена
- Открытые порты: 80, 443, 8090

## 🔧 Шаг 1: Подготовка сервера

### Установка зависимостей
```bash
# Обновление системы
sudo apt update && sudo apt upgrade -y

# Установка Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Установка MySQL
sudo apt install mysql-server -y
sudo mysql_secure_installation

# Установка Nginx
sudo apt install nginx -y

# Установка PHP и зависимостей для тестирования
sudo apt install php php-curl php-json php-mbstring -y
```

### Настройка MySQL
```sql
-- Создание базы данных
CREATE DATABASE portaldata CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Создание пользователя
CREATE USER 'portaldata'@'localhost' IDENTIFIED BY 'secure_password';
GRANT ALL PRIVILEGES ON portaldata.* TO 'portaldata'@'localhost';
FLUSH PRIVILEGES;
```

## 🔧 Шаг 2: Развертывание приложения

### Клонирование репозитория
```bash
cd /var/www
git clone https://github.com/ostkru/go-mod.git
cd go-mod
```

### Настройка конфигурации
```bash
# Создание конфигурационного файла
cat > config.json << EOF
{
    "database": {
        "host": "localhost",
        "port": 3306,
        "user": "portaldata",
        "password": "secure_password",
        "name": "portaldata"
    },
    "server": {
        "port": 8090,
        "host": "0.0.0.0"
    }
}
EOF
```

### Сборка приложения
```bash
# Установка зависимостей
go mod download

# Сборка бинарного файла
go build -o app cmd/api/main.go

# Создание systemd сервиса
sudo tee /etc/systemd/system/portaldata-api.service > /dev/null << EOF
[Unit]
Description=PortalData Trading API
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/go-mod
ExecStart=/var/www/go-mod/app
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# Активация сервиса
sudo systemctl daemon-reload
sudo systemctl enable portaldata-api
sudo systemctl start portaldata-api
```

## 🔧 Шаг 3: Настройка Nginx

### SSL сертификаты
```bash
# Установка Certbot для Let's Encrypt
sudo apt install certbot python3-certbot-nginx -y

# Получение SSL сертификата
sudo certbot --nginx -d api.portaldata.ru
```

### Конфигурация Nginx
```bash
# Копирование конфигурации
sudo cp nginx-trading.conf /etc/nginx/sites-available/portaldata-trading

# Активация сайта
sudo ln -s /etc/nginx/sites-available/portaldata-trading /etc/nginx/sites-enabled/

# Удаление дефолтного сайта
sudo rm /etc/nginx/sites-enabled/default

# Проверка конфигурации
sudo nginx -t

# Перезапуск nginx
sudo systemctl reload nginx
```

## 🔧 Шаг 4: Настройка файрвола

```bash
# Установка UFW
sudo apt install ufw -y

# Настройка правил
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443
sudo ufw allow 8090

# Активация файрвола
sudo ufw enable
```

## 🔧 Шаг 5: Тестирование

### Базовое тестирование
```bash
# Проверка статуса сервисов
sudo systemctl status portaldata-api
sudo systemctl status nginx
sudo systemctl status mysql

# Проверка портов
netstat -tlnp | grep -E ':(80|443|8090)'

# Тестирование локального API
curl -H "Authorization: Bearer f428fbc16a97b9e2a55717bd34e97537ec34cb8c04a5f32eeb4e88c9ee998a53" \
     http://localhost:8090/api/v1/products
```

### Тестирование продакшена
```bash
# Запуск тестов продакшена
./test_production.sh

# Полное тестирование
php comprehensive_api_test.php
```

## 🔧 Шаг 6: Мониторинг и логирование

### Настройка логирования
```bash
# Создание директории для логов
sudo mkdir -p /var/log/portaldata
sudo chown www-data:www-data /var/log/portaldata

# Настройка logrotate
sudo tee /etc/logrotate.d/portaldata-api > /dev/null << EOF
/var/log/portaldata/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 www-data www-data
    postrotate
        systemctl reload portaldata-api
    endscript
}
EOF
```

### Настройка мониторинга
```bash
# Установка htop для мониторинга
sudo apt install htop -y

# Создание скрипта мониторинга
cat > monitor.sh << 'EOF'
#!/bin/bash
echo "=== PortalData API Status ==="
echo "Time: $(date)"
echo ""

echo "Services:"
systemctl is-active portaldata-api
systemctl is-active nginx
systemctl is-active mysql

echo ""
echo "Ports:"
netstat -tlnp | grep -E ':(80|443|8090)'

echo ""
echo "Memory usage:"
free -h

echo ""
echo "Disk usage:"
df -h /

echo ""
echo "Recent logs:"
tail -5 /var/log/portaldata/app.log
EOF

chmod +x monitor.sh
```

## 🔧 Шаг 7: Автоматическое обновление

### Скрипт деплоя
```bash
cat > deploy.sh << 'EOF'
#!/bin/bash

echo "🚀 Starting deployment..."

# Остановка сервиса
sudo systemctl stop portaldata-api

# Обновление кода
git pull origin main

# Сборка нового бинарного файла
go build -o app cmd/api/main.go

# Обновление документации
redoc-cli bundle openapi.json -o redoc-documentation.html

# Запуск сервиса
sudo systemctl start portaldata-api

# Проверка статуса
sleep 5
sudo systemctl status portaldata-api

echo "✅ Deployment completed!"
EOF

chmod +x deploy.sh
```

## 🔧 Шаг 8: Резервное копирование

### Скрипт бэкапа
```bash
cat > backup.sh << 'EOF'
#!/bin/bash

BACKUP_DIR="/var/backups/portaldata"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# Бэкап базы данных
mysqldump -u portaldata -p portaldata > $BACKUP_DIR/db_backup_$DATE.sql

# Бэкап конфигурации
tar -czf $BACKUP_DIR/config_backup_$DATE.tar.gz \
    config.json \
    nginx-trading.conf \
    /etc/systemd/system/portaldata-api.service

# Бэкап логов
tar -czf $BACKUP_DIR/logs_backup_$DATE.tar.gz /var/log/portaldata/

# Очистка старых бэкапов (оставляем последние 7 дней)
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete

echo "Backup completed: $BACKUP_DIR"
EOF

chmod +x backup.sh
```

## 🔧 Шаг 9: Проверка безопасности

### Проверка SSL
```bash
# Проверка SSL сертификата
openssl s_client -connect api.portaldata.ru:443 -servername api.portaldata.ru

# Проверка через SSL Labs (онлайн)
echo "Проверьте SSL на https://www.ssllabs.com/ssltest/analyze.html?d=api.portaldata.ru"
```

### Проверка безопасности
```bash
# Проверка открытых портов
nmap -sT -O localhost

# Проверка процессов
ps aux | grep -E '(app|nginx|mysql)'

# Проверка прав доступа
ls -la /var/www/go-mod/
ls -la /etc/nginx/sites-enabled/
```

## 📊 Мониторинг производительности

### Настройка мониторинга
```bash
# Установка инструментов мониторинга
sudo apt install sysstat iotop -y

# Настройка sar для сбора статистики
sudo sed -i 's/ENABLED="false"/ENABLED="true"/' /etc/default/sysstat
sudo systemctl enable sysstat
sudo systemctl start sysstat
```

## 🚨 Устранение неполадок

### Частые проблемы

1. **Сервис не запускается**
```bash
sudo systemctl status portaldata-api
sudo journalctl -u portaldata-api -f
```

2. **Nginx не работает**
```bash
sudo nginx -t
sudo systemctl status nginx
sudo tail -f /var/log/nginx/error.log
```

3. **Проблемы с базой данных**
```bash
sudo systemctl status mysql
mysql -u portaldata -p portaldata -e "SHOW TABLES;"
```

4. **Проблемы с SSL**
```bash
sudo certbot certificates
sudo certbot renew --dry-run
```

## 📞 Поддержка

При возникновении проблем:

1. Проверьте логи: `sudo journalctl -u portaldata-api -f`
2. Проверьте статус сервисов: `./monitor.sh`
3. Создайте issue в репозитории: https://github.com/ostkru/go-mod/issues

## 🔄 Обновления

Для обновления системы:

```bash
# Обновление системы
sudo apt update && sudo apt upgrade -y

# Обновление приложения
./deploy.sh

# Проверка после обновления
./test_production.sh
``` 