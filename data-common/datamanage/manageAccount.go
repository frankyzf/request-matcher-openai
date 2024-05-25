package datamanage

import (
	"errors"
	"request-matcher-openai/data-mydb/dmtable"

	"golang.org/x/crypto/bcrypt"

	"request-matcher-openai/data-mydb/mydb"
)

func (p *DataManager) GetAccountList(param mydb.RequestParam, bMaskPassword bool) ([]mydb.Account, int, error) {
	tableFilter := ""
	tableParams := []interface{}{}
	data, count, err2 := p.DMManage.GetList("account_user", param, tableFilter, tableParams)
	if err2 != nil {
		return []mydb.Account{}, count, err2
	}
	data1, err3 := dmtable.AccountConversion{}.ConvertItemsWithType(data)
	if err3 != nil {
		return []mydb.Account{}, count, err3
	}
	data2 := []mydb.Account{}
	if bMaskPassword {
		data2 = p.cleanAccountsPassword(data1)
	} else {
		data2 = data1
	}
	return data2, count, nil
}

func (p *DataManager) GetAccountListByDB(param mydb.RequestParam, bMaskPassword bool) ([]mydb.Account, int, error) {
	tableFilter := ""
	tableParams := []interface{}{}
	data, count, err2 := p.DMManage.GetListByDB("account_user", param, tableFilter, tableParams)
	if err2 != nil {
		return []mydb.Account{}, count, err2
	}
	data1, err3 := dmtable.AccountConversion{}.ConvertItemsWithType(data)
	if err3 != nil {
		return []mydb.Account{}, count, err3
	}
	data2 := []mydb.Account{}
	if bMaskPassword {
		data2 = p.cleanAccountsPassword(data1)
	} else {
		data2 = data1
	}
	return data2, count, nil
}

func (p *DataManager) GetFullAccountListByDB(param mydb.RequestParam, group string) ([]mydb.AccountShort, int, error) {
	tableFilter := ""
	tableParams := []interface{}{}
	if group != "" {
		tableFilter += " and user_group=? "
		tableParams = append(tableParams, group)
	}
	return p.GetFullAccountListWithFilter(param, tableFilter, tableParams)
}

func (p *DataManager) GetFullAccountListWithFilter(param mydb.RequestParam,
	tableFilter string, tableParams []interface{}) ([]mydb.AccountShort, int, error) {
	data, count, err2 := p.DMManage.GetFullItemListByDB("account_user", param, tableFilter, tableParams, "", []interface{}{})
	if err2 != nil {
		return []mydb.AccountShort{}, count, err2
	}
	data1, err3 := dmtable.AccountConversion{}.ConvertFullItemsWithType(data)
	if err3 != nil {
		return []mydb.AccountShort{}, count, err3
	}
	return data1, count, nil
}

func (p *DataManager) GetOneAccount(id string, bMaskPassword bool) (mydb.Account, error) {
	item, err := p.DMManage.GetOneItem("account_user", id)
	if err != nil {
		return mydb.Account{}, err
	}
	data2, err2 := dmtable.AccountConversion{}.ConvertOneItemWithType(item)
	if err2 != nil {
		return mydb.Account{}, err2
	}
	if bMaskPassword {
		data2 = p.cleanAccountPassword(data2)
	}
	return data2, nil
}

func (p *DataManager) GetOneAccountByDB(id string, bMaskPassword bool) (mydb.Account, error) {
	item, err := p.DMManage.GetOneItemByDB("account_user", id)
	if err != nil {
		return mydb.Account{}, err
	}
	data2, err2 := dmtable.AccountConversion{}.ConvertOneItemWithType(item)
	if err2 != nil {
		return mydb.Account{}, err2
	}
	if bMaskPassword {
		data2 = p.cleanAccountPassword(data2)
	}
	return data2, nil
}

func (p *DataManager) GetOneAccountByDBWithFilter(tableFilter string, tableParams []interface{}, bMaskPassword bool) (mydb.Account, error) {
	item, err := p.DMManage.GetOneItemByDBWithFilter("account_user", tableFilter, tableParams)
	if err != nil {
		return mydb.Account{}, err
	}
	data2, err2 := dmtable.AccountConversion{}.ConvertOneItemWithType(item)
	if err2 != nil {
		return mydb.Account{}, err2
	}
	if bMaskPassword {
		data2 = p.cleanAccountPassword(data2)
	}
	return data2, nil
}

func (p *DataManager) GetOneFullAccount(id string) (mydb.AccountShort, error) {
	item, err := p.DMManage.GetOneFullItem("account_user", id)
	if err != nil {
		return mydb.AccountShort{}, err
	}
	return dmtable.AccountConversion{}.ConvertOneFullItemWithType(item)
}

func (p *DataManager) GetOneFullAccountByDB(id string) (mydb.AccountShort, error) {
	item, err := p.DMManage.GetOneFullItemByDB("account_user", id)
	if err != nil {
		return mydb.AccountShort{}, err
	}
	return dmtable.AccountConversion{}.ConvertOneFullItemWithType(item)
}

func (p *DataManager) CheckAccountByEmail(email string) bool {
	item, err2 := p.DMManage.GetOneItemByDBWithFilter("account_user", ` and email=? `, []interface{}{email})
	if err2 != nil {
		return false
	}
	data, err3 := dmtable.AccountConversion{}.ConvertOneItemWithType(item)
	if err3 != nil {
		return false
	}
	if err3 == nil && data.Email == email {
		return true
	}
	return false
}

func (p *DataManager) CheckAccountAndVerifyPassword(email string, password string) (mydb.Account, error) {
	account, err2 := p.GetOneAccountByDBWithFilter(` and email=? `, []interface{}{email}, false)
	if err2 != nil {
		return mydb.Account{}, err2
	}
	saltedPassword, _ := p.GenerateSaltedPassword(password)
	if saltedPassword != account.Password {
		return mydb.Account{}, errors.New("password is not matched")
	}
	return account, nil
}

func (p *DataManager) GetOneAccountByEmail(email string, bMaskPassword bool) (mydb.Account, error) {
	item, err2 := p.DMManage.GetOneItemByDBWithFilter("account_user", ` and email=? `, []interface{}{email})
	if err2 != nil {
		return mydb.Account{}, err2
	}
	data, err3 := dmtable.AccountConversion{}.ConvertOneItemWithType(item)
	if bMaskPassword {
		data = p.cleanAccountPassword(data)
	}
	return data, err3
}

func (p *DataManager) GetOneAccountByContactEmail(email string, bMaskPassword bool) (mydb.Account, error) {
	item, err2 := p.DMManage.GetOneItemByDBWithFilter("account_user", ` and contact_email=? `, []interface{}{email})
	if err2 != nil {
		return mydb.Account{}, err2
	}
	data, err3 := dmtable.AccountConversion{}.ConvertOneItemWithType(item)
	if bMaskPassword {
		data = p.cleanAccountPassword(data)
	}
	return data, err3
}

func (p *DataManager) updateAccountPasswordByEmail(email string, password string) error {
	account := mydb.Account{}
	err := p.myDbConn.First(&account, `email=? `, email).Error
	if err != nil {
		return err
	}
	return p.accountUpdatePassword(account, password)
}

func (p *DataManager) UpdateAccountPassword(id string, password string) error {
	account := mydb.Account{}
	err := p.myDbConn.First(&account, "id=?", id).Error
	if err != nil {
		return err
	}
	return p.accountUpdatePassword(account, password)
}

func (p *DataManager) accountUpdatePassword(account mydb.Account, password string) error {
	err := p.saveAccountNewPassword(account.ID, password)
	if err == nil {
		//p.sendAccountUpdateMessage(account)
	}
	return err
}

func (p *DataManager) CreateOneAccount(msg mydb.AccountMessage) (mydb.Account, error) {
	// use internal method
	account := mydb.Account{}
	saltedPassword, err := p.GenerateSaltedPassword(msg.Password)
	if err != nil {
		return mydb.Account{}, err
	}

	account = mydb.Account{
		Name:               msg.Name,
		Password:           saltedPassword,
		PasswordUpdateTime: mydb.GetNow(),
		NeedChangePassword: true,
		Email:              msg.Email,
		Phone:              msg.Phone,
	}
	if msg.ContactEmail != nil {
		account.ContactEmail = *msg.ContactEmail
	}
	account, err = p.AddOrUpdateAccountItem(account)
	if err != nil {
		return mydb.Account{}, err
	}
	p.mylogger.Infof("create account with email:%v and create account with phone:%v and pwd:%v", account.Email, account.Phone, msg.Password)
	//p.sendAccountWelcomeMessage(account, msg.Password)
	return account, nil
}

func (p *DataManager) UpdateOneAccount(id string, msg mydb.AccountMessage) (mydb.Account, error) {
	// use internal method
	var err error

	account := mydb.Account{
		ID:    id,
		Name:  msg.Name,
		Email: msg.Email,
		Phone: msg.Phone,
	}
	if msg.Password != "" {
		account.Password, err = p.GenerateSaltedPassword(msg.Password)
		if err != nil {
			return mydb.Account{}, err
		}
		account.PasswordUpdateTime = mydb.GetNow()
	}
	account, err = p.AddOrUpdateAccountItem(account)
	if err != nil {
		return mydb.Account{}, err
	}
	if msg.ContactEmail != nil {
		account.ContactEmail = *msg.ContactEmail
		err = p.DMManage.UpdateOneItemByForce("account_user", id, map[string]interface{}{"contact_email": account.ContactEmail})
	}
	if msg.Password != "" {
		p.DMManage.UpdateOneItemFlag("account_user", account.ID, []string{"need_change_password"}, []bool{false})
		//p.sendAccountUpdateMessage(account)
	}
	return account, nil
}

func (p *DataManager) AddOrUpdateAccountItem(account mydb.Account) (mydb.Account, error) {
	// use internal method
	item, err2 := p.DMManage.AddOrUpdateOneItem(account, "id=?", []interface{}{account.ID})
	if err2 != nil {
		return mydb.Account{}, err2
	}
	data, err2 := dmtable.AccountConversion{}.ConvertOneItemWithType(item)
	if err2 != nil {
		return mydb.Account{}, err2
	}

	return p.GetOneAccountByDB(data.ID, true)
}

func (p *DataManager) DeleteOneAccount(id string) error {
	err := p.DMManage.DeleteOneItem("account_user", id)
	return err
}

func (p *DataManager) UpdateAccountEnableFlag(id string, enable bool) error {
	return p.UpdateAccountFields(id, map[string]interface{}{"enable": enable})
}

func (p *DataManager) UpdateAccountFields(id string, fields map[string]interface{}) error {
	err2 := p.DMManage.UpdateOneItemByForceWithFilter("account_user", " id=? ", []interface{}{id}, fields)
	return err2
}

func (p *DataManager) cleanAccountPassword(data mydb.Account) mydb.Account {
	data.Password = ""
	return data
}

func (p *DataManager) cleanAccountsPassword(data []mydb.Account) []mydb.Account {
	res := []mydb.Account{}
	for _, item := range data {
		res = append(res, p.cleanAccountPassword(item))
	}
	return res
}

func (p *DataManager) AccountLoginHandler(userMsg mydb.UserLoginMessage, clientIP string) (mydb.UserShort, error) {
	myAccount, err := p.AccountLogin(userMsg, clientIP)
	if err != nil {
		return mydb.UserShort{}, err
	}
	item := mydb.ConvertBaseAccountToUserShort(mydb.ConvertAccountToBaseAccount(myAccount))
	return item, nil
}

func (p *DataManager) AccountLogin(userMsg mydb.UserLoginMessage, clientIP string) (mydb.Account, error) {
	var err error
	if userMsg.Email == "" || userMsg.Password == "" {
		return mydb.Account{}, errors.New("incorrect email or password")
	}
	password := []byte(userMsg.Password)
	email := userMsg.Email
	if userMsg.ForceSkipCaptureVerify == false { // if  true - skip captcha verify, case: tcp login
		err = p.myAuth.VerifyCaptcha(clientIP, userMsg.Captcha)
		if err != nil {
			return mydb.Account{}, err
		}
	}

	myAccount, err1 := p.GetOneAccountByContactEmail(email, false)
	if err1 != nil {
		myAccount, err1 = p.GetOneAccountByEmail(email, false)
		if err1 != nil {
			return mydb.Account{}, errors.New("email is not registered")
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(myAccount.Password), password)
	if err != nil {
		return mydb.Account{}, errors.New("incorrect email or password ")
	}
	return myAccount, err
}
