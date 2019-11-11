package main

import (
  "bytes"
  "path"
  "log"
	"github.com/dchest/captcha"
	"github.com/valyala/fasthttp"
	"unsafe"
)

func getCaptchaId(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("content-type", "application/json") //·µ»ي�½ˇjson
	ctx.Write([]byte("{\"CaptchaId\":" + "\""+captcha.New()+"\"}") )
	return
}

func verifyCaptchaId(ctx *fasthttp.RequestCtx)  {
	if !captcha.VerifyString( string(ctx.FormValue("captchaId")), string(ctx.FormValue("captchaSolution"))) {
		ctx.Write([]byte("{\"Result\" :" + "false"))
		return 
	} else {
		ctx.Write([]byte("{\"Result\" :" + "true"))
		return 
	}
}

func verifyCaptcha( captchaId string, captchaSolution string  )  bool  {
	if !captcha.VerifyString( captchaId, captchaSolution   ) {
		return false
	} else {
		return true
	}
}

func servePnghttp(ctx *fasthttp.RequestCtx, id, ext string )  {
	ctx.Response.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Response.Header.Set("Pragma", "no-cache")
	ctx.Response.Header.Set("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		log.Println("case png:")
		ctx.Response.Header.Set("Content-Type", "image/png")
		captcha.WriteImage(&content, id, captcha.StdWidth, captcha.StdHeight)
	default:
		 log.Println("ctx.NotFound()")
		 ctx.NotFound()
		 return
	}

	ctx.Write(content.Bytes())
}

func ServePNG(ctx *fasthttp.RequestCtx) {
	log.Println("ServePNG(ctx *fasthttp.RequestCtx):")
	params_map := make(map[string]string)
		
	ctx.URI().QueryArgs().VisitAll(func(k []byte, value []byte) {
	str := (*string)(unsafe.Pointer(&k))
	params_map[*str] = string(value)
	log.Println(params_map[*str])
	})
	
	file := params_map["id"]
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || id == "" {
		ctx.NotFound()
		return
	}
	if string( ctx.FormValue("reload") ) != "" {
		captcha.Reload(id)
	}
  servePnghttp(ctx, id, ext) 
}
