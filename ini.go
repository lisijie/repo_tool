package main

import (
    "bufio"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    Section
    filename string
    sections map[string]*Section
}

// 解析指定ini文件并创建Config对象
func NewConfig(fileName string) (*Config, error) {
    c := &Config{
        filename: fileName,
        sections: make(map[string]*Section),
    }
    c.data = make(map[string]string)
    if err := c.init(); err != nil {
        return nil, err
    }
    return c, nil
}

func (c *Config) init() error {
    file, err := os.Open(c.filename)
    if err != nil {
        return err
    }
    defer file.Close()

    bs := bufio.NewScanner(file)
    section := ""
    for bs.Scan() {
        line := strings.TrimSpace(bs.Text())
        if len(line) == 0 || line[0:1] == "#" || line[0:1] == ";" {
            continue
        }
        if line[0:1] == "[" {
            section = line[1:len(line)-1]
            continue
        }
        ss := strings.SplitN(line, "=", 2)
        if len(ss) != 2 {
            continue
        }
        k := strings.TrimSpace(ss[0])
        if k != "" {
            v := strings.TrimSpace(ss[1])
            if section == "" {
                c.data[k] = v
            } else {
                if _, ok := c.sections[section]; !ok {
                    c.sections[section] = newSection()
                }
                c.sections[section].data[k] = v
            }
        }
    }
    return nil
}

func (c *Config) GetSection(name string) *Section {
    return c.sections[name]
}

type Section struct {
    data map[string]string
}

func newSection() *Section {
    return &Section{data: make(map[string]string)}
}

// 获取所有配置项map
func (s *Section) GetAll() map[string]string {
    return s.data
}

// 获取一个字符串值的配置项，如果值不存在返回空字符串
func (s *Section) GetString(key string, def ...string) string {
    if v, ok := s.data[key]; ok {
        return v
    }
    if len(def) > 0 {
        return def[0]
    }
    return ""
}

// 获取一个整型配置项，如果出错或值不存在返回0
func (s *Section) GetInt(key string, def ...int) int {
    if v, ok := s.data[key]; ok {
        i, _ := strconv.Atoi(v)
        return i
    }
    if len(def) > 0 {
        return def[0]
    }
    return 0
}
