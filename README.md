# Буровичок

Десктопное приложение для хакатона от РосНефти по анализу месторождений

Для получения exe файла выполните

```Bash
go build -o burovichok.exe cmd/main.go
```

Содержит зависимость PostgreSQL и minIO

Для локального запуска можно воспользоваться `docker-compose` образом

Для запуска
```bash
docker-compose -f build/docker-compose/docker-compose.yaml up
```
