package main

import (
	"log"

	"github.com/robertkrimen/otto"
)

func main() {
	vm := otto.New()

	vm.Set("sayHello", func(call otto.FunctionCall) otto.Value {
		log.Printf("Hello, %s.\n", call.Argument(0).String())
		return otto.Value{}
	})

	vm.Set("twoPlus", func(call otto.FunctionCall) otto.Value {
		right, _ := call.Argument(0).ToInteger()
		result, _ := vm.ToValue(2 + right)
		return result
	})

	testExample, _ := vm.Run(`
    	sayHello("Xyzzy");      // Hello, Xyzzy.
		sayHello();             // Hello, undefined

		testExample = twoPlus(2.0); // 4
	`)

	log.Println(testExample)
}

func checkNilErr(err any) {
	if err != nil {
		// log.Fatalln("Error:\n%v\n", err)
		log.Fatalln(err)
	}
}
