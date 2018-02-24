package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/apex/log"
	lcli "github.com/apex/log/handlers/cli"
	"github.com/fatih/color"
	"github.com/fitzix/go-udp/server/models"
	"github.com/fitzix/go-udp/server/reader"
	"github.com/urfave/cli"
)

var (
	VERSION   = "DEV"
	COMMIT    = "UNKNOWN"
	DATE      = "UNKNOWN"
	GOVERSION = "UNKNOWN"

	bold = color.New(color.Bold)
)

func init() {
	log.SetHandler(lcli.Default)
}

func main() {
	var configPath string

	App := cli.NewApp()

	App.Name = "GO UDP LOG SERVER"
	App.Version = fmt.Sprintf("%v, commit %v, built at %v, %v", VERSION, COMMIT, DATE, GOVERSION)
	App.Usage = "GO UDP 高并发日志收集工具"

	App.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config,file,c,f",
			Usage:       "Load configuration from `FILE`",
			Value:       "config.toml",
			Destination: &configPath,
		},
	}

	App.Action = func(c *cli.Context) error {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.WithError(err).Error(color.RedString("未找到配置文件,请先使用init方法生成配置文件"))
			return cli.NewExitError("", 1)
		}
		log.Infof(bold.Sprint("开始解析配置文件"))
		ServerConf := &reader.ServerConf
		if _, err := toml.DecodeFile(configPath, ServerConf); err != nil {
			log.WithError(err).Error(color.RedString("解析配置文件失败,请检查配置文件格式"))
			return cli.NewExitError("", 1)
		}
		reader.Start()
		return nil
	}

	App.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "generate config.toml",
			Action: func(c *cli.Context) error {
				file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					log.WithError(err).Error("创建配置文件失败")
					return cli.NewExitError("\n", 1)
				}
				defer file.Close()
				if _, err = file.WriteString(models.DefaultServerConf); err != nil {
					log.WithError(err).Error("写入配置文件失败")
					return cli.NewExitError("\n", 1)
				}
				log.Infof(bold.Sprint("配置文件生成成功"))
				return nil
			},
		},
	}

	if err := App.Run(os.Args); err != nil {
		log.WithError(err).Fatal("failed")
	}
}
