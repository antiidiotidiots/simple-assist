package main

import (
	"fmt"
	"log"

	v8 "rogchap.com/v8go"
)

func main() {
	iso := v8.NewIsolate()

	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		fmt.Printf("%v", info.Args())
		return nil
	})

	global := v8.NewObjectTemplate(iso)
	global.Set("print", printfn)

	ctx := v8.NewContext(iso, global)
	ctx.RunScript("print('foo')", "print.js")
}

func checkNilErr(err any) {
	if err != nil {
		// log.Fatalln("Error:\n%v\n", err)
		log.Fatalln(err)
	}
}
