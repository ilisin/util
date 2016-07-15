package paramParser

import (
	"testing"
)

type IntOrNull struct {
	IsNull bool  `thrift:"IsNull,1" json:"IsNull"`
	Value  int32 `thrift:"Value,2" json:"Value"`
}

type StringOrNull struct {
	IsNull bool   `thrift:"IsNull,1" json:"IsNull"`
	Value  string `thrift:"Value,2" json:"Value"`
}

type BoolOrNull struct {
	IsNull bool `thrift:"IsNull,1" json:"IsNull"`
	Value  bool `thrift:"Value,2" json:"Value"`
}

type CommonParam struct {
	Numbers map[string]*IntOrNull    `thrift:"Numbers,1" json:"Numbers"`
	Strings map[string]*StringOrNull `thrift:"Strings,2" json:"Strings"`
	Bools   map[string]*BoolOrNull   `thrift:"Bools,3" json:"Bools"`
}

func init() {
	RegisterMap("test", []ParamOption{
		ParamOption{
			Key:          "CI_KeyInt1",
			Required:     false,
			DefaultValue: IntOrNull{false, 0},
		},
		ParamOption{
			Key:          "CI_KeyInt2",
			Required:     true,
			DefaultValue: IntOrNull{true, 0},
		},
		ParamOption{
			Key:          "CS_KeyString1",
			Required:     true,
			DefaultValue: StringOrNull{false, "st1"},
		},
		ParamOption{
			Key:          "CS_KeyString2",
			Required:     false,
			DefaultValue: StringOrNull{true, ""},
		},
		ParamOption{
			Key:          "CB_KeyBool1",
			Required:     false,
			DefaultValue: BoolOrNull{false, true},
		},
		ParamOption{
			Key:          "CB_KeyBool2",
			Required:     true,
			DefaultValue: BoolOrNull{true, false},
		},
	})
}

func TestParse(t *testing.T) {
	param := CommonParam{
		Numbers: make(map[string]*IntOrNull),
		Strings: make(map[string]*StringOrNull),
	}
	param.Numbers["CI_KeyInt1"] = &IntOrNull{false, 14}
	v, err := Parse("test", "CI_KeyInt1", param)
	if err != nil {
		t.Error(err)
	} else {
		if in, ok := v.(*IntOrNull); ok {
			if in.Value != int32(14) {
				t.Errorf("取值错误 ：%v", in.Value)
			}
		} else {
			t.Errorf("%v 解析错误", in)
		}
	}

	v, err = Parse("test", "CI_KeyIntt1", param)
	if err == nil {
		t.Error(err)
	}

	v, err = Parse("test", "CI_KeyInt2", param)
	if err != ErrRequired {
		t.Errorf("必须的参数验证错误")
	}

	param.Strings["CS_KeyString1"] = &StringOrNull{false, "xxxxxxx"}
	v, err = Parse("test", "CS_KeyString1", param)
	if err != nil {
		t.Error(err)
	} else {
		if in, ok := v.(*StringOrNull); ok {
			if in.Value != "xxxxxxx" {
				t.Errorf("取值错误 ：%v", in.Value)
			}
		} else {
			t.Errorf("%v 解析错误", in)
		}
	}

	v, err = Parse("test", "CS_KeyString2", param)
	if err != nil {
		t.Errorf("%v", err)
	} else {
		t.Logf("%v", v)
	}

}
