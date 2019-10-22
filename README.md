# Chakala 使用文档
## 一、概述
### 1.1 背景
大量的项目，例如：大屏、BI类、图表展示类、演示类项目工作以展示为主。可以说整个项目除了登录以外，底层就是大量的 **select** 语句。这种项目时间紧，工作量大，但难度低。还按照以往 java web mvc 三层的开发方式，不能快速的完成，响应客户的要求。
另外，迁移进司X平台后，部署维护变的繁杂，按照以往的方式，每次更新打jar 包，再打docker镜像，过于耗时，消耗人力。
如果能够只是改下数据库中配置表的配置项，就可以配出一个接口，可以快速响应客户需求，同时节省部署消耗。
**本项目目的，只要会 sql 和 json, 就可以开发简单的接口。**

### 1.2 不适合使用的场景
    不适合应用于业务逻辑复杂的项目，无法像编程语言一样灵活。
    不适合应用于大数据量的场景，不是大数据处理程序。
    不适合直接应用于外网的场景，没有健全的权限控制等。

### 1.3 技术选型
**开发语言**   ：**Go**语言 
**orm框架**   ：**gorm**   github.com/jinzhu/gorm 

**数据库驱动**： github.com/jinzhu/gorm/dialects/postgres 

**Http框架**  ：**fasthttp** github.com/valyala/fasthttp 

**Router路由** ：github.com/buaazp/fasthttprouter 

**Session框架**：github.com/phachon/fasthttpsession github.com/phachon/fasthttpsession/memory

**JWT框架**：github.com/dgrijalva/jwt-go

**Redis-Cluster 连接池**： github.com/wuxibin89/redis-go-cluster 的改造版 github.com/zdy02004/redis-go-cluster 

**Memcache 连接池**：github.com/bradfitz/gomemcache/memcache 

**权限管理框架 casbin**：github.com/casbin/casbin 

**casbin 策略持久化库**：github.com/casbin/gorm-adapter 

**javascript 执行解释器**：github.com/robertkrimen/otto 

**图片验证码框架**：github.com/dchest/captcha 

### 1.4 功能概述
    自定义登录，采用单机内存中的 Session 方式，适合浏览器端使用,支持带图片验证码的登录.
	基于casbin实现的RBAC的接口权限管控
	支持 JWT 的 Token 方式登录，适合后端使用.
	支持执行可配置 sql
	支持执行可配置 redis 命令,目前只支持带密码验证的 Redis Cluster 集群
	支持执行可配置 memcache 命令,目前只支持单机模式
	支持执行可配置本地 Shell 命令
	支持执行可配置远端 SSH 命令
	支持AES加密接口服务
	支持AES解密接口服务
	支持发起可配置 GET 请求
	支持发起可配置 POST 请求
	支持发起采用表单方式的文件上传
	支持发起文件下载
	支持代理静态页面,应用于前后端物理分离部署
	支持内嵌 javascript 脚本，进行简单的逻辑处理
## 二、概念
### 2.1 术语
**HTTP请求类型：**  服务端接收响应的http请求类型 **method_type**，1 GET 2 POST

**入参校验：**      **valid**， 对接口调用传入的入参按配置规则校验

**取值：**         **get_value**,根据配置规则，执行 sql、redis、接口调用等方式获得接口返回值

**输出：**         **out_put** 根据配置规则，将 get_value 中的返回值填入返回报文中

### 2.2 执行流程
	根据配置文件配置，连接数据库
	扫描配置表中的所有记录，针对每一条配置项，注册生成对应的接口
	配置表不存在可自动建表
	启动 http 服务，监听配置文件中配置的端口
	接收到接口请求，处理
	入参校验
	取值并声明到内部变量中
	调用javascript 解释器，对内部变量进行逻辑处理
	生成出参报文
### 2.3 配置 json 的一般组成
```json
{
"动作类型"：{
		"变量名"：{
			"属性":"值或取值方式"
			 }
	    }   
}
```
**动作类型：** 指程序的功能，例如，执行sql，执行 redis 命令，执行本地命令等。

**变量：** 自定义变量，将动作的执行结果赋值给该变量。

**属性：** 变量的属性。

**值或取值方式：**   属性的值，或获得该值的方法。

## 三、部署
### 3.1 下载部署
**下载程序 centOS 版，https://github.com/zdy02004/chakala/release/chakala_go.zip

解压 unzip chakala_go.zip 

填写配置文件 conf.cfg 

```bash
#程序所在目录，可不配
program_name=/home/project/chakala_go/chakala_go
#数据库类型
db_type=postgres
#数据库主机
db_host=10.1.xxx.xxx
#数据库接口
db_port=5432
#数据库用户
db_user=
#库名
db_name=public
#库名
db_database=public
#数据库密码
db_password=
#http 监听端口
server_port=xxxx
#http 监听ip
server_ip=10.19.xx.xx
#是否使用 redis
is_use_redis=1
# redis cluster 连接地址
redis_cluster_server=10.1.xxx.xxx:6001;10.1.xxx.xxx:6001;10.1.xxx.xxx:6002;10.1.xxx.xxx:6001;10.1.235.xxx.xxx;127.0.0.1:6001
# redis cluster 密码
redis_cluster_auth=xxxx
# Aes 加密服务密钥
aeskey=3204R3u9y8Y2Uwfl
# 日志是否输出到文件
is_log_file=1
# memcache 连接地址
memcache_server=10.1.xxx.xxx:11211
# JWT 签名
SigningKey=c48iuheO9_l
# JWT token 有效期时间
expire_times=15000
```
```shell
sqlite3 session.db
```

```sql
CREATE TABLE session (
 session_id varchar(64) NOT NULL DEFAULT '',
 contents TEXT NOT NULL,
last_active int(10) NOT NULL DEFAULT '0',
 PRIMARY KEY (session_id)
);

create index last_active on session (last_active);
```
然后进入 postgresql 中执行

```sql
CREATE TABLE chakala_config (
    id INT8 NOT NULL DEFAULT unique_rowid(),   --  id
    created_at TIMESTAMPTZ NULL,               --  创建时间
    updated_at TIMESTAMPTZ NULL,               --  修改时间
    deleted_at TIMESTAMPTZ NULL,               --  删除时间
    name STRING NULL,                          --  url二级域名
    method_type INT8 NULL,                     --  请求类型：1 GET，2 POST,3 下载文件
    valid STRING NULL,                         --  入参校验 json 配置串
    get_value STRING NULL,                      --  取值 json 配置串
    js_script STRING NULL,                      --  javascript 交互脚本
    out_put STRING NULL,                        --  出参报文 json 模板
    is_use INT8 NULL,                           --  是否在用，1 是 ，0 否
    re_mark STRING NULL,                        -- 备注
    CONSTRAINT "primary" PRIMARY KEY (id ASC),  --  主键 id
    INDEX idx_chakala_config_deleted_at (deleted_at ASC),
    UNIQUE INDEX chakala_config_name_key (name ASC), -- 唯一约束 name
    FAMILY "primary" (id, created_at, updated_at, deleted_at, name, method_type, valid, get_value, out_put, is_use，re_mark)
);

CREATE UNIQUE INDEX if not exists ON chakala_config (name);
```

启动程序 nohup ./restart.sh & 

### 3.2 编译环境搭建

安装 go 语言环境参考 https://www.jianshu.com/p/b2222fc04f47

安装依赖库
```shell
go get github.com/jinzhu/gorm
go get github.com/jinzhu/gorm/dialects/postgres
go get github.com/valyala/fasthttp
go get github.com/buaazp/fasthttprouter
go get github.com/phachon/fasthttpsession github.com/phachon/fasthttpsession/memory
go get github.com/dgrijalva/jwt-go
go get github.com/zdy02004/redis-go-cluster
go get github.com/bradfitz/gomemcache/memcache
go get github.com/casbin/casbin
go get github.com/casbin/gorm-adapter
go get github.com/robertkrimen/otto
go get github.com/dchest/captcha
```
下载程序包
git clone https://github.com/zdy02004/chakala.git
然后进入源码文件夹,执行编译
```shell
go build
```

## 四、数据库增删改查

### 4.1 配置表

```sql
CREATE TABLE chakala_config (
    id INT8 NOT NULL DEFAULT unique_rowid(),   --  id
    created_at TIMESTAMPTZ NULL,               --  创建时间
    updated_at TIMESTAMPTZ NULL,               --  修改时间
    deleted_at TIMESTAMPTZ NULL,               --  删除时间
    name STRING NULL,                          --  url二级域名
    method_type INT8 NULL,                     --  请求类型：1 GET，2 POST,3 下载文件
    valid STRING NULL,                         --  入参校验 json 配置串
    get_value STRING NULL,                      --  取值 json 配置串
    js_script STRING NULL,                      --  javascript 交互脚本
    out_put STRING NULL,                        --  出参报文 json 模板
    is_use INT8 NULL,                           --  是否在用，1 是 ，0 否
    re_mark STRING NULL,                        -- 备注
    CONSTRAINT "primary" PRIMARY KEY (id ASC),  --  主键 id
    INDEX idx_chakala_config_deleted_at (deleted_at ASC),
    UNIQUE INDEX chakala_config_name_key (name ASC), -- 唯一约束 name
    FAMILY "primary" (id, created_at, updated_at, deleted_at, name, method_type, valid, get_value, out_put, is_use，re_mark)
);

CREATE UNIQUE INDEX if not exists ON chakala_config (name);

```
### 4.2 配置执行sql，接口为GET方式,不传参数

配置表
id 自行配置
name = 'chakala_list'
method_type = 1
valid 配置 ''
get_value  配置如下：
```json
{
"sql":{
"sql_result":"select '{\"id\": '||'\"'|| cast(id as varchar)||'\",' ||'\"name\": '||'\"'|| cast(name as varchar)||'\",' ||'\"method_type\": '||'\"'|| cast(method_type as varchar)||'\",' ||'\"valid\": '|| cast(valid as varchar)||',' ||'\"get_value\": '|| cast(get_value as varchar)||',' ||'\"out_put\": '||'\"'|| replace(cast(out_put as varchar),'\"','\"')||'\"'||',\"is_use\":'||'\"'|| cast(is_use as varchar)||'\"'||',\"re_mark\":'||'\"'|| cast(COALESCE(re_mark,'') as varchar)||'\"}' as out  from public.chakala_config"
}}
```
**注意：sql语句中，select 的部分要有 as out**

这里说明一下，配置结构
![](http://10.1.235.103:4999/server/../Public/Uploads/2019-09-19/5d82638e7565f.png)

**out_put**  配置如下：

```json
{"sql_result":${sql_result}}
```
is_use 配置为1

#### 参考配置sql

```sql
INSERT INTO chakala_config
SELECT -3, now(), now(), now(), 'chakala_list'
	, 1, '{
"method_type":{
"val":"${method_type}",
"notBlank":"false",
"isNum":"false"
}
}', '{
"sql":{
"sql_result":"select ''''{\"id\": ''''||''''\"''''|| cast(id as varchar)||''''\",'''' ||''''\"name\": ''''||''''\"''''|| cast(name as varchar)||''''\",'''' ||''''\"method_type\": ''''||''''\"''''|| cast(method_type as varchar)||''''\",'''' ||''''\"valid\": ''''|| cast(valid as varchar)||'''','''' ||''''\"get_value\": ''''|| cast(get_value as varchar)||'''','''' ||''''\"out_put\": ''''||''''\"''''|| replace(cast(out_put as varchar),''''\"'''',''''\\\"'''')||''''\"''''||'''',\"is_use\":''''||''''\"''''|| cast(is_use as varchar)||''''\"}'''' as out  from chakala_config"
}
}','', '{"sql_result":${sql_result}}',
-- re_mark
're_mark'
;
```

重启程序后测试执行 ` curl 10.19.xxx.xxx:8000/chakala_list`
![](http://10.1.235.103:4999/server/../Public/Uploads/2019-09-18/5d820e358f2cc.png)

### 4.3 配置执行sql，接口为POST方式传参

name = 'channel'  二级域名为channel
method_type = 2        post方式
valid 配置如下：
```json
{
"userid":{
"val":"${userid}",
"notBlank":"true",
"isNum":"false"
},
"channel_id":{
"val":"${channel_id}",
"notBlank":"true",
"isNum":"true"
},
"channel_id2":{
"val":"${channel_id2}",
"notBlank":"true",
"isNum":"true"
}
}
```
**校验入参 userid 非空，校验入参channel_id和channel_id2非空且为数字**

get_value  配置如下：
```json
{
"sql":{
 "user": " select '{\"name\": '||'\"'|| cast(name as varchar)||'\"' || '}'  as out from public.chakala_config where name = '${userid}'",
 "channel_name":  " select '{\"channel_name\": '||'\"'|| cast(channel_name as varchar)||'\"' || '}'  as out from public.channel_info where id >= '${channel_id}' and id <= '${channel_id2}'"
}
}
```

**注： post中的入参为 userid、channel_id 这两个，会替换get_value中被${}包裹的同名变量**

**out_put**  配置如下：
```json
{
"userid":"${user}",
"channel_names":"${channel_name}"
}
```

**注： output 中这两个变量 ${user}、${channel_name} ，会被get_value 中定义的变量值替换，而不是被入参的的变量替换**

#### 参考配置sql

```json
INSERT INTO chakala_config
SELECT 2, now(), now(), now(), 'channel',2,
'{
"userid":{
"val":"${userid}",
"notBlank":"true",
"isNum":"false"
},
"channel_id":{
"val":"${channel_id}",
"notBlank":"true",
"isNum":"true"
},
"channel_id2":{
"val":"${channel_id2}",
"notBlank":"true",
"isNum":"true"
}
}'
，
'{
"sql":{
 "user": " select ''{\"name\": ''||''\"''|| cast(name as varchar)||''\"'' || ''}''  as out from public.chakala_config where name = ''${userid}''",
 "channel_name":  " select ''{\"channel_name\": ''||''\"''|| cast(channel_name as varchar)||''\"'' || ''}''  as out from public.channel_info where id >= ''${channel_id}'' and id <= ''${channel_id2}''"
}
}'
,'',
'{
"userid":"${user}",
"channel_names":"${channel_name}"
}',1，'';
```

**注：如果使用高版本PostgreSql的话，可以使用 row_2_json 函数将一行结果集转换为json，简化取值中sql的复杂度。参考4.7 的例子**

重启程序后测试执行，如果使用脚本 test.sh 启动的程序则会自动重启

```shell
curl 10.19.xxx.xxx:8000/channel -X POST -H "Content-Type:application/json" -d '{"userid":"name" , "channel_id":"471578910","channel_id2":"471579910"}'
```

![](http://10.1.235.103:4999/server/../Public/Uploads/2019-09-19/5d82617b8c89a.png)

### 4.4 配置执行Insert语句，接口为POST方式传参

name = 'chakala_insert'  二级域名为 chakala_insert
method_type = 2        post方式
valid 配置如下：
```json
{
"method_type":{
"val":"${method_type}",
"notBlank":"true",
"isNum":"true",
"<":"2",
">":"0"
},

"is_use":{
"val":"${is_use}",
"notBlank":"true",
"isNum":"true"
}
}
```

**校验入参 method_type、is_use 非空、且为数字，其中 method_type取值范围大于0且小于2**

get_value  配置如下：
```json
{
"sql":{
"sql_result":"insert into chakala_config select nextval('seq_chakala') ,now(),now(),now(),'${name}','${method_type}','${valid}','${get_value}','${js_script}','${out_put}','${is_use}','${re_mark}'"
}
}
```

**out_put**  配置如下：

```json
{
"sql_result":"${sql_result}"
}
```

#### 参考配置sql

```sql
insert into chakala_config
select -1,now(),now(),now(),'chakala_insert',2,
-- valid:
'{
"method_type":{
"val":"${method_type}",
"notBlank":"true",
"isNum":"true",
"<":"2",
">":"0"
},

"is_use":{
"val":"${is_use}",
"notBlank":"true",
"isNum":"true"
}
}',
-- get_value:
'{
"sql":{
"sql_result":"insert into chakala_config select nextval(''seq_chakala'') ,now(),now(),now(),''${name}'',''${method_type}'',''${valid}'',''${get_value}'',''${js_script}'',''${out_put}'',''${is_use}'',''${re_mark}''"
}
}',
-- js_script:
'',
-- out_put:
'{
"sql_result":"${sql_result}"
}',
1,'chakala_insert'
```
重启程序后测试执行，如果使用脚本 test.sh 启动的程序则会自动重启

**注意：程序会检查前缀为 chakala 的二级域名的接口，一旦这样的接口执行，则会优雅的关闭自己，然后由 test.sh 脚本自动拉起。**

POST http://10.19.xxx.xxx:8000/chakala_insert
```json
{
    "name": "try_insert7",
    "method_type": "1",
    "is_use": "1",
    "valid": "{ \"userid\": { \"val\": \"${userid}\",\"notBlank\": \"true\", \"isNum\": \"true\",\"<\": \"13\", \">\": \"4\", \">=\": \"4\",\"<=\": \"7\",\"==\": \"\", \"isPast\": \"\",    \"isFuture\": \"\",\"pattern\": \"\"   }}",
    "get_value": "{ \"sql\": { \"user\": \"select ''{\\\"name\\\": ''||''\\\"''|| name||''\\\"'' || ''}''  as out from public.chakala_config where id = ''${userid}''\" } }",
"js_script":"",
    "out_put": "{\"userid\": \"${user}\" }",
    "re_mark": "try_insert7"
}
```
再查询数据库，可以验证 try_insert7的接口添加成功。
curl 10.19.xx.xxx:8000/try_insert7?userid=5

**注意：json value 中的双引号要用 \ 转义；value 中的单引号要用两个单引号转义；sql 语句中的\ 转义要用三个\ 转义，转义后为 \" 形式插入配置表，如果插入的不是自己的配置表，而是其他表则不需要 三个\ 转义。**

### 4.5 配置执行 delete 语句，接口为 GET 方式传参

name = 'chakala_delete'  二级域名为 chakala_delete
method_type = 1        get 方式
valid 配置如下：
```json
{
"method_type":{
"val":"${method_type}",
"notBlank":"true",
"isNum":"false"
}
}
```
get_value  配置如下：
```json
{
"sql":{
"sql_result":"delete from chakala_config where name='${name}'"
}
}
```

#### 参考配置sql

```sql
insert into chakala_config
select -2,now(),now(),now(),'chakala_delete',1,
-- valid:
'{
"method_type":{
"val":"${method_type}",
"notBlank":"true",
"isNum":"false"
}
}',
-- get_value:
'{
"sql":{
"sql_result":"delete from chakala_config where name=''${name}''"
}
}','',
-- out_put:
'{"sql_result":"${sql_result}"}',
1,'';
```
重启程序后测试执行，如果使用脚本 test.sh 启动的程序则会自动重启
```shell
curl "10.19.xxx.xxx:8000/chakala_delete?name=try_insert31"
```
再查询数据库，可以验证 try_insert31 接口被删除。

### 4.6 配置执行 update 语句，接口为 POST 方式传参

name = 'chakala_upsert'  二级域名为 chakala_upsert
method_type = 2        post 方式
valid 配置如下：
```json
{
"method_type":{
"val":"${method_type}",
"notBlank":"true",
"isNum":"true",
"<=":"2",
">":"0"
},

"is_use":{
"val":"${is_use}",
"notBlank":"true",
"isNum":"true"
}
}
```
get_value  配置如下：
```json
{
"sql":{
"sql_result":"insert into chakala_config (id,created_at,updated_at,deleted_at,name,method_type,valid,get_value,out_put,is_use) values ( nextval('seq_chakala') ,now(),now(),now(),'${name}','${method_type}','${valid}','${get_value}','${out_put}','${is_use}' ) on  CONFLICT(name) do update set name=excluded.name,method_type=excluded.method_type,valid=excluded.valid,get_value=excluded.get_value,out_put=excluded.out_put,is_use=excluded.is_use,re_mark=excluded.re_mark"
}}
```

#### 参考配置sql

```sql
insert into chakala_config
select -4,now(),now(),now(),'chakala_upsert',2,
-- valid:
'{
"method_type":{
"val":"${method_type}",
"notBlank":"true",
"isNum":"false"
}
}',
-- get_value:
'{"sql":{
"sql_result":"insert into chakala_config (id,created_at,updated_at,deleted_at,name,method_type,valid,get_value,js_script,out_put,is_use,re_mark) values ( nextval(''seq_chakala'') ,now(),now(),now(),''${name}'',''${method_type}'',''${valid}'',''${get_value}'',''${js_script}'',''${out_put}'',''${is_use}'',''${re_mark}'' ) on  CONFLICT(name) do update set name=excluded.name,method_type=excluded.method_type,valid=excluded.valid,get_value=excluded.get_value,js_script=excluded.js_script,out_put=excluded.out_put,is_use=excluded.is_use,re_mark=excluded.re_mark"
}}',
-- js_script:
'',
-- out_put:
'{"sql_result":"${sql_result}"}',
1,
-- remark
'chakala_upsert'
```
重启程序后测试执行，如果使用脚本 test.sh 启动的程序则会自动重启
 post url

    10.19.xx.xxx:8000/chakala_upsert
 post 报文
```json
{
    "name": "try_upsert7",
    "method_type": "1",
    "is_use": "1",
    "valid": "{ \"userid\": { \"val\": \"${userid}\",\"notBlank\": \"true\", \"isNum\": \"true\",\"<\": \"13\", \">\": \"4\", \">=\": \"4\",\"<=\": \"7\",\"==\": \"\", \"isPast\": \"\",    \"isFuture\": \"\",\"pattern\": \"\"   }}",
    "get_value": "{ \"sql\": { \"user\": \"select ''{\\\"name\\\": ''||''\\\"''|| name||''\\\"'' || ''}''  as out from public.chakala_config where id = ''${userid}''\" } }",
    "js_script":"",
    "out_put": "{\"userid\": \"${user}\" }",
    "re_mark": "try_insert711"
}
```

重启程序后测试执行，如果使用脚本 test.sh 启动的程序则会自动重启
 post url

    10.19.xxx.xxx:8000/chakala_upsert
 post 报文
```json
{
    "name": "try_upsert7",
    "method_type": "1",
    "is_use": "1",
    "valid": "{ \"userid\": { \"val\": \"${userid}\",\"notBlank\": \"true\", \"isNum\": \"true\",\"<\": \"13\", \">\": \"4\", \">=\": \"4\",\"<=\": \"7\",\"==\": \"\", \"isPast\": \"\",    \"isFuture\": \"\",\"pattern\": \"\"   }}",
    "get_value": "{ \"sql\": { \"user\": \"select ''{\\\"name\\\": ''||''\\\"''|| name||''\\\"'' || ''}''  as out from public.chakala_config where id = ''${userid}''\" } }",
    "js_script":"",
    "out_put": "{\"userid\": \"${user}\" }",
    "re_mark": "try_insert711"
}
```

**注：这里使用了 PostGreSql 的 upsert 特性，不存在记录则插入，存在记录则update。不了解请自行百度 upsert 。**

### 4.7 row_2_json 案例

name = 'test_row_2_json'  二级域名为 test_row_2_json
method_type = 1        get 方式
valid 配置如下：
```json
{ "userid": { "val": "${userid}","notBlank": "false" }}
```
get_value  配置如下：
```json
{ "sql": { "user": "select row_to_json(t)::text as out from ( select id,name,method_type,valid,get_value,js_script,out_put,is_use,re_mark from  public.chakala_config ) t " } }
```
out_put 配置如下：
```json
{"userid": ${user} }
```

**注意： row_to_json(chakala_config)::text 析出的是字符串，所以 out_put 中的value两边不能加双引号**

#### 参考配置sql

```sql
insert into chakala_config
select 24,now(),now(),now(),'test_row_2_json',1,
'{ "userid": { "val": "${userid}","notBlank": "false" }}', -- valid
'{ "sql": { "user": "select row_to_json(t)::text as out from ( select id,name,method_type,valid,get_value,js_script,out_put,is_use,re_mark from  public.chakala_config ) t " } }', -- get_value
'', -- js_script
'{"userid": ${user} }', -- out_put
1,
'test_row_2_json'
```

## 五、入参校验 

上一章中已经见到几个入参校验的例子。入参校验即可校验get请求中url带的参数，也可校验post请求中 json报文中带的参数，无区别。本章系统的介绍。

### 支持的校验方式

  非空校验
	是否为数字
	小于某个值
	小于等于某个值
	大于某个值
	大于等于某个值
	等于某个值
	是否是过去
	是否是未来
	是否是疑似低级别sql注入
	是否是疑似高级别sql注入
	是否匹配某个正则表达式
	校验是否登录
	校验是否有权限访问
	是否校验 token
  
### json配置结构

```json
{
"userid":{
"val":"${userid}",    -- 入参变量名
"notBlank":"true",	-- 入参是否非空
"isNum":"true",       -- 入参是否位数字
"<":"13",             -- 入参小于某个值
">":"4",              -- 入参大于某个值
">=":"4",			 -- 入参大于等于某个值
"<=":"10",			-- 入参小于等于某个值
"==":"",              -- 入参等于某个值
"isPast":"false",     -- 入参是否是过去
"isFuture":"false",   -- 入参是否是未来
"Sql_low_inject":"true"  -- 入参是否疑似低级别sql注入
"Sql_high_inject":"true"  -- 入参是否疑似高级别sql注入
"pattern":"",         -- 入参是否满足某个正则表达式，正则语法参考go标准库的正则语法
"islogin":"true"      -- 本接口是否必须登录后才能访问
"ispermission":"true" -- 本接口登录后是否有权限访问
"isjwt":"true"        -- 本接口是否校验 token
}}
```

#####本章不再另行举例，例子参考第四章使用。

## 六、执行 redis 指令

### 6.1 配置执行 redis 指令，接口为GET方式

name = 'test_redis'  二级域名为 test_redis
method_type = 1        get 方式
valid 配置如下：
```json
{
"method_type":{
"val":"${method_type}",
"notBlank":"false",
"isNum":"false"
}}
```
get_value  配置如下：
```json
{
"Redis":{
    "user":["GET", "${userid}"]
}}
```

**注：其中 Redis 为 Redis 动作类型
使用方法和 sql 取值类似，只是命令中，用 json 数组 [] 将多个命令连起来。**

out_put 配置如下：
```json
    { "userid" : "${user}" }
```
重启程序后测试执行，
```shell
    curl "10.19.xxx.xxx:8000/try_insert?userid=1"
```
观察 返回值中 GET 的返回值

## 七、执行 Memcache 指令

### 7.1 配置执行 Memcache put 指令，接口为GET方式

name = 'memcache_put'  二级域名为 memcache_put
method_type = 1        get 方式
valid 配置如下：
```json
{ "method_type":{ "val":"${method_type}","notBlank":"false","isNum":"false"}}
```
get_value  配置如下：
```json
{ "Memcache":{ "userid":{"set":"${userid}"}}}
```
使用方法和 sql 取值类似，只是命令中，用 json 数组 [] 将多个命令连起来。
out_put 配置如下：
```json
{"userid": "${userid}" }
```
重启程序后测试执行，
```shell
    curl "10.19.xxx.xxx:8000/memcache_put?userid=1"
```
观察 返回值中 GET 的返回值

### 7.2 配置执行 Memcache get 指令，接口为GET方式

name = 'memcache_get'  二级域名为 memcache_get
method_type = 1        get 方式
valid 配置如下：
```json
{ "method_type":{ "val":"${method_type}","notBlank":"false","isNum":"false"}}
```
get_value  配置如下：
```json
{ "Memcache":{ "userid":{"get":"${userid}"}}}
```
使用方法和 sql 取值类似，只是命令中，用 json 数组 [] 将多个命令连起来。
out_put 配置如下：
```json
{"userid": "${userid}" }
```
重启程序后测试执行，
```shell
    curl "10.19.xxx.xxx:8000/memcache_get?userid=1"
```
观察 返回值中 GET 的返回值

### 7.3 get_value 支持的 Memcache 指令

```json
{
"Memcache":{
"user":{
"get":"${userid}",      -- 或者
"set":"${userid}",      -- 或者
"add":"${userid}",    -- 或者
"replace":"${userid}",  -- 或者
"delete":"${userid}",  -- 或者
"incrby":"${userid}"   -- 或者
}]}}
```

**注：其中 Memcache 为 Memcache 动作类型
每个变量下只能支持一种 get 或 set 或其他，各指令不能同时存在**

## 八、执行本地 Shell 指令

### 8.1 配置执行 本地 Shell 指令，接口为GET方式

name = 'shell'  二级域名为 shell
method_type = 1        get 方式
valid 配置如下：
```json
{
"method_type":{
"val":"${method_type}",
"notBlank":"false",
"isNum":"false"
}}
```
get_value  配置如下：
```json
{
"Shell":{
"user":"${userid}"
}}
```
out_put 配置如下：
```json
{ "userid" : "${user}" }
```
重启程序后测试执行，使用浏览器访问 
   http://10.19.xxx.xxx:8000/shell?userid=ls -l
观察返回值中
![](http://10.1.235.103:4999/server/../Public/Uploads/2019-09-19/5d8328d1ce791.png)

## 九、执行远端 ssh 指令

### 9.1 配置执行，连接远端 ssh，执行指令，接口为GET方式

name = 'ssh_test'  二级域名为 ssh_test
method_type = 1        get 方式
valid 配置如下：
```json
{ "method_type":{ "val":"${method_type}","notBlank":"false","isNum":"false"}}
```
get_value  配置如下：
```json
{
"SSH":{
"user":{
"user":"work",
"pwd": "XXXXXSWXXXIDFXXXXj/PbgXX",
"addr": "10.1.xxx.xxx:22",
"cmd":"${userid}"
}
}}
```

**注：其中 SSH 为 ssh 动作类型
	  "user"："work" 是配置 ssh 登录用户名
     "pwd": "XXXXXSWXXXIDFXXXXj/PbgXX" 是配置ssh登录密码，密码是在服务端会用 aes 算法解密，密钥为配置文件中的 aeskey
     "addr": "10.1.xxx.xxx:22" 是配置 ssh 登录主机ip和端口
     "cmd":"${userid}" 是在远端执行的shell命令**

out_put 配置如下：
```json
{"userid" : "${user}" }
```
重启程序后测试执行，使用浏览器访问
   http://10.19.xxx.xxx:8000/ssh_test?userid=free -g
观察返回值

## 十、AES 加密接口服务

### 10.1 配置执行，加密接口服务，接口为 POST 方式

name = 'test_aesEncrypt'  二级域名为 test_aesEncrypt
method_type = 2           post 方式
valid 配置如下：
```json
{ "method_type":{ "val":"${method_type}","notBlank":"false","isNum":"false"}}
```
get_value  配置如下：
```json
{
"aesEncrypt":{
"user":"${userid}"
}}
```

**注：其中 aesEncrypt 表示加密动作类型
     "user"： 为变量
     ${userid}： 为要加密的值
	 **

out_put 配置如下：
```json
{"userid" : "${user}" }
```
重启程序后测试执行，使用 post 访问

## 十一、AES 解密接口服务

### 11.1 配置执行，解密接口服务，接口为 POST 方式

name = 'test_aesDecrypt'  二级域名为 test_aesDecrypt
method_type = 2           post 方式
valid 配置如下：
```json
{ "method_type":{ "val":"${method_type}","notBlank":"false","isNum":"false"}}
```
get_value  配置如下：
```json
{
"aesDecrypt":{
"user":"${userid}"
}}
```

**注：其中 aesDecrypt 表示加密动作类型
     "user"： 为变量
     ${userid}： 为要解密的值
	 **

out_put 配置如下：
```json
{"userid" : "${user}" }
```
重启程序后测试执行，使用 post 访问

## 十二、服务端调用其他 GET 接口

### 12.1 配置执行，服务端调用其他 GET 接口，GET 方式

name = 'get_test'  二级域名为 get_test
method_type = 1           get 方式
valid 配置如下：
```json
{
"method_type":{
"val":"${method_type}",
"notBlank":"false",
"isNum":"false"
}}
```
get_value  配置如下：
```json
{
"Get":{
 "user": "http://10.19.xxx.xxx:8001/name?userid=${userid}"
}}
```

**注：其中 Get 表示调用其他get接口服务
     "user"： 为变量
     值为其他的接口地址
	 **

out_put 配置如下：
```json
{ "userid" : "${user}" }
```
重启程序后测试执行，使用浏览器访问
http://10.19.xxx.xxx:8000/get_test?userid=1

## 十三、服务端调用其他 POST 接口

### 13.1 配置执行，服务端调用其他 POST 接口，POST 方式

name = 'test_post'  二级域名为 test_post
method_type = 2           post 方式
valid 配置如下：
```json
{
"method_type":{
"val":"${method_type}",
"notBlank":"false",
"isNum":"false"
}}
```
get_value  配置如下：
```json
{
"Post":{
 "user": {
     "url" : "http://10.19.xxx.xxx:8001/channel",
     "body" :
	 "{
	 \"userid\":\"${userid}\",
	 \"channel_id\":\"${channel_id}\",
	 \"channel_id2\":\"${channel_id2}\"
	 }"
     }
}}
```

**注：其中 Post 表示调用其他 Post 接口服务
url 为其他 POST 接口的地址
body 为其他 POST 接口的入参报文
	 **

out_put 配置如下：
```json
{
"userid" : "${user}",
"channel_names" : "${channel_name}"                                                 }
```
重启程序后测试执行，使用post访问
```shell
 curl 10.19.xxx.xxx:8000/test_post -X POST -H "Content-Type:application/json" -d '{"userid":"name" , "channel_id":"471578910","channel_id2":"471579910"}'
```

## 十四、代理静态 html 页面

### 14.1 配置代理静态html页面

name = 'test_html'  二级域名为 test_html
method_type = 1           get 方式
valid 配置为'' 空
get_value  配置如下：
```json
{
"html":{
"user":"/home/project/chakala_go/html"
}}
```

**注：其中 html 表示代理静态html页面
     "user"： 为变量
     值为代理静态页面所在的目录
	 **
out_put 配置 '' 空


在服务器 /home/project/chakala_go/html 下编写html文件 aa.html
```html
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>try1</title>
</head>
<body>
<h2>Norwegian Mountain Trip</h2>
<img border="0" src="https://www.runoob.com/images/pulpit.jpg" alt="Pulpit rock" width="304" height="228">
</body>
</html>
```

重启程序后测试执行，使用浏览器访问
10.19.xxx.xxx:8000/test_html/aa.html

## 十五、登录

目前支持基本的登录的操作,采用 session 方式
在 get_value 中，变量名叫 #login 的变量一单初始化为非空，程序会给客户端申请一个 session
变量名叫 #logout 的变量一单初始化为非空，程序会释放给客户端申请一个 session

### 15.1 登入

name = 'test_login'  二级域名为 test_login
method_type = 1           get 方式
valid 配置为
```json
{
"userid":
{ "val": "${name}",
"notBlank": "true",
"isNum": "false",
"<": "13",
">": "4",
">=": "4",
"<=": "7",
"==": "",
"isPast": "",
"isFuture": ""
"pattern": ""
}}
```
get_value  配置如下：
```json
{
"sql": {
"#login": "select cast(coalesce(id,0) as varchar)  as out from \"User\" where name = '${name}' and member_number ='${phone}' "
}}
```

**注：其中 #login 表示系统规定的登录变量,该变量返回值会存入 mem_session 中，作为权限验证的角色.
其中 #sys_login 表示系统规定的特权登录变量,该变量返回值会存入sqlit3 数据库中的 session.db中作为可持久化 session ，作为权限验证的角色
	 **
out_put 配置
```json
{"user": "${#login}"}
```
重启程序后测试执行，使用浏览器访问
```xml
10.19.xxx.xxx:8000/test_login?name=a&phone=13893654380
```
打开浏览器开发者模式查看 cookie


### 15.2 登出
name = 'test_logout'  二级域名为 test_logout
method_type = 1           get 方式
valid 配置为
```json
{
"method_type":{
"val":"${method_type}",
"notBlank":"false",
"isNum":"false"
}}
```
get_value  配置如下：
```json
{
"sql": {
"#logout": "select 1 "
}}
```

**注：其中 #logout 表示系统规定的登录变量.
其中 #sys_logout 表示系统规定的特权登录变量.
	 **
   
out_put 配置
```json
{"user": "${#logout}"}
```

### 15.3 校验接口登录过才可以访问
valid 配置为
```json
{
"method_type":{
"val":"${method_type}",
"islogin":"true",
}}
```
**注：其中 valid 中变量的属性 islogin 为"true",接口会在入参校验阶段检查客户端是否登录过，登录过可继续调用，未登录则不准调用。**
### 15.4 带验证码的登录
name = 'test_login_with_png'  二级域名为 test_logout
method_type = 1           get 方式
valid 配置为
```json
{
"method_type":{
"val":"${method_type}",
"notBlank":"false",
"isNum":"false"
}}
```
get_value  配置如下：
```json
{ "sql": { "#login_with_png": "select coalesce(name,'0')  as out from chakala_config where name = '${name}' " } }
```
**注：其中 #login_with_png 表示系统规定的带验证码的登录变量.
其中 #sys_login_with_png 表示系统规定的特权带验证码的登录变量.
	 **	
out_put 配置
```json
	{"user": "${#login_with_png}"}
```

测试方法
1 访问 http://10.19.xx.xxx:8000/getCaptchaId 获得图片id
{"CaptchaId" :"nOb1mEy06EE0jcEeBV63"}
2 显示图片
http://10.19.xx.xxx:8000/png?id=nOb1mEy06EE0jcEeBV63.png
3 登录的时候带上 id captchaId，和 图片读出的验证码 captchaSolution
http://10.19.xx.xxx:8000/test_login_with_png?name=test_login_with_png&captchaId=nOb1mEy06EE0jcEeBV63&captchaSolution=000048


## 十六、上传文件
name = 'try_upload'  二级域名为 try_upload
method_type = 3      表单上传方式
valid 配置为
```json
{"method_type":{ "val":"${method_type}","notBlank":"false","isNum":"false"}}
```
get_value 配置为
```json
{"Upload":{"key":"key","path":"/home/project/chakala_go"}}
```
**注：其中 Upload 表示上传动作类型，key为form表单汇总的key值，path 为上传文件的存储目录。**
Content-Type 为 application/x-www-form-urlencoded
form-data 为 file 类型



**注：最大支持上传 200M 大小的文件**
## 十七、下载文件
name = 'test_donwload'  二级域名为 test_donwload
method_type = 1      GET 请求
valid 配置为
```json
{ "method_type":{ "val":"${filename}","notBlank":"true","isNum":"false"}}
```
get_value 配置为
```json
{
"Download":{
"key":"${filename}",
"path":"/home/project/chakala_go/"
}}
```
参考配置SQL
```sql
insert into public.chakala_config select 17,now(),now(),now(),'test_donwload',1,'{ "method_type":{ "val":"${filename}","notBlank":"true","isNum":"false"}}','{
"Download":{
"key":"${filename}",
"path":"/home/project/chakala_go/"
}}','','',1，''；
```
**注：其中 Download 表示下载动作类型，key为要下载的文件名，path 为下载文件的存储目录。**
response 的 Content-Disposition 为 "attachment; filename=xxxx"
##### 下载成功！

## 十八、JWT鉴权
### 18.1 申请 token
name = 'request_token'  二级域名为 request_token
method_type = 1      GET 请求
valid 配置为
```json
{ "userid": { "val": "${name}","notBlank": "true"  }}
```
get_value 配置为
```json
{ "sql":
{ "#request_jwt": "select cast(coalesce(id,0) as varchar)  as out from \"User\" 
where name = '${name}' and member_number ='${phone}' " 
}}
```
**注：其中 #request_jwt 表示系统规定的申请token专用变量
	 **
out_put 配置为
```json
{"token": "${#request_jwt}"}
```
参考配置sql
```sql
insert into public.chakala_config
select 18,
now(),
now(),
now(),
'request_token',
1,
'{ "userid": { "val": "${name}","notBlank": "true", "isNum": "false","<": "13", ">": "4", ">=": "4","<=": "7","==": "", "isPast": "",    "isFuture": "","pattern": ""   }}',
'{ "sql": { "#request_jwt": "select cast(coalesce(id,0) as varchar)  as out from \"User\" where name = ''${name}'' and member_number =''${phone}'' " } }','',
'{"token": "${#request_jwt}"}',
1,'';
```
重启程序后测试执行，使用GET访问


### 18.2 token 鉴权
name = 'check_token'  二级域名为 check_token
method_type = 1      GET 请求
valid 配置为
```json
{ "userid": { "val": "${token}","notBlank": "true", "isjwt": "true"}}
```

**注：其中 #isjwt 表示系统规定的鉴权 token 专用属性**

get_value 配置为
```json
{ "sql": { "request": "select 1 as out from \"User\"  limit 1 " } }
```
out_put 配置为
```json
{"result": "${request}"}
```
参考配置sql
```sql
insert into public.chakala_config
 select
 19,
 now(),
 now(),
 now(),
 'check_token',
 1,
 '{ "userid": { "val": "${token}","notBlank": "true", "isjwt": "true"}}',
 '{ "sql": { "request": "select 1 as out from \"User\"  limit 1 " } }','',
 '{"result": "${request}"}',
 1, '';
```
重启程序后测试执行，使用GET访问


## 十九、接口访问权限

### 19.1 权限表

权限表 casbin_rule 会自行创建。
```sql
 CREATE TABLE casbin_rule (
     p_type VARCHAR(100) NULL,  -- 类型 ‘p’ 代表赋权 ，‘g’代表组定义
     v0 VARCHAR(100) NULL,	--用户名 或 角色 或 组
     v1 VARCHAR(100) NULL, --资源
     v2 VARCHAR(100) NULL, -- 动作
     v3 VARCHAR(100) NULL,
     v4 VARCHAR(100) NULL,
     v5 VARCHAR(100) NULL,
     FAMILY "primary" (p_type, v0, v1, v2, v3, v4, v5, rowid)；
```
权限表配置参考casbin 策略的配置 https://casbin.org/docs/en/rbac
这里配置 shell 接口要验证权限

name = 'shell'  二级域名为 shell
method_type = 1      GET 请求
valid 配置为
```json
{
"method_type":{
"islogin":"true",
"ispermission":"true",
"val":"${method_type}",
"notBlank":"false",
"isNum":"false"
}}
```
get_value 配置为
```json
{
"Shell":{
"user":"${userid}"
}
}
```
out_put 配置为
```json
{ "userid" : "${user}" }
```
参考权限配置sql
```sql
insert into casbin_rule
select 'p','a','shell','read',NULL,NULL union all -- 用户 a 可以访问 shell
select 'p','data2_admin','data2','read',NULL,NULL union all  -- 组 data2_admin 可以读 read
select 'p','data2_admin','data2','write',NULL,NULL union all -- 组 data2_admin 可以写 write
select 'g','a','data2_admin',NULL,NULL ,NULL;  -- a 属于组 data2_admin
```
注： 更改权限表后应重启程序，目前不支持动态加载

重启程序后测试执行，先登录
http://10.19.xx.xx:xxxx/test_login?name=a&phone=13893654380
再访问接口shell
http://10.19.xx.xx:xxxx/shell?userid=ls
重启浏览器后，修改 ispermission为false 再 登录后访问接口shell

### 19.2 添加权限

name = 'add_policy'  二级域名为 add_policy
method_type = 1      GET 请求
valid 配置为
```json
{ "param1": { "val": "${param1}","notBlank": "true"},"param2": { "val": "${param2}","notBlank": "true"},"param3": { "val": "${param3}","notBlank": "true"} }
```

get_value 配置为
```json
{"add_policy":{"user":["${param1}", "${param2}", "${param3}"]}}
```

**注：其中 add_policy 表示系统规定的添加权限动作**
out_put 配置为
```json
{"user": "${user}"}
```
参考配置sql
```sql
insert into public.chakala_config 
select 20,now(),now(),now(),
'add_policy',
1,
'{ "param1": { "val": "${param1}","notBlank": "true"},"param2": { "val": "${param2}","notBlank": "true"},"param3": { "val": "${param3}","notBlank": "true"} }',
'{"add_policy":{"user":["${param1}", "${param2}", "${param3}"]}}',
'{"user": "${user}"}' ,'',
1,
'add_policy'
;
```

### 19.3 删除权限

name = 'delete_policy'  二级域名为 delete_policy
method_type = 1      GET 请求
valid 配置为
```json
{ "param1": { "val": "${param1}","notBlank": "true"},"param2": { "val": "${param2}","notBlank": "true"},"param3": { "val": "${param3}","notBlank": "true"} }
```
get_value 配置为
```json
{"delete_policy":{"user":["${param1}", "${param2}", "${param3}"]}}
```

**注：其中 delete_policy 表示系统规定的删除权限动作**

out_put 配置为
```json
{"user": "${user}"}
```
参考配置sql
```sql
insert into public.chakala_config
select 21,now(),now(),now(),
'delete_policy',
1,
'{ "param1": { "val": "${param1}","notBlank": "true"},"param2": { "val": "${param2}","notBlank": "true"},"param3": { "val": "${param3}","notBlank": "true"} }',
'{"delete_policy":{"user":["${param1}", "${param2}", "${param3}"]}}','',
'{"user": "${user}"}' ,
1,
'delete_policy'
;
```

### 19.4 重新加载权限

当不是通过接口调整权限，而是通过表修改权限时，需要重现加载表
name = 'load_policy'  二级域名为 load_policy
method_type = 1      GET 请求
valid 配置为
```json
{ "param1": { "val": "${param1}","notBlank": "false"} }
```
get_value 配置为
```json
{"load_policy":{"aa":"1"}}
```

**注：其中 load_policy 表示系统规定的删除权限动作**

out_put 配置为
```json
{"load_policy": "${load_policy}"}
```
参考配置sql
```sql
insert into public.chakala_config 
select 22,now(),now(),now(),
'load_policy',
1,
'{ "param1": { "val": "${param1}","notBlank": "false"} }',
'{"load_policy":{"aa":"1"}}',
'{"load_policy": "${load_policy}"}' ,'',
1,
'load_policy'
;
```

## 二十、javascript 嵌入脚本

### 20.1 javascript 嵌入脚本

name = 'test_js'  二级域名为 test_js
method_type = 1      GET 请求
valid 配置为
```json
{ "js_script": { "val": "${js_script}","notBlank": "true", "isNum": "false","<": "13", ">": "4", ">=": "4","<=": "7","==": "", "isPast": "",    "isFuture": "","pattern": ""   }}
```
get_value 配置为
```json
{ "sql": { "name": "select cast(coalesce(name,'0') as varchar)  as out from \"User\" where name = '${name}' and member_number ='${phone}'  " } }
```
js_script 配置为
```json
function main(tab) {
    result = {}
    tab["name"] += ":http://www.baidu.com"
    return result
}
```

**注：其中 main 为主函数，必须写，tab 代表 get_value 中各种取值变量 **

out_put 配置为
```json
{"user": "${name}"}
```
参考配置sql
```sql
insert into public.chakala_config
 select 23,now(),now(),now(),
 'test_js',1,
 '{ "js_script": { "val": "${js_script}","notBlank": "true", "isNum": "false","<": "13", ">": "4", ">=": "4","<=": "7","==": "", "isPast": "",    "isFuture": "","pattern": ""   }}',
 '{ "sql": { "name": "select cast(coalesce(name,''0'') as varchar)  as out from \"User\" where name = ''${name}'' and member_number =''${phone}''  " } }',
 'function main(tab) {
    result = {}
    tab["name"] += ":http://www.baidu.com"
    return result
}',
 '{"user": "${name}"}' ,1,'js_test';
```
测试： http://10.19.xxX.XXX:8000/test_js?name=a&phone=13893654380

## 二十一、设置报文响应头

### 21.1 设置报文响应头

name = 'test_head'  二级域名为 test_head
method_type = 1      GET 请求
valid 配置为
```json
{ "userid":{ "val":"${name}","notBlank":"true", "isNum":"false","<":"13", ">":"4", ">=":"4","<=":"7","==":"", "isPast":"",    "isFuture":"","pattern":""   }}
```
get_value 配置为
```json
{ "Header":{ "head1":{"key":"key1","val":"val1" },"head2":{"key":"key2","val":"val2" }} }
```

**注：其中 Header 表示要设置响应头，变量里面的 key 代表 http header中的 key，变量里面的 val 代表 http heade r中的 val **
out_put 配置为
```json
{"head1": "${head1}"}
```

参考配置sql
```sql
insert into public.chakala_config select 27,
now(),
now(),
now(),
'test_head',
1,
'{ "userid":{ "val":"${name}","notBlank":"true", "isNum":"false","<":"13", ">":"4", ">=":"4","<=":"7","==":"", "isPast":"",    "isFuture":"","pattern":""   }}',
'{ "Header":{ "head1":{"key":"key1","val":"val1" },"head2":{"key":"key2","val":"val2" }} }',
'',
'{"head1":"${head1}"',
1,'';
```
测试： http://10.19.xxX.XXX:8000/test_head

