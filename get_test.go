package simpledi_test

import (
	"testing"

	"github.com/eerzho/simpledi"
)

func Test_Get_Success(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
	})
	simpledi.Resolve()

	assertNoPanic(t, func() {
		_ = simpledi.Get[string]("yeast")
	})
}

func Test_Get_ID_Not_Found(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
	})
	simpledi.Resolve()

	assertPanic(t, func() {
		_ = simpledi.Get[string]("bread")
	}, simpledi.ErrIDNotFound)
}

func Test_Get_Type_Mismatch(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
	})
	simpledi.Resolve()

	assertPanic(t, func() {
		_ = simpledi.Get[int]("yeast")
	}, simpledi.ErrTypeMismatch)
}

func Test_Get_Empty_String_ID(t *testing.T) {
	defer simpledi.Close()
	simpledi.Resolve()

	assertPanic(t, func() {
		simpledi.Get[string]("")
	}, simpledi.ErrIDRequired)
}

func Test_Get_After_Close(t *testing.T) {
	defer simpledi.Close()
	simpledi.Set(simpledi.Definition{
		ID: "yeast",
		New: func() any {
			return "yeast"
		},
	})
	simpledi.Resolve()
	simpledi.Close()

	assertPanic(t, func() {
		simpledi.Get[string]("yeast")
	}, simpledi.ErrIDNotFound)
}

func Test_Get_Same_Instance_Returned(t *testing.T) {
	defer simpledi.Close()
	type service struct {
		data string
	}
	simpledi.Set(simpledi.Definition{
		ID: "service",
		New: func() any {
			return &service{data: "some data"}
		},
	})
	simpledi.Resolve()

	first := simpledi.Get[*service]("service")
	second := simpledi.Get[*service]("service")
	assertSameInstance(t, first, second)
}

func Test_Get_During_Resolve(t *testing.T) {
	defer simpledi.Close()

	type Repository struct {
		name string
	}

	type Service struct {
		repo *Repository
	}

	type Controller struct {
		service *Service
	}

	simpledi.Set(simpledi.Definition{
		ID: "repository",
		New: func() any {
			return &Repository{name: "user_db"}
		},
	})

	simpledi.Set(simpledi.Definition{
		ID:   "service",
		Deps: []string{"repository"},
		New: func() any {
			repo := simpledi.Get[*Repository]("repository")
			return &Service{repo: repo}
		},
	})

	simpledi.Set(simpledi.Definition{
		ID:   "controller",
		Deps: []string{"service"},
		New: func() any {
			svc := simpledi.Get[*Service]("service")
			return &Controller{service: svc}
		},
	})

	simpledi.Resolve()

	assertNoPanic(t, func() {
		repo := simpledi.Get[*Repository]("repository")
		if repo.name != "user_db" {
			t.Errorf("got: %s, want: user_db", repo.name)
		}

		svc := simpledi.Get[*Service]("service")
		if svc.repo == nil {
			t.Errorf("got: nil, want: *Repository")
		}
		if svc.repo.name != "user_db" {
			t.Errorf("got: %s, want: user_db", svc.repo.name)
		}

		ctrl := simpledi.Get[*Controller]("controller")
		if ctrl.service == nil {
			t.Errorf("got: nil, want: *Service")
		}
		if ctrl.service.repo == nil {
			t.Errorf("got: nil, want: *Repository")
		}
	})
}
