package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var fakeJSON = `[{
	"name": "LiLei",
	"sex": "M",
	"height": 150,
	"class": 1
}, {
	"name": "HanMeimei",
	"sex": "F",
	"height": 181,
	"class": 2
}]`

var fakeMap = []map[string]interface{}{{
	"name":   "LiLei",
	"sex":    "M",
	"height": 150,
	"class":  1,
}, {
	"name":   "HanMeimei",
	"sex":    "F",
	"height": 181,
	"class":  2,
},
}

func Test_parseFile(t *testing.T) {
	type args struct {
		fn      string
		content string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]map[string]interface{}
		wantErr bool
	}{
		{
			name: "1_ok",
			args: args{fn: "abc.json", content: fakeJSON},
			want: map[string][]map[string]interface{}{
				"abc": fakeMap,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := ioutil.TempFile("", tt.args.fn)
			if err != nil {
				t.Errorf("create temp file failed, Err: %+v", err)
			}
			defer os.Remove(f.Name())
			_, err = f.WriteString(tt.args.content)
			if err != nil {
				t.Errorf("write temp file failed, Err: %+v", err)
			}
			got, err := parseFile(f.Name())
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if deepEqual(got, tt.want) {
				t.Errorf("parseFile() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func deepEqual(a, b map[string][]map[string]interface{}) bool {
	for key, value := range a {
		v, ok := b[key]
		if !ok {
			return false
		}
		for i, item := range value {
			for itemK, itemV := range item {
				bv, ok := v[i][itemK]
				if !ok {
					return false
				}
				if !reflect.DeepEqual(itemV, bv) {
					return false
				}
			}
		}
	}
	return true
}

func Test_parseDir(t *testing.T) {
	type args struct {
		fns     []string
		content []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]map[string]interface{}
		wantErr bool
	}{
		{
			name: "1_ok",
			args: args{fns: []string{"abc.json", "edf.json"}, content: []string{fakeJSON, fakeJSON}},
			want: map[string][]map[string]interface{}{
				"abc": fakeMap,
				"edf": fakeMap,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "")
			if err != nil {
				t.Errorf("create temp dir failed, Err: %+v", err)
			}
			defer os.RemoveAll(dir)
			for i, file := range tt.args.fns {
				f, err := ioutil.TempFile(dir, file)
				if err != nil {
					t.Errorf("create temp file failed, Err: %+v", err)
				}
				_, err = f.WriteString(tt.args.content[i])
				if err != nil {
					t.Errorf("write temp file failed, Err: %+v", err)
				}
			}

			got, err := parseDir(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if deepEqual(got, tt.want) {
				t.Errorf("parseDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
