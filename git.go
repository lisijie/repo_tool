package main

import "strings"

type GitRepo struct {
    path    string
    lastCmd string
}

func NewGitRepo(path string) *GitRepo {
    return &GitRepo{path: path}
}

func (r *GitRepo) GetBranches() ([]string, error) {
    r.lastCmd = "git branch -r --no-color"
    cmd := NewCommand(r.lastCmd)
    if err := cmd.RunInDir(r.path); err != nil {
        return nil, err
    }
    out := string(cmd.Stdout())
    lines := strings.Split(out, "\n")
    branches := make([]string, 0, len(lines))
    prefixLen := len("origin/")
    for k, v := range lines {
        if k == 0 {
            continue
        }
        if len(v) > prefixLen {
            branches = append(branches, strings.TrimSpace(v)[prefixLen:])
        }
    }
    return branches, nil
}

func (r *GitRepo) GetCurrentBranch() (string, error) {
    r.lastCmd = "git branch --no-color"
    cmd := NewCommand(r.lastCmd)
    if err := cmd.RunInDir(r.path); err != nil {
        return "", err
    }
    out := string(cmd.Stdout())
    lines := strings.Split(out, "\n")
    for _, v := range lines {
        if len(v) > 0 && v[0:1] == "*" {
            return strings.TrimSpace(v[1:]), nil
        }
    }
    return "", nil
}

func (r *GitRepo) Pull() (string, error) {
    r.lastCmd = "git pull"
    cmd := NewCommand(r.lastCmd)
    if err := cmd.RunInDir(r.path); err != nil {
        return "", err
    }
    return string(cmd.Stdout()), nil
}

func (r *GitRepo) Checkout(branch string) (string, error) {
    r.lastCmd = "git reset --hard HEAD && git checkout " + branch + " && git pull"
    cmd := NewCommand(r.lastCmd)
    if err := cmd.RunInDir(r.path); err != nil {
        return "", err
    }
    return string(cmd.Stdout()), nil
}

func (r *GitRepo) Clean() (string, error) {
    r.lastCmd = "git remote prune origin"
    cmd := NewCommand(r.lastCmd)
    if err := cmd.RunInDir(r.path); err != nil {
        return "", err
    }
    return string(cmd.Stdout()), nil
}

func (r *GitRepo) LastCmd() string {
    return r.lastCmd
}
