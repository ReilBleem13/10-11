# 10-11

# Добавление сервисов для проверки
curl --location 'http://localhost:8080/check' \
    --header 'Content-Type: application/json' \
    --data '{
        "links": [
            "google.com",
            "vk.ru"
        ]
    }'

# Получение PDF-отчета для набора 1
curl --location 'http://localhost:8080/report' \
    --header 'Content-Type: application/json' \
    --data '{
        "links_list": [1] 
    }' \
    --output report_set_1.pdf


Проверенные наборы ссылок хранятся в формате JSON в директории data.