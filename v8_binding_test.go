package v8

import "testing"
import "runtime"

func Test_Bind_Variadic(t *testing.T) {
	template := engine.NewObjectTemplate()

	template.Bind("Call", func(arg1, arg2 string, args ...string) *Value {
		val := engine.NewObject()
		obj := val.ToObject()
		obj.SetProperty("a1", engine.NewString(arg1), PA_None)
		obj.SetProperty("a2", engine.NewString(arg2), PA_None)
		array := engine.NewArray(len(args))
		arrayObj := array.ToObject()
		for i, arg := range args {
			arrayObj.SetElement(i, engine.NewString(arg))
		}
		obj.SetProperty("as", array, PA_None)
		return val
	})

	engine.NewContext(template).Scope(func(cs ContextScope) {
		script := engine.Compile([]byte(`
		a = Call("aaa", "bbb");
		if (a.a1 != "aaa" || a.a2 != "bbb") {
			throw "value should be {\"a1\":\"aaa\",\"a2\":\"bbb\"} not " + JSON.stringify(a);
		}
		a = Call("aaa", "bbb", "ccc");
		if (a.a1 != "aaa" || a.a2 != "bbb" || a.as.length != 1 || a.as[0] != "ccc") {
			throw "value should be {\"a1\":\"aaa\",\"a2\":\"bbb\",\"as\":[\"ccc\"]} not " + JSON.stringify(a);
		}
		a = Call("aaa", "bbb", "ccc", "ddd");
		if (a.a1 != "aaa" || a.a2 != "bbb" || a.as.length != 2 || a.as[0] != "ccc" || a.as[1] != "ddd") {
			throw "value should be {\"a1\":\"aaa\",\"a2\":\"bbb\",\"as\":[\"ccc\",\"ddd\"]} not " + JSON.stringify(a);
		}
		"ok"
		`), nil)

		var retVal *Value
		if err := cs.TryCatch(func() {
			retVal = cs.Run(script)
		}); err != nil {
			t.Fatal(err)
		}
		if !retVal.IsString() || retVal.ToString() != "ok" {
			t.Fatalf("value should be \"ok\" not %s", ToJSON(retVal))
		}
	})

	runtime.GC()
}