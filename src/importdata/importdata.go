package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type options struct {
	// file or dir path
	path string

	dsn string
}

func gatherOption(fs *flag.FlagSet, args ...string) options {
	var o options
	fs.StringVar(&o.path, "path", "", "data file path")
	fs.StringVar(&o.dsn, "dsn", "", "data source name")
	fs.Parse(args)
	return o
}

func (o *options) validate() error {
	if o.path == "" {
		return errors.New("empty file path")
	}

	if o.dsn == "" {
		return errors.New("empty database")
	}
	return nil
}

func parsePath(path string) (map[string][]map[string]interface{}, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return parseDir(path)
	}

	return parseFile(path)
}

func parseDir(path string) (map[string][]map[string]interface{}, error) {
	m := make(map[string][]map[string]interface{})
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		// depth is 1
		if p == path {
			return nil
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		fm, err := parseFile(p)
		if err != nil {
			return err
		}
		for k, v := range fm {
			m[k] = v
		}

		return nil
	})
	return m, err
}

func parseFile(path string) (map[string][]map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	data := []map[string]interface{}{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	m := make(map[string][]map[string]interface{})
	fn := filepath.Base(path)
	m[fn[0:len(fn)-len(filepath.Ext(fn))]] = data

	return m, nil
}

func kvs(m map[string]interface{}) (ks []string, vs []interface{}) {
	for k, v := range m {
		ks = append(ks, k)
		vs = append(vs, v)
	}
	return
}

func insertDB(dsn string, m map[string][]map[string]interface{}) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	for table, data := range m {
		for _, v := range data {
			clumns, values := kvs(v)
			sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
				table, strings.Join(clumns, ","), strings.Repeat(", ?", len(values))[1:])
			if err := db.Exec(sql, values...).Error; err != nil {
				return err
			}
		}
	}

	log.Println("insert successfully")
	return nil
}

func main() {
	o := gatherOption(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.validate(); err != nil {
		log.Fatalf("invalid options, Err: %+v", err)
	}

	m, err := parsePath(o.path)
	if err != nil {
		log.Fatalf("read path failed, Err: %+v", err)
	}

	if err := insertDB(o.dsn, m); err != nil {
		log.Fatalf("insert data into db failed, Err: %+v", err)
	}
}
