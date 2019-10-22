#!/bin/sh
curl http://10.19.13.27:8000/name?userid=name

curl 10.19.13.27:8000/name -X POST -H "Content-Type:application/json" -d '{"userid":1}'
curl 10.19.13.27:8000/channel  -X POST -H "Content-Type:application/json" -d '{"userid":"name" , "channel_id":"471578910","channel_id2":"472578996"}'


curl "http://10.19.13.27:8000/channel?userid=name&channel_id=471578910&channel_id2=472599996"




update public.chakala_config set get_value =
'{
"sql":{
 "user": " select ''{\"name\": ''||''\"''|| cast(name as varchar)||''\"'' || ''}''  as out from public.chakala_config where name = ''${userid}''",
 "channel_name":  " select ''{\"channel_name\": ''||''\"''|| cast(channel_name as varchar)||''\"'' || ''}''  as out from public.channel_info where id >= ''${channel_id}'' and id <= ''${channel_id2}''"
}
}' where id = 2;


update public.chakala_config set out_put =
 '{ "userid" : "${user}",
    "channel_names" : "${channel_name}"
  }' where id =2;


http://www.bejson.com/



curl "http://10.19.13.27:8000/get_test?userid=name"
curl 10.19.13.27:8000/test_post -X POST -H "Content-Type:application/json" -d '{"userid":"name" , "channel_id":"471578910","channel_id2":"472578996"}'


curl "http://10.19.13.27:8000/test_redis?userid=foo1"
curl "http://10.19.13.27:8000/valid?userid=7"

update chakala_config set get_value ='{                                                                                                                                         
 "sql":{                                                                                                                                  
  "user": " select ''{\"name\": ''||''\"''|| cast(name as varchar)||''\"'' || ''}''  as out from public.chakala_config where id = ''${userid}''"  
 }                                                                                                                                        
 }' where id = 7;
