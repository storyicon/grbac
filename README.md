# GRBAC [![CircleCI](https://circleci.com/gh/storyicon/grbac/tree/master.svg?style=svg)](https://circleci.com/gh/storyicon/grbac/tree/master) [![Go Report Card](https://goreportcard.com/badge/github.com/storyicon/grbac)](https://goreportcard.com/report/github.com/storyicon/grbac)  [![Build Status](https://travis-ci.org/storyicon/grbac.svg?branch=master)](https://travis-ci.org/storyicon/grbac) [![GoDoc](https://godoc.org/github.com/storyicon/grbac?status.svg)](https://godoc.org/github.com/storyicon/grbac) [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/storyicon/Lobby)

![grbac](https://raw.githubusercontent.com/storyicon/grbac/master/docs/screenshot/grbac.png)

[中文文档](https://github.com/storyicon/grbac/blob/master/docs/README-chinese.md)

Grbac is a fast, elegant and concise [RBAC](https://en.wikipedia.org/wiki/Role-based_access_control) framework. It supports [enhanced wildcards](#4-enhanced-wildcards) and matches HTTP requests using [Radix](https://en.wikipedia.org/wiki/Radix) trees. Even more amazing is that you can easily use it in any existing database and data structure.        

What grbac does is ensure that the specified resource can only be accessed by the specified role. Please note that grbac is not responsible for the storage of rule configurations and "what roles the current request initiator has". It means you should configure the rule information first and provide the roles that the initiator of each request has.        

grbac treats the combination of `Host`, `Path`, and `Method` as a `Resource`, and binds the `Resource` to a set of role rules (called `Permission`). Only users who meet these rules can access the corresponding `Resource`.        

The component that reads the rule information is called `Loader`. grbac presets some loaders, you can also customize a loader by implementing `func()(grbac.Rules, error)` and load it via `grbac.WithLoader`.        

- [1. Most Common Use Case](#1-most-common-use-case)
- [2. Concept](#2-concept)
    - [2.1. Rule](#21-rule)
    - [2.2. Resource](#22-resource)
    - [2.3. Permission](#23-permission)
    - [2.4. Loader](#24-loader)
- [3. Other Examples](#3-other-examples)
    - [3.1. gin && grbac.WithJSON](#31-gin--grbacwithjson)
    - [3.2. echo && grbac.WithYaml](#32-echo--grbacwithyaml)
    - [3.3. iris && grbac.WithRules](#33-iris--grbacwithrules)
    - [3.4. ace && grbac.WithAdvancedRules](#34-ace--grbacwithadvancedrules)
    - [3.5. gin && grbac.WithLoader](#35-gin--grbacwithloader)
- [4. Enhanced wildcards](#4-enhanced-wildcards)
- [5. BenchMark](#5-benchmark)
- [6. Production](#6-production)

## 1. Most Common Use Case

Below is the most common use case, which uses `gin` and wraps `grbac` as a middleware. With this example, you can easily know how to use `grbac` in other http frameworks(like `echo`, `iris`, `ace`, etc):

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/storyicon/grbac"
    "net/http"
    "time"
)

func LoadAuthorizationRules() (rules grbac.Rules, err error) {
    // Implement your logic here
    // ...
    // You can load authorization rules from database or file
    // But you need to return your authentication rules in the form of grbac.Rules
    // tips: You can also bind this function to a golang struct
    return
}

func QueryRolesByHeaders(header http.Header) (roles []string,err error){
    // Implement your logic here
    // ...
    // This logic maybe take a token from the headers and
    // query the user's corresponding roles from the database based on the token.
    return roles, err
}

func Authorization() gin.HandlerFunc {
    // Here, we use a custom Loader function via "grbac.WithLoader"
    // and specify that this function should be called every minute to update the authentication rules.
    // Grbac also offers some ready-made Loaders:
    // grbac.WithYAML
    // grbac.WithRules
    // grbac.WithJSON
    // ...
    rbac, err := grbac.New(grbac.WithLoader(LoadAuthorizationRules, time.Minute))
    if err != nil {
        panic(err)
    }
    return func(c *gin.Context) {
        roles, err := QueryRolesByHeaders(c.Request.Header)
        if err != nil {
            c.AbortWithError(http.StatusInternalServerError, err)
            return
        }
        state, _ := rbac.IsRequestGranted(c.Request, roles)
        if !state.IsGranted() {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
    }
}

func main(){
    c := gin.New()
    c.Use(Authorization())

    // Bind your API here
    // ...

    c.Run(":8080")
}
```
## 2. Concept

Here are some concepts about `grbac`. It's very simple, you may only need three minutes to understand.

### 2.1. Rule

```go
// Rule is used to define the relationship between "resource" and "permission"
type Rule struct {
    // The ID controls the priority of the rule.
    // The higher the ID means the higher the priority of the rule.
    // When a request is matched to more than one rule,
    // then authentication will only use the permission configuration for the rule with the highest ID value.
    // If there are multiple rules that are the largest ID, then one of them will be used randomly.
    ID int `json:"id"`
    *Resource
    *Permission
}
```

As you can see, the `Rule` consists of three parts: `ID`, `Resource`, and `Permission`.      
The `ID` determines the priority of the Rule.      
When a request meets multiple rules at the same time (such as in a wildcard), 
`grbac` will select the one with the highest ID, then authenticate with its Permission definition.       
If multiple rules of the same ID are matched at the same time, grbac will randomly select one from them.      

Here is a very simple example:

```yaml
#Rule
- id: 0
  # Resource
  host: "*"
  path: "**"
  method: "*"
  # Permission
  authorized_roles:
  - "*"
  forbidden_roles: []
  allow_anyone: false

#Rule 
- id: 1
  # Resource
  host: domain.com
  path: "/article"
  method: "{DELETE,POST,PUT}"
  # Permission
  authorized_roles:
  - editor
  forbidden_roles: []
  allow_anyone: false
```

In this configuration file written in yaml format, the rule with `ID=0` states that all resources can be accessed by anyone with any role.
But the rule with `ID=1` states that only the `editor` can operate on the article.          
Then, except that the operation of the article can only be accessed by the `editor`, all other resources can be accessed by anyone with any role.         

### 2.2. Resource

```go
// Resource defines resources
type Resource struct {
    // Host defines the host of the resource, allowing wildcards to be used.
    Host string `json:"host"`
    // Path defines the path of the resource, allowing wildcards to be used.
    Path string `json:"path"`
    // Method defines the method of the resource, allowing wildcards to be used.
    Method string `json:"method"`
}
```

Resource is used to describe which resources a rule applies to. 
When `IsRequestGranted(c.Request, roles)` is executed, grbac first matches the current `Request` with the `Resources` in all `Rule`s.

Each field of Resource supports [enhanced wildcards](#4-enhanced-wildcards)

### 2.3. Permission

```go
// Permission is used to define permission control information
type Permission struct {
    // AuthorizedRoles defines roles that allow access to specified resource
    // Accepted type: non-empty string, *
    //      *: means any role, but visitors should have at least one role,
    //      non-empty string: specified role
    AuthorizedRoles []string `json:"authorized_roles"`
    // ForbiddenRoles defines roles that not allow access to specified resource
    // ForbiddenRoles has a higher priority than AuthorizedRoles
    // Accepted type: non-empty string, *
    //      *: means any role, but visitors should have at least one role,
    //      non-empty string: specified role
    //
    ForbiddenRoles []string `json:"forbidden_roles"`
    // AllowAnyone has a higher priority than ForbiddenRoles/AuthorizedRoles
    // If set to true, anyone will be able to pass authentication.
    // Note that this will include people without any role.
    AllowAnyone bool `json:"allow_anyone"`
}
```

`Permission` is used to define the authorization rules of the `Resource` to which it is bound.
That's understandable. When the roles of the requester meets the definition of `Permission`, he will be allowed access, otherwise he will be denied access.

For faster speeds, fields in `Permission` do not support `enhanced wildcards`.
Only `*` is allowed in `AuthorizedRoles` and `ForbiddenRoles` to indicate `all`.

### 2.4. Loader

Loader is used to load authorization rules. grbac presets some loaders, you can also customize a loader by implementing `func()(grbac.Rules, error)` and load it via `grbac.WithLoader`.        

| method | description |
| --- | --- |
| WithJSON(path, interval)  | periodically load rules configuration from `json` file  |
| WithYaml(path, interval) | periodically load rules configuration from `yaml` file |
| WithRules(Rules) | load rules configuration from `grbac.Rules` |
| WithAdvancedRules(loader.AdvancedRules) | load advanced rules from `loader.AdvancedRules`| 
| WithLoader(loader func()(Rules, error), interval) | periodically load rules with custom functions |
     
`interval` defines the reload period of the authentication rule.     
When `interval < 0`, `grbac` will abandon periodically loading the configuration file;     
When `interval∈[0,1s)`, `grbac` will automatically set the `interval` to `5s`;     

## 3. Other Examples

Here are some simple examples to make it easier to understand how `grbac` works.     
Although `grbac` works well in most http frameworks, I am sorry that I only use gin now, so if there are some flaws in the example below, please let me know.     
### 3.1. gin && grbac.WithJSON

If you want to write the configuration file in a `JSON` file, you can load it via `grbac.WithJSON(file, interval)`, `file` is your json file path, and grbac will reload the file every `interval`.     

```json
[
    {
        "id": 0,
        "host": "*",
        "path": "**",
        "method": "*",
        "authorized_roles": [
            "*"
        ],
        "forbidden_roles": [
            "black_user"
        ],
        "allow_anyone": false
    },
    {
        "id":1,
        "host": "domain.com",
        "path": "/article",
        "method": "{DELETE,POST,PUT}",
        "authorized_roles": ["editor"],
        "forbidden_roles": [],
        "allow_anyone": false
    }
]
```
The above is an example of authentication rule in `JSON` format. It's structure is based on [grbac.Rules](#21-rule).

```go

func QueryRolesByHeaders(header http.Header) (roles []string,err error){
    // Implement your logic here
    // ...
    // This logic maybe take a token from the headers and
    // query the user's corresponding roles from the database based on the token.
    return roles, err
}

func Authentication() gin.HandlerFunc {
    rbac, err := grbac.New(grbac.WithJSON("config.json", time.Minute * 10))
    if err != nil {
        panic(err)
    }
    return func(c *gin.Context) {
        roles, err := QueryRolesByHeaders(c.Request.Header)
        if err != nil {
            c.AbortWithError(http.StatusInternalServerError, err)
            return
        }

        state, err := rbac.IsRequestGranted(c.Request, roles)
        if err != nil {
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }

        if !state.IsGranted() {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
    }
}

func main(){
    c := gin.New()
    c.Use(Authentication())

    // Bind your API here
    // ...
    
    c.Run(":8080")
}

```

### 3.2. echo && grbac.WithYaml

If you want to write the configuration file in a `YAML` file, you can load it via `grbac.WithYAML(file, interval)`, `file` is your yaml file path, and grbac will reload the file every `interval`.

```yaml
#Rule
- id: 0
  # Resource
  host: "*"
  path: "**"
  method: "*"
  # Permission
  authorized_roles:
  - "*"
  forbidden_roles: []
  allow_anyone: false

#Rule 
- id: 1
  # Resource
  host: domain.com
  path: "/article"
  method: "{DELETE,POST,PUT}"
  # Permission
  authorized_roles:
  - editor
  forbidden_roles: []
  allow_anyone: false
```

The above is an example of authentication rule in `YAML` format. It's structure is based on [grbac.Rules](#21-rule).

```go
func QueryRolesByHeaders(header http.Header) (roles []string,err error){
    // Implement your logic here
    // ...
    // This logic maybe take a token from the headers and
    // query the user's corresponding roles from the database based on the token.
    return roles, err
}

func Authentication() echo.MiddlewareFunc {
    rbac, err := grbac.New(grbac.WithYAML("config.yaml", time.Minute * 10))
    if err != nil {
            panic(err)
    }
    return func(echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            roles, err := QueryRolesByHeaders(c.Request().Header)
            if err != nil {
                    c.NoContent(http.StatusInternalServerError)
                    return nil
            }
            state, err := rbac.IsRequestGranted(c.Request(), roles)
            if err != nil {
                    c.NoContent(http.StatusInternalServerError)
                    return nil
            }
            if state.IsGranted() {
                    return nil
            }
            c.NoContent(http.StatusUnauthorized)
            return nil
        }
    }
}

func main(){
    c := echo.New()
    c.Use(Authentication())

    // Implement your logic here
    // ...
}
```

### 3.3. iris && grbac.WithRules

If you want to write the authentication rules directly in the code, `grbac.WithRules(rules)` provides this way, you can use it like this:

```go

func QueryRolesByHeaders(header http.Header) (roles []string,err error){
    // Implement your logic here
    // ...
    // This logic maybe take a token from the headers and
    // query the user's corresponding roles from the database based on the token.
    return roles, err
}

func Authentication() iris.Handler {
    var rules = grbac.Rules{
        {
            ID: 0,
            Resource: &grbac.Resource{
                Host: "*",
                Path: "**",
                Method: "*",
            },
            Permission: &grbac.Permission{
                AuthorizedRoles: []string{"*"},
                ForbiddenRoles: []string{"black_user"},
                AllowAnyone: false,
            },
        },
        {
            ID: 1,
            Resource: &grbac.Resource{
                Host: "domain.com",
                Path: "/article",
                Method: "{DELETE,POST,PUT}",
            },
            Permission: &grbac.Permission{
                AuthorizedRoles: []string{"editor"},
                ForbiddenRoles: []string{},
                AllowAnyone: false,
            },
        },
    }
    rbac, err := grbac.New(grbac.WithRules(rules))
    if err != nil {
        panic(err)
    }
    return func(c context.Context) {
        roles, err := QueryRolesByHeaders(c.Request().Header)
        if err != nil {
            c.StatusCode(http.StatusInternalServerError)
            c.StopExecution()
            return
        }
        state, err := rbac.IsRequestGranted(c.Request(), roles)
        if err != nil {
            c.StatusCode(http.StatusInternalServerError)
            c.StopExecution()
            return
        }
        if !state.IsGranted() {
            c.StatusCode(http.StatusUnauthorized)
            c.StopExecution()
            return
        }
    }
}

func main(){
    c := iris.New()
    c.Use(Authentication())

    // Implement your logic here
    // ...
}
```

### 3.4. ace && grbac.WithAdvancedRules

If you want to write the authentication rules directly in the code, `grbac.WithAdvancedRules(rules)` provides this way, you can use it like this:

```go

func QueryRolesByHeaders(header http.Header) (roles []string,err error){
    // Implement your logic here
    // ...
    // This logic maybe take a token from the headers and
    // query the user's corresponding roles from the database based on the token.
    return roles, err
}

func Authentication() ace.HandlerFunc {
    var advancedRules = loader.AdvancedRules{
        {
            Host: []string{"*"},
            Path: []string{"**"},
            Method: []string{"*"},
            Permission: &grbac.Permission{
                AuthorizedRoles: []string{},
                ForbiddenRoles: []string{"black_user"},
                AllowAnyone: false,
            },
        },
        {
            Host: []string{"domain.com"},
            Path: []string{"/article"},
            Method: []string{"PUT","DELETE","POST"},
            Permission: &grbac.Permission{
                AuthorizedRoles: []string{"editor"},
                ForbiddenRoles: []string{},
                AllowAnyone: false,
            },
        },
    }
    auth, err := grbac.New(grbac.WithAdvancedRules(advancedRules))
    if err != nil {
        panic(err)
    }
    return func(c *ace.C) {
        roles, err := QueryRolesByHeaders(c.Request.Header)
        if err != nil {
        c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
        state, err := auth.IsRequestGranted(c.Request, roles)
        if err != nil {
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
        if !state.IsGranted() {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
    }
}
func main(){
    c := ace.New()
    c.Use(Authentication())

    // Implement your logic here
    // ...
}

```

`loader.AdvancedRules` attempts to provide a simpler way to define authentication rules than `grbac.Rules`.


### 3.5. gin && grbac.WithLoader

```go

func QueryRolesByHeaders(header http.Header) (roles []string,err error){
    // Implement your logic here
    // ...
    // This logic maybe take a token from the headers and
    // query the user's corresponding roles from the database based on the token.
    return roles, err
}

type MySQLLoader struct {
    session *gorm.DB
}

func NewMySQLLoader(dsn string) (*MySQLLoader, error) {
    loader := &MySQLLoader{}
    db, err := gorm.Open("mysql", dsn)
    if err  != nil {
        return nil, err
    }
    loader.session = db
    return loader, nil
}

func (loader *MySQLLoader) LoadRules() (rules grbac.Rules, err error) {
    // Implement your logic here
    // ...
    // Extract data from the database, assemble it into grbac.Rules and return
    return
}

func Authentication() gin.HandlerFunc {
    loader, err := NewMySQLLoader("user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
    if err != nil {
        panic(err)
    }
    rbac, err := grbac.New(grbac.WithLoader(loader.LoadRules, time.Second * 5))
    if err != nil {
        panic(err)
    }
    return func(c *gin.Context) {
        roles, err := QueryRolesByHeaders(c.Request.Header)
        if err != nil {
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
            
        state, err := rbac.IsRequestGranted(c.Request, roles)
        if err != nil {
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
        if !state.IsGranted() {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
    }
}

func main(){
    c := gin.New()
    c.Use(Authorization())

    // Bind your API here
    // ...

    c.Run(":8080")
}
```

## 4. Enhanced wildcards

`Wildcard` supported syntax:        
```text
pattern:
  { term }
term:
  '*'         matches any sequence of non-path-separators
  '**'        matches any sequence of characters, including
              path separators.
  '?'         matches any single non-path-separator character
  '[' [ '^' ] { character-range } ']'
        character class (must be non-empty)
  '{' { term } [ ',' { term } ... ] '}'
  c           matches character c (c != '*', '?', '\\', '[')
  '\\' c      matches character c

character-range:
  c           matches character c (c != '\\', '-', ']')
  '\\' c      matches character c
  lo '-' hi   matches character c for lo <= c <= hi
```

## 5. BenchMark

```go
➜ gos test -bench=. 
goos: linux
goarch: amd64
pkg: github.com/storyicon/grbac/pkg/tree
BenchmarkTree_Query                   2000           541397 ns/op
BenchmarkTree_Foreach_Query           2000           1360719 ns/op
PASS
ok      github.com/storyicon/grbac/pkg/tree     13.182s
```
The test case contains 1000 random rules, and the `BenchmarkTree_Query` and `BenchmarkTree_Foreach_Query` functions test four requests separately, after calculation:

```
541397/(4*1e9)=0.0001s
```

When there are 1000 rules, the average verification time per request is `0.0001s`.

## 6. Production      

`grbac` has been used in the `production` environment by the following companies:    

![wallstreetcn](https://raw.githubusercontent.com/storyicon/grbac/master/docs/screenshot/wallstreetcn.png)