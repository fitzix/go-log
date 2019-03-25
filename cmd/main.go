package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/fitzix/go-log/models"
	"github.com/fitzix/go-log/reader"
	"github.com/urfave/cli"
)

var (
	VERSION   = "DEV"
	COMMIT    = "UNKNOWN"
	DATE      = "UNKNOWN"
	GOVERSION = "UNKNOWN"
)

func main() {
	var configPath string

	App := cli.NewApp()

	App.Name = "GO TCP/UDP LOG SERVER"
	App.Version = fmt.Sprintf("%v, commit %v, built at %v, %v", VERSION, COMMIT, DATE, GOVERSION)
	App.Usage = "GO TCP/UDP 高并发日志收集工具"
	App.Author = "Fitzix"
	App.Email = "caojunkaiv@gmail.com"

	App.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config,file,c,f",
			Usage:       "Load configuration from `FILE`",
			Value:       "config.toml",
			Destination: &configPath,
		},
	}

	App.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "generate config.toml",
			Action: func(c *cli.Context) error {
				file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					log.Println("创建配置文件失败", err)
					return cli.NewExitError("\n", 1)
				}
				defer file.Close()
				if _, err = file.WriteString(models.DefaultServerConf); err != nil {
					log.Println("写入配置文件失败", err)
					return cli.NewExitError("\n", 1)
				}
				log.Println("配置文件生成成功")
				return nil
			},
		},
	}

	App.Action = func(c *cli.Context) error {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Println("未找到配置文件,请先使用init方法生成配置文件")
			return cli.NewExitError("", 1)
		}
		log.Println("开始解析配置文件")
		ServerConf := new(models.SenderConf)
		if _, err := toml.DecodeFile(configPath, ServerConf); err != nil {
			log.Println("解析配置文件失败,请检查配置文件格式")
			return cli.NewExitError("", 1)
		}
		return reader.Start()

	}

	if err := App.Run(os.Args); err != nil {
		log.Println("failed", err)
	}
}
