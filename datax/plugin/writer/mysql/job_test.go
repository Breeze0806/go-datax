package mysql

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

func TestJob_Init(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		j       *Job
		args    args
		conf    *config.JSON
		jobConf *config.JSON
		wantErr bool
	}{
		{
			name: "1",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return &mockExecer{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: testJSONFromString(`{
				"job":{
					"content":[
						{
							"writer": {
								"name":"mysqlwriter",
								"parameter":{

								}
							}
						}
					]
				}
			}`),
		},
		{
			name: "2",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return &mockExecer{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromString(`{}`),
			jobConf: testJSONFromString(`{
				"job":{
					"content":[
						{
							"writer": {
								"name":"mysqlwriter",
								"parameter":{

								}
							}
						}
					]
				}
			}`),
			wantErr: true,
		},
		{
			name: "3",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return &mockExecer{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: testJSONFromString(`{
				"job":{
					"content":[
						{
							"writer": {
								"name":"mysqlwriter"
							}
						}
					]
				}
			}`),
			wantErr: true,
		},
		{
			name: "4",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return &mockExecer{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: testJSONFromString(`{
				"job":{
					"content":[
						{
							"writer": {
								"name":"mysqlwriter",
								"parameter":{
									"username": 1
								}
							}
						}
					]
				}
			}`),
			wantErr: true,
		},
		{
			name: "5",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return nil, errors.New("mock error")
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: testJSONFromString(`{
				"job":{
					"content":[
						{
							"writer": {
								"name":"mysqlwriter",
								"parameter":{
								}
							}
						}
					]
				}
			}`),
			wantErr: true,
		},
		{
			name: "6",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				newExecer: func(name string, conf *config.JSON) (Execer, error) {
					return &mockExecer{
						queryErr: errors.New("mock error"),
					}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: testJSONFromString(`{
				"job":{
					"content":[
						{
							"writer": {
								"name":"mysqlwriter",
								"parameter":{
								}
							}
						}
					]
				}
			}`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.j.SetPluginConf(tt.conf)
			tt.j.SetPluginJobConf(tt.jobConf)
			if err := tt.j.Init(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Job.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJob_Destroy(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		j       *Job
		args    args
		wantErr bool
	}{
		{
			name: "1",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				execer:  &mockExecer{},
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.j.Destroy(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Job.Destroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJob_Split(t *testing.T) {
	type args struct {
		ctx    context.Context
		number int
	}
	tests := []struct {
		name    string
		j       *Job
		args    args
		jobConf *config.JSON
		want    []*config.JSON
		wantErr bool
	}{
		{
			name: "1",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
			},
			args: args{
				ctx:    context.TODO(),
				number: 1,
			},
			jobConf: testJSONFromString(`{
				"job":{
					"content":[
						{
							"reader": {
								"name":"mysqlreader",
								"parameter":{
								}
							}
						}
					]
				}
			}`),
			want: []*config.JSON{
				testJSONFromString(`{
				"job":{
					"content":[
						{
							"reader": {
								"name":"mysqlreader",
								"parameter":{
								}
							}
						}
					]
				}
			}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.j.SetPluginJobConf(tt.jobConf)
			got, err := tt.j.Split(tt.args.ctx, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("Job.Split() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.Split() = %v, want %v", got, tt.want)
			}
		})
	}
}
