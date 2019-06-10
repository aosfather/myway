package core

import "testing"

func TestAuthFieldValue_Validate(t *testing.T) {
	f := AuthField{"test1", "测试", FT_INT}

	v := AuthFieldValue{&f, VT_SINGLE, "12", nil}

	t.Log(v.Validate("21"))
	t.Log(v.Validate("12"))
}

func TestAuthobject_Validate(t *testing.T) {

	f1 := AuthField{"test1", "测试", FT_INT}
	f2 := AuthField{"test2", "测试", FT_INT}
	c := AuthClass{"T", make(map[string]*AuthField)}
	c.Fields["test1"] = &f1
	c.Fields["test2"] = &f2

	os := AuthobjectSet{&c, []*Authobject{}, nil}

	v1 := AuthFieldValue{&f1, VT_SINGLE, "12", nil}
	v2 := AuthFieldValue{&f2, VT_SINGLE, "12", nil}
	o := Authobject{false, []*AuthFieldValue{&v1, &v2}}

	os.Objects = append(os.Objects, &o)

	m := make(map[string]string)
	m["test1"] = "21"
	m["test2"] = "12"

	t.Log(o.Validate(m))
	t.Log(os.Validate("T", m))

	m["test1"] = "12"
	m["test2"] = "12"

	t.Log(o.Validate(m))

	t.Log(os.Validate("T", m))
}

func TestDefaultRoleMetaManager_AddRole(t *testing.T) {
	dm := DefaultRoleMetaManager{}
	sr := Role{nil, "super", make(map[string]*AuthobjectSet)}
	r := Role{&sr, "myrole1", make(map[string]*AuthobjectSet)}
	dm.AddRole(&r)
	dm.AddRole(&sr)

	t.Log(dm.FindRole("myrole1"))

	t.Log(dm.FindRole("super"))
}

func TestRole_Validate(t *testing.T) {
	f1 := AuthField{"test1", "测试", FT_INT}
	f2 := AuthField{"test2", "测试", FT_INT}
	c := AuthClass{"T", make(map[string]*AuthField)}
	c.Fields["test1"] = &f1
	c.Fields["test2"] = &f2

	os := AuthobjectSet{&c, []*Authobject{}, nil}

	v1 := AuthFieldValue{&f1, VT_SINGLE, "12", nil}
	v2 := AuthFieldValue{&f2, VT_SINGLE, "12", nil}
	o := Authobject{false, []*AuthFieldValue{&v1, &v2}}

	os.Objects = append(os.Objects, &o)

	r := Role{nil, "myrole1", make(map[string]*AuthobjectSet)}

	r.setMap["T"] = &os

	m := make(map[string]string)
	m["test1"] = "12"
	m["test2"] = "21"
	t.Log(r.Validate("T", m, 0))
	t.Log(r.Validate("T", m, 11))

	m["test2"] = "12"
	t.Log(r.Validate("T", m, 1))
}
