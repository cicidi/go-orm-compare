package util

import "log"

//Created by Walter Chen / walter.chen@byton.com
//Date: 2019-04-09
//Time: 17:08

func CheckErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
