package mysql

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/storage/database"
)

func testReader(filename string) *Reader {
	reader, err := NewReader(filename)
	if err != nil {
		panic(err)
	}
	return reader
}

func TestReader_Job(t *testing.T) {
	tests := []struct {
		name string
		r    *Reader
		want reader.Job
		conf *config.JSON
	}{
		{
			name: "1",
			r:    testReader(filepath.Join("resources", "plugin.json")),
			want: &Job{
				BaseJob: plugin.NewBaseJob(),
				newQuerier: func(name string, conf *config.JSON) (Querier, error) {
					return database.Open(name, conf)
				},
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.SetPluginConf(tt.conf)
			if got := tt.r.Job(); !reflect.DeepEqual(got.PluginConf(), tt.want.PluginConf()) {
				t.Errorf("Reader.Job() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_Task(t *testing.T) {
	tests := []struct {
		name string
		r    *Reader
		want reader.Task
		conf *config.JSON
	}{
		{
			name: "1",
			r:    testReader(filepath.Join("resources", "plugin.json")),
			want: &Task{
				BaseTask: plugin.NewBaseTask(),
				newQuerier: func(name string, conf *config.JSON) (Querier, error) {
					return database.Open(name, conf)
				},
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.SetPluginConf(tt.conf)
			if got := tt.r.Task(); !reflect.DeepEqual(got.PluginConf(), tt.want.PluginConf()) {
				t.Errorf("Reader.Task() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewReader(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantR   *Reader
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				filename: filepath.Join("resources", "plugin.json"),
			},
			wantR: testReader(filepath.Join("resources", "plugin.json")),
		},
		{
			name: "2",
			args: args{
				filename: filepath.Join("tmpresources", "tmpplugin.json"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := NewReader(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("NewReader() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
