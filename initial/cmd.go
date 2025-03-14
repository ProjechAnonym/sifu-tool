package initial

import (
	"os"

	"github.com/alecthomas/kingpin/v2"
)

func Initkingpin() (*string, *string, *string, *[]string){
	app := kingpin.New("sifu-tool", "A application for managing the ddns and ssl certificate.")
	app.Version("1.0.0")
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')
	run := app.Command("run", "Boot up the application.")
	environment := run.Flag("env", "Path of the working directory").Short('e').PlaceHolder("/opt/sifutool/lib").Required().String()
	config := run.Flag("config", "Path of the config file").Short('c').PlaceHolder("/opt/sifutool/config/config.yaml").String()
	address := run.Flag("address", "Address to listen on").Short('a').PlaceHolder(":8080").Default(":8080").String()
	domains := run.Arg("origins", "The origins which allowed to access the application").PlaceHolder("http://example1.com http://example2.com").Default("*").Strings()
	run.Help(`Boot up the application. The flags "--config" and "--env" are required. If the flag "--address" is empty, it will take ":8080" by default.If "origins" args are empty, it will allow any domains to access the application.`)
	kingpin.MustParse(app.Parse(os.Args[1:]))
	return environment, config, address, domains
}