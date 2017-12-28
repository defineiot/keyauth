package application_test

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"

	appctr "openauth/pkg/application"
	"openauth/store/application"
	appmock "openauth/store/application/mock"
	usrmock "openauth/store/user/mock"
)

func NewAppController() *appctr.Controller {
	log := logrus.New()
	as := new(appmock.AppStore)
	as.RegistrationFn = func(userID, name, redirectURI, clientType, description, website string) (*application.Application, error) {
		client := application.Client{
			ClientID:     "aabbxx",
			ClientSecret: "aabbxx",
			RedirectURI:  redirectURI,
			ClientType:   clientType,
		}
		app := application.Application{
			ID:          "unit-test-app-id",
			UserID:      userID,
			Name:        name,
			Website:     website,
			Description: description,
			CreateAt:    time.Now().Unix(),
			Client:      &client,
		}
		return &app, nil
	}
	as.UnregistrationFn = func(id string) error {
		if id != "unit-test-app-id" {
			return errors.New("unit-test-app-id not find")
		}
		return nil
	}
	as.GetUserAppsFn = func(userID string) ([]*application.Application, error) {
		client := application.Client{
			ClientID:     "aabbxx",
			ClientSecret: "aabbxx",
			RedirectURI:  "redirect-uri",
			ClientType:   "public",
		}
		app := application.Application{
			ID:       "unit-test-id",
			UserID:   "user-exist-user",
			Name:     "app01",
			CreateAt: time.Now().Unix(),
			Client:   &client,
		}
		apps := []*application.Application{&app}
		return apps, nil
	}

	us := new(usrmock.UserStore)
	us.CheckUserIsExistByIDFn = func(userID string) (bool, error) {
		if userID == "unit-exist-user" {
			return true, nil
		}
		return false, nil
	}

	appctr.InitController(log, as, us)

	appsvr, err := appctr.GetController()
	if err != nil {
		panic(err)
	}
	return appsvr
}
