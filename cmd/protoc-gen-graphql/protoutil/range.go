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

func RangeMessages(
	ds protoreflect.MessageDescriptors,
	f func(protoreflect.MessageDescriptor) error,
) error {
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

func RangeEnums(
	ds protoreflect.EnumDescriptors,
	f func(protoreflect.EnumDescriptor) error,
) error {
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

func RangeEnumValues(
	parent protoreflect.EnumDescriptor,
	f func(protoreflect.EnumValueDescriptor) error,
) error {
	ds := parent.Values()
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

func TypeDFS(d protoreflect.Descriptor) ([]protoreflect.Descriptor, error) {
	var resp []protoreflect.Descriptor
	if d, ok := d.(protoreflect.MessageDescriptor); ok {
		resp = append(resp, d)
	}
	if d, ok := d.(interface {
		Messages() protoreflect.MessageDescriptors
	}); ok {
		err := RangeMessages(d.Messages(), func(m protoreflect.MessageDescriptor) error {
			typs, err := TypeDFS(m)
			if err != nil {
				return err
			}
			resp = append(resp, typs...)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	if d, ok := d.(protoreflect.EnumDescriptor); ok {
		resp = append(resp, d)
	}
	if d, ok := d.(interface {
		Enums() protoreflect.EnumDescriptors
	}); ok {
		err := RangeEnums(d.Enums(), func(m protoreflect.EnumDescriptor) error {
			typs, err := TypeDFS(m)
			if err != nil {
				return err
			}
			resp = append(resp, typs...)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}
