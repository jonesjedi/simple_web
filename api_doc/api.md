## 域名
https://www.onb.io/

## 接口汇总

### 1.注册接口

#### 接口地址:
/api/user/register
#### 请求方式：
POST
#### 请求参数：
|  参数名   | 类型  | 是否必须   | 说明 |
|  ----  | ----  | ----  | ----  |
| user_name  | string | Y  | 用户名 |
| email  | string |Y  | 邮箱 |
| user_pwd  | string |Y  | 密码 |
#### 返回
```
{
    "ret":0,
    "msg":"succ",
    "data":""
}
```

### 2.发送注册邮件接口

#### 接口地址:
/api/user/send_validate_email
#### 请求方式：
POST
#### 请求参数：
|  参数名   | 类型  | 是否必须   | 说明 |
|  ----  | ----  | ----  | ----  |
|  user_email  | string  | Y  | 用户邮箱  |

#### 返回
```json
{
    "ret":0,
    "msg":"succ",
    "data":""
}
```

### 3.验证邮箱接口

#### 接口地址:
/api/user/validate_email
#### 请求方式：
POST
#### 请求参数：
|  参数名   | 类型  | 是否必须   | 说明 |
|  ----  | ----  | ----  | ----  |
|  code  | string  | 是  | 验证code  |
|  user_id  | string  | 是  | 用户名  |

#### 返回
```json
{
    "ret":0,
    "msg":"succ",
    "data":""
}
```

### 4.登录接口

#### 接口地址:
/api/user/login
#### 请求方式：
Form 表单提交
#### 请求参数：
|  参数名   | 类型  | 是否必须   | 说明 |
|  ----  | ----  | ----  | ----  |
|  user_name  | string  | 是  | 用户名  |
|  user_pwd  | string  | 是  | 密码  |

#### 返回
```json
跳转个人主页
```

### 5.发送密码重置邮件

#### 接口地址:
/api/user/reset_pwd_email
#### 请求方式：
POST
#### 请求参数：
|  参数名   | 类型  | 是否必须   | 说明 |
|  ----  | ----  | ----  | ----  |
|  user_email  | string  | 是  | 用户邮箱  |

#### 返回
```json
{
    "ret":0,
    "msg":"succ",
    "data":""
}
```

### 6.密码重置请求

#### 接口地址:
/api/user/reset_pwd
#### 请求方式：
POST
#### 请求参数：
|  参数名   | 类型  | 是否必须   | 说明 |
|  ----  | ----  | ----  | ----  |
|  user_name  | string  | 是  | 用户名称  |
|  code  | string  | 是  | 上一个请求会在链接里埋一个code，在这里带过来，验证是否来自正常的重置流程  |
|  new_pwd  | string  | 是  | 新密码  |

#### 返回
```json
{
    "ret":0,
    "msg":"succ",
    "data":""
}
```

### 7.个人主页

#### 接口地址:
/api/user/index
#### 请求方式：
GET
#### 请求参数：
|  参数名   | 类型  | 是否必须   | 说明 |
|  ----  | ----  | ----  | ----  |
不需要参数，从登录态获取用户信息

#### 返回
```json
{
    "ret":0,
    "msg":"succ",
    "data":{
        "user_name":"kris",
        "user_avatar":"用户头像",
        "user_link":"用户主页",
        "is_confirmed":false,
        "user_extra":"其他附加信息"
    }
}
```

### 8.更新个人头像

#### 接口地址:
/api/user/update_info
#### 请求方式：
GET
#### 请求参数：
|  参数名   | 类型  | 是否必须   | 说明 |
|  ----  | ----  | ----  | ----  |
|  user_img  | string  | 是  | 新头像链接  |

#### 返回
```json
{
    "ret":0,
    "msg":"succ",
    "data":""
}
```


### 9.个人链接列表

#### 接口地址:
/api/link/link_list
#### 请求方式：
GET
#### 请求参数：
|  参数名   | 类型  | 是否必须   | 说明 |
|  ----  | ----  | ----  | ----  |
不需要参数，从登录态获取用户信息

#### 返回
```json
{
    "ret":0,    
    "msg":"succ",
    "data":[
        {
            "link_url":"http://www.qq.com",
            "link_desc":"etss",
            "user_img":"首图链接"
        }
    ]
}
```


