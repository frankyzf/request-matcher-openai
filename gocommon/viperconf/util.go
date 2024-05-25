package viperconf

import "github.com/spf13/viper"

func GetDefaultString(vp *viper.Viper, name string, defaultValue string) string {
	if vp != nil && vp.IsSet(name) == true {
		return vp.GetString(name)
	}
	return defaultValue
}

func GetDefaultBool(vp *viper.Viper, name string, defaultValue bool) bool {
	if vp != nil && vp.IsSet(name) == true {
		return vp.GetBool(name)
	}
	return defaultValue
}

func GetDefaultInt(vp *viper.Viper, name string, defaultValue int) int {
	if vp != nil && vp.IsSet(name) == true {
		return vp.GetInt(name)
	}
	return defaultValue
}

func GetDefaultFloat64(vp *viper.Viper, name string, defaultValue float64) float64 {
	if vp != nil && vp.IsSet(name) == true {
		return vp.GetFloat64(name)
	}
	return defaultValue
}

func GetDefaultStringSlice(vp *viper.Viper, name string, defaultValue []string) []string {
	if vp != nil && vp.IsSet(name) == true {
		value := vp.GetStringSlice(name)
		return value
	}
	return defaultValue
}
