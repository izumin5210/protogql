package protoprocessor

import (
	"context"
	"io"
	"io/ioutil"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/pluginpb"
)

type GenerateFunc func(ctx context.Context, fd protoreflect.FileDescriptor, types *Types) (*pluginpb.CodeGeneratorResponse_File, error)

func (f GenerateFunc) Generate(ctx context.Context, fd protoreflect.FileDescriptor, types *Types) (*pluginpb.CodeGeneratorResponse_File, error) {
	return f(ctx, fd, types)
}

type Generator interface {
	Generate(context.Context, protoreflect.FileDescriptor, *Types) (*pluginpb.CodeGeneratorResponse_File, error)
}

func New(g Generator) *Processor {
	return &Processor{g: g}
}

type Processor struct {
	g Generator
}

func (p *Processor) Process(ctx context.Context, r io.Reader, w io.Writer) error {
	req, err := p.readReq(r)
	if err != nil {
		return err
	}

	resp, err := p.processEach(ctx, req)
	if err != nil {
		return err
	}

	err = p.writeResp(w, resp)
	return err
}

func (p *Processor) readReq(r io.Reader) (*pluginpb.CodeGeneratorRequest, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var req pluginpb.CodeGeneratorRequest
	if err = proto.Unmarshal(buf, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (p *Processor) processEach(ctx context.Context, req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	files := new(protoregistry.Files)
	for _, fd := range req.GetProtoFile() {
		rfd, err := protodesc.NewFile(fd, files)
		if err != nil {
			return nil, err
		}
		files.RegisterFile(rfd)
		if err != nil {
			return nil, err
		}
	}

	types := NewTypes()
	types.RegisterFromFiles(files)

	var resp pluginpb.CodeGeneratorResponse
	var errs []error
	for _, file := range req.GetFileToGenerate() {
		fd, err := files.FindFileByPath(file)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		out, err := p.g.Generate(ctx, fd, types)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		resp.File = append(resp.File, out)
	}

	if len(errs) > 0 {
		// TODO: handling
		return nil, errs[0]
	}

	return &resp, nil
}

func (p *Processor) writeResp(w io.Writer, resp *pluginpb.CodeGeneratorResponse) error {
	buf, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}
