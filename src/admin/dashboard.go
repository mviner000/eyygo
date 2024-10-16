package admin

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/src/auth"
	"github.com/mviner000/eyymi/src/http"
)

const SiteName = "Eyygo Administration"

func Dashboard(c *fiber.Ctx) error {
	userID, _, err := auth.GetSessionFromDB(c)
	if err != nil {
		return http.HttpResponseRedirect("/login", false).Render(c)
	}
	uintUserID := uint(userID)

	user, err := auth.GetUserByID(uintUserID)
	if err != nil {
		return http.HttpResponseServerError("Error retrieving user information", nil).Render(c)
	}

	log.Printf("User data: %+v", user)

	return http.HttpResponseHTMX(fiber.Map{
		"User":      user,
		"MetaTitle": "Dashboard | " + SiteName,
	}, "src/admin/templates/dashboard.html", "src/admin/templates/layout.html").Render(c)
}
