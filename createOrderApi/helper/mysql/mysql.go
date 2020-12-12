package mysql

import (
	"database/sql"
	"fmt"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Output   output
	Database database
}

type database struct {
	Server    string
	Port      string
	Database  string
	User      string
	Password1 string
}

type output struct {
	Directory string
	Format    string
}

func Connect() {
	var conf Config
	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", conf)

	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.Database.User, conf.Database.Password1, conf.Database.Server, conf.Database.Port, conf.Database.Database)

	db, err := sql.Open("mysql", connString)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer db.Close()

}
