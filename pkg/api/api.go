package api

import (
	"encoding/json"
	"github.com/ademaxweb/mfa-go-core/pkg/handler"
	"io"
	"net/http"
	"users/pkg/db"
)

type Api struct {
	writer io.Writer
	db     db.Interface
}

func New(writer io.Writer, db db.Interface) *Api {
	return &Api{
		writer: writer,
		db:     db,
	}
}

func SendResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func (a *Api) write(s string) {
	if a.writer == nil {
		return
	}
	a.writer.Write([]byte(s))
}

func (a *Api) RegisterRoutes(h *handler.Handler) {
	h.BaseMiddleware(a.logMiddleware, a.contentTypeMiddleware)

	h.Handle(
		handler.Route{
			Path:    "/users",
			Func:    a.getUsers,
			Methods: handler.Methods{http.MethodGet},
		},
		handler.Route{
			Path:    "/users",
			Func:    a.createUser,
			Methods: handler.Methods{http.MethodPost},
		},
		handler.Route{
			Path:    "/users/{id:[0-9]+}",
			Func:    a.getUser,
			Methods: handler.Methods{http.MethodGet},
		},
		handler.Route{
			Path:    "/users/{id:[0-9]+}",
			Func:    a.deleteUser,
			Methods: handler.Methods{http.MethodDelete},
		},
		handler.Route{
			Path:    "/users/{id:[0-9]+}",
			Func:    a.updateUser,
			Methods: handler.Methods{http.MethodPut},
		},
		handler.Route{
			Path:    "/users/{email:[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}}",
			Func:    a.getUserByEmail,
			Methods: handler.Methods{http.MethodGet},
		},
	)
}
