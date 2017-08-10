package api

import (
	"database/sql"
)

type AppServicer interface {
}

type AppService struct {
	DB *sql.DB
}

// func (app *AppService) VerifyUser(name, pass string) (interface{}, error) {
// 	user := models.User{}
// 	err := user.UserByName(app.DB, name)
// 	if err != nil {
// 		return nil, errors.New(err.Error() + "Invalid username")
// 	}

// 	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(pass)); err != nil {
// 		return nil, errors.New(err.Error() + "Invalid Password")
// 	}
// 	return user, nil
// }
