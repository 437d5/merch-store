# Тестовое задание Avito Backend 2025 зима.

В директории `api` находится документация API.

Стек:

  * Go
  * Posgtres
  * Docker

Внутри go использовался Gin и pxgpool. Реализована чистая архитектура с разделением на слои.

В директории `test` лежат Python файлы с двумя простыми сценариями e2e тестов и файл для нагрузочного тестирования с использованием Locust.

# Запуск 

### Запуск API
```
git clone https://github.com/437d5/merch-store.git
```
```
cd merch-store/
```
```
docker compose up --build
```

### Запуск E2E

Нужно перейти в директорию test/e2e_test
```
python3 -m venv venv
```
```
source venv/bin/activate
```
```
pip install -r requirements.txt
```
```
pytest e2e.py
```

### Запуск нагрузочного тестирования

Нужно перейти в директорию test/load_test
```
python3 -m venv venv
```
```
source venv/bin/activate
```
```
pip install -r requirements.txt
```
```
locust -f locustfile.py --host http://localhost:8080
```
