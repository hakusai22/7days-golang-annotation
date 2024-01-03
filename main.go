package main

import (
	"fmt"
	"net/http"
	"net/url"
)

/*
   --idea
   -*- coding: utf-8 -*-
   @Author  : hakusai22
   @Time    : 2024/01/03 14:11
*/

func main() {
	fmt.Println("test")
	var link = "https://www.baidu.com?network_name=Apple+Search+Ads&app_name=com.btang.cbt&os_name=ios&partner_name={partner_name}&userId=315294050&language=1&debug_env_ignore_auth=true&app_flag=128&platform_type=2&app_version=2.10.0&ab_param=a&localTimezone=Asia/Shanghai"

	r2, err := http.Get(link)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r2.Request.Form)
	fmt.Println(r2.Body)
	fmt.Println(r2.StatusCode)

	u, err := url.Parse(link)
	q := u.Query()
	u.RawQuery = q.Encode()
	fmt.Println(u.String())
	res, err := http.Get(u.String())
	if err != nil {
		fmt.Println(err)
		return
	}
	if err != nil {
		fmt.Println("0")
		return
	}
	resCode := res.StatusCode
	if err != nil {
		fmt.Println("0")
		return
	}
	fmt.Printf("%d", resCode)
	fmt.Printf("%s", res.Body)

}
