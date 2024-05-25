package datamanage

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"request-matcher-openai/data-mydb/dmtable"
	"request-matcher-openai/data-mydb/mydb"
)

func (p *DataManager) GetUserList(param mydb.RequestParam, tableFilter string, tableParams []interface{}) ([]mydb.User, int, error) {
	data, count, err2 := p.DMManage.GetList("user", param, tableFilter, tableParams)
	if err2 != nil {
		return []mydb.User{}, count, err2
	}
	data1, err3 := dmtable.UserConversion{}.ConvertItemsWithType(data)
	if err3 != nil {
		return []mydb.User{}, count, err3
	}
	return data1, count, nil
}

func (p *DataManager) GetFullUserList(param mydb.RequestParam, tableFilter string, tableParams []interface{}) ([]mydb.UserShort, int, error) {
	data, count, err2 := p.DMManage.GetFullItemList("user", param, tableFilter, tableParams, "", []interface{}{})
	if err2 != nil {
		return []mydb.UserShort{}, count, err2
	}
	data1, err3 := dmtable.UserConversion{}.ConvertFullItemsWithType(data)
	if err3 != nil {
		return []mydb.UserShort{}, count, err3
	}
	return data1, count, nil
}

func (p *DataManager) GetOneUser(id string, bMaskPassword bool) (mydb.User, error) {
	item, err := p.DMManage.GetOneItem("user", id)
	if err != nil {
		return mydb.User{}, err
	}
	data2, err2 := dmtable.UserConversion{}.ConvertOneItemWithType(item)
	if err2 != nil {
		return mydb.User{}, err2
	}
	if bMaskPassword {
		data2 = p.cleanUserPassword(data2)
	}
	return data2, nil
}

func (p *DataManager) GetOneFullUser(id string) (mydb.UserShort, error) {
	item, err := p.DMManage.GetOneFullItem("user", id)
	if err != nil {
		return mydb.UserShort{}, err
	}
	return dmtable.UserConversion{}.ConvertOneFullItemWithType(item)
}

func (p *DataManager) GetOneUserByDBWithFilter(tableFilter string, tableParams []interface{}, bMaskPassword bool) (mydb.User, error) {
	item, err := p.DMManage.GetOneItemByDBWithFilter("user", tableFilter, tableParams)
	if err != nil {
		return mydb.User{}, err
	}
	data2, err2 := dmtable.UserConversion{}.ConvertOneItemWithType(item)
	if err2 != nil {
		return mydb.User{}, err2
	}
	if bMaskPassword {
		data2 = p.cleanUserPassword(data2)
	}
	return data2, nil
}

func (p *DataManager) GetOneFullUserByDBWithFilter(tableFilter string, tableParams []interface{}) (mydb.UserShort, error) {
	item, err2 := p.DMManage.GetOneFullItemByDBWithFilter("user", tableFilter, tableParams, "", []interface{}{})
	if err2 != nil {
		return mydb.UserShort{}, err2
	}
	return dmtable.UserConversion{}.ConvertOneFullItemWithType(item)
}

func (p *DataManager) CreateOneUser(req mydb.UserMessage) (mydb.User, error) {
	data := mydb.ConvertMessageToUser(req)
	item, err := p.DMManage.CreateOneItem(data)
	if err != nil {
		return mydb.User{}, err
	}
	return dmtable.UserConversion{}.ConvertOneItemWithType(item)
}

func (p *DataManager) UpdateOneUser(id string, req mydb.UserMessage) (mydb.User, error) {
	data := mydb.ConvertMessageToUser(req)
	item, err2 := p.DMManage.UpdateOneItem(data, id)
	if err2 != nil {
		return mydb.User{}, err2
	}
	return dmtable.UserConversion{}.ConvertOneItemWithType(item)
}

func (p *DataManager) DeleteOneUser(id string) error {
	return p.DMManage.DeleteOneItem("user", id)
}

func (p *DataManager) cleanUserPassword(data mydb.User) mydb.User {
	data.Password = ""
	return data
}

func (p *DataManager) CheckUserExistByEmail(email string) bool {
	tableFilter := " and email=? "
	tableParams := []interface{}{email}
	item, err3 := p.DMManage.GetOneItemByDBWithFilter("user", tableFilter, tableParams)
	if err3 != nil {
		return false
	}
	data, err2 := dmtable.UserConversion{}.ConvertOneItemWithType(item)
	if err2 != nil {
		return false
	}
	if data.Email == email {
		return true
	}
	return false
}

func (p *DataManager) CheckUserExistByPhone(phone string) bool {
	tableFilter := " and phone=? "
	tableParams := []interface{}{phone}
	item, err3 := p.DMManage.GetOneItemByDBWithFilter("user", tableFilter, tableParams)
	if err3 != nil {
		return false
	}
	data, err2 := dmtable.UserConversion{}.ConvertOneItemWithType(item)
	if err2 != nil {
		return false
	}
	if err2 == nil && data.Phone == phone {
		return true
	}
	return false
}

func (p *DataManager) UserLoginHandler(userMsg mydb.UserLoginMessage) (mydb.UserShort, error) {
	item := mydb.UserShort{}
	if userMsg.Password == "" {
		return mydb.UserShort{}, errors.New("Password cannot be empty.")
	}
	myUser, err := p.GetOneUserByDBWithFilter(" and email=? ", []interface{}{userMsg.Email}, false)
	if err != nil {
		p.mylogger.Errorf("failed to get with email:%v, err:%v", userMsg.Email, err)
		return mydb.UserShort{}, errors.New("The email entered is not registered. If you are a new user, please signup first.")
	}
	password := []byte(userMsg.Password)
	err = bcrypt.CompareHashAndPassword([]byte(myUser.Password), password)
	if err != nil {
		return mydb.UserShort{}, errors.New("Your account or password is incorrect. If you don’t remember your password, you may reset it below.")
	}
	item, err = p.GetOneFullUser(myUser.ID)
	return item, err
}

func (p *DataManager) UserLoginWithPhoneHandler(userMsg mydb.UserLoginMessage) (mydb.UserShort, error) {
	item := mydb.UserShort{}
	if userMsg.Password == "" {
		return mydb.UserShort{}, errors.New("Password cannot be empty.")
	}
	myUser, err := p.GetOneUserByDBWithFilter(" and phone=? ", []interface{}{userMsg.Phone}, false)
	if err != nil {
		p.mylogger.Errorf(" uknown user with phone:%v, err:%v", userMsg.Phone, err)
		return mydb.UserShort{}, errors.New("The mobile number entered is not registered. If you are a new user, please signup first.")
	}
	err = bcrypt.CompareHashAndPassword([]byte(myUser.Password), []byte(userMsg.Password))
	if err != nil {
		return mydb.UserShort{}, errors.New("Your account or password is incorrect. If you don’t remember your password, you may reset it below.")
	}
	item, err = p.GetOneFullUser(myUser.ID)
	return item, err
}

func (p *DataManager) UserSignupHandler(msg mydb.UserSignupMessage) (mydb.User, error) {
	var err error
	data := mydb.User{}

	if msg.Email != "" && p.CheckUserExistByEmail(msg.Email) {
		return data, errors.New(`User account already exists. If you would like to reset your password, please user "Forget Password" instead.`)
	}
	if msg.Phone != "" && p.CheckUserExistByPhone(msg.Phone) {
		return data, errors.New(`User account already exists. If you would like to reset your password, please user "Forget Password" instead.`)
	}

	data, err = p.SignupUser(msg)
	return data, err
}

func (p *DataManager) SignupUser(userMsg mydb.UserSignupMessage) (mydb.User, error) {
	if userMsg.Email != "" && p.CheckUserExistByEmail(userMsg.Email) {
		return mydb.User{}, fmt.Errorf("this email is already registered")
	} else if userMsg.Phone != "" && p.CheckUserExistByPhone(userMsg.Phone) {
		return mydb.User{}, fmt.Errorf("this mobile is already registered")
	}
	item, err := p.CreateOneUserWithoutNotifyButMessageHook(userMsg)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (p *DataManager) CreateOneUserWithoutNotifyButMessageHook(userMsg mydb.UserSignupMessage) (mydb.User, error) {
	//initial a return usershort item
	generatePassword := ""
	if userMsg.Password == "" {
		generatePassword = p.GenerateRandomCode()
		userMsg.Password = generatePassword
	}
	item, err := p.CreateOneUserWithoutNotify(userMsg) //already notify
	if err != nil {
		return item, err
	}
	//p.sendWelcomeEmailForNewUser(item, generatePassword)
	if generatePassword == "" {
		p.UpdateOneUserAnyFlagWithoutHook(item.ID, "need_change_password", false)
	}
	return item, err
}

func (p *DataManager) CreateOneUserWithoutNotify(userMsg mydb.UserSignupMessage) (mydb.User, error) {
	tmp := mydb.ConvertUserSignupMessageToUser(userMsg, loc)
	var err error
	tmp.Password, err = p.GenerateSaltedPassword(tmp.Password)
	if err != nil {
		return mydb.User{}, err
	}
	tmp.PasswordUpdateTime = mydb.GetNow()

	item, err2 := p.DMManage.CreateOneItem(tmp)
	if err2 != nil {
		return mydb.User{}, err2
	}
	p.mylogger.Infof("create user:%T, item:%v", item, item)
	data2, err3 := p.GetOneUser(item.GetID(), true)
	return data2, err3
}

func (p *DataManager) UpdateOneUserAnyFlagWithoutHook(id string, flag string, enable bool) error {
	err3 := p.DMManage.UpdateOneItemFlag("user", id, []string{flag}, []bool{enable})
	return err3
}

func (p *DataManager) UserHandleUpdatePassword(userID string, userMsg mydb.UserUpdatePasswordMessage) (mydb.UserShort, error) {
	var err error
	item := mydb.UserShort{}
	tmp, err := p.GetOneUser(userID, false)
	if err != nil {
		return item, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(tmp.Password), []byte(userMsg.OldPassword))
	if err != nil {
		p.mylogger.Errorf("the old password is incorrect, err:%v", err)
		return item, errors.New("Your old password is incorrect. Your password has not been changed. Please re-enter your password.")
	}
	err = p.updateUserPassword(userID, userMsg.NewPassword)
	if err != nil {
		return item, err
	}
	item, err = p.GetOneFullUser(userID)
	return item, err
}

func (p *DataManager) updateUserPassword(userID string, password string) error {
	_, err := p.GetOneUser(userID, true)
	if err != nil {
		return errors.New("unfound user")
	}
	err = p.saveUserNewPassword(userID, password)
	if err != nil {
		return err
	}

	return nil
}

func (p *DataManager) UserHandleSelfDeleteWithPhone(userID string, message mydb.UserSelfDeleteMessage) error {
	user, err := p.GetOneUser(userID, true)
	if err != nil {
		return errors.New("not existing user")
	}
	if user.Phone == "" {
		return errors.New("empty mobile")
	}

	_, err = p.DMManage.UpdateOneItem(mydb.User{
		ID: user.ID, Remark: "self delete",
	}, user.ID)
	if err != nil {
		return err
	}

	err = p.DeleteOneUser(user.ID)

	return err
}

func (p *DataManager) UserHandleSelfDeleteWithEmail(userID string, message mydb.UserSelfDeleteMessage) error {
	user, err := p.GetOneUser(userID, true)
	if err != nil {
		return errors.New("not existing user")
	}
	if user.Email == "" {
		return errors.New("empty email address")
	}

	_, err = p.DMManage.UpdateOneItem(mydb.User{
		ID: user.ID, Remark: "self delete",
	}, user.ID)
	if err != nil {
		return err
	}

	err = p.DeleteOneUser(user.ID)

	return err
}

func (p *DataManager) HandleUserUpdateMessage(data mydb.UserUpdateMessage) (mydb.User, error) {
	user, err := p.GetOneUser(data.UserID, false)
	if err != nil {
		return mydb.User{}, err
	}
	if data.Email != "" && data.Email != user.Email && p.CheckUserExistByEmail(data.Email) {
		return mydb.User{}, fmt.Errorf("this email is already registered")
	}

	if data.ExpiredAt == "" && user.ExpiredAt != nil {
		data.ExpiredAt = user.ExpiredAt.Format("2006-01-02 15:04:05")
	}
	tmp := mydb.ConvertUserUpdateMessageToUser(data, loc)
	tmp.Phone = ""
	tmp.Password = ""
	item, err3 := p.DMManage.UpdateOneItem(tmp, data.UserID)
	if err3 != nil {
		return mydb.User{}, err3
	}
	data2, err2 := dmtable.UserConversion{}.ConvertOneItemWithType(item)
	if err2 != nil {
		return mydb.User{}, err2
	}
	data2, _ = p.GetOneUser(data2.GetID(), true)
	return data2, nil
}

func (p *DataManager) UpdateUserPhoneWithoutHookByForce(id string, phone string) error {
	tmp, err := p.GetOneUser(id, false)
	if err != nil || tmp.ID == "" {
		return errors.New("not existing user")
	}
	if phone != "" && phone != tmp.Phone && p.CheckUserExistByPhone(phone) {
		return fmt.Errorf("this mobile is already registered")
	}
	p.mylogger.Infof("update user:%v new phone:%v old phone:%v", id, phone, tmp.Phone)
	return p.DMManage.UpdateOneItemByForce("user", id, map[string]interface{}{"phone": phone})
}

func (p *DataManager) UpdateOneUserEnableFlag(id string, enable bool) error {
	err3 := p.DMManage.UpdateOneItemFlag("user", id, []string{"enable"}, []bool{enable})
	if err3 != nil {
		return err3
	}

	return nil
}
