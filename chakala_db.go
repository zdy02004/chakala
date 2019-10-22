package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/buaazp/fasthttprouter"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/phachon/fasthttpsession"
	"github.com/phachon/fasthttpsession/memory"
	"github.com/phachon/fasthttpsession/sqlite3"
	"github.com/robertkrimen/otto"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"log"
	"github.com/zdy02004/redis-go-cluster"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
	//"sync"
	//"syscall"
	"flag"
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	"net"
	"runtime"
)

// 默认的 session 全局配置
var session = fasthttpsession.NewSession(fasthttpsession.NewDefaultConfig())
var sys_session = fasthttpsession.NewSession(fasthttpsession.NewDefaultConfig())

var gloable_program_name string

type Config struct {
	gorm.Model
	id         int
	Name       string
	MethodType int
	Valid      string
	GetValue   string
	JsScript   string
	OutPut     string
	IsUse      int
	ReMark     string
}

//配置表执行结果集
type record struct {
	Id         int
	Name       string
	MethodType int
	Valid      string
	GetValue   string
	JsScript   string
	OutPut     string
	IsUse      int
	ReMark     string
}

//裸sql执行结果集
type RawSql struct {
	Out string
}

//校验信息结构体
type err_info struct {
	Key             string `json:key`
	Val             string `json:val`
	IsBlank         bool   `json:isBlank`
	IsNum           bool   `json:isNum`
	Gt              string `json:gt`
	Lt              string `json:lt`
	Ge              string `json:ge`
	Le              string `json:le`
	Eq              string `json:eq`
	IsPast          bool   `json:isPast`
	IsFuture        bool   `json:isFuture`
	Pattern         string `json:pattern`
	Sql_low_inject  string `json:sql_low_inject`
	Sql_high_inject string `json:sql_high_inject`
	IsPass          bool   `json:isPass`
	IsLogin         bool   `json:islogin`
	IsSys           bool   `json:issys`
	IsJwt           bool   `json:isjwt`
	IsPermission    bool   `json:ispermission`
}

//配置json列表 chakala_list 接口专用
type ChakalaList struct {
	Id         int
	Name       string
	MethodType int
	Valid      map[string]map[string]interface{}
	GetValue   map[string]map[string]interface{}
	JsScript   string
	OutPut     map[string]interface{}
	IsUse      int
	ReMark     string
}

//配置json列表 chakala_upsert 接口专用
type ChakalaUpsert struct {
	Id         int
	Name       string
	MethodType int
	Valid      map[string]map[string]interface{}
	GetValue   map[string]map[string]interface{}
	JsScript   string
	OutPut     map[string]string
	IsUse      int
	ReMark     string
}

// 设置环境变量
func setupEnv() {
	// 设置最大的CPU数量
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 设置进程名称
	//gspt.SetProcTitle(constants.Name)
}

var later int64
var later_ string

func init() {

	flag.StringVar(&later_, "later", "defualt", "log in user")

}

const ShellToUse = "bash"

//本地执行命令
func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// request handler
func loginHandle(ctx *fasthttp.RequestCtx, user_name string) (err error) {
	// start session
	sessionStore, err := session.Start(ctx)
	if err != nil {
		//ctx.SetBodyString(err.Error())
		return err
	}
	// 必须 defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionStore.Set("sid", "chakala")
	sessionStore.Set("user_name", user_name)
	return nil
}

func logoutHandle(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := session.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionStore.Delete("sid")
	sessionStore.Delete("user_name")

	s := sessionStore.Get("sid")
	if s == nil {
		return
	}
	ctx.Write([]byte("{\"Error\":\"logout failed\"}\n"))

}

// request handler
func sysloginHandle(ctx *fasthttp.RequestCtx, user_name string) (err error) {
	// start session
	sessionStore, err := sys_session.Start(ctx)
	if err != nil {
		//ctx.SetBodyString(err.Error())
		return err
	}
	// 必须 defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionStore.Set("sid", "chakala")
	sessionStore.Set("user_name", user_name)
	return nil
}

func syslogoutHandle(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := sys_session.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionStore.Delete("sid")
	sessionStore.Delete("user_name")

	s := sessionStore.Get("sid")
	if s == nil {
		return
	}
	ctx.Write([]byte("{\"Error\":\"logout failed\"}\n"))

}

//实际主函数
func real_main() {

	f, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() //设置日志输出到 f

	log.Println("Start real_main() ...... ")
	//读取配置文件 conf.cfg
	configMap := InitConfig("conf.cfg")
	//读取启动命令———绝对路径
	program_name := configMap["program_name"]
	gloable_program_name = program_name

	//是否开启独立日志文件
	is_log_file := configMap["is_log_file"]
	if is_log_file == "1" {
		log.SetOutput(f)
	}
	//读取数据库配置
	host := configMap["db_host"]
	db_type := configMap["db_type"]
	db_port := configMap["db_port"]
	db_user := configMap["db_user"]
	db_database := configMap["db_database"]
	password := configMap["db_password"]
	//读取加解密秘钥
	aeskey := configMap["aeskey"]
	//读取 memcache 连接池
	memcache_server := configMap["memcache_server"]
	//读取 jwt 签名和实效时间
	SigningKey := configMap["SigningKey"]
	expire_times := configMap["expire_times"]
	// http 服务监听
	server_port := configMap["server_port"]
	server_ip := configMap["server_ip"]

	listen_url := server_ip + ":" + server_port
	log.Println(listen_url)

	ln, err := net.Listen("tcp4", listen_url)
	if err != nil {
		log.Println("start server fail,because ", err)
		//f.Close()
		time.Sleep(time.Duration(1) * time.Second)
		//c <- 1
		return
		//os.Exit(-5)
	}
	// 优雅的关闭配置
	// 最大等待3秒
	Listener := NewGracefulListener(ln, time.Second*3)
	// 读取 redis 配置
	is_use_redis := configMap["is_use_redis"]
	redis_cluster_server := configMap["redis_cluster_server"]
	redis_cluster_auth := configMap["redis_cluster_auth"]

	var GeneralCluster *redis.Cluster
	//连接 redis cluster
	if is_use_redis == "1" {
		Nodes := strings.Split(redis_cluster_server, ";")
		cluster, errr := redis.NewCluster(
			&redis.Options{
				StartNodes:   Nodes,
				ConnTimeout:  50 * time.Millisecond,
				ReadTimeout:  50 * time.Millisecond,
				WriteTimeout: 50 * time.Millisecond,
				KeepAlive:    16,
				AliveTime:    60 * time.Second,
			}, redis_cluster_auth)

		if errr != nil {
			log.Println("Redis Connect err is ", errr, "auth is ", redis_cluster_auth)
			return
		}
		GeneralCluster = cluster
	}
	//连接DB
	constr := "host=" + host + " port=" + db_port + " user=" + db_user + " dbname=" + db_database + " password=" + password
	db, err := gorm.Open(db_type, constr)
	//defer	db.Close()

	if err != nil {
		log.Println("gorm.Open() failed! beacuse ", err, " try to use sslmode=disable")
		db, err = gorm.Open(db_type, constr+" sslmode=disable")
		if err != nil {
			log.Println("gorm.Open(", constr+" sslmode=disable", " failed! ", err)
			return
		} else {
			log.Println("gorm.Open() success with  sslmode=disable")
		}
	}

	log.Println("************************ DB.Open success *************************")

	//自动建表
	if !db.HasTable("chakala_config") {
		//建表
		db.Table("chakala_config").CreateTable(&Config{})
		db.Raw("CREATE UNIQUE INDEX if not exists ON chakala_config (name)")
	}

	// casbin 权限中间件初始化,模型为 rbac
	adapter, _ := gormadapter.NewAdapterByDB(db)
	effect, err51 := casbin.NewEnforcer("./rbac_model.conf", adapter)
	if err51 != nil {
		log.Println("casbin.NewEnforcer(/rbac_model.conf) failed because ", err51)
	}
	effect.LoadPolicy()

	var records []record
	//获得配置记录
	db.Table("chakala_config").Select("id,name,method_type,valid,get_value,js_script,out_put,is_use").Where("is_use > ?", -1).Scan(&records)

	// chakala_list 接口专用数据结构
	length := len(records)
	var chakala_list []ChakalaList = make([]ChakalaList, length)
	for line_index, one_line := range records {
		chakala_list[line_index].Id = one_line.Id
		chakala_list[line_index].Name = one_line.Name
		chakala_list[line_index].MethodType = one_line.MethodType
		chakala_list[line_index].JsScript = one_line.JsScript
		chakala_list[line_index].IsUse = one_line.IsUse
		chakala_list[line_index].ReMark = one_line.ReMark
		err = json.Unmarshal([]byte(one_line.Valid), &(chakala_list[line_index].Valid))
		err = json.Unmarshal([]byte(one_line.GetValue), &(chakala_list[line_index].GetValue))
		err = json.Unmarshal([]byte(one_line.OutPut), &(chakala_list[line_index].OutPut))

	}

	router := fasthttprouter.New()

	//针对每一个配置项
	for line_index, one_line := range records {
		log.Print("record:", line_index)
		//log.Println("=====================================================================================\n")
		//log.Println(one_line)
		//log.Println("=====================================================================================\n")

		var method_type int
		var valid string
		var get_value string
		var out_put_ string
		var jscript string

		method_type = one_line.MethodType
		valid = one_line.Valid
		get_value = one_line.GetValue
		out_put_ = one_line.OutPut
		jscript = one_line.JsScript

		//入参字符串转map
		html_map := make(map[string]map[string]interface{})
		err02 := json.Unmarshal([]byte(get_value), &html_map)
		if err02 == nil {
			exec_html_map := html_map["html"]
			if len(exec_html_map) > 0 {
				log.Println("-----------------------1.-2、Exec html --------------------\n")
				log.Println("exec_html_map :==>> ", exec_html_map)
				for k, v := range exec_html_map {
					log.Println("k:v=", k, ":", v)
					router.ServeFiles("/"+one_line.Name+"/*filepath", v.(string))
				}
			}
		} else {
			log.Println("Unmarshal html_map err:==>> ", err02)
		}

		//http.HandleFunc("/"+one_line.Name, func(w http.ResponseWriter, r *http.Request) {
		var real_method string
		if one_line.MethodType == 1 {
			real_method = "GET"
		} else {
			real_method = "POST"
		}

		name_route := one_line.Name
		log.Println(name_route)
		//log.Println(one_line.OutPut)
		//注册Controller
		router.Handle(real_method, "/"+name_route, func(ctx *fasthttp.RequestCtx) {

			valid_map := make(map[string]map[string]interface{})
			get_value_map := make(map[string]map[string]interface{})
			params_map := make(map[string]string)
			out_put := out_put_

			ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
			ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type") //header的类型
			ctx.Response.Header.Set("content-type", "application/json")             //返回数据格式是json
			log.Println(one_line.Name)
			log.Println(method_type)
			log.Println("======================== 1、Get Input Params ================================\n")
			//解析get中入参，存为 params_map
			if method_type == 1 {
				ctx.URI().QueryArgs().VisitAll(func(k []byte, value []byte) {
					str := (*string)(unsafe.Pointer(&k))
					params_map[*str] = string(value)
					log.Println(params_map[*str])
				})
			} else if method_type == 2 {
				//解析post中的入参，存为 params_map
				get_params_map := make(map[string]interface{})

				err4 := json.Unmarshal(ctx.PostBody(), &get_params_map)
				if err4 != nil {
					log.Println("Post Body is not json")
					ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
					ctx.Write([]byte("{\"Error\":\"Post Body is not json\"}\n"))
					return
				}
				//处理创建配置表的情况——前缀带有 chakala 的，将 value 字符串化，避免深度解析
				for k, v := range get_params_map {
					if strings.HasPrefix(name_route, "chakala") {

						_str_v, err := json.Marshal(v)
						if err != nil {
							log.Println("json.Marshal(", v, ") is bad")
							ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
							ctx.Write([]byte("{\"Error\":\"json to string failed\"}\n"))
							return
						}
						str_len := len(_str_v)
						str_v := string(_str_v)[1 : str_len-1]
						log.Println("Marshal ", k, ": ", str_v)
						params_map[k] = str_v
					} else {
						var istrans bool
						params_map[k], istrans = v.(string)
						//转型失败
						if istrans == false {
							log.Printf(params_map[k], " assagned failed beacuse of ", err4)
							ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
							ctx.Write([]byte("{\"Error\":\"" + params_map[k] + "  assagned failed beacuse of " + err4.Error() + "\"}\n"))
							return
						}
					}

				}
			}

			log.Println("params_map is :==>> ", params_map)

			//if strings.HasPrefix(name_route, "chakala_upsert") {
			//log.Println("chakala_upsert fix ")
			//var chakala_one ChakalaUpsert
			//err4 := json.Unmarshal(ctx.PostBody(), &chakala_one)
			//	if err4 != nil {
			//		log.Println("Post Body is not json")
			//		ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
			//		ctx.Write([]byte("{\"Error\":\"Post Body is not json\"}\n"))
			//		return
			//	}
			//log.Println("chakala_one: ",chakala_one)
			// valid_ ,err41 := json.Marshal(  chakala_one.Valid )
			//		if err41 != nil {
			//		log.Println("chakala_one.Valid is not json")
			//		ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
			//		ctx.Write([]byte("{\"Error\":\"chakala_one.Valid is not json\"}\n"))
			//		return
			//	}
			// get_value_ ,err42 := json.Marshal(  chakala_one.GetValue )
			//		if err42 != nil {
			//		log.Println("chakala_one.GetValue is not json")
			//		ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
			//		ctx.Write([]byte("{\"Error\":\"chakala_one.GetValue is not json\"}\n"))
			//		return
			//	}
			//	log.Println("chakala_one.GetValue: ", chakala_one.GetValue )
			// out_put_ ,err43 := json.Marshal(  chakala_one.OutPut )
			//	if err43 != nil {
			//		log.Println("chakala_one.OutPut is not json")
			//		ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
			//		ctx.Write([]byte("{\"Error\":\"chakala_one.OutPut is not json\"}\n"))
			//		return
			//	}
			//	log.Println("chakala_one.OutPut: ", chakala_one.OutPut )
			//
			//
			//	params_map["id"] = strconv.Itoa(chakala_one.Id)
			//	params_map["Name"] =  chakala_one.Name
			//	params_map["MethodType"] = strconv.Itoa(chakala_one.MethodType)
			//
			//  params_map["js_script"] = chakala_one.JsScript
			//	params_map["is_use"] =  strconv.Itoa(chakala_one.IsUse)
			//	params_map["re_mark"] = chakala_one.ReMark
			//
			//  params_map["valid"] = string(valid_)
			//	params_map["get_value"] = string(get_value_)
			//	params_map["out_put"] = string(out_put_)
			//log.Println("After chakala_upsert fix,the params_map is :==>> ", params_map)
			//}

			get_value_replace := get_value
			valid_replace := valid
			//入参值替换执行和校验脚本
			for k, v := range params_map {
				//入参值替换 get_value
				get_value_replace = strings.Replace(get_value_replace, "${"+k+"}", v, -1)
				//入参值替换 valid
				valid_replace = strings.Replace(valid_replace, "${"+k+"}", v, -1)
			}

			log.Println("valid_replace: ", valid_replace)
			//校验json字符串转map
			err1 := json.Unmarshal([]byte(valid_replace), &valid_map)
			if err1 != nil {
				log.Println("valid_map is not json:==>> ", err1)
				//break
			}
			log.Println("get_value_replace: " + get_value_replace)
			//入参字符串转map
			err2 := json.Unmarshal([]byte(get_value_replace), &get_value_map)

			if err2 != nil {
				log.Println("get_value is not json:==>> ", err2)
				ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
				ctx.Write([]byte("{\"Error\":\"get_value is not json\"}\n"))
				return
			}

			log.Println("get_value_map is ==>>  ", get_value_map)
			log.Println("valid_map is ==>>  ", valid_map)

			log.Println("======================= 1.2、Exec Valid  ================================\n")
			Err_map_info_map := make(map[string](*err_info))
			//遍历校验 map，依次执行校验逻辑
			for k, v := range valid_map {
				Err_map_info := new(err_info)

				log.Println("k:v=", k, ":", v)
				val := v["val"].(string)
				Err_map_info.Val = val
				Err_map_info.Key = k
				Err_map_info.IsPass = true

				//判断是否系统登录
				if v["issys"] != nil && v["issys"] == "true" {
					log.Println(" check if syslogin ")
					// start session
					sessionStore, err := sys_session.Start(ctx)
					if err != nil {
						log.Println("Start sys_session failed because ", err.Error())
						ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
						ctx.Write([]byte("{\"Error\":\"Start sys_session failed\"}\n"))
						return
					}

					//defer sessionStore.Save(ctx)
					s := sessionStore.Get("sid")
					if s == nil {
						log.Println("Please SysLogin first ")
						ctx.Write([]byte("{\"Error\":\"Please SysLogin first\"}\n"))
						ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
						return
					}
					log.Println(" check syslogin Ok")
					//判断是否有访问权限
					if v["ispermission"] != nil && v["ispermission"] == "true" {
						log.Println(" check if permission ")
						user_name := sessionStore.Get("user_name")
						if user_name == nil {
							log.Println("Please Login first ")
							ctx.Write([]byte("{\"Error\":\"Please Login first\"}\n"))
							ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
							return
						}
						ispermit, _ := effect.Enforce(user_name, name_route, "read")
						if ispermit != true {
							log.Println("You are not permissinoned to ", name_route)
							ctx.Write([]byte("{\"Error\":\"You are not permissinoned to " + name_route + "\"}\n"))
							ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
							return
						}

						log.Println(" Permissinoned Ok")
					}

				}

				//判断是否需要登录
				if v["islogin"] != nil && v["islogin"] == "true" {
					log.Println(" check if login ")
					// start session
					sessionStore, err := session.Start(ctx)
					if err != nil {
						log.Println("Start session failed because ", err.Error())
						ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
						ctx.Write([]byte("{\"Error\":\"Start session failed\"}\n"))
						return
					}

					//defer sessionStore.Save(ctx)
					s := sessionStore.Get("sid")
					if s == nil {
						log.Println("Please Login first ")
						ctx.Write([]byte("{\"Error\":\"Please Login first\"}\n"))
						ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
						return
					}
					log.Println(" check login Ok")

					//判断是否有访问权限
					if v["ispermission"] != nil && v["ispermission"] == "true" {
						log.Println(" check if permission ")
						user_name := sessionStore.Get("user_name")
						if user_name == nil {
							log.Println("Please Login first ")
							ctx.Write([]byte("{\"Error\":\"Please Login first\"}\n"))
							ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
							return
						}
						ispermit, _ := effect.Enforce(user_name, name_route, "read")
						if ispermit != true {
							log.Println("You are not permissinoned to ", name_route)
							ctx.Write([]byte("{\"Error\":\"You are not permissinoned to " + name_route + "\"}\n"))
							ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
							return
						}

						log.Println(" Permissinoned Ok")
					}

				}

				//判断是否需要验证 token
				if v["isjwt"] != nil && v["isjwt"] == "true" {
					if len(SigningKey) <= 0 {
						log.Println("len(SigningKey) <= 0 ")
						ctx.Write([]byte("{\"Error\":\"len(SigningKey) <= 0\"}\n"))
						return
					}
					log.Println(" check if token ")
					token := params_map["token"]
					ok, err := jwt_parse(token, SigningKey)

					if err != nil {
						log.Println("Please request token first ")
						ctx.Write([]byte("{\"Error\":\"Please request token first\"}\n"))
						ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
						return
					}
					if ok != true {
						log.Println(" token Invalid ")
						ctx.Write([]byte("{\"Error\":\"token Invalid\"}\n"))
						ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
						return
					}
					log.Println(" check token Ok")
				}

				//判断非空
				if v["notBlank"] != nil && v["notBlank"] == "true" {
					Length := len(val)
					if !(Length > 0) {
						Err_map_info.IsPass = false
						log.Println("key is Blank ")
					}
				}
				//判断是否是数字类型
				if v["isNum"] != nil && v["isNum"] == "true" {
					Pattern := "\\d+"
					resuLt, _ := regexp.MatchString(Pattern, val)
					if !(resuLt) {
						Err_map_info.IsNum = false
						Err_map_info.IsPass = false
						log.Println("IsNum = false")
					} else {
						//获取需要比较的左值
						lleft, _ := strconv.ParseFloat(val, 64)
						//判断左值是否小于某个值
						if v["<"] != nil && len(v["<"].(string)) > 0 {
							right := v["<"].(string)
							lright, _ := strconv.ParseFloat(right, 64)
							if !(lleft < lright) {
								Err_map_info.IsPass = false
								log.Println("key:", val, " is not < ", right)
							}
						}
						//判断左值是否大于某个值
						if v[">"] != nil && len(v[">"].(string)) > 0 {
							right := v[">"].(string)
							lright, _ := strconv.ParseFloat(right, 64)
							if !(lleft > lright) {
								Err_map_info.IsPass = false
								log.Println("key:", val, " is not > ", right)
							}
						}
						//判断左值是否小于等于某个值
						if v["<="] != nil && len(v["<="].(string)) > 0 {
							right := v["<="].(string)
							lright, _ := strconv.ParseFloat(right, 64)
							if !(lleft <= lright) {
								Err_map_info.IsPass = false
								log.Println("key:", val, " is not <= ", right)
							}
						}
						//判断左值是否大于等于某个值
						if v[">="] != nil && len(v[">="].(string)) > 0 {
							right := v[">="].(string)
							lright, _ := strconv.ParseFloat(right, 64)
							if !(lleft >= lright) {
								Err_map_info.IsPass = false
								log.Println("key:", val, " is not >= ", right)
							}
						}
						//判断左值是否等于某个值
						if v["=="] != nil && len(v["=="].(string)) > 0 {
							right := v["=="].(string)
							lright, _ := strconv.ParseFloat(right, 64)
							if lleft != lright {
								Err_map_info.IsPass = false
								log.Println("key:", val, " is not == ", right)
							}
						}
					}
				} else {
					right := time.Now().Unix()
					time_template := "20060102150405"                                 //外部传入的时间字符串模板
					_lleft, _ := time.ParseInLocation(time_template, val, time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
					lleft := _lleft.Unix()

					//判断时间是否是过去
					if v["isPast"] != nil && v["isPast"] == "true" {
						if lleft < right {
						} else {
							Err_map_info.IsPass = false
							log.Println("key:", val, " is not Past ")
							log.Println("left:", lleft, ", right:", right)
						}
					}
					//判断时间是否是未来
					if v["isFuture"] != nil && v["isFuture"] == "true" {
						if lleft > right {
						} else {
							Err_map_info.IsPass = false
							log.Println("key:", val, " is not Future ")
							log.Println("left:", lleft, ", right:", right)

						}
					}
					//判断正则条件是否满足
					if v["pattern"] != nil && len(v["pattern"].(string)) > 0 {
						Pattern := v["pattern"].(string)
						resuLt, _ := regexp.MatchString(Pattern, val)
						if !(resuLt) {
							Err_map_info.IsPass = false
							log.Println("pattern [", Pattern, "] is not Satisfy ")
						}
					}
					//判断是否是低级别sql注入
					if v["sql_low_inject"] != nil && len(v["sql_low_inject"].(string)) > 0 {
						resuLt := sql_low_inject(val)
						if resuLt != "ok" {
							Err_map_info.IsPass = false
							Err_map_info.Sql_low_inject = resuLt
							log.Println(val, " sql_low_inject ", resuLt)
						}
					}
					//判断是否是高级别sql注入
					if v["sql_high_inject"] != nil && len(v["sql_high_inject"].(string)) > 0 {
						resuLt := sql_high_inject(val)
						if resuLt != "ok" {
							Err_map_info.IsPass = false
							Err_map_info.Sql_high_inject = resuLt
							log.Println(val, " sql_high_inject ", resuLt)
						}
					}
				}
				//汇总校验是否满足所有条件
				if Err_map_info.IsPass == false {
					log.Println("Valid: ", k, " is not Pass")
					result, err := json.Marshal(&Err_map_info)
					if err != nil {
						log.Println(err)
					}
					ctx.Write([]byte("{\"Error\":\"Valid is not Pass\","))
					ctx.Write([]byte("{\"Valid\":\""))
					ctx.Write([]byte(string(result)))
					ctx.Write([]byte("\"}"))
					ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
					return
				} else {
					Err_map_info_map[k] = Err_map_info
				}

			}

			log.Println("======================= 1.3、Get chakala_list  ================================\n")
			if strings.HasPrefix(name_route, "chakala_list") {
				//result_, err := json.Marshal(&records)
				result, err := json.Marshal(&chakala_list)
				ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json

				if err != nil {
					log.Println(err)
					ctx.Write([]byte("{\"Error\":\""))
					ctx.Write([]byte(err.Error()))
					ctx.Write([]byte("\"}"))
					return
				} else {
					ctx.Write([]byte(string(result)))
					return
				}

			}

			get_value_result := make(map[string]string)

			log.Println("======================= 1.4、Get chakala_turn  ================================\n")
			if strings.HasPrefix(name_route, "chakala_turn") {
				get_params_map := make(map[string]string)
				err4 := json.Unmarshal(ctx.PostBody(), &get_params_map)
				if err4 != nil {
					log.Println("Chakala_turn Post Body is not json")
					ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
					ctx.Write([]byte("{\"Error\":\"Chakala_turn Post Body is not json\"}\n"))
					return
				}
				i := 0

				for k, v := range get_params_map {
					log.Println("key: ", k, "val: ", v)

					if v == "1" || v == "0" {
						db.Table("chakala_config").Where("name = ? ", k).Update("is_use", v)
						i = i + 1
					}

				}

				log.Println("ture on/off succuess: ", i)
				get_value_result["chakala_turn"] = strconv.Itoa(i)
				//ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
				//ctx.Write([]byte("{\"ture on/off succuess\": "+ strconv.Itoa(i) +"}\n"))
				//return

			}

			// 获得 get_value 中的值：依次执行map，根据标签执行 SQL、redis cluster、memcahcae、本地 shell 、
			// http get、http set、ssh 、加解密等，获得 get value 中的所有真实值
			log.Println("=======================2、Get Value ================================\n")
			//执行 sql
			log.Println("-----------------------2.1、Exec sql--------------------------------\n")
			var exec_sql_map map[string]interface{}
			exec_sql_map = get_value_map["sql"]

			if exec_sql_map != nil {
				for k, v := range exec_sql_map {
					var rawSqlResult []RawSql
					var thisOneResult string
					//执行裸sql
					log.Println("exec sql:", v)
					err5 := db.Raw(v.(string)).Scan(&rawSqlResult)

					if err5 != nil {
						//			log.Fatal("exec sql bad")
						//			log.Fatal(err5)
						//		break
					}
					//多行返回值拼接json数组
					resultlen := len(rawSqlResult)

					if resultlen > 1 {
						thisOneResult = "["
					}
					if rawSqlResult != nil && resultlen != 0 {
						for index, value := range rawSqlResult {

							thisOneResult = thisOneResult + value.Out
							if resultlen == 1 {
								break
							}
							if index != resultlen-1 {
								thisOneResult = thisOneResult + ","

							} else {
								thisOneResult = thisOneResult + "]"
							}
						}
					} else {
						thisOneResult = "sql excuted!"
					}

					get_value_result[k] = thisOneResult
				}

			}

			log.Println("-----------------------2.2、Exec curl.Get() --------------------------\n")

			var exec_get_map map[string]interface{}
			exec_get_map = get_value_map["Get"]

			if exec_get_map != nil {
				for k, v := range exec_get_map {

					log.Println("Get:", v)
					_, resp, _ := fasthttp.Get(nil, v.(string))

					// if status != fasthttp.StatusOK {
					//  log.Println(v," 请求没有成功:", status)
					//  return
					// }

					get_value_result[k] = string(resp)
					log.Print("curl.Get().Body():==>> ", string(resp))
				}
			}

			log.Println("-----------------------2.3、Exec curl.Post() ------------------------\n")

			exec_post_map := get_value_map["Post"]
			log.Println("Post :==>> ", exec_post_map)

			if exec_post_map != nil {
				for k, v := range exec_post_map {
					log.Println("Post object v is :==>> ", v)

					v_url := v.(map[string]interface{})["url"]
					v_body := v.(map[string]interface{})["body"]
					v_head := v.(map[string]interface{})["head"]

					log.Println("Post Url:  ", v_url)
					log.Println("Post Body: ", v_body)
					log.Println("Post Head: ", v_head)

					req := &fasthttp.Request{}
					req.SetRequestURI(v_url.(string))

					body := []byte(v_body.(string))
					req.SetBody(body)

					req.Header.SetContentType("application/json")

					if v_head != nil {
						for head_key, head_value := range v_head.(map[string]string) {
							req.Header.Set(head_key, head_value)
							log.Println("Post Head: ", head_key, ":", head_value)

						}
					}

					req.Header.SetMethod("POST")

					resp := &fasthttp.Response{}

					client := &fasthttp.Client{}
					log.Println("Post:", req)
					//log.Println("\nPost body:",resp )
					if err := client.Do(req, resp); err != nil {
						log.Println(v_url, " 请求失败:", err.Error())
						ctx.Write([]byte("{\"Error\" :" + string(v_url.(string)) + " \"请求失败\"}\n" + err.Error()))
						return
						//break
					}

					b := resp.Body()
					log.Println("result:\r\n", string(b))
					get_value_result[k] = string(b)
				}
			}

			if is_use_redis == "1" {
				log.Println("-----------------------2.4、Exec Redis Cluster --------------------\n")
				exec_redis_map := get_value_map["Redis"]
				log.Println("exec_redis_map :==>> ", exec_redis_map)

				if exec_redis_map != nil {
					for k, v := range exec_redis_map {
						log.Println("k:v=", k, ":", v)
						log.Println(len(v.([]interface{})))
						args := make([]interface{}, len(v.([]interface{})))
						//拼接 redis 入参
						for i, vv := range v.([]interface{}) {
							args[i] = interface{}(vv)
						}
						//执行redis 命令
						rep, err := redis.String(GeneralCluster.Do((args[0]).(string), args[1:]...))
						log.Println("Redis rep is ", rep)

						if err != nil {
							log.Println("Redis exec err is ", err)
							ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
							ctx.Write([]byte("{\"Error\" : \"Redis exec err is \"" + err.Error()))
							return
						}
						get_value_result[k] = string(rep)
					}
				}
			}
			log.Println("-----------------------2.5、Exec Shell --------------------\n")
			exec_shell_map := get_value_map["Shell"]
			log.Println("exec_shell_map :==>> ", exec_shell_map)

			if exec_shell_map != nil {
				for k, v := range exec_shell_map {
					log.Println("k:v=", k, ":", v)
					//执行 shell 命令
					err11, out, errout := Shellout(v.(string))

					if err11 != nil {
						log.Println("Shell exec err is ", err11)
						ctx.Write([]byte("{\"Error\" : \"Shell exec err is " + errout + "\"}"))
						ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
						return
					}
					get_value_result[k] = string(out)
				}
			}
			log.Println("-----------------------2.6、Exec SSH --------------------\n")
			exec_ssh_map := get_value_map["SSH"]
			log.Println("exec_ssh_map :==>> ", exec_shell_map)

			if exec_ssh_map != nil {
				for k, v := range exec_ssh_map {
					log.Println("k:v=", k, ":", v)
					user_ := v.(map[string]interface{})["user"].(string)
					pwd__ := v.(map[string]interface{})["pwd"].(string)
					addr_ := v.(map[string]interface{})["addr"].(string)
					cmd_ := v.(map[string]interface{})["cmd"].(string)

					pwd_, err14 := Decrypt(pwd__, aeskey)

					if err14 != nil {
						log.Println("aesDecrypt exec err is ", err14)
						ctx.Write([]byte("{\"Error\" : \"aesDecrypt exec err is " + fmt.Sprintf("%s", err14) + "\"}"))
						ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
						return
					}

					//执行 shell 命令
					cli := Cli{
						user: user_,
						pwd:  pwd_,
						addr: addr_,
					}
					out, err13 := cli.Run(cmd_)

					if err13 != nil {
						log.Println("SSH exec err is ", err13)
						ctx.Write([]byte("{\"Error\" : \"SSH exec err is " + fmt.Sprintf("%s", err13) + "\"}"))
						ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
						return
					}
					get_value_result[k] = string(out)
				}
			}
			log.Println("-----------------------2.7、Exec aesEncrypt --------------------\n")
			exec_encrypt_map := get_value_map["aesEncrypt"]
			log.Println("exec_encrypt_map :==>> ", exec_shell_map)

			if exec_encrypt_map != nil {
				for k, v := range exec_encrypt_map {
					log.Println("k:v=", k, ":", v)
					xpass, err := Encrypt(v.(string), aeskey)

					if err != nil {
						log.Println("aesEncrypt exec err is ", err)
						ctx.Write([]byte("{\"Error\" : \"aesEncrypt exec err is " + fmt.Sprintf("%s", err) + "\"}"))
						ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
						return
					}
					get_value_result[k] = string(xpass)
				}
			}

			log.Println("-----------------------2.8、Exec aesDecrypt --------------------\n")
			exec_decrypt_map := get_value_map["aesDecrypt"]
			log.Println("exec_decrypt_map :==>> ", exec_shell_map)

			if exec_decrypt_map != nil {
				for k, v := range exec_decrypt_map {
					log.Println("k:v=", k, ":", v)
					xpass, err := Decrypt(v.(string), aeskey)

					if err != nil {
						log.Println("aesDecrypt exec err is ", err)
						ctx.Write([]byte("{\"Error\" : \"aesDecrypt exec err is " + fmt.Sprintf("%s", err) + "\"}"))
						ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
						return
					}
					get_value_result[k] = string(xpass)
				}
			}
			log.Println("-----------------------2.9、Exec Memcache --------------------\n")
			exec_memcache_map := get_value_map["Memcache"]
			log.Println("exec_memcache_map :==>> ", exec_memcache_map)

			if exec_memcache_map != nil {
				mc := memcache.New(memcache_server)
				if mc == nil {
					log.Println("memcache New failed")
					ctx.Write([]byte("{\"Error\" : \"memcache New failed \"}"))
					ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
					return
				}
				for k, v := range exec_memcache_map {
					log.Println("k:v=", k, ":", v)
					//val := v.(map[string]string)
					for k1, v1 := range v.(map[string]interface{}) {
						var err error
						var rep string
						//执行memcache 命令
						switch k1 {
						case "set":
							err = mc.Set(&memcache.Item{Key: k, Value: []byte(v1.(string))})
							if err != nil {
								log.Println("Set ", k, " ", v1.(string), " failed,because ", err)
								ctx.Write([]byte("{\"Error\" : \"Set " + k + " " + v1.(string) + " failed,because " + fmt.Sprintf("%s", err) + " \"}"))
								ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
								return
							} else {
								rep = "ok"
							}
						case "get":
							rep_, err21 := mc.Get(k)
							if err21 != nil {
								log.Println("Get ", k, " ", v1.(string), " failed,because ", err21)
								ctx.Write([]byte("{\"Error\" : \"Set " + k + " " + v1.(string) + " failed,because " + fmt.Sprintf("%s", err21) + " \"}"))
								ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
								return
							} else {
								rep = string(rep_.Value)
								log.Println("memcache get key  ", k, " ", rep)
							}
						case "add":
							err = mc.Add(&memcache.Item{Key: k, Value: []byte(v1.(string))})
							if err != nil {
								log.Println("Add ", k, " ", v1.(string), " failed,because ", err)
								ctx.Write([]byte("{\"Error\" : \"Add " + k + " " + v1.(string) + " failed,because " + fmt.Sprintf("%s", err) + " \"}"))
								ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
								return
							} else {
								rep = "ok"
							}
						case "replace":
							err = mc.Replace(&memcache.Item{Key: k, Value: []byte(v1.(string))})
							if err != nil {
								log.Println("Replace ", k, " ", v1.(string), " failed,because ", err)
								ctx.Write([]byte("{\"Error\" : \"Replace " + k + " " + v1.(string) + " failed,because " + fmt.Sprintf("%s", err) + " \"}"))
								ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
								return
							} else {
								rep = "ok"
							}
						case "delete":
							err = mc.Delete(k)
							if err != nil {
								log.Println("Delete ", k, " ", v1.(string), " failed,because ", err)
								ctx.Write([]byte("{\"Error\" : \"Delete " + k + " " + v1.(string) + " failed,because " + fmt.Sprintf("%s", err) + " \"}"))
								ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
								return
							} else {
								rep = "ok"
							}
						case "incrby":
							err = mc.Set(&memcache.Item{Key: k, Value: []byte(v1.(string))})
							if err != nil {
								log.Println("Incrby ", k, " ", v1.(string), " failed,because ", err)
								ctx.Write([]byte("{\"Error\" : \"Incrby " + k + " " + v1.(string) + " failed,because " + fmt.Sprintf("%s", err) + " \"}"))
								ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
								return
							} else {
								rep = "ok"
							}
						//case "decrby":
						//rep, err := mc.Increment( k, []byte(v1.(string)) )
						//if err != nil {
						//	log.Println("Decrby ", k, " ", v1.(string), " failed,because ", err)
						//	ctx.Write([]byte("{\"Error\" : \"Decrby " + k + " " + v1.(string) + " failed,because " + fmt.Sprintf("%s", err) + " \"}"))
						//	return
						//} else {
						//	rep = "ok"
						//}
						default:
							log.Println("memcache ops not support")
							ctx.Write([]byte("{\"Error\" : \"memcache ops not support \"}"))
							ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
							return
						}
						get_value_result[k] = string(rep)
					}

				}
			}

			log.Println("-----------------------2.10、Exec upload --------------------\n")
			exec_upload_path_ := get_value_map["Upload"]["path"]
			var exec_upload_path string
			var exec_upload_key string

			if exec_upload_path_ != nil {
				exec_upload_path = exec_upload_path_.(string)
			}

			exec_upload_key_ := get_value_map["Upload"]["key"]
			if exec_upload_key_ != nil {
				exec_upload_key = exec_upload_key_.(string)
			}

			log.Println("exec_upload_path :==>> ", exec_upload_path)
			log.Println("exec_upload_key  :==>> ", exec_upload_key)

			if exec_upload_path_ != nil && exec_upload_key_ != nil && exec_upload_key != "" && exec_upload_path != "" {

				//根据参数名获取上传的文件
				fileHeader, err := ctx.FormFile(exec_upload_key)
				if err != nil {
					ctx.Write([]byte("{\"Error\" : \"FormFile err" + err.Error() + " \"}"))
					ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
					return
				}
				//打开上传的文件
				file, err := fileHeader.Open()
				if err != nil {
					ctx.Write([]byte("{\"Error\" : \"fileHeader.Open() err" + err.Error() + " \"}"))
					ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
					return
				}
				//使用完关闭文件
				defer file.Close()
				newFile := exec_upload_path + "/" + fileHeader.Filename

				is_exists, _ := PathExists(exec_upload_path)
				if is_exists == false {
					ctx.WriteString("{\"upload " + fileHeader.Filename + " Not Exists!\"} ")
					ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
					return
				}

				//新建一个文件，此处使用默认的txt格式
				nf, err := os.OpenFile(newFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
				if err != nil {
					ctx.WriteString("{\"os.OpenFile: " + fileHeader.Filename + " Failed ! because " + err.Error() + "\"} ")
					ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
					return
				}
				//使用完需要关闭
				defer nf.Close()
				//复制文件内容
				_, err = io.Copy(nf, file)
				if err != nil {
					ctx.WriteString("{\"io.Copy: " + fileHeader.Filename + " Failed ! because " + err.Error() + "\"} ")
					ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
					return
				}
				ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
				ctx.WriteString("{\"upload " + fileHeader.Filename + " success\"} ")
				return
			}

			log.Println("-----------------------2.11、Exec download --------------------\n")
			exec_download_path_ := get_value_map["Download"]["path"]
			var exec_download_path string
			var file_name string

			if exec_download_path_ != nil {
				exec_download_path = exec_download_path_.(string)
			}

			file_name_ := get_value_map["Download"]["key"]
			if file_name_ != nil {
				file_name = file_name_.(string)
			}

			log.Println("exec_download_path :==>> ", exec_download_path)
			log.Println("file_name  :==>> ", file_name)

			if exec_download_path_ != nil && file_name_ != nil && file_name != "" && exec_download_path != "" {
				download_file := exec_download_path + file_name
				is_exists, _ := PathExists(download_file)
				if is_exists == false {
					ctx.WriteString("{\"download " + file_name + " Not Exists!\"} ")
					return
				}
				f, err := os.Open(download_file)
				if err != nil {
					ctx.WriteString("{\"download filed because " + file_name + " not exists \"} ")
					return
				}
				defer f.Close()
				// 将文件读取出来
				data, err20 := ioutil.ReadAll(f)
				if err20 != nil {
					ctx.WriteString("{\"Read file filed because " + err20.Error() + " not exists \"} ")
					return
				}
				// 设置头信息：Content-Disposition ，消息头指示回复的内容该以何种形式展示，
				// 是以内联的形式（即网页或者页面的一部分），还是以附件的形式下载并保存到本地
				// Content-Disposition: inline
				// Content-Disposition: attachment
				// Content-Disposition: attachment; filename="filename.后缀"
				// 第一个参数或者是inline（默认值，表示回复中的消息体会以页面的一部分或者
				// 整个页面的形式展示），或者是attachment（意味着消息体应该被下载到本地；
				// 大多数浏览器会呈现一个“保存为”的对话框，将filename的值预填为下载后的文件名，
				// 假如它存在的话）。
				ctx.Response.Header.Set("Content-Disposition", "attachment; filename="+file_name)
				ctx.Write(data)
				return

			}

			log.Println("======================= 2.12、request_jwt ================================\n")
			//执行jwt登录
			str_jwt := get_value_result["#request_jwt"]
			if str_jwt != "" && str_jwt != "0" && str_jwt != "sql excuted!" {
				if len(SigningKey) <= 0 {
					log.Println("len(SigningKey) <= 0 ")
					ctx.Write([]byte("{\"Error\":\"len(SigningKey) <= 0\"}\n"))
					return
				}

				log.Print("jwt_request ... ")
				expire_times_int, _ := strconv.ParseInt(expire_times, 10, 64)
				token_string, err := jwt_request(str_jwt, SigningKey, expire_times_int)
				if err != nil {
					log.Println("jwt_request failed,because ", err.Error())
					ctx.Write([]byte("{\"Error\":\"jwt request failed,because " + err.Error() + "\"}\n"))
				}
				get_value_result["#request_jwt"] = string(token_string)
			}
			if str_jwt == "sql excuted!" {
				get_value_result["#request_jwt"] = "request token failed"
			}

			log.Println("-----------------------2.13、Exec add_policy --------------------\n")
			exec_add_policy_map := get_value_map["add_policy"]
			log.Println("exec_add_policy_map :==>> ", exec_add_policy_map)

			if exec_add_policy_map != nil {
				for k, v := range exec_add_policy_map {
					log.Println("k:v=", k, ":", v)
					args := make([]string, len(v.([]interface{})))
					//拼接 add_policy 入参
					for i, vv := range v.([]interface{}) {
						args[i] = vv.(string)
					}
					//执行add_policy 命令
					rep, err := effect.AddPolicy(args[0], args[1], args[2])

					if err != nil {
						log.Println("add_policy exec err is ", err)
						ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
						ctx.Write([]byte("{\"Error\" : \"add_policy exec err is \"" + err.Error()))
						return
					}

					get_value_result[k] = strconv.FormatBool(rep)
				}
			}

			log.Println("-----------------------2.14、Exec delete_policy --------------------\n")
			exec_delete_policy_map := get_value_map["delete_policy"]
			log.Println("exec_delete_policy_map :==>> ", exec_delete_policy_map)

			if exec_delete_policy_map != nil {
				for k, v := range exec_delete_policy_map {
					log.Println("k:v=", k, ":", v)
					args := make([]string, len(v.([]interface{})))
					//拼接 delete_policy 入参
					for i, vv := range v.([]interface{}) {
						args[i] = vv.(string)
					}
					//执行delete_policy 命令
					rep, err := effect.RemovePolicy(args[0], args[1], args[2])

					if err != nil {
						log.Println("delete_policy exec err is ", err)
						ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
						ctx.Write([]byte("{\"Error\" : \"delete_policy exec err is \"" + err.Error()))
						return
					}

					get_value_result[k] = strconv.FormatBool(rep)
				}
			}

			log.Println("-----------------------2.15、Exec load_policy --------------------\n")
			exec_load_policy_map := get_value_map["load_policy"]
			if exec_load_policy_map != nil {
				log.Println("exec_load_policy_map :==>> ", exec_load_policy_map)
				effect.LoadPolicy()
				get_value_result["load_policy"] = strconv.FormatBool(true)

			}
			
			log.Println("-----------------------2.16、Set Header --------------------\n")
			exec_header_map := get_value_map["Header"]
			log.Println("exec_header_map :==>> ", exec_header_map)

			if exec_header_map != nil {
				for k, v := range exec_header_map {
					log.Println("k:v=", k, ":", v)
					key_ := v.(map[string]interface{})["key"].(string)
					val_ := v.(map[string]interface{})["val"].(string)
					ctx.Response.Header.Set(key_, val_) //设置响应头
				}
			}
			
			log.Println("======================= 3、Run JavaScript ========================\n")
			log.Println("get_value_result: ", get_value_result)
			if len(jscript) > 0 {
				//javascript 解释器初始化
				vm := otto.New()
				//解析 js 脚本
				_, err := vm.Run(jscript)
				if err != nil {
					log.Println("Run javascript failed because of ", err)
				} else {
					// get_value_result 绑定为 js 中的变量
					jsa, err := vm.ToValue(get_value_result)
					if err != nil {
						log.Println("Javascript VM get map failed because of ", err)
					} else {
						// 执行 js 中的 main 函数
						result, err := vm.Call("main", nil, jsa)
						if err != nil {
							log.Println("vm.Call(main, nil, jsa) failed because of ", err)
						} else {
							//  js中的 main 函数执行结果导出
							tmpR, err := result.Export()
							if err != nil {
								log.Println("result.Export() failed because of ", err)
							} else {
								//get_value_result = tmpR
								for key, val := range tmpR.(map[string]interface{}) {
									get_value_result[key] = val.(string)
								}
								log.Println("get_value_result after js: ", get_value_result)
							}
						}
					}
				}

			} else {
				log.Println("len(jscript) < 0 ,So No Need to run JavaScript VM!")
			}
			log.Println("======================= 4、OutPut ================================\n")
			//打印执行结果map
			log.Print("get_value_result: ", get_value_result)
			//替换执行结果body
			//log.Print("before set out_put :"  )
			log.Println(out_put)
			for k, v := range get_value_result {
				out_put = strings.Replace(out_put, "${"+k+"}", v, -1)
			}
			log.Println("out_put: ", out_put)

			//执行登录
			str_login := get_value_result["#login"]
			if str_login != "" && str_login != "0" && str_login != "sql excuted!" {
				log.Print("Login ... ")
				err := loginHandle(ctx, str_login)
				if err != nil {
					log.Println("login failed ")
				}
			}
			//执行登出
			str_logout := get_value_result["#logout"]
			if str_logout != "" && str_logout != "0" && str_login != "sql excuted!" {
				log.Print("Logout ... ")
				logoutHandle(ctx)

			}

			//执行系统登录
			sys_str_login := get_value_result["#sys_login"]
			if sys_str_login != "" && sys_str_login != "0" && sys_str_login != "sql excuted!" {
				log.Print("SysLogin ... ")
				err := sysloginHandle(ctx, sys_str_login)
				if err != nil {
					log.Println("login failed because ", err)
				}
			}
			//执行系统登出
			sys_str_login = get_value_result["#sys_logout"]
			if sys_str_login != "" && sys_str_login != "0" && sys_str_login != "sql excuted!" {
				log.Print("SysLogout ... ")
				syslogoutHandle(ctx)
			}

			//执行带验证码的登录 需要传 captchaId 和 captchaSolution 连个变量
			if len(params_map["captchaId"]) > 0 && len(params_map["captchaSolution"]) > 0 {
				if verifyCaptcha(params_map["captchaId"], params_map["captchaSolution"]) {
					str_login_with_png := get_value_result["#login_with_png"]
					if str_login_with_png != "" && str_login_with_png != "0" && str_login_with_png != "sql excuted!" {
						log.Print("Login ... ")
						err := loginHandle(ctx, str_login)
						if err != nil {
							log.Println("login failed ")
						}
					}
					//执行系统登录
					sys_str_login_with_png := get_value_result["#sys_login_with_png"]
					if sys_str_login_with_png != "" && sys_str_login_with_png != "0" && sys_str_login_with_png != "sql excuted!" {
						log.Print("SysLogin ... ")
						err := sysloginHandle(ctx, sys_str_login)
						if err != nil {
							log.Println("login failed because ", err)
						}
					}

				} else {
					log.Print("Incorrect verification code  ")
					ctx.Write([]byte("{\"Error\" : \"Incorrect verification code \""))
					return
				}
			} 

			ctx.Response.Header.Set("content-type", "application/json") //返回数据格式是json
			ctx.Write([]byte(out_put))

			//在运行期间通过 接口 修改配置表，则重启
			if strings.HasPrefix(name_route, "chakala") && !strings.HasPrefix(name_route, "chakala_list") && !strings.HasPrefix(name_route, "chakala_import") {
				//server.Listener.Close()
				//c <- 1

				go func() {
					log.Println("Listener.Close() .... ")

					err = Listener.Close()

					for {

						if err != nil {
							log.Println("Listener.Closed err ", err)
							err = Listener.Close()
						} else {
							break
						}
					}
					log.Println("Listener.Closed .... ")

				}()

				return
			}

		})
	}

	router.Handle("GET", "/getCaptchaId", getCaptchaId)
	router.Handle("POST", "/verifyCaptchaId", verifyCaptchaId)
	router.Handle("GET", "/png", ServePNG)

	//fasthttp.ListenAndServe(listen_url, router.Handler)
	server := &fasthttp.Server{
		Handler:            router.Handler,
		Name:               "Chakala SERVER",
		MaxRequestBodySize: 200 * 1024 * 1024,
	}
	//err = fasthttp.Serve(Listener, router.Handler)
	log.Println("************************** server.Serve(Listener) success ***********************")

	err = server.Serve(Listener)
	//defer db.Close()

	if err != nil {
		log.Println("start server fail")
		os.Exit(-5)
	}

}

func main() {
	flag.Parse() //暂停获取参数
	later, _ = strconv.ParseInt(later_, 10, 64)
	log.Println("later ", later, " s")
	setupEnv()
	time.Sleep(time.Duration(later) * time.Second)

	// 必须在使用之前指定 session 的存储
	err := session.SetProvider("memory", &memory.Config{})
	if err != nil {
		log.Println("session init: ", err.Error())
		os.Exit(1)
	}

	err = sys_session.SetProvider("sqlite3", sqlite3.NewConfigWith("session.db", "session"))
	if err != nil {
		log.Println("sys_session init: ", err.Error())
		os.Exit(1)
	}

	//c := make(chan int)
	//defer close(c)
	//real_main()
	//for{
	////k := <- c
	//log.Println("Restart " )
	//for {
	real_main()
	//}
	//log.Println("nohup "+gloable_program_name+" -later 1 &")
	//Shellout("nohup "+gloable_program_name+" -later 1 &")
	//os.Exit(-5)
	// }
}
