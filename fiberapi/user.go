package fiberapi

import (
	"errors"
	"fmt"
	"ginbar/api/models"
	"ginbar/api/utils"
	"ginbar/mysql/db"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// --------------------
// Structs
// --------------------

type SessionUser struct {
	Name  string
	ID    int32
	Level int32
}

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
	session.Set("userID", fmt.Sprintf("%v", user.ID))
	session.Set("userlevel", fmt.Sprintf("%v", user.Level))
	err = session.Save()

	if err != nil {
		return err
	}

	return c.JSON(userJSON)
}

func (server *FiberServer) Logout(c *fiber.Ctx) error {
	session, err := server.Sessions.Get(c)
	if err != nil {
		return nil
	}

	// Todo: look into probably redundant code
	session.Set("user", "")     // this
	session.Set("userID", 0)    // this
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

func (server *FiberServer) CreateUser(c *fiber.Ctx) error {
	// Parse HTML Queries
	form := new(userRegisterForm)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	// Get Session Information
	user, err := server.GetUserFromSession(c)
	if err != nil {
		return err
	}

	if user.ID <= 0 {
		return fmt.Errorf("upload post: user data could not be loaded from session [userid]")
	}

	// Hash Password
	hash, err := utils.CreatePasswordHash([]byte(form.Password))
	if err != nil {
		return err
	}

	createParams := db.CreateUserParams{
		Name:     form.Name,
		Email:    form.Email,
		Password: hash,
	}

	if len(createParams.Name) < 4 || !utils.IsEmailValid(createParams.Email) {
		return errors.New("username or Email not valid")
	}

	// Insert User
	err = server.store.CreateUser(c.Context(), createParams)
	if err != nil {
		return err
	}

	return nil
}

func (server *FiberServer) GetUserFromSession(c *fiber.Ctx) (*SessionUser, error) {
	session, err := server.Sessions.Get(c)
	if err != nil {
		return nil, err
	}

	strUserName, _ := session.Get("user").(string)
	strUserID, _ := session.Get("userID").(string)
	strLevel, _ := session.Get("userlevel").(string)

	userID, _ := strconv.ParseInt(strUserID, 10, 64)
	userLevel, _ := strconv.ParseInt(strLevel, 10, 64)

	user := &SessionUser{
		Name:  strUserName,
		ID:    int32(userID),
		Level: int32(userLevel),
	}

	return user, nil
}
