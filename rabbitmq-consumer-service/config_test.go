package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "ok_config",
			args: args{configPath: "test_config.yaml"},
			want: &Config{
				RabbitMQ: struct {
					URL            string `yaml:"url"`
					QueueName      string `yaml:"queue_name"`
					PrefetchCount  int    `yaml:"prefetch_count"`
					ConsumerCount  int    `yaml:"consumer_count"`
					WorkerPoolSize int    `yaml:"worker_pool_size"`
				}{
					URL: "amqp://guest:guest@localhost:5672/",
					QueueName: "test_queue",
					PrefetchCount: 10,
					ConsumerCount: 5,
					WorkerPoolSize: 20,
				},
			},
			wantErr: false,
		},
		{
			name:    "non_existent_config",
			args:    args{configPath: "test_non_existent_config.yaml"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid_format_config",
			args:    args{configPath: "test_invalid_format_config.yaml"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 如果是有效的文件，创建文件
			if tt.name == "ok_config" {
				// 创建有效配置文件
				configData := `
rabbitmq:
  url: amqp://guest:guest@localhost:5672/
  queue_name: test_queue
  prefetch_count: 10
  consumer_count: 5
  worker_pool_size: 20
`
				err := ioutil.WriteFile(tt.args.configPath, []byte(configData), 0644)
				assert.NoError(t, err)
				defer os.Remove(tt.args.configPath)
			} else if tt.name == "invalid_format_config" {
				// 创建无效的 YAML 配置文件
				invalidData := `
rabbitmq:
  url: amqp://guest:guest@localhost:5672/
  queue_name: test_queue
  prefetch_count: invalid_number
`
				err := ioutil.WriteFile(tt.args.configPath, []byte(invalidData), 0644)
				assert.NoError(t, err)
				defer os.Remove(tt.args.configPath)
			}

			got, err := LoadConfig(tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
