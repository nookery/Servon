package deploy

type Project struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`           // 项目名称
    Type        string `json:"type"`           // 项目类型：static/node/go/python
    GitRepo     string `json:"git_repo"`       // Git仓库地址
    Branch      string `json:"branch"`         // 分支
    BuildCmd    string `json:"build_cmd"`      // 构建命令
    OutputDir   string `json:"output_dir"`     // 构建输出目录
    Domain      string `json:"domain"`         // 域名
    Path        string `json:"path"`           // URL路径
    Port        int    `json:"port"`           // 服务端口(非静态项目)
    Status      string `json:"status"`         // 状态：running/stopped/failed
    LastDeploy  string `json:"last_deploy"`    // 最后部署时间
    Environment []Env  `json:"environment"`    // 环境变量
}

type Env struct {
    Key   string `json:"key"`
    Value string `json:"value"`
} 