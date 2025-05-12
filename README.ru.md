# SimpleDI

[English](README.md) | [Русский](README.ru.md)

SimpleDI - это легковесный контейнер внедрения зависимостей для приложений на Go. Он предоставляет простой способ управления зависимостями и их жизненным циклом в вашем приложении.

### Возможности

- Простой и интуитивно понятный API
- Разрешение зависимостей с автоматическим упорядочиванием
- Обнаружение циклических зависимостей
- Типобезопасное внедрение зависимостей
- Отсутствие внешних зависимостей

### Установка

```bash
go get github.com/eerzho/simpledi@latest
```

### Быстрый старт

```go
package main

import "github.com/eerzho/simpledi"

func main() {
    c := simpledi.NewContainer()

    // Регистрация зависимостей
    c.Register("db", nil, func() any {
        return &DB{DSN: "example"}
    })

    c.Register("repo", []string{"db"}, func() any {
        db := c.Get("db").(*DB)
        return &Repo{DB: db}
    })

    c.Register("service", []string{"repo"}, func() any {
        repo := c.Get("repo").(*Repo)
        return &Service{Repo: repo}
    })

    // Разрешение всех зависимостей
    if err := c.Resolve(); err != nil {
        panic(err)
    }
}
```

### Справочник API

#### NewContainer()

Создает новый контейнер внедрения зависимостей.

#### Register(name string, deps []string, constructor func() any)

Регистрирует новую зависимость в контейнере.
- `name`: Уникальный идентификатор зависимости
- `deps`: Список имен зависимостей, от которых зависит этот компонент
- `constructor`: Функция, создающая экземпляр зависимости

#### Get(name string) any

Получает разрешенный экземпляр зависимости по его имени.

#### Resolve() error

Разрешает все зарегистрированные зависимости в правильном порядке. Возвращает ошибку при наличии циклических или отсутствующих зависимостей.

### Лицензия

MIT License - подробности в файле [LICENSE](LICENSE)
