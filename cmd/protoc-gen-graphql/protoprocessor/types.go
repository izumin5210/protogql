package protoprocessor

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type Types struct {
	files          map[string]*descriptor.FileDescriptorProto
	msgByFullname  map[typePath]*descriptor.DescriptorProto
	svcByFullname  map[typePath]*descriptor.ServiceDescriptorProto
	enumByFullname map[typePath]*descriptor.EnumDescriptorProto
}

func NewTypes() *Types {
	return &Types{
		files:          map[string]*descriptor.FileDescriptorProto{},
		msgByFullname:  map[typePath]*descriptor.DescriptorProto{},
		svcByFullname:  map[typePath]*descriptor.ServiceDescriptorProto{},
		enumByFullname: map[typePath]*descriptor.EnumDescriptorProto{},
	}
}

type typePath string

func (p typePath) Join(path string) typePath {
	return p + "." + typePath(path)
}

func (t *Types) FindFile(file string) *descriptor.FileDescriptorProto {
	return t.files[file]
}

func (t *Types) FindMessage(path string) *descriptor.DescriptorProto {
	return t.msgByFullname[typePath(path)]
}

func (t *Types) FindEnum(path string) *descriptor.EnumDescriptorProto {
	return t.enumByFullname[typePath(path)]
}

func (t *Types) AddFile(fd *descriptor.FileDescriptorProto) {
	t.files[fd.GetName()] = fd

	path := typePath("." + fd.GetPackage())

	t.addSvcs(fd.GetService(), path)
	t.addMsgs(fd.GetMessageType(), path)
	t.addEnums(fd.GetEnumType(), path)
}

func (t *Types) addSvcs(sds []*descriptor.ServiceDescriptorProto, path typePath) {
	for _, sd := range sds {
		t.svcByFullname[path.Join(sd.GetName())] = sd
	}
}

func (t *Types) addMsgs(mds []*descriptor.DescriptorProto, path typePath) {
	for _, md := range mds {
		path := path.Join(md.GetName())

		t.msgByFullname[path] = md
		t.addMsgs(md.GetNestedType(), path)
		t.addEnums(md.GetEnumType(), path)
	}
}

func (t *Types) addEnums(eds []*descriptor.EnumDescriptorProto, path typePath) {
	for _, ed := range eds {
		t.enumByFullname[path.Join(ed.GetName())] = ed
	}
}
