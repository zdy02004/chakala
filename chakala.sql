--
-- PostgreSQL database dump
--

-- Dumped from database version 12.0
-- Dumped by pg_dump version 12.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: casbin_rule; Type: TABLE; Schema: public; Owner: javauser
--

CREATE TABLE public.casbin_rule (
    p_type character varying(100),
    v0 character varying(100),
    v1 character varying(100),
    v2 character varying(100),
    v3 character varying(100),
    v4 character varying(100),
    v5 character varying(100)
);


ALTER TABLE public.casbin_rule OWNER TO javauser;

--
-- Name: chakala_config; Type: TABLE; Schema: public; Owner: javauser
--

CREATE TABLE public.chakala_config (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text,
    method_type integer,
    valid text,
    get_value text,
    js_script text,
    out_put text,
    is_use integer,
    re_mark text
);


ALTER TABLE public.chakala_config OWNER TO javauser;

--
-- Name: chakala_config_id_seq; Type: SEQUENCE; Schema: public; Owner: javauser
--

CREATE SEQUENCE public.chakala_config_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.chakala_config_id_seq OWNER TO javauser;

--
-- Name: chakala_config_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: javauser
--

ALTER SEQUENCE public.chakala_config_id_seq OWNED BY public.chakala_config.id;


--
-- Name: seq_chakala; Type: SEQUENCE; Schema: public; Owner: javauser
--

CREATE SEQUENCE public.seq_chakala
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.seq_chakala OWNER TO javauser;

--
-- Name: chakala_config id; Type: DEFAULT; Schema: public; Owner: javauser
--

ALTER TABLE ONLY public.chakala_config ALTER COLUMN id SET DEFAULT nextval('public.chakala_config_id_seq'::regclass);


--
-- Data for Name: casbin_rule; Type: TABLE DATA; Schema: public; Owner: javauser
--

COPY public.casbin_rule (p_type, v0, v1, v2, v3, v4, v5) FROM stdin;
g	a	data2_admin	\N	\N	\N	\N
p	bob	data2	write	\N	\N	\N
p	a	shell2	read	\N	\N	\N
p	data2_admin	data2	read	\N	\N	\N
p	data2_admin	data2	write	\N	\N	\N
\.


--
-- Data for Name: chakala_config; Type: TABLE DATA; Schema: public; Owner: javauser
--

COPY public.chakala_config (id, created_at, updated_at, deleted_at, name, method_type, valid, get_value, js_script, out_put, is_use, re_mark) FROM stdin;
-3	2019-10-10 14:36:35.843283+08	2019-10-10 14:36:35.843283+08	2019-10-10 14:36:35.843283+08	chakala_list	1	{\n"method_type":{\n"val":"${method_type}",\n"notBlank":"false",\n"isNum":"false"\n}\n}	{\n"sql":{\n"sql_result":"select '{\\"id\\": '||'\\"'|| cast(id as varchar)||'\\",' ||'\\"name\\": '||'\\"'|| cast(name as varchar)||'\\",' ||'\\"method_type\\": '||'\\"'|| cast(method_type as varchar)||'\\",' ||'\\"valid\\": '|| cast(valid as varchar)||',' ||'\\"get_value\\": '|| cast(get_value as varchar)||','||'\\"js_script\\": '|| cast(js_script as varchar)||',' ||'\\"out_put\\": '||'\\"'|| replace(cast(out_put as varchar),'\\"','\\"')||'\\"'||',\\"is_use\\":'||'\\"'|| cast(is_use as varchar)||'\\"'||',\\"remark\\":'||'\\"'|| cast(COALESCE(remark,'') as varchar)||'\\"}' as out  from public.chakala_config"\n}}		{"sql_result":${sql_result}}	1	remark
-2	2019-09-08 03:09:15.973062+08	2019-09-08 03:09:15.973062+08	2019-09-08 03:09:15.973062+08	chakala_delete	1	{\n"method_type":{\n"val":"${method_type}",\n"notBlank":"true",\n"isNum":"false"\n}\n}	{\n"sql":{\n"sql_result":"delete from chakala_config where name='${name}'"\n}\n}		{\n"sql_result":"${sql_result}"\n}	1	\N
-1	2019-10-10 00:07:21.361531+08	2019-10-10 00:07:21.361531+08	2019-10-10 00:07:21.361531+08	chakala_insert	2	{\n"method_type":{\n"val":"${method_type}",\n"notBlank":"true",\n"isNum":"true",\n"<":"2",\n">":"0"\n},\n\n"is_use":{\n"val":"${is_use}",\n"notBlank":"true",\n"isNum":"true"\n}\n}	{\n"sql":{\n"sql_result":"insert into chakala_config select nextval('seq_chakala') ,now(),now(),now(),'${name}','${method_type}','${valid}','${get_value}','${js_script}','${out_put}','${is_use}','${re_mark}'"\n}\n}		{\n"sql_result":"${sql_result}"\n}	1	chakala_insert
11	2019-09-10 14:50:58.080783+08	2019-09-10 14:50:58.080783+08	2019-09-10 14:50:58.080783+08	test_aesDecrypt	2	{"method_type":{"isNum":"false","notBlank":"false","val":"2"}}	{"aesDecrypt":{"user":"${userid}"}}		{"userid":"${user}"}	1	
5	2019-08-28 18:35:06.290214+08	2019-08-28 18:35:06.290214+08	2019-08-28 18:35:06.290214+08	test_redis	2	{"method_type":{"isNum":"false","notBlank":"false","val":"${method_type}"}}	{"Redis":{"user":["GET","${userid}"]}}		{"userid":"${user}"}	1	
2	2019-08-18 22:57:24.478618+08	2019-08-18 22:57:24.478618+08	2019-08-18 22:57:24.478618+08	channel	2	{"channel_id":{"isNum":"true","notBlank":"true","val":"${channel_id}"},"channel_id2":{"isNum":"true","notBlank":"true","val":"${channel_id2}"},"userid":{"isNum":"false","notBlank":"true","val":"${userid}"}}	{"sql":{"channel_name":" select '{\\"channel_name\\": '||'\\"'|| cast(channel_name as varchar)||'\\"' || '}'  as out from public.channel_info where id >= '${channel_id}' and id <= '${channel_id2}'","user":" select '{\\"name\\": '||'\\"'|| cast(name as varchar)||'\\"' || '}'  as out from public.chakala_config where name = '${userid}'"}}		{"channel_names":"${channel_name}","userid":"${user}"}	1	
9	2019-09-10 11:58:13.506976+08	2019-09-10 11:58:13.506976+08	2019-09-10 11:58:13.506976+08	ssh_test	1	{"method_type":{"isNum":"false","notBlank":"false","val":"${method_type}"}}	{"SSH":{"user":{"addr":"10.1.235.103:22022","cmd":"${userid}","pwd":"XXXXXXXXXXXXXXXXX","user":"work"}}}		{"userid":"${user}"}	1	
4	2019-08-25 14:38:33.434025+08	2019-08-25 14:38:33.434025+08	2019-08-25 14:38:33.434025+08	test_post	2	{"method_type":{"isNum":"false","notBlank":"false","val":"2"}}	{"Post":{"user":{"body":"{\\"userid\\":\\"${userid}\\" , \\"channel_id\\":\\"${channel_id}\\",\\"channel_id2\\":\\"${channel_id2}\\"}","url":"http://10.19.13.27:8001/channel"}}}		{"channel_names":"${channel_name}","userid":"${user}"}	1	
15	2019-09-14 21:00:05.942284+08	2019-09-14 21:00:05.942284+08	2019-09-14 21:00:05.942284+08	test_html	1	null	{"html":{"user":"/home/project/chakala_go/html"}}		null	1	
7	2019-09-03 16:54:27.0311+08	2019-09-03 16:54:27.0311+08	2019-09-03 16:54:27.0311+08	valid	1	{\n"method_type":{\n"islogin":"true",\n"val":"${method_type}",\n"notBlank":"false",\n"isNum":"false"\n}}	{                                                                                                                                         \n "sql":{                                                                                                                                  \n  "user": " select '{\\"name\\": '||'\\"'|| cast(name as varchar)||'\\"' || '}'  as out from public.chakala_config where id = '${userid}'"  \n }                                                                                                                                        \n }		{ "userid" : "${user}" }	1	\N
10	2019-09-10 14:50:22.376676+08	2019-09-10 14:50:22.376676+08	2019-09-10 14:50:22.376676+08	test_aesEncrypt	2	{"method_type":{"isNum":"false","notBlank":"false","val":"2"}}	{"aesEncrypt":{"user":"${userid}"}}		{"userid":"${user}"}	1	
12	2019-09-11 18:58:44.204912+08	2019-09-11 18:58:44.204912+08	2019-09-11 18:58:44.204912+08	memcache_put	1	{"method_type":{"isNum":"false","notBlank":"false","val":"1"}}	{"Memcache":{"userid":{"set":"${userid}"}}}		{"userid":"${userid}"}	1	
13	2019-09-11 20:20:30.910795+08	2019-09-11 20:20:30.910795+08	2019-09-11 20:20:30.910795+08	memcache_get	1	{"method_type":{"isNum":"false","notBlank":"false","val":"1"}}	{"Memcache":{"userid":{"get":"${userid}"}}}		{"userid":"${userid}"}	1	
8	2019-09-03 22:43:05.050223+08	2019-09-03 22:43:05.050223+08	2019-09-03 22:43:05.050223+08	test_past	1	{\n"userid":{\n"val":"${userid}",\n"notBlank":"true",\n"isNum":"",\n"<":"13",\n">":"4",\n">=":"4",\n"<=":"7",\n"==":"",\n"isPast":"true",\n"isFuture":"",\n"pattern":""\n}\n}	{                                                                                                                                         \n "sql":{                                                                                                                                  \n  "user": " select '{\\"name\\": '||'\\"'|| cast(name as varchar)||'\\"' || '}'  as out from public.chakala_config where id = '${userid}'"  \n }                                                                                                                                        \n }		{ "userid" : "${user}" }	1	\N
3	2019-08-25 11:38:53.86097+08	2019-08-25 11:38:53.86097+08	2019-08-25 11:38:53.86097+08	get_test	2	{"method_type":{"isNum":"false","notBlank":"false","val":"2"}}	{"Get":{"user":"http://10.19.13.27:8001/name?userid=${userid}"}}		{"userid":"${user}"}	1	
14	2019-09-12 13:36:32.996582+08	2019-09-12 13:36:32.996582+08	2019-09-12 13:36:32.996582+08	try_upload	3	{"method_type":{"isNum":"false","notBlank":"false","val":"${method_type}"}}	{"Upload":{"key":"key","path":"/home/project/chakala_go"}}		{"userid":"${userid}"}	1	
17	2019-09-19 23:45:56.104437+08	2019-09-19 23:45:56.104437+08	2019-09-19 23:45:56.104437+08	test_donwload	1	{ "method_type":{ "val":"${filename}","notBlank":"true","isNum":"false"}}	{\n"Download":{\n"key":"${filename}",\n"path":"/home/project/chakala_go/"\n}}			1	\N
18	2019-09-20 17:05:59.045205+08	2019-09-20 17:05:59.045205+08	2019-09-20 17:05:59.045205+08	request_token	1	{ "userid": { "val": "${name}","notBlank": "true", "isNum": "false","<": "13", ">": "4", ">=": "4","<=": "7","==": "", "isPast": "",    "isFuture": "","pattern": ""   }}	{ "sql": { "#request_jwt": "select cast(coalesce(id,0) as varchar)  as out from \\"User\\" where name = '${name}' and member_number ='${phone}' " } }		{"token": "${#request_jwt}"}	1	\N
19	2019-09-20 17:17:59.388452+08	2019-09-20 17:17:59.388452+08	2019-09-20 17:17:59.388452+08	check_token	1	{ "userid": { "val": "${token}","notBlank": "true", "isjwt": "true"}}	{ "sql": { "request": "select 1 as out from \\"User\\"  limit 1 " } }		{"result": "${request}"}	1	\N
20	2019-10-07 14:51:49.727701+08	2019-10-07 14:51:49.727701+08	2019-10-07 14:51:49.727701+08	add_policy	1	{ "param1": { "val": "${param1}","notBlank": "true"},"param2": { "val": "${param2}","notBlank": "true"},"param3": { "val": "${param3}","notBlank": "true"} }	{"add_policy":{"user":["${param1}", "${param2}", "${param3}"]}}		{"user": "${user}"}	1	add_policy
21	2019-10-07 15:12:17.64992+08	2019-10-07 15:12:17.64992+08	2019-10-07 15:12:17.64992+08	delete_policy	1	{ "param1": { "val": "${param1}","notBlank": "true"},"param2": { "val": "${param2}","notBlank": "true"},"param3": { "val": "${param3}","notBlank": "true"} }	{"delete_policy":{"user":["${param1}", "${param2}", "${param3}"]}}		{"user": "${user}"}	1	delete_policy
22	2019-10-07 16:05:09.99732+08	2019-10-07 16:05:09.99732+08	2019-10-07 16:05:09.99732+08	load_policy	1	{ "param1": { "val": "${param1}","notBlank": "false"} }	{"load_policy":{"aa":"1"}}		{"load_policy": "${load_policy}"}	1	load_policy
23	2019-10-10 02:09:45.723196+08	2019-10-10 02:09:45.723196+08	2019-10-10 02:09:45.723196+08	test_js	1	{ "userid": { "val": "${name}","notBlank": "true", "isNum": "false","<": "13", ">": "4", ">=": "4","<=": "7","==": "", "isPast": "",    "isFuture": "","pattern": ""   }}	{ "sql": { "name": "select cast(coalesce(name,'0') as varchar)  as out from \\"User\\" where name = '${name}' and member_number ='${phone}'  " } }	function main(tab) {\n    result = {}\n    tab["name"] += ":http://www.baidu.com"\n    return result\n}	{"user": "${name}"}	1	js_test
-4	2019-10-16 16:47:14.097295+08	2019-10-16 16:47:14.097295+08	2019-10-16 16:47:14.097295+08	chakala_upsert	2	{\n"method_type":{\n"val":"${method_type}",\n"notBlank":"true",\n"isNum":"false"\n}\n}	{"sql":{\n"sql_result":"insert into chakala_config (id,created_at,updated_at,deleted_at,name,method_type,valid,get_value,js_script,out_put,is_use,re_mark) values ( nextval('seq_chakala') ,now(),now(),now(),'${name}','${method_type}','${valid}','${get_value}','${js_script}','${out_put}',${is_use},'${re_mark}' ) on  CONFLICT(name) do update set name=excluded.name,method_type=excluded.method_type,valid=excluded.valid,get_value=excluded.get_value,js_script=excluded.js_script,out_put=excluded.out_put,is_use=excluded.is_use,re_mark=excluded.re_mark"\n}}		{"sql_result":"${sql_result}"}	1	chakala_upsert
16	2019-09-18 13:59:27.658041+08	2019-09-18 13:59:27.658041+08	2019-09-18 13:59:27.658041+08	test_login	1	{ "userid": { "val": "${name}","notBlank": "true", "isNum": "false","<": "13", ">": "4", ">=": "4","<=": "7","==": "", "isPast": "",    "isFuture": "","pattern": ""   }}	{ "sql": { "#login": "select cast(coalesce(name,'0') as varchar)  as out from \\"User\\" where name = '${name}' and member_number ='${phone}' " } }		{"user": "${user}"}	1	\N
26	2019-10-19 21:20:11.459697+08	2019-10-19 21:20:11.459697+08	2019-10-19 21:20:11.459697+08	test_login_with_png	1	{ "userid": { "val": "${name}","notBlank": "true", "isNum": "false","<": "13", ">": "4", ">=": "4","<=": "7","==": "", "isPast": "",    "isFuture": "","pattern": ""   }}	{ "sql": { "#login_with_png": "select coalesce(name,'0')  as out from chakala_config where name = '${name}' " } }		{"user": "${#login_with_png}"}	1	\N
24	2019-10-17 13:58:40.470855+08	2019-10-17 13:58:40.470855+08	2019-10-17 13:58:40.470855+08	test_row_2_json	1	{ "userid": { "val": "${userid}","notBlank": "false" }}	{ "sql": { "user": "select row_to_json(t)::text as out from ( select id,name,method_type,valid,get_value,js_script,out_put,is_use,re_mark from  public.chakala_config ) t " } }		{"userid": ${user} }	1	test_row_2_json
-5	2019-10-17 12:01:53.330681+08	2019-10-17 12:01:53.330681+08	2019-10-17 12:01:53.330681+08	chakala_turn	2	{\n"method_type":{\n"val":"${method_type}",\n"notBlank":"false",\n"isNum":"false"\n}}	{"ok":{"ok":"1"}}		{"turned":${chakala_turn}}	1	chakala_turn
-6	2019-10-17 13:31:24.18745+08	2019-10-17 13:31:24.18745+08	2019-10-17 13:31:24.18745+08	chakala_import	2	{\n"method_type":{\n"val":"${method_type}",\n"notBlank":"true",\n"isNum":"false"\n}\n}	{"sql":{\n"sql_result":"insert into chakala_config (id,created_at,updated_at,deleted_at,name,method_type,valid,get_value,js_script,out_put,is_use,re_mark) values ( nextval('seq_chakala') ,now(),now(),now(),'${name}','${method_type}','${valid}','${get_value}','${js_script}','${out_put}',${is_use},'${re_mark}' ) on  CONFLICT(name) do update set name=excluded.name,method_type=excluded.method_type,valid=excluded.valid,get_value=excluded.get_value,js_script=excluded.js_script,out_put=excluded.out_put,is_use=excluded.is_use,re_mark=excluded.re_mark"\n}}		{"sql_result":"${sql_result}"}	1	chakala_import
25	2019-10-17 15:05:12.581606+08	2019-10-17 15:05:12.581606+08	2019-10-17 15:05:12.581606+08	test_syslogin	1	{ "userid": { "val": "${name}","notBlank": "true", "isNum": "false","<": "13", ">": "4", ">=": "4","<=": "7","==": "", "isPast": "",    "isFuture": "","pattern": ""   }}	{ "sql": { "#sys_login": "select coalesce(name,'0')  as out from chakala_config where name = '${name}' " } }		{"user": "${user}"}	1	\N
27	2019-10-22 14:03:12.323473+08	2019-10-22 14:03:12.323473+08	2019-10-22 14:03:12.323473+08	test_head	1	{ "userid":{ "val":"${name}","notBlank":"true", "isNum":"false","<":"13", ">":"4", ">=":"4","<=":"7","==":"", "isPast":"",    "isFuture":"","pattern":""   }}	{ "Header":{ "head1":{"key":"key1","val":"val1" },"head2":{"key":"key2","val":"val2" }} }		{"head1":"${head1}","head2":"${head2}"}	1	\N
\.


--
-- Name: chakala_config_id_seq; Type: SEQUENCE SET; Schema: public; Owner: javauser
--

SELECT pg_catalog.setval('public.chakala_config_id_seq', 1, false);


--
-- Name: seq_chakala; Type: SEQUENCE SET; Schema: public; Owner: javauser
--

SELECT pg_catalog.setval('public.seq_chakala', 303, true);


--
-- Name: chakala_config chakala_config_pkey; Type: CONSTRAINT; Schema: public; Owner: javauser
--

ALTER TABLE ONLY public.chakala_config
    ADD CONSTRAINT chakala_config_pkey PRIMARY KEY (id);


--
-- Name: chakala_config_name_idx; Type: INDEX; Schema: public; Owner: javauser
--

CREATE UNIQUE INDEX chakala_config_name_idx ON public.chakala_config USING btree (name);


--
-- Name: idx_chakala_config_deleted_at; Type: INDEX; Schema: public; Owner: javauser
--

CREATE INDEX idx_chakala_config_deleted_at ON public.chakala_config USING btree (deleted_at);


--
-- PostgreSQL database dump complete
--

