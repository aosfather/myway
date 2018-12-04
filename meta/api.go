package meta

//api的描述
type Api struct {
	NameSpace      string          //所在的模块，用于空间的管理
	Url            string          //api的url
	Desc           string          //描述
	Method         []HttpMethod    //允许的访问方法
	Cluster        *ServerCluster  //对应的集群
	Access         IPAccessControl //ip访问控制
	ServerUrl      string          //服务对应url
	Domain         string          //对应的域名
	Status         Status
	AuthFilter     string
	MatchRule      MatchRule //url匹配规则
	CircuitBreaker *CircuitBreaker
	MaxQPS         int64
}

func (this *Api) Key() string {
	return this.NameSpace + "/" + this.Url
}

// 通用状态
type Status byte

const (
	Down    Status = 0
	Up      Status = 1
	Unknown Status = 2
)

type Validation struct {
	Parameter Parameter
	Required  bool
	Rules     []ValidationRule
}

type ValidationRule struct {
	RuleType   RuleType
	Expression string
}

type RuleType byte

const (
	RuleRegexp RuleType = 0
)

//参数来源
type ParameterSource byte

const (
	PS_QueryString ParameterSource = 0
	PS_FormData    ParameterSource = 1
	PS_JSONBody    ParameterSource = 2
	PS_Header      ParameterSource = 3
	PS_Cookie      ParameterSource = 4
	PS_PathValue   ParameterSource = 5
)

// Parameter is a parameter from a http request
type Parameter struct {
	Name   string          //参数名
	Source ParameterSource //参数来源
	Index  int32           //序号
}

//ip权限控制
type IPAccessControl struct {
	Whitelist []string
	Blacklist []string
}
