package fireworks

import (
	"reflect"
	"testing"
)

func Test_parseChunk(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want parsedChunk
	}{
		{name: "Empty chunk", args: args{data: []byte("\n")}, want: parsedChunk{Empty: true}},
		{name: "Final Chunk", args: args{data: []byte("data: [DONE]\n")}, want: parsedChunk{LastMessage: true}},
		{name: "Empty Delta", args: args{data: []byte("data: {\"model\":\"accounts/fireworks/models/mixtral-8x7b-instruct\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\"},\"finish_reason\":null}]}\n")}, want: parsedChunk{Delta: ""}},
		{name: "Filled Delta", args: args{data: []byte("data: {\"model\":\"accounts/fireworks/models/mixtral-8x7b-instruct\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"Hello world\"},\"finish_reason\":null}]}\n")}, want: parsedChunk{Delta: "Hello world"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseChunk(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseChunk() = %v, want %v", got, tt.want)
			}
		})
	}
}
