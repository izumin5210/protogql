package protoutil

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/pluginpb"
)

type ProtocGen struct {
	Generate func(context.Context, *ProtocGenRequest) error
}

type ProtocGenRequest struct {
	Request         *pluginpb.CodeGeneratorRequest
	FileRegistry    *protoregistry.Files
	FilesToGenerate []protoreflect.FileDescriptor

	genFiles []*GeneratedFile
}

func (r *ProtocGenRequest) GenerateFile(name string) *GeneratedFile {
	f := &GeneratedFile{name: name}
	r.genFiles = append(r.genFiles, f)
	return f
}

type GeneratedFile struct {
	name string
	buf  bytes.Buffer
}

func (f *GeneratedFile) Write(p []byte) (n int, err error) {
	return f.buf.Write(p)
}

func (g *ProtocGen) Run(ctx context.Context, r io.Reader, w io.Writer) error {
	req, err := g.readReq(r)
	if err != nil {
		return err
	}

	resp, err := g.processEach(ctx, req)
	if err != nil {
		return err
	}

	err = g.writeResp(w, resp)
	return err
}

func (g *ProtocGen) readReq(r io.Reader) (*pluginpb.CodeGeneratorRequest, error) {
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

func (g *ProtocGen) processEach(ctx context.Context, rawReq *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	req := &ProtocGenRequest{
		Request:      rawReq,
		FileRegistry: new(protoregistry.Files),
	}

	for _, fd := range rawReq.GetProtoFile() {
		rfd, err := protodesc.NewFile(fd, req.FileRegistry)
		if err != nil {
			return nil, err
		}
		req.FileRegistry.RegisterFile(rfd)
		if err != nil {
			return nil, err
		}
	}

	for _, file := range rawReq.GetFileToGenerate() {
		fd, err := req.FileRegistry.FindFileByPath(file)
		if err != nil {
			return nil, err
		}
		req.FilesToGenerate = append(req.FilesToGenerate, fd)
	}

	err := g.Generate(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := new(pluginpb.CodeGeneratorResponse)
	for _, f := range req.genFiles {
		resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(f.name),
			Content: proto.String(f.buf.String()),
		})
	}

	return resp, nil
}

func (g *ProtocGen) writeResp(w io.Writer, resp *pluginpb.CodeGeneratorResponse) error {
	buf, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}
