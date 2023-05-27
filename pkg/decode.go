package pkg

import (
	"encoding/json"

	"github.com/graphql-go/graphql"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

func ParseInto(src any, dst any) error {
	jsonByte, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonByte, dst)
}

func ParseGraphQLData(src *graphql.Result, dst any, key string) error {
	jsonByte, err := json.Marshal(src.Data)
	if err != nil {
		return err
	}

	var data map[string]json.RawMessage
	if err = json.Unmarshal(jsonByte, &data); err != nil {
		return err
	}

	return json.Unmarshal(data[key], dst)
}

func ProtoAnyToMap(in *anypb.Any) (map[string]interface{}, error) {
	out := structpb.Struct{}
	if err := in.UnmarshalTo(&out); err != nil {
		return nil, err
	}

	return out.AsMap(), nil
}

func ToProtoAny(in any) (*anypb.Any, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := structpb.Struct{}
	if err = protojson.Unmarshal(data, &out); err != nil {
		return nil, err
	}

	return anypb.New(&out)
}

func ParseProtoAnyInto(src *anypb.Any, dst any) error {
	mp, err := ProtoAnyToMap(src)
	if err != nil {
		return err
	}

	return ParseInto(mp, dst)
}
