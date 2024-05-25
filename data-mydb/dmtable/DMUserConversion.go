package dmtable

import (
	"encoding/json"
	"errors"
	"request-matcher-openai/data-mydb/mydb"
)

type UserConversion struct {
}

func (p UserConversion) ParseOneItem(buf string) (mydb.DataItem, error) {
	data := mydb.User{}
	err := json.Unmarshal([]byte(buf), &data)
	return data, err
}

func (p UserConversion) MarshalOneItem(item interface{}) (string, error) {
	if data, ok := item.(mydb.User); ok {
		buf, err := json.Marshal(data)
		return string(buf), err
	} else {
		return "", errors.New("failed to convert to User")
	}
}

func (p UserConversion) ParseOneFullItem(buf string) (mydb.DataItem, error) {
	data := mydb.UserShort{}
	err := json.Unmarshal([]byte(buf), &data)
	return data, err
}

func (p UserConversion) MarshalOneFullItem(item interface{}) (string, error) {
	if data, ok := item.(mydb.UserShort); ok {
		buf, err := json.Marshal(data)
		return string(buf), err
	} else {
		return "", errors.New("failed to convert to UserShort")
	}
}

func (p UserConversion) ParseItems(buf string) ([]mydb.DataItem, error) {
	res := []mydb.DataItem{}
	data := []mydb.User{}
	err := json.Unmarshal([]byte(buf), &data)
	if err != nil {
		return res, err
	}
	for _, item := range data {
		res = append(res, item)
	}
	return res, nil
}

func (p UserConversion) MarshalItems(data []mydb.DataItem) (string, error) {
	res, err := p.ConvertItemsWithType(data)
	if err != nil {
		return "", err
	}
	buf, err := json.Marshal(res)
	return string(buf), err
}

func (p UserConversion) ParseFullItems(buf string) ([]mydb.DataItem, error) {
	res := []mydb.DataItem{}
	data := []mydb.UserShort{}
	err := json.Unmarshal([]byte(buf), &data)
	if err != nil {
		return res, err
	}
	for _, item := range data {
		res = append(res, item)
	}
	return res, nil
}

func (p UserConversion) MarshalFullItems(data []mydb.DataItem) (string, error) {
	res, err := p.ConvertFullItemsWithType(data)
	if err != nil {
		return "", err
	}
	buf, err := json.Marshal(res)
	return string(buf), err
}

func (p UserConversion) ConvertItems(data interface{}) ([]mydb.DataItem, error) {
	if res, ok := data.([]mydb.User); ok {
		res2 := []mydb.DataItem{}
		for _, item := range res {
			res2 = append(res2, item)
		}
		return res2, nil
	}
	return []mydb.DataItem{}, errors.New("failed to convert to User")
}

func (p UserConversion) ConvertFullItems(data interface{}) ([]mydb.DataItem, error) {
	if res, ok := data.([]mydb.UserShort); ok {
		res2 := []mydb.DataItem{}
		for _, item := range res {
			res2 = append(res2, item)
		}
		return res2, nil
	}
	return []mydb.DataItem{}, errors.New("failed to convert to UserShort")
}

func (p UserConversion) ConvertOneItemWithType(data interface{}) (mydb.User, error) {
	res2, ok := data.(mydb.User)
	if ok {
		return res2, nil
	}
	return mydb.User{}, errors.New("failed to convert to User")
}

func (p UserConversion) ConvertOneFullItemWithType(data interface{}) (mydb.UserShort, error) {
	res2, ok := data.(mydb.UserShort)
	if ok {
		return res2, nil
	}
	return mydb.UserShort{}, errors.New("failed to convert to UserShort")
}

func (p UserConversion) ConvertItemsWithType(data interface{}) ([]mydb.User, error) {
	if res, ok := data.([]mydb.DataItem); ok {
		res2 := []mydb.User{}
		for _, item := range res {
			item2, ok2 := item.(mydb.User)
			if ok2 {
				res2 = append(res2, item2)
			} else {
				return []mydb.User{}, errors.New("failed to convert to User")
			}
		}
		return res2, nil
	}
	return []mydb.User{}, errors.New("failed to convert to User")
}

func (p UserConversion) ConvertFullItemsWithType(data interface{}) ([]mydb.UserShort, error) {
	if res, ok := data.([]mydb.DataItem); ok {
		res2 := []mydb.UserShort{}
		for _, item := range res {
			item2, ok2 := item.(mydb.UserShort)
			if ok2 {
				res2 = append(res2, item2)
			} else {
				return []mydb.UserShort{}, errors.New("failed to convert to UserShort")
			}
		}
		return res2, nil
	}
	return []mydb.UserShort{}, errors.New("failed to convert to UserShort")
}
