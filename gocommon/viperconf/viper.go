package viperconf

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func LoadConfiguration(env string, v *viper.Viper) *viper.Viper {
	// 1 config.json
	// 2 config_xxx_template.json
	// 3 config_xxx_part-${ENV}.json
	// 4 config-${ENV}.json

	// 1
	v.SetConfigType("json")
	v.AddConfigPath("config")
	v.SetConfigName("config")
	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("reading config.json failed : %v\n", err)
	}

	dir, err := ioutil.ReadDir("config")
	if err != nil {
		fmt.Printf("find configuration files error: %v\n", err)
		return v
	}
	var files []string
	var templateFiles []string
	var templateEnvFiles []string
	var customerFiles []string

	configEnvFile := "config-" + env
	for _, fi := range dir {
		fileName := strings.TrimSuffix(fi.Name(), ".json")
		//skip the fixed config file
		if fileName == "config" || fileName == configEnvFile || fileName == ".license" {
			fmt.Printf("skip the configure file:%v\n", fileName)
			continue
		}
		if fileName == fi.Name() {
			fmt.Printf("the file:%v does not end with json, skip it\n", fileName)
			continue
		}

		//match template file
		reg1 := regexp.MustCompile(`config_(.*?)_template`)

		//match template env file
		tmpsuffix := "_part-" + env
		reg2 := regexp.MustCompile(`config_(.*?)` + tmpsuffix)
		if reg1.Match([]byte(fileName)) {
			templateFiles = append(templateFiles, fileName)
		} else if reg2.Match([]byte(fileName)) {
			templateEnvFiles = append(templateEnvFiles, fileName)
		} else {
			fmt.Printf("error to load config file because the name is illegal: %v\n", fileName)
		}

		files = append(files, fileName)
	}
	//read custom dir
	dir2, err2 := ioutil.ReadDir("config/customer")
	if err2 != nil {
		fmt.Printf("skip to load customer configuration files for error: %v\n", err2)
	} else {
		for _, fi := range dir2 {
			fileName := strings.TrimSuffix(fi.Name(), ".json")
			if fileName == fi.Name() {
				fmt.Printf("the file:%v does not end with json, skip it\n", fileName)
				continue
			}
			customerFiles = append(customerFiles, "customer/"+fileName)
		}
	}

	fmt.Printf("find configure files number: %d and files:%v\n", len(files), files)
	fmt.Printf("start merge config files \n")

	// 2
	for _, fileitem := range templateFiles {
		v.SetConfigName(fileitem)
		err = v.MergeInConfig()
		if err != nil {
			fmt.Printf("error in merge the config file %s, error: %v\n", fileitem, err)
		} else {
			fmt.Printf("merge the config file %s, \n", fileitem)
		}
	}

	// 3
	for _, fileitem := range templateEnvFiles {
		v.SetConfigName(fileitem)
		err = v.MergeInConfig()
		if err != nil {
			fmt.Printf("error in merge the config file %s, error: %v\n", fileitem, err)
		} else {
			fmt.Printf("merge the config file %s, \n", fileitem)
		}
	}

	// 3
	for _, fileitem := range customerFiles {
		v.SetConfigName(fileitem)
		err = v.MergeInConfig()
		if err != nil {
			fmt.Printf("error in merge the config file %s, error: %v\n", fileitem, err)
		} else {
			fmt.Printf("merge the config file %s, \n", fileitem)
		}
	}

	// 4
	v.SetConfigName(configEnvFile)
	err = v.MergeInConfig()
	if err != nil {
		fmt.Printf("error in merge the config file %s, error: %v\n", configEnvFile, err)
	} else {
		fmt.Printf("merge the config file %s, \n", configEnvFile)
	}
	return v
}

func MergeConfiguration(viper *viper.Viper) {
	pflag.VisitAll(func(flg *pflag.Flag) {
		if flg.Changed {
			return
		}
		viperValue := viper.Get(flg.Name)
		if viperValue != nil {
			strValue, err1 := cast.ToStringE(viperValue)
			if err1 == nil {
				err1 = pflag.Set(flg.Name, strValue)
				if err1 == nil {
					fmt.Printf("set pflag %s from viper\n", flg.Name)
				} else {
					fmt.Printf("err set pflag %s from viper err: %s\n", flg.Name, err1)
				}
				//flg.Value.Set(viperValue.(string))
			} else {
				fmt.Printf("%s cast.ToStringE err %s\n", flg.Name, err1)
			}
		}
	})
	viper.BindPFlags(pflag.CommandLine)
}

func DumpConfiguration(viper *viper.Viper, dpath string) {
	err := viper.WriteConfigAs(dpath)
	if err == nil {
		fmt.Printf("dump config:%v\n", dpath)
	} else {
		fmt.Printf("error happened when dumping config：%s\n", err.Error())
	}
}
