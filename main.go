package main

import snowflake "frame/library/util"

func main() {
	//config.InitConfig("./")
	//fmt.Println(viper.Get("mysql.port"))
	id := snowflake.GetId(1,1,false)
	snowflake.Parse(id)
}
