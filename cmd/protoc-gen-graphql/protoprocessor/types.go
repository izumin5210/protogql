package protoprocessor

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type Types struct {
	descriptorByName map[protoreflect.FullName]interface{}
}

func NewTypes() *Types {
	return &Types{
		descriptorByName: map[protoreflect.FullName]interface{}{},
	}
}

func (t *Types) RegisterFromFiles(files *protoregistry.Files) {
	files.RangeFiles(func(f protoreflect.FileDescriptor) bool {
		t.RegisterFromFile(f)
		return true
	})
}

func (t *Types) RegisterFromFile(fd protoreflect.FileDescriptor) {
	t.registerDescriptors(fd)
}

func (t *Types) registerDescriptors(d interface {
	Enums() protoreflect.EnumDescriptors
	Messages() protoreflect.MessageDescriptors
	Extensions() protoreflect.ExtensionDescriptors
}) {
	t.registerMessages(d.Messages())
	t.registerEnums(d.Enums())
}

func (t *Types) registerEnums(ds protoreflect.EnumDescriptors) {
	n := ds.Len()
	for i := 0; i < n; i++ {
		d := ds.Get(i)
		t.descriptorByName[d.FullName()] = d
	}
}

func (t *Types) registerMessages(ds protoreflect.MessageDescriptors) {
	n := ds.Len()
	for i := 0; i < n; i++ {
		d := ds.Get(i)
		t.descriptorByName[d.FullName()] = d
		t.registerDescriptors(d)
	}
}

func (t *Types) registerExtensions(ds protoreflect.ExtensionDescriptors) {
	n := ds.Len()
	for i := 0; i < n; i++ {
		d := ds.Get(i)
		t.descriptorByName[d.FullName()] = d
	}
}

func (t *Types) FindEnumByName(name protoreflect.FullName) (protoreflect.EnumDescriptor, error) {
	d, ok := t.descriptorByName[name]
	if !ok {
		return nil, fmt.Errorf("enum %s is not found", name)
	}
	ed, ok := d.(protoreflect.EnumDescriptor)
	if !ok {
		return nil, fmt.Errorf("enum %s is not found", name)
	}
	return ed, nil
}

func (t *Types) FindMessageByName(name protoreflect.FullName) (protoreflect.MessageDescriptor, error) {
	d, ok := t.descriptorByName[name]
	if !ok {
		return nil, fmt.Errorf("enum %s is not found", name)
	}
	md, ok := d.(protoreflect.MessageDescriptor)
	if !ok {
		return nil, fmt.Errorf("enum %s is not found", name)
	}
	return md, nil
}

func (t *Types) FindExtensionByName(name protoreflect.FullName) (protoreflect.ExtensionDescriptor, error) {
	d, ok := t.descriptorByName[name]
	if !ok {
		return nil, fmt.Errorf("extension %s is not found", name)
	}
	ed, ok := d.(protoreflect.ExtensionDescriptor)
	if !ok {
		return nil, fmt.Errorf("extension %s is not found", name)
	}
	return ed, nil
}
