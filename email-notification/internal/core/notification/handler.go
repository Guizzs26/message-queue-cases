package notification

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type NotificationHTTPHandler struct {
	svc EmailService
}

func NewNotificationHTTPHandler(s EmailService) *NotificationHTTPHandler {
	return &NotificationHTTPHandler{svc: s}
}

type sendEmailRequest struct {
	To      string `json:"to" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

func (nh *NotificationHTTPHandler) SendEmailHandler(c echo.Context) error {
	var req sendEmailRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	if err := nh.svc.SendEmail(
		c.Request().Context(),
		req.To,
		req.Subject,
		req.Body,
	); err != nil {
		log.Printf("ERROR: failed to send email: %+v\n", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not process request"})
	}

	return c.JSON(http.StatusAccepted, echo.Map{"message": "request accepted and is being processed"})
}
