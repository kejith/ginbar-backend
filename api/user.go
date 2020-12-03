package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"ginbar/api/models"
	"ginbar/api/utils"
	"ginbar/mysql/db"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
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

// GetUsers retrieves a user with a given id or username from the database
func (server *Server) GetUsers(context *gin.Context) {
	// TODO: Abstract Database Result Struct into a custom JSON Model
	user, err := server.store.GetUsers(context)
	if err != nil {
		context.Error(err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   user,
	})

}

// GetUser retrieves a user with a given id or username from the database
func (server *Server) GetUser(context *gin.Context) {
	userName := context.Param("id")
	userID, _ := strconv.ParseInt(context.Param("id"), 10, 64)

	// If UserID != 0 we use this to get the User
	// Else we try to use the Username instead
	if userID != 0 {
		// TODO: Abstract Database Result Struct into a custom JSON Model
		user, err := server.store.GetUser(context, int32(userID))
		if err != nil {
			context.Error(err)
			context.Status(http.StatusInternalServerError)
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   user,
		})
	} else {
		// TODO: Abstract Database Result Struct into a custom JSON Model
		user, err := server.store.GetUserByName(context, userName)
		if err != nil {
			context.Error(err)
			context.Status(http.StatusInternalServerError)
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   user,
		})
	}

}

// GetUserByName retrieves a user with a given id from the database
func (server *Server) GetUserByName(context *gin.Context) {
	// TODO GetUserByName Should get a User by Name not ID
	userID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.Error(err)
		context.Status(http.StatusInternalServerError)
		return
	}

	user, err := server.store.GetUser(context, int32(userID))
	context.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   user,
	})
}

// CreateUser saves a user to the database
func (server *Server) CreateUser(context *gin.Context) {
	// Read POST Data and write it into registerForm
	var registerForm userRegisterForm
	err := context.ShouldBind(&registerForm)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	// Hash Password
	hash, err := utils.CreatePasswordHash([]byte(registerForm.Password))
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	createParams := db.CreateUserParams{
		Name:     registerForm.Name,
		Email:    registerForm.Email,
		Password: hash,
	}

	if len(createParams.Name) < 4 || !utils.IsEmailValid(createParams.Email) {
		context.Status(http.StatusUnprocessableEntity)
		context.Error(errors.New("Username or Email not valid"))
		return
	}

	// Insert User
	err = server.store.CreateUser(context, createParams)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	context.Status(http.StatusOK)
}

// Login handles a request from the client to login a user and create
// a jwt token which is set in the Cookie
func (server *Server) Login(context *gin.Context) {
	// Read POST data and write it into userLoginForm
	var form userLoginForm
	err := context.ShouldBind(&form)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	// find user with Given Name
	user, err := server.store.GetUserByName(context, form.Name)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	userJSON := models.PublicUserJSON{}
	userJSON.Populate(user)

	// Compare Passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	// Compare given Username with Username from Database
	if user.Name != form.Name {
		context.Status(http.StatusInternalServerError)
		return
	}

	// Write Session Information
	session := sessions.Default(context)
	session.Set("user", user.Name)
	session.Set("userid", user.ID)
	if err = session.Save(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	context.JSON(http.StatusOK, userJSON)
}

// Logout sets the current MaxAge of the Session to -1 so it expires
// and the user is logged out
func (server *Server) UserLogout(context *gin.Context) {
	fmt.Println("Logout")
	session := sessions.Default(context)
	session.Set("user", "")
	session.Set("userid", 0)
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	session.Save()
}

// Me returns the Data of the currently logged in User
func (server *Server) Me(context *gin.Context) {
	// Read Session Information
	session := sessions.Default(context)
	userName, ok := session.Get("user").(string)

	if ok {
		user, err := server.store.GetUserByName(context, userName)
		if err != nil {
			context.Error(err)
			context.Status(http.StatusInternalServerError)
			return
		}

		// TODO: maybe send more informative Information
		// remove sensible Data
		u := models.PublicUserJSON{}
		u.Populate(user)

		// context.JSON(http.StatusOK, gin.H{
		// 	"status": http.StatusOK,
		// 	"data":   user,
		// })

		context.JSON(http.StatusOK, u)
	} else {
		context.Status(http.StatusUnauthorized)
	}
}
