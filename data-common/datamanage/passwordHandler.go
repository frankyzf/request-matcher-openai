package datamanage

import (
	"errors"
	"request-matcher-openai/data-mydb/mydb"

	"golang.org/x/crypto/bcrypt"
)

func (p *DataManager) GenerateSaltedPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("empty password")
	}
	if len(password) > 255 {
		return "", errors.New("password length overflow") //bsd bcrypt bug
	}
	salted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(salted), err
}

func (p *DataManager) saveUserNewPassword(id string, newPassword string) error {
	salted, err := p.GenerateSaltedPassword(newPassword)
	if err != nil {
		return err
	}
	_, err = p.DMManage.UpdateOneItem(mydb.User{
		ID:                 id,
		Password:           string(salted),
		PasswordUpdateTime: mydb.GetNow(),
	}, id)
	if err != nil {
		return err
	}
	//here, we set need_change_password to false because it only applies to update
	p.DMManage.UpdateOneItemFlag("user", id, []string{"need_change_password"}, []bool{false})
	return nil
}

func (p *DataManager) saveAccountNewPassword(id string, newPassword string) error {
	salted, err := p.GenerateSaltedPassword(newPassword)
	if err != nil {
		return err
	}
	_, err = p.DMManage.UpdateOneItem(mydb.Account{
		ID:                 id,
		Password:           string(salted),
		PasswordUpdateTime: mydb.GetNow(),
	}, id)
	if err != nil {
		return err
	}
	//here, we set need_change_password to false because it only applies to update
	p.DMManage.UpdateOneItemFlag("account_user", id, []string{"need_change_password"}, []bool{false})
	return nil
}
