package oauth2

import (
	"fmt"
	"net/http"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Registers a new user
func register(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB) {
	user := User{}
	if err := r.DecodeJsonPayload(&user); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Case insensitive search, usernames will probably be emails and
	// foo@bar.com is identical to FOO@BAR.com
	if db.Where("LOWER(username) = LOWER(?)", user.Username).First(&User{}).RowsAffected > 0 {
		rest.Error(w, fmt.Sprintf("%s already taken", user.Username), http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 3)
	if err != nil {
		rest.Error(w, "Bcrypt error", http.StatusInternalServerError)
		return
	}

	user.Password = string(passwordHash)
	if err := db.Create(&user).Error; err != nil {
		rest.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteJson(map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	})
}