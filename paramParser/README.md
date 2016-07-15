## parse为通用参数解析器 ##

**在使用之前需要先注册参数默认值**

> 注册方法如下

	func init() {
		paramParser.RegisterMap("myService", []paramParser.ParamOption{
			paramParser.ParamOption{
				Key:          "CI_KeyInt1",
				Required:     false,
				DefaultValue: IntOrNull{false, 0},
			},
			paramParser.ParamOption{
				Key:          "CI_KeyInt2",
				Required:     true,
				DefaultValue: IntOrNull{true, 0},
			},
			paramParser.ParamOption{
				Key:          "CS_KeyString1",
				Required:     true,
				DefaultValue: "DEF_V",
			},
			paramParser.ParamOption{
				Key:          "CB_KeyBool1",
				Required:     true,
				DefaultValue: true,
			},
		})
	}
Key申明参数标示名，而Required则标示参数是否必须，如果设置为true而参数值未指定则会报**ErrRequired**异常

DefaultValue为对应默认值,支持以下几种数据类型

- int,int8,int16,int32,int64
- 数值类型的封装，例如

		type IntOrNull struct {
			IsNull bool  `thrift:"IsNull,1" json:"IsNull"`
			Value  int32 `thrift:"Value,2" json:"Value"`
		}
其中需要包含bool类型的IsNull字段和数值(int8,int18,int32,int64,ing)类型的Value字段
- string
- 字符串类型的数据封装，例如

		type StringOrNull struct {
			IsNull bool   `thrift:"IsNull,1" json:"IsNull"`
			Value  string `thrift:"Value,2" json:"Value"`
		}
其中bool类型的IsNull字段和string类型的Value字段是必须的
- bool
- bool类型的数据封装,例如

		type BoolOrNull struct {
			IsNull bool `thrift:"IsNull,1" json:"IsNull"`
			Value  bool `thrift:"Value,2" json:"Value"`
		}
其中bool类型的IsNull和bool类型的Value是必须的
> 满足数值类型或者数据类型的封装的条件可以通过以下函数解析参数值

	func ParseInt(pname string, paramKey string, param interface{}) (int, bool, error)

> 而满足字符串类型或者字符串类型封装的条件可以通过以下函数解析参数

	func ParseString(pname string, paramKey string, param interface{}) (string, bool, error)

> 满足bool类型的封装的条件可以通过以下函数解析参数

	func ParseBool(pname string, paramKey string, param interface{}) (bool, bool, error)

如果以上三种都不满足则可以使用

	func Parse(pname string, paramKey string, param interface{}) (interface{}, error) 
来解析
