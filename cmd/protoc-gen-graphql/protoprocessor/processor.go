package protoprocessor

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

type GenerateFunc func(ctx context.Context, file string, types *Types) (*plugin.CodeGeneratorResponse_File, error)

func (f GenerateFunc) Generate(ctx context.Context, file string, types *Types) (*plugin.CodeGeneratorResponse_File, error) {
	return f(ctx, file, types)
}

type Generator interface {
	Generate(context.Context, string, *Types) (*plugin.CodeGeneratorResponse_File, error)
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

func (p *Processor) readReq(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var req plugin.CodeGeneratorRequest
	if err = proto.Unmarshal(buf, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (p *Processor) processEach(ctx context.Context, req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	types := NewTypes()
	for _, f := range req.ProtoFile {
		types.AddFile(f)
	}

	var resp plugin.CodeGeneratorResponse
	var errs []error
	for _, file := range req.GetFileToGenerate() {
		out, err := p.g.Generate(ctx, file, types)
		if err != nil {
			errs = append(errs, err)
		}
		resp.File = append(resp.File, out)
	}

	if len(errs) > 0 {
		// TODO: handling
		return nil, errs[0]
	}

	return &resp, nil
}

func (p *Processor) writeResp(w io.Writer, resp *plugin.CodeGeneratorResponse) error {
	buf, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}
