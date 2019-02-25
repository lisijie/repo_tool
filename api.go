package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
)

type APIHandler func(http.ResponseWriter, *http.Request, httprouter.Params) (interface{}, error)

// 获取项目状态
func apiStatus(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
    // 重新加载配置
    if c, err := NewConfig(*iniFile); err == nil {
        config = c
    }
    // 项目信息
    projects := make([]map[string]interface{}, 0)
    if config.GetSection("projects") != nil {
        data := config.GetSection("projects").GetAll()
        for name, path := range data {
            repo := NewGitRepo(path)
            bs, _ := repo.GetBranches()
            info := make(map[string]interface{})
            info["name"] = name
            info["current"], _ = repo.GetCurrentBranch()
            info["branches"] = bs
            projects = append(projects, info)
        }
    }
    // 命令列表
    commands := make(map[string]string)
    if config.GetSection("commands") != nil {
        commands = config.GetSection("commands").GetAll()
    }
    out := map[string]interface{}{
        "projects":   projects,
        "commands":   commands,
        "version":    Version,
        "build_time": BuildTime,
    }
    return out, nil
}

// 获取分支列表
func apiBranches(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
    project := r.URL.Query().Get("p")
    path := config.GetSection("projects").GetString(project)
    repo := NewGitRepo(path)
    return repo.GetBranches()
}

// 更新代码
func apiPull(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
    project := r.URL.Query().Get("p")
    path := config.GetSection("projects").GetString(project)
    repo := NewGitRepo(path)
    ret, err := repo.Pull()
    if err != nil {
        return "", nil
    }
    out := make(map[string]interface{})
    out["cmd"] = repo.LastCmd()
    out["out"] = ret
    return out, nil
}

// 切换分支
func apiCheckout(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
    project := r.URL.Query().Get("p")
    branch := r.URL.Query().Get("branch")
    path := config.GetSection("projects").GetString(project)
    repo := NewGitRepo(path)
    ret, err := repo.Checkout(branch)
    if err != nil {
        return "", nil
    }
    out := make(map[string]interface{})
    out["cmd"] = repo.LastCmd()
    out["out"] = ret
    return out, nil
}

// 清理无效追踪分支
func apiClean(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
    project := r.URL.Query().Get("p")
    path := config.GetSection("projects").GetString(project)
    repo := NewGitRepo(path)
    ret, err := repo.Clean()
    if err != nil {
        return "", nil
    }
    out := make(map[string]interface{})
    out["cmd"] = repo.LastCmd()
    out["out"] = ret
    return out, nil
}

// 执行命令
func apiCommand(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
    name := r.URL.Query().Get("name")
    commands := config.GetSection("commands")
    if commands == nil || commands.GetString(name) == "" {
        return nil, errors.New("command not found")
    }
    cmdStr := commands.GetString(name)
    cmd := NewCommand(cmdStr)
    if err := cmd.Run(); err != nil {
        return "", err
    }
    out := make(map[string]string)
    out["cmd"] = cmdStr
    out["out"] = string(cmd.Stdout())
    return out, nil
}

func decorate(f APIHandler) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        code := 200
        msg := ""
        data, err := f(w, r, ps)
        if err != nil {
            code = 500
            msg = err.Error()
        }
        out := make(map[string]interface{})
        if code == 200 {
            out["code"] = 0
        } else {
            out["code"] = code
            out["msg"] = msg
        }
        out["data"] = data
        response, err := json.Marshal(out)
        if err != nil {
            code = 500
            response = []byte(fmt.Sprintf(`{"code":500, "msg":"%s", "data":null}`, err.Error()))
        }
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        w.WriteHeader(200)
        w.Write(response)
    }
}
