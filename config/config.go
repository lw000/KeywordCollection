package config

import (
	"errors"
	"fmt"
	"github.com/Unknwon/goconfig"
	tymysql "github.com/lw000/gocommon/db/mysql"
	tyrdsex "github.com/lw000/gocommon/db/rdsex"
	"strconv"
)

// IniConfig ini配置
type IniConfig struct {
	RdsCfg   *tyrdsex.JsonConfig
	MysqlCfg *tymysql.JsonConfig
	TLS      struct {
		Enable   bool
		CertFile string
		KeyFile  string
	}
	Port     int64
	Debug    int64
	SplitLog int
}

// NewIniConfig ...
func NewIniConfig() *IniConfig {
	return &IniConfig{
		RdsCfg:   &tyrdsex.JsonConfig{},
		MysqlCfg: &tymysql.JsonConfig{},
	}
}

// LoadIniConfig ...
func LoadIniConfig(file string) (*IniConfig, error) {
	cfg := NewIniConfig()
	err := cfg.Load(file)
	return cfg, err
}

// Load ...
func (c *IniConfig) Load(file string) error {
	var (
		err error
		f   *goconfig.ConfigFile
	)

	f, err = goconfig.LoadConfigFile(file)
	if err != nil {
		return fmt.Errorf("配置文件读取失败 [%s]", file)
	}

	err = c.readMainCfg(f)
	if err != nil {
		return err
	}

	err = c.readTlsCfg(f)
	if err != nil {
		return err
	}

	err = c.readMysqlCfg(f)
	if err != nil {
		return err
	}

	// err = c.readRdsCfg(f)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (c *IniConfig) readMainCfg(f *goconfig.ConfigFile) error {
	var (
		err      error
		port     string
		debug    string
		splitlog string
	)

	section := "main"

	port, err = f.GetValue(section, "port")
	if err != nil {
		return fmt.Errorf("获取键值(%s): %s", "port", err.Error())
	}

	c.Port, err = strconv.ParseInt(port, 10, 64)
	if err != nil {
		return errors.New(err.Error())
	}

	debug, err = f.GetValue(section, "debug")
	if err != nil {
		return fmt.Errorf("获取键值(%s): %s", "debug", err.Error())
	}

	c.Debug, err = strconv.ParseInt(debug, 10, 64)
	if err != nil {
		return errors.New(err.Error())
	}

	splitlog, err = f.GetValue(section, "splitlog")
	if err != nil {
		return fmt.Errorf("获取键值(%s): %s", "splitlog", err.Error())
	}

	c.SplitLog, err = strconv.Atoi(splitlog)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (c *IniConfig) readMysqlCfg(f *goconfig.ConfigFile) error {
	var (
		err          error
		maxOdleConns string
		maxOpenConns string
		section      string
	)
	if c.Debug == 1 {
		section = "dev_mysql"
	} else {
		section = "prod_mysql"
	}

	c.MysqlCfg.Username, err = f.GetValue(section, "username")
	if err != nil {
		return fmt.Errorf("get key(%s): %s", "username", err.Error())
	}

	c.MysqlCfg.Password, err = f.GetValue(section, "password")
	if err != nil {
		return fmt.Errorf("get key(%s): %s", "password", err.Error())
	}

	c.MysqlCfg.Host, err = f.GetValue(section, "host")
	if err != nil {
		return fmt.Errorf("get key(%s): %s", "host", err.Error())
	}

	c.MysqlCfg.Database, err = f.GetValue(section, "database")
	if err != nil {
		return fmt.Errorf("get key(%s): %s", "database", err.Error())
	}

	maxOdleConns, err = f.GetValue(section, "MaxOdleConns")
	if err != nil {
		return fmt.Errorf("get key(%s): %s", "MaxOdleConns", err.Error())
	}
	c.MysqlCfg.MaxOdleConns, err = strconv.ParseInt(maxOdleConns, 10, 64)
	if err != nil {
		return fmt.Errorf("get key(%s): %s", "MaxOdleConns", err.Error())
	}

	maxOpenConns, err = f.GetValue(section, "MaxOpenConns")
	if err != nil {
		return err
	}
	c.MysqlCfg.MaxOpenConns, err = strconv.ParseInt(maxOpenConns, 10, 64)
	if err != nil {
		return err
	}

	return nil
}

func (c *IniConfig) readRdsCfg(f *goconfig.ConfigFile) error {
	var (
		err          error
		Db           string
		PoolSize     string
		MinIdleConns string
	)

	section := "redis"
	c.RdsCfg.Host, err = f.GetValue(section, "host")
	if err != nil {
		return fmt.Errorf("get key(%s): %s", "host", err.Error())
	}

	c.RdsCfg.Psd, err = f.GetValue(section, "psd")
	if err != nil {
		return fmt.Errorf("get key(%s): %s", "psd", err.Error())
	}

	Db, err = f.GetValue(section, "db")
	if err != nil {
		return fmt.Errorf("get key(%s): %s", "db", err.Error())
	}
	c.RdsCfg.Db, err = strconv.ParseInt(Db, 10, 64)
	if err != nil {
		return err
	}

	PoolSize, err = f.GetValue(section, "poolSize")
	if err != nil {
		return fmt.Errorf("get key(%s): %s", "poolSize", err.Error())
	}
	c.RdsCfg.PoolSize, err = strconv.ParseInt(PoolSize, 10, 64)
	if err != nil {
		return err
	}

	MinIdleConns, err = f.GetValue(section, "minIdleConns")
	if err != nil {
		return fmt.Errorf("get key(%s): %s", "minIdleConns", err.Error())
	}
	c.RdsCfg.MinIdleConns, err = strconv.ParseInt(MinIdleConns, 10, 64)
	if err != nil {
		return err
	}
	return nil
}

func (c *IniConfig) readTlsCfg(f *goconfig.ConfigFile) error {
	var (
		err    error
		enable string
	)

	section := "tls"

	enable, err = f.GetValue(section, "enable")
	if err != nil {
		return fmt.Errorf("get key(%s): %s", "port", err.Error())
	}
	c.TLS.Enable, err = strconv.ParseBool(enable)
	if err != nil {
		return errors.New(err.Error())
	}

	if c.TLS.Enable {
		c.TLS.CertFile, err = f.GetValue(section, "certFile")
		if err != nil {
			return fmt.Errorf("get key(%s): %s", "certFile", err.Error())
		}
		if c.TLS.CertFile == "" {
			return errors.New("cretFile is empty")
		}

		c.TLS.KeyFile, err = f.GetValue(section, "keyFile")
		if err != nil {
			return fmt.Errorf("get key(%s): %s", "keyFile", err.Error())
		}

		if c.TLS.KeyFile == "" {
			return errors.New("keyFile is empty")
		}
	}

	return nil
}

func (c IniConfig) String() string {
	return fmt.Sprintf("{%v, %v}", c.MysqlCfg, c.RdsCfg)
}
