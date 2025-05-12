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
package simpledi

import "github.com/eerzho/simpledi"

func main() {
	c := simpledi.NewContainer()

	// Register dependencies
	c.Register("db", nil, func() any {
		return &DB{DSN: "example"}
	})

	c.Register("repo1", []string{"db"}, func() any {
		return &Repo1{
			DB: c.Get("db").(*DB),
		}
	})

	c.Register("repo2", []string{"db"}, func() any {
		return &Repo2{
			DB: c.Get("db").(*DB),
		}
	})

	c.Register("service", []string{"repo1", "repo2"}, func() any {
		return &Service{
			Repo1: c.Get("repo1").(*Repo1),
			Repo2: c.Get("repo2").(*Repo2),
		}
	})

	// Resolve all dependencies
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
