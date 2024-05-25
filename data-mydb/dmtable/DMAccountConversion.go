package dmtable

import (
	"encoding/json"
	"errors"
	"request-matcher-openai/data-mydb/mydb"
)

type AccountConversion struct {
}

func (p AccountConversion) ParseOneItem(buf string) (mydb.DataItem, error) {
	data := mydb.Account{}
	err := json.Unmarshal([]byte(buf), &data)
	return data, err
}

func (p AccountConversion) MarshalOneItem(item interface{}) (string, error) {
	if data, ok := item.(mydb.Account); ok {
		buf, err := json.Marshal(data)
		return string(buf), err
	} else {
		return "", errors.New("failed to convert to Account")
	}
}

func (p AccountConversion) ParseOneFullItem(buf string) (mydb.DataItem, error) {
	data := mydb.AccountShort{}
	err := json.Unmarshal([]byte(buf), &data)
	return data, err
}

func (p AccountConversion) MarshalOneFullItem(item interface{}) (string, error) {
	if data, ok := item.(mydb.AccountShort); ok {
		buf, err := json.Marshal(data)
		return string(buf), err
	} else {
		return "", errors.New("failed to convert to AccountShort")
	}
}

func (p AccountConversion) ParseItems(buf string) ([]mydb.DataItem, error) {
	res := []mydb.DataItem{}
	data := []mydb.Account{}
	err := json.Unmarshal([]byte(buf), &data)
	if err != nil {
		return res, err
	}
	for _, item := range data {
		res = append(res, item)
	}
	return res, nil
}

func (p AccountConversion) MarshalItems(data []mydb.DataItem) (string, error) {
	res, err := p.ConvertItemsWithType(data)
	if err != nil {
		return "", err
	}
	buf, err := json.Marshal(res)
	return string(buf), err
}

func (p AccountConversion) ParseFullItems(buf string) ([]mydb.DataItem, error) {
	res := []mydb.DataItem{}
	data := []mydb.AccountShort{}
	err := json.Unmarshal([]byte(buf), &data)
	if err != nil {
		return res, err
	}
	for _, item := range data {
		res = append(res, item)
	}
	return res, nil
}

func (p AccountConversion) MarshalFullItems(data []mydb.DataItem) (string, error) {
	res, err := p.ConvertFullItemsWithType(data)
	if err != nil {
		return "", err
	}
	buf, err := json.Marshal(res)
	return string(buf), err
}

func (p AccountConversion) ConvertItems(data interface{}) ([]mydb.DataItem, error) {
	if res, ok := data.([]mydb.Account); ok {
		res2 := []mydb.DataItem{}
		for _, item := range res {
			res2 = append(res2, item)
		}
		return res2, nil
	}
	return []mydb.DataItem{}, errors.New("failed to convert to Account")
}

func (p AccountConversion) ConvertFullItems(data interface{}) ([]mydb.DataItem, error) {
	if res, ok := data.([]mydb.AccountShort); ok {
		res2 := []mydb.DataItem{}
		for _, item := range res {
			res2 = append(res2, item)
		}
		return res2, nil
	}
	return []mydb.DataItem{}, errors.New("failed to convert to AccountShort")
}

func (p AccountConversion) ConvertOneItemWithType(data interface{}) (mydb.Account, error) {
	res2, ok := data.(mydb.Account)
	if ok {
		return res2, nil
	}
	return mydb.Account{}, errors.New("failed to convert to Account")
}

func (p AccountConversion) ConvertOneFullItemWithType(data interface{}) (mydb.AccountShort, error) {
	res2, ok := data.(mydb.AccountShort)
	if ok {
		return res2, nil
	}
	return mydb.AccountShort{}, errors.New("failed to convert to AccountShort")
}

func (p AccountConversion) ConvertItemsWithType(data interface{}) ([]mydb.Account, error) {
	if res, ok := data.([]mydb.DataItem); ok {
		res2 := []mydb.Account{}
		for _, item := range res {
			item2, ok2 := item.(mydb.Account)
			if ok2 {
				res2 = append(res2, item2)
			} else {
				return []mydb.Account{}, errors.New("failed to convert to Account")
			}
		}
		return res2, nil
	}
	return []mydb.Account{}, errors.New("failed to convert to Account")
}

func (p AccountConversion) ConvertFullItemsWithType(data interface{}) ([]mydb.AccountShort, error) {
	if res, ok := data.([]mydb.DataItem); ok {
		res2 := []mydb.AccountShort{}
		for _, item := range res {
			item2, ok2 := item.(mydb.AccountShort)
			if ok2 {
				res2 = append(res2, item2)
			} else {
				return []mydb.AccountShort{}, errors.New("failed to convert to AccountShort")
			}
		}
		return res2, nil
	}
	return []mydb.AccountShort{}, errors.New("failed to convert to AccountShort")
}
