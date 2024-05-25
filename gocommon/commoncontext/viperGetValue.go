package commoncontext

func IsFieldSet(name string) bool {
	vp := GetInstance().VP
	return vp.IsSet(name)
}

func GetDefaultString(name string, defaultValue string) string {
	vp := GetInstance().VP
	if vp != nil && vp.IsSet(name) == true {
		return vp.GetString(name)
	}
	return defaultValue
}

func GetDefaultBool(name string, defaultValue bool) bool {
	vp := GetInstance().VP
	if vp != nil && vp.IsSet(name) == true {
		return vp.GetBool(name)
	}
	return defaultValue
}

func GetDefaultInt(name string, defaultValue int) int {
	vp := GetInstance().VP
	if vp != nil && vp.IsSet(name) == true {
		return vp.GetInt(name)
	}
	return defaultValue
}

func GetDefaultFloat64(name string, defaultValue float64) float64 {
	vp := GetInstance().VP
	if vp != nil && vp.IsSet(name) == true {
		return vp.GetFloat64(name)
	}
	return defaultValue
}

func GetDefaultStringSlice(name string, defaultValue []string) []string {
	vp := GetInstance().VP
	if vp != nil && vp.IsSet(name) == true {
		value := vp.GetStringSlice(name)
		return value
	}
	return defaultValue
}

func GetDefaultStringMap(name string, defaultValue map[string]string) map[string]string {
	vp := GetInstance().VP
	if vp != nil && vp.IsSet(name) == true {
		value := vp.GetStringMapString(name)
		return value
	}
	return defaultValue
}

func GetDefaultInterfaceMap(name string, defaultValue map[string]interface{}) map[string]interface{} {
	vp := GetInstance().VP
	if vp != nil && vp.IsSet(name) == true {
		value := vp.GetStringMap(name)
		return value
	}
	return defaultValue
}

func GetDefaultNode(name string, defaultValue interface{}) interface{} {
	vp := GetInstance().VP
	if vp != nil && vp.IsSet(name) == true {
		value := vp.Get(name)
		return value
	}
	return defaultValue
}
