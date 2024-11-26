package get

import (
	"log/slog"
	"net/http"

	resp "github.com/AlexMickh/PFly2/internal/lib/api/response"
	"github.com/AlexMickh/PFly2/internal/models"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UserGetter interface {
	GetUserByEmail(email string) (models.User, error)
}

type Request struct {
	Email string `json:"email" validate:"required,email"`
}

type Response struct {
	resp.Response
	Id          int      `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Email       string   `json:"email,omitempty"`
	Password    []byte   `json:"password,omitempty"`
	ImageUrl    string   `json:"image_url,omitempty"`
	Description string   `json:"description,omitempty"`
	Interests   []string `json:"interests,omitempty"`
}

func New(log *slog.Logger, userGettet UserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		user, err := userGettet.GetUserByEmail(req.Email)
		if err != nil {
			log.Error("failed to get user by email", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to get user by email"))

			return
		}

		render.JSON(w, r, &user)
	}
}
