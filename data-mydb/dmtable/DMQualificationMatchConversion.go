package dmtable

import (
	"encoding/json"
	"errors"
	"request-matcher-openai/data-mydb/mydb"
)

type QualificationMatchConversion struct {
}

func (p QualificationMatchConversion) ParseOneItem(buf string) (mydb.DataItem, error) {
	data := mydb.QualificationMatch{}
	err := json.Unmarshal([]byte(buf), &data)
	return data, err
}

func (p QualificationMatchConversion) MarshalOneItem(item interface{}) (string, error) {
	if data, ok := item.(mydb.QualificationMatch); ok {
		buf, err := json.Marshal(data)
		return string(buf), err
	} else {
		return "", errors.New("failed to convert to QualificationMatch")
	}
}

func (p QualificationMatchConversion) ParseOneFullItem(buf string) (mydb.DataItem, error) {
	data := mydb.QualificationMatchShort{}
	err := json.Unmarshal([]byte(buf), &data)
	return data, err
}

func (p QualificationMatchConversion) MarshalOneFullItem(item interface{}) (string, error) {
	if data, ok := item.(mydb.QualificationMatchShort); ok {
		buf, err := json.Marshal(data)
		return string(buf), err
	} else {
		return "", errors.New("failed to convert to QualificationMatchShort")
	}
}

func (p QualificationMatchConversion) ParseItems(buf string) ([]mydb.DataItem, error) {
	res := []mydb.DataItem{}
	data := []mydb.QualificationMatch{}
	err := json.Unmarshal([]byte(buf), &data)
	if err != nil {
		return res, err
	}
	for _, item := range data {
		res = append(res, item)
	}
	return res, nil
}

func (p QualificationMatchConversion) MarshalItems(data []mydb.DataItem) (string, error) {
	res, err := p.ConvertItemsWithType(data)
	if err != nil {
		return "", err
	}
	buf, err := json.Marshal(res)
	return string(buf), err
}

func (p QualificationMatchConversion) ParseFullItems(buf string) ([]mydb.DataItem, error) {
	res := []mydb.DataItem{}
	data := []mydb.QualificationMatchShort{}
	err := json.Unmarshal([]byte(buf), &data)
	if err != nil {
		return res, err
	}
	for _, item := range data {
		res = append(res, item)
	}
	return res, nil
}

func (p QualificationMatchConversion) MarshalFullItems(data []mydb.DataItem) (string, error) {
	res, err := p.ConvertFullItemsWithType(data)
	if err != nil {
		return "", err
	}
	buf, err := json.Marshal(res)
	return string(buf), err
}

func (p QualificationMatchConversion) ConvertItems(data interface{}) ([]mydb.DataItem, error) {
	if res, ok := data.([]mydb.QualificationMatch); ok {
		res2 := []mydb.DataItem{}
		for _, item := range res {
			res2 = append(res2, item)
		}
		return res2, nil
	}
	return []mydb.DataItem{}, errors.New("failed to convert to QualificationMatch")
}

func (p QualificationMatchConversion) ConvertFullItems(data interface{}) ([]mydb.DataItem, error) {
	if res, ok := data.([]mydb.QualificationMatchShort); ok {
		res2 := []mydb.DataItem{}
		for _, item := range res {
			res2 = append(res2, item)
		}
		return res2, nil
	}
	return []mydb.DataItem{}, errors.New("failed to convert to QualificationMatchShort")
}

func (p QualificationMatchConversion) ConvertOneItemWithType(data interface{}) (mydb.QualificationMatch, error) {
	res2, ok := data.(mydb.QualificationMatch)
	if ok {
		return res2, nil
	}
	return mydb.QualificationMatch{}, errors.New("failed to convert to QualificationMatch")
}

func (p QualificationMatchConversion) ConvertOneFullItemWithType(data interface{}) (mydb.QualificationMatchShort, error) {
	res2, ok := data.(mydb.QualificationMatchShort)
	if ok {
		return res2, nil
	}
	return mydb.QualificationMatchShort{}, errors.New("failed to convert to QualificationMatchShort")
}

func (p QualificationMatchConversion) ConvertItemsWithType(data interface{}) ([]mydb.QualificationMatch, error) {
	if res, ok := data.([]mydb.DataItem); ok {
		res2 := []mydb.QualificationMatch{}
		for _, item := range res {
			item2, ok2 := item.(mydb.QualificationMatch)
			if ok2 {
				res2 = append(res2, item2)
			} else {
				return []mydb.QualificationMatch{}, errors.New("failed to convert to QualificationMatch")
			}
		}
		return res2, nil
	}
	return []mydb.QualificationMatch{}, errors.New("failed to convert to QualificationMatch")
}

func (p QualificationMatchConversion) ConvertFullItemsWithType(data interface{}) ([]mydb.QualificationMatchShort, error) {
	if res, ok := data.([]mydb.DataItem); ok {
		res2 := []mydb.QualificationMatchShort{}
		for _, item := range res {
			item2, ok2 := item.(mydb.QualificationMatchShort)
			if ok2 {
				res2 = append(res2, item2)
			} else {
				return []mydb.QualificationMatchShort{}, errors.New("failed to convert to QualificationMatchShort")
			}
		}
		return res2, nil
	}
	return []mydb.QualificationMatchShort{}, errors.New("failed to convert to QualificationMatchShort")
}
