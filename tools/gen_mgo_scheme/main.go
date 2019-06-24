package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/simplejia/utils"
)

var (
	env string
)

type SchemeIndex struct {
	Fields     []string `json:"fields,omitempty"`
	Unique     bool     `json:"unique,omitempty"`
	Background bool     `json:"background,omitempty"`
}

type Scheme struct {
	Dsn      string         `json:"dsn,omitempty"`
	AuthUser string         `json:"auth_user,omitempty"`
	Db       string         `json:"db,omitempty"`
	Table    string         `json:"table,omitempty"`
	DbNum    int            `json:"db_num,omitempty"`
	TableNum int            `json:"table_num,omitempty"`
	Indices  []*SchemeIndex `json:"indices,omitempty"`
}

// Conf 定义配置参数
type Conf struct {
	Schemes []*Scheme `json:"schemes,omitempty"`
}

func getSchemes() (schemes []*Scheme, err error) {
	file := filepath.Join("mongo", "scheme.json")

	if _, err = os.Stat(file); err != nil {
		if !os.IsNotExist(err) {
			return
		}
		err = nil
		return
	}

	fcontent, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	fcontent = utils.RemoveAnnotation(fcontent)
	var envs map[string]*Conf
	if err = json.Unmarshal(fcontent, &envs); err != nil {
		return
	}

	conf := envs[env]
	if conf == nil {
		return
	}

	schemes = conf.Schemes
	return
}

func runSchemeCommand(scheme *Scheme) (err error) {
	dbTables := [][2]string{}
	if tableNum := scheme.TableNum; tableNum > 1 {
		dbNum := scheme.DbNum
		if dbNum == 0 {
			dbNum = 1
		}

		for i := 0; i < tableNum; i++ {
			db := fmt.Sprintf("%s_%d", scheme.Db, i%scheme.DbNum)
			table := fmt.Sprintf("%s_%d", scheme.Table, i)
			dbTables = append(dbTables, [...]string{db, table})
		}
	} else {
		dbTables = append(dbTables, [...]string{scheme.Db, scheme.Table})
	}

	session, err := mgo.Dial(scheme.Dsn)
	if err != nil {
		return
	}

	if authUser := scheme.AuthUser; authUser != "" {
		for _, dbTable := range dbTables {
			cmd := bson.D{
				{
					"grantRolesToUser",
					authUser,
				},
				{
					"roles",
					[]interface{}{
						bson.M{
							"role": "readWrite",
							"db":   dbTable[0],
						},
					},
				},
			}

			if err = session.Run(cmd, nil); err != nil {
				return
			}
		}
	}

	if indices := scheme.Indices; len(indices) != 0 {
		for _, dbTable := range dbTables {
			for _, index := range scheme.Indices {
				col := session.DB(dbTable[0]).C(dbTable[1])
				if err = col.EnsureIndex(mgo.Index{
					Key:        index.Fields,
					Unique:     index.Unique,
					Background: index.Background,
				}); err != nil {
					return
				}
			}
		}
	}

	return
}

func exit(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	println()
	println("Failed!")
	os.Exit(-1)
}

func main() {
	flag.StringVar(&env, "env", "prod", "set env")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "A tiny tool, used to generate project's db scheme\n")
		fmt.Fprintf(os.Stderr, "version: 1.11, Created by simplejia [12/2018]\n\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	println("Begin generate scheme")

	schemes, err := getSchemes()
	if err != nil {
		exit("get schemes err: %v", err)
	}

	log.Printf("env: %s\nconf: %s\n", env, utils.Iprint(schemes))

	for _, scheme := range schemes {
		if err := runSchemeCommand(scheme); err != nil {
			exit("run scheme command err: %v", err)
		}
	}

	println("Success!")
	return
}
