package claude

import (
	"reflect"
	"testing"
)

func Test_parseChunk_delta(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want parsedChunk
	}{
		{name: "Empty chunk", args: args{data: []byte(`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":""}}`)}, want: parsedChunk{Empty: true}},
		{name: "Filled Delta", args: args{data: []byte(`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"foo"}}`)}, want: parsedChunk{Empty: false, Delta: "foo"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseChunk(tt.args.data, "content_block_delta"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseChunk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseChunk_start(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want parsedChunk
	}{
		{name: "Empty chunk", args: args{data: []byte(`{"type":"content_block_delta","index":0,"content_block":{"type":"text_delta","text":""}}`)}, want: parsedChunk{Empty: true}},
		{name: "Filled Delta", args: args{data: []byte(`{"type":"content_block_delta","index":0,"content_block":{"type":"text_delta","text":"foo"}}`)}, want: parsedChunk{Empty: false, Delta: "foo"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseChunk(tt.args.data, "content_block_start"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseChunk() = %v, want %v", got, tt.want)
			}
		})
	}
}
