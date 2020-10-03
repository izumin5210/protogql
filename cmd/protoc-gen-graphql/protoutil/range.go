package protoutil

import (
	"errors"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var BreakRange = errors.New("break range")

func RangeServices(
	parent interface {
		Services() protoreflect.ServiceDescriptors
	},
	f func(protoreflect.ServiceDescriptor) error,
) error {
	ds := parent.Services()
	n := ds.Len()
	for i := 0; i < n; i++ {
		if err := f(ds.Get(i)); err != nil {
			if err == BreakRange {
				return nil
			}
			return err
		}
	}
	return nil
}

func RangeMethods(
	parent interface {
		Methods() protoreflect.MethodDescriptors
	},
	f func(protoreflect.MethodDescriptor) error,
) error {
	ds := parent.Methods()
	n := ds.Len()
	for i := 0; i < n; i++ {
		if err := f(ds.Get(i)); err != nil {
			if err == BreakRange {
				return nil
			}
			return err
		}
	}
	return nil
}

func RangeFields(
	parent interface {
		Fields() protoreflect.FieldDescriptors
	},
	f func(protoreflect.FieldDescriptor) error,
) error {
	ds := parent.Fields()
	n := ds.Len()
	for i := 0; i < n; i++ {
		if err := f(ds.Get(i)); err != nil {
			if err == BreakRange {
				return nil
			}
			return err
		}
	}
	return nil
}
