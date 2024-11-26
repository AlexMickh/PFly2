package save

import (
	"log/slog"
	"net/http"

	resp "github.com/AlexMickh/PFly2/internal/lib/api/response"
	"github.com/AlexMickh/PFly2/internal/models"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserSaver interface {
	SaveUser(user models.User) (int64, error)
}

type Request struct {
	Name        string   `json:"name" validate:"required"`
	Email       string   `json:"email" validate:"required,email"`
	Password    string   `json:"password" validate:"required,min=8"`
	ImageUrl    string   `json:"image_url" validate:"url"`
	Description string   `json:"description"`
	Interests   []string `json:"interests"`
}

type Response struct {
	resp.Response
	Id int64 `json:"id,omitempty"`
}

func New(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"

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

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to hash password", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to hash password"))

			return
		}

		userData := models.User{
			Id:          -1,
			Name:        req.Name,
			Email:       req.Email,
			Password:    hash,
			ImageUrl:    req.ImageUrl,
			Description: req.Description,
			Interests:   req.Interests,
		}

		// TODO: add validation for case when user already exists
		id, err := userSaver.SaveUser(userData)
		if err != nil {
			log.Error("failed to save user", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to save user"))

			return
		}

		log.Info("user saved", slog.Int64("id", id))

		responseOK(w, r, id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id int64) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Id:       id,
	})
}
