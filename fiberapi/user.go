package fiberapi

import (
	"errors"
	"fmt"
	"ginbar/api/models"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// --------------------
// FORMS
// --------------------

type userRegisterForm struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

type userLoginForm struct {
	Name     string `form:"name"`
	Password string `form:"password"`
}

func (server *FiberServer) GetUsers(c *fiber.Ctx) error {
	userList, err := server.store.GetUsers(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": http.StatusOK,
		"data":   userList,
	})
}

func (server *FiberServer) GetUser(c *fiber.Ctx) error {

	userID, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	// Todo retrieve User by Username
	userName := ""

	if userID != 0 {
		// TODO: Abstract Database Result Struct into a custom JSON Model
		user, err := server.store.GetUser(c.Context(), int32(userID))
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"status": http.StatusOK,
			"data":   user,
		})
	}

	if userName != "" {
		// TODO: Abstract Database Result Struct into a custom JSON Model
		user, err := server.store.GetUserByName(c.Context(), userName)
		if err != nil {
			return nil
		}

		return c.JSON(fiber.Map{
			"status": http.StatusOK,
			"data":   user,
		})
	}

	return fmt.Errorf(
		"input error: Could not retrieve user with given Information UserID: %v, UserName: %s",
		userID,
		userName,
	)
}

func (server *FiberServer) Login(c *fiber.Ctx) error {
	session, err := server.Sessions.Get(c)
	if err != nil {
		return nil
	}

	// Parse HTML Queries
	loginForm := &userLoginForm{}
	if err := c.BodyParser(loginForm); err != nil {
		return err
	}

	user, err := server.store.GetUserByName(c.Context(), loginForm.Name)
	if err != nil {
		return err
	}

	userJSON := models.PublicUserJSON{}
	userJSON.Populate(user)

	// Compare Passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password))
	if err != nil {
		return err
	}

	// Compare given Username with Username from Database
	if user.Name != loginForm.Name {
		return errors.New("login: username doesn't Match")
	}

	session.Set("user", user.Name)
	session.Set("userID", user.ID)
	session.Set("userlevel", user.Level)

	return session.Save()
}

func (server *FiberServer) Logout(c *fiber.Ctx) error {
	session, err := server.Sessions.Get(c)
	if err != nil {
		return nil
	}

	// Todo: look into probably redundant code
	session.Set("user", "")     // this
	session.Set("userid", 0)    // this
	session.Set("userlevel", 0) // this

	return session.Destroy()
}

func (server *FiberServer) Me(c *fiber.Ctx) error {
	session, err := server.Sessions.Get(c)
	if err != nil {
		return err
	}

	userName, ok := session.Get("user").(string)
	if ok {
		user, err := server.store.GetUserByName(c.Context(), userName)
		if err != nil {
			return err
		}

		u := models.PublicUserJSON{}
		u.Populate(user)

		return c.JSON(u)
	} else {
		return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
	}
}
