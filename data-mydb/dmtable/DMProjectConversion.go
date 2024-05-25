package dmtable

import (
	"encoding/json"
	"errors"
	"request-matcher-openai/data-mydb/mydb"
)

type ProjectConversion struct {
}

func (p ProjectConversion) ParseOneItem(buf string) (mydb.DataItem, error) {
	data := mydb.Project{}
	err := json.Unmarshal([]byte(buf), &data)
	return data, err
}

func (p ProjectConversion) MarshalOneItem(item interface{}) (string, error) {
	if data, ok := item.(mydb.Project); ok {
		buf, err := json.Marshal(data)
		return string(buf), err
	} else {
		return "", errors.New("failed to convert to Project")
	}
}

func (p ProjectConversion) ParseOneFullItem(buf string) (mydb.DataItem, error) {
	data := mydb.ProjectShort{}
	err := json.Unmarshal([]byte(buf), &data)
	return data, err
}

func (p ProjectConversion) MarshalOneFullItem(item interface{}) (string, error) {
	if data, ok := item.(mydb.ProjectShort); ok {
		buf, err := json.Marshal(data)
		return string(buf), err
	} else {
		return "", errors.New("failed to convert to ProjectShort")
	}
}

func (p ProjectConversion) ParseItems(buf string) ([]mydb.DataItem, error) {
	res := []mydb.DataItem{}
	data := []mydb.Project{}
	err := json.Unmarshal([]byte(buf), &data)
	if err != nil {
		return res, err
	}
	for _, item := range data {
		res = append(res, item)
	}
	return res, nil
}

func (p ProjectConversion) MarshalItems(data []mydb.DataItem) (string, error) {
	res, err := p.ConvertItemsWithType(data)
	if err != nil {
		return "", err
	}
	buf, err := json.Marshal(res)
	return string(buf), err
}

func (p ProjectConversion) ParseFullItems(buf string) ([]mydb.DataItem, error) {
	res := []mydb.DataItem{}
	data := []mydb.ProjectShort{}
	err := json.Unmarshal([]byte(buf), &data)
	if err != nil {
		return res, err
	}
	for _, item := range data {
		res = append(res, item)
	}
	return res, nil
}

func (p ProjectConversion) MarshalFullItems(data []mydb.DataItem) (string, error) {
	res, err := p.ConvertFullItemsWithType(data)
	if err != nil {
		return "", err
	}
	buf, err := json.Marshal(res)
	return string(buf), err
}

func (p ProjectConversion) ConvertItems(data interface{}) ([]mydb.DataItem, error) {
	if res, ok := data.([]mydb.Project); ok {
		res2 := []mydb.DataItem{}
		for _, item := range res {
			res2 = append(res2, item)
		}
		return res2, nil
	}
	return []mydb.DataItem{}, errors.New("failed to convert to Project")
}

func (p ProjectConversion) ConvertFullItems(data interface{}) ([]mydb.DataItem, error) {
	if res, ok := data.([]mydb.ProjectShort); ok {
		res2 := []mydb.DataItem{}
		for _, item := range res {
			res2 = append(res2, item)
		}
		return res2, nil
	}
	return []mydb.DataItem{}, errors.New("failed to convert to ProjectShort")
}

func (p ProjectConversion) ConvertOneItemWithType(data interface{}) (mydb.Project, error) {
	res2, ok := data.(mydb.Project)
	if ok {
		return res2, nil
	}
	return mydb.Project{}, errors.New("failed to convert to Project")
}

func (p ProjectConversion) ConvertOneFullItemWithType(data interface{}) (mydb.ProjectShort, error) {
	res2, ok := data.(mydb.ProjectShort)
	if ok {
		return res2, nil
	}
	return mydb.ProjectShort{}, errors.New("failed to convert to ProjectShort")
}

func (p ProjectConversion) ConvertItemsWithType(data interface{}) ([]mydb.Project, error) {
	if res, ok := data.([]mydb.DataItem); ok {
		res2 := []mydb.Project{}
		for _, item := range res {
			item2, ok2 := item.(mydb.Project)
			if ok2 {
				res2 = append(res2, item2)
			} else {
				return []mydb.Project{}, errors.New("failed to convert to Project")
			}
		}
		return res2, nil
	}
	return []mydb.Project{}, errors.New("failed to convert to Project")
}

func (p ProjectConversion) ConvertFullItemsWithType(data interface{}) ([]mydb.ProjectShort, error) {
	if res, ok := data.([]mydb.DataItem); ok {
		res2 := []mydb.ProjectShort{}
		for _, item := range res {
			item2, ok2 := item.(mydb.ProjectShort)
			if ok2 {
				res2 = append(res2, item2)
			} else {
				return []mydb.ProjectShort{}, errors.New("failed to convert to ProjectShort")
			}
		}
		return res2, nil
	}
	return []mydb.ProjectShort{}, errors.New("failed to convert to ProjectShort")
}
