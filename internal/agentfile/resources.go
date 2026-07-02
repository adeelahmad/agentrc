package agentfile

import (
	"fmt"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// PopulateLocalResources parses f.CleanedSource with the real Dockerfile
// instruction parser and fills f.LocalResources from every COPY/ADD
// instruction (local — remote ADDs never reach CleanedSource, they were
// extracted into f.RemoteAdds) whose destination lands under /mnt.
func PopulateLocalResources(f *File) error {
	res, err := parser.Parse(strings.NewReader(string(f.CleanedSource)))
	if err != nil {
		return fmt.Errorf("parsing cleaned Dockerfile: %w", err)
	}
	stages, _, err := instructions.Parse(res.AST, nil)
	if err != nil {
		return fmt.Errorf("parsing cleaned Dockerfile instructions: %w", err)
	}

	var resources []LocalResource
	for _, stage := range stages {
		for _, cmd := range stage.Commands {
			dest, sources, line, ok := destAndSources(cmd)
			if !ok {
				continue
			}
			kind := ResourceKindForDest(dest)
			if kind == KindOther {
				continue
			}
			name := ""
			if kind != KindSOP {
				name = pathBase(dest)
			}
			src := ""
			if len(sources) > 0 {
				src = sources[0]
			}
			resources = append(resources, LocalResource{Source: src, Dest: dest, Kind: kind, Name: name, Line: line})
		}
	}
	f.LocalResources = resources
	return reconcileFileBackedSOP(f)
}

// reconcileFileBackedSOP handles SOP's third documented form: authored as
// `COPY ./sop.md /mnt/SOP` or `ADD --remote <url> /mnt/SOP` rather than the
// SOP keyword (spec/index.md §7). Extract only ever sets f.SOP from the
// keyword forms, so this synthesizes it from whichever COPY/ADD instruction
// targets /mnt/SOP — and rejects an Agentfile that declares an SOP more than
// once, whether via the keyword, a local COPY, or a remote ADD.
func reconcileFileBackedSOP(f *File) error {
	var fileBackedLine int
	found := false

	for _, r := range f.LocalResources {
		if r.Kind != KindSOP {
			continue
		}
		if found {
			return fmt.Errorf("line %d: SOP declared more than once (already declared at line %d)", r.Line, fileBackedLine)
		}
		found, fileBackedLine = true, r.Line
	}
	for _, ra := range f.RemoteAdds {
		cleaned, ok := CleanMntPath(ra.Dest)
		if !ok || cleaned != "/mnt/SOP" {
			continue
		}
		if found {
			return fmt.Errorf("line %d: SOP declared more than once (already declared at line %d)", ra.Line, fileBackedLine)
		}
		found, fileBackedLine = true, ra.Line
	}

	if !found {
		return nil
	}
	if f.SOP != nil {
		return fmt.Errorf("line %d: SOP declared more than once (already declared at line %d)", fileBackedLine, f.SOP.Line)
	}
	f.SOP = &SOP{FileBacked: true, Line: fileBackedLine}
	return nil
}

// BaseImageRef returns the first stage's FROM reference from f.CleanedSource
// (as authored — not resolved to a digest).
func BaseImageRef(f *File) (string, error) {
	res, err := parser.Parse(strings.NewReader(string(f.CleanedSource)))
	if err != nil {
		return "", fmt.Errorf("parsing cleaned Dockerfile: %w", err)
	}
	stages, _, err := instructions.Parse(res.AST, nil)
	if err != nil {
		return "", fmt.Errorf("parsing cleaned Dockerfile instructions: %w", err)
	}
	if len(stages) == 0 {
		return "", fmt.Errorf("no FROM instruction found")
	}
	return stages[0].BaseName, nil
}

func destAndSources(cmd instructions.Command) (dest string, sources []string, line int, ok bool) {
	switch c := cmd.(type) {
	case *instructions.CopyCommand:
		dest, sources = c.DestPath, c.SourcePaths
	case *instructions.AddCommand:
		dest, sources = c.DestPath, c.SourcePaths
	default:
		return "", nil, 0, false
	}
	if loc := cmd.Location(); len(loc) > 0 {
		line = loc[0].Start.Line
	}
	return dest, sources, line, true
}

func pathBase(p string) string {
	idx := strings.LastIndex(p, "/")
	if idx < 0 {
		return p
	}
	return p[idx+1:]
}
