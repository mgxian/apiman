package setting

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
)

var (
	Conf       string
	CustomConf string
	Cfg        *ini.File
	AppPath    string
)

func execPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}

func WorkDir() (string, error) {
	AppPath, err := execPath()
	if err != nil {
		println("app path error")
	}
	AppPath = strings.Replace(AppPath, "\\", "/", -1)

	i := strings.LastIndex(AppPath, "/")
	if i == -1 {
		return AppPath, nil
	}
	return AppPath[:i], nil
}

func NewConfig() {
	workdir, err := WorkDir()
	if err != nil {
		println("workdir error in newconfig")
	}
	if len(Conf) == 0 {
		//println("conf len 0")
		Conf = workdir + "/conf/app.conf"
	}
	println(Conf)
	Cfg, err = ini.Load(Conf)
	if err != nil {
		println("ini load error")
	}

	if len(CustomConf) == 0 {
		//println("customconf len 0")
		CustomConf = workdir + "/conf/custom.app.conf"
	}
	println(CustomConf)
	Cfg.Append(CustomConf)
	//sec := Cfg.Section("database")
	//names := sec.Key("type").String()
	//println(names)
}
