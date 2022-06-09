package core

import (
	"errors"
	"fmt"
	"reflect"
)

type _AspectJointPoint struct {
	Aspect     reflect.Type
	JointPoint reflect.Value
}

type AspectJointPointTab map[string]*_AspectJointPoint

func (aspectJPTab *AspectJointPointTab) Analysis(aop interface{}) error {
	vAOP := reflect.ValueOf(aop).Elem()
	newAspectJPTab := map[string]*_AspectJointPoint{}

	for i := 0; i < vAOP.NumField(); i++ {
		vAspectJP := vAOP.Field(i)
		tAspectJP := vAspectJP.Type()

		if tAspectJP.Kind() != reflect.Func {
			continue
		}

		if !vAspectJP.IsValid() || vAspectJP.IsNil() {
			return fmt.Errorf("analysis aspect join point [%s] failed, nil join point", tAspectJP.String())
		}

		if tAspectJP.NumOut() <= 0 {
			return fmt.Errorf("analysis aspect join point [%s] failed, aspect invalid", tAspectJP.String())
		}

		tAspect := tAspectJP.Out(0)

		newAspectJPTab[tAspect.String()] = &_AspectJointPoint{
			Aspect:     tAspect,
			JointPoint: vAspectJP,
		}
	}

	*aspectJPTab = newAspectJPTab

	return nil
}

func (aspectJPTab *AspectJointPointTab) NewAspect(aspect string, args []reflect.Value) (reflect.Value, error) {
	if *aspectJPTab == nil {
		return reflect.Value{}, errors.New("new aspect failed, no analysis")
	}

	ajp, ok := (*aspectJPTab)[aspect]
	if !ok {
		return reflect.Value{}, fmt.Errorf("new aspect [%s] failed, aspect not found", aspect)
	}

	tjp := ajp.JointPoint.Type()

	for i := 0; i < tjp.NumIn(); i++ {
		targ := tjp.In(i)

		if tjp.IsVariadic() {
			if i >= len(args) {
				break
			}

			if i >= tjp.NumIn()-1 {
				targ = targ.Elem()
			}

		} else {
			if i >= len(args) {
				return reflect.Value{}, fmt.Errorf("new aspect [%s] failed, aspect args not matching", aspect)
			}
		}

		if !args[i].Type().AssignableTo(targ) {
			return reflect.Value{}, fmt.Errorf("new aspect [%s] failed, aspect args not matching", aspect)
		}
	}

	return ajp.JointPoint.Call(args)[0], nil
}

func (aspectJPTab *AspectJointPointTab) IsAspect(aspect string) bool {
	if *aspectJPTab == nil {
		return false
	}

	_, ok := (*aspectJPTab)[aspect]
	return ok
}
