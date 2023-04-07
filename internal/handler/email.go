package handler

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/umalmyha/go-neo4j/internal/model"
	"github.com/umalmyha/go-neo4j/internal/storage"
)

type EmailHandler struct {
	str *storage.Neo4jEmailStorage
}

func NewEmailHandler(str *storage.Neo4jEmailStorage) *EmailHandler {
	return &EmailHandler{str: str}
}

func (h *EmailHandler) CreateUser(c echo.Context) error {
	var u model.User
	if err := c.Bind(&u); err != nil {
		return err
	}

	u.ID = uuid.NewString()
	if err := h.str.CreateUser(c.Request().Context(), &u); err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(http.StatusCreated, &u)
}

func (h *EmailHandler) FindByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "path parameter id must be provided")
	}

	u, err := h.str.FindUserByID(c.Request().Context(), id)
	if err != nil {
		log.Println(err)
		return err
	}

	if u == nil {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, u)
}

func (h *EmailHandler) CreateEmail(c echo.Context) error {
	var email model.Email
	if err := c.Bind(&email); err != nil {
		return err
	}

	if email.To == "" || email.From == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "both email sender and receiver must be specified")
	}

	if err := h.str.CreateEmail(c.Request().Context(), &email); err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(http.StatusCreated, &email)
}
