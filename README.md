# Воркшоп по тестированию в Go ...

Данный проект предназначен для практического изучения методов тестирования в языке программирования Go. В ходе воркшопа студенты познакомятся с различными техниками тестирования, включая модульное тестирование, интеграционное тестирование и тестирование производительности.

## Инструкции

1. Склонируйте репозиторий:
   ```
   git clone https://github.com/your-username/ya-praktukum-go-testing-workshop.git
   cd ya-praktukum-go-testing-workshop
   ```

2. Установите зависимости:
   ```
   go mod tidy
   ```

3. Изучите код проекта.

4. На первом этапе попробуйте написать unit тесты не меняя дизайн приложения и не используя сторонние библиотеки.

5. Запустите тесты:

    ```
    make test
    ```

6. Послее реализации первого этапа перепишите тесты с использованием знаний, которые получили на вебинаре. Используйте принципы и инструменты о которых мы поговорили. На этом эапе можно использлвать для тестирования сторонние assert библиотеки. Попробуйте реализовать табличные тесты или тест сьюты.

5. На последнем этапе переработайте дизайн приложения так чтобы можно было использовать моки. Если вы пользуетесь кодогенерацией, сгенерируйте моки.


## Сборка проекта

Для сборки проекта выполните следующую команду:

```
make build
```

Это создаст исполняемый файл `shortener` в корневой директории проекта.

## Запуск приложения

Можно запустить приложение командой:

```
make run
```

или используя docker compose:

```
make setup
```

## Дополнительные ресурсы

- [Go testing package](https://pkg.go.dev/testing)
- [Fuzzing in Go](https://go.dev/doc/security/fuzz)
- [How to write benchmarks in Go](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
- [Testify package](https://github.com/stretchr/testify)
- [Ginkgo package](https://github.com/onsi/ginkgo)

Удачи в изучении тестирования в Go!