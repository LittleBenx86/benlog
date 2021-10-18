package convertor

import "testing"

func Test_Map2Struct(t *testing.T) {
	mapIns := make(map[string]interface{})
	mapIns["Name"] = "Ben"
	mapIns["Age"] = 16

	type Man struct {
		Name string
		Age  int
	}
	var m Man
	err := Map2Struct(mapIns, &m)
	if err != nil {
		t.Logf("[error] %s\n", err)
	}
	t.Logf("%+v\n", m)
}

func Test_Struct2Map(t *testing.T) {
	type Career struct {
		Name string
	}
	type Man struct {
		Name   string
		Age    int
		Career Career
	}
	m := Man{
		Name: "Ben",
		Age:  16,
	}

	tar := Struct2Map(m)
	t.Logf("%+v\n", tar)
}
