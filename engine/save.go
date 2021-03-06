package engine

import (
	"bufio"
	"fmt"
	"github.com/mumax/3/cuda"
	"github.com/mumax/3/data"
	"github.com/mumax/3/draw"
	"github.com/mumax/3/oommf"
	"github.com/mumax/3/util"
	"os"
	"path"
	"strings"
)

func init() {
	DeclFunc("Save", Save, "Save space-dependent quantity once, with auto filename")
	DeclFunc("SaveAs", SaveAs, "Save space-dependent with custom filename")
	DeclVar("FilenameFormat", &FilenameFormat, "printf formatting string for output filenames.")
	DeclVar("OutputFormat", &outputFormat, "Format for data files: OVF1_TEXT, OVF1_BINARY, OVF2_TEXT or OVF2_BINARY")
	DeclROnly("OVF1_BINARY", OVF1_BINARY, "OutputFormat = OVF1_BINARY sets binary OVF1 output")
	DeclROnly("OVF2_BINARY", OVF2_BINARY, "OutputFormat = OVF2_BINARY sets binary OVF2 output")
	DeclROnly("OVF1_TEXT", OVF1_TEXT, "OutputFormat = OVF1_TEXT sets text OVF1 output")
	DeclROnly("OVF2_TEXT", OVF2_TEXT, "OutputFormat = OVF2_TEXT sets text OVF2 output")
	DeclFunc("Snapshot", Snapshot, "Save image of quantity")
	DeclVar("SnapshotFormat", &SnapshotFormat, "Image format for snapshots: jpg, png or gif.")
}

var (
	FilenameFormat = "%s%06d"    // formatting string for auto filenames.
	SnapshotFormat = "jpg"       // user-settable snapshot format
	outputFormat   = OVF2_BINARY // user-settable output format
)

// Save once, with auto file name
func Save(q Quantity) {
	fname := autoFname(q.Name(), autonum[q])
	SaveAs(q, fname)
	autonum[q]++
}

// Save under given file name (transparent async I/O).
func SaveAs(q Quantity, fname string) {
	if !path.IsAbs(fname) && !strings.HasPrefix(fname, OD) {
		fname = path.Clean(OD + "/" + fname)
	}
	if path.Ext(fname) == "" {
		fname += ".ovf"
	}
	buffer, recylce := q.Slice()
	if recylce {
		defer cuda.Recycle(buffer)
	}
	info := data.Meta{Time: Time, Name: q.Name(), Unit: q.Unit(), CellSize: q.Mesh().CellSize()}
	data := buffer.HostCopy() // must be copy (async io)
	SaveQue <- func() { saveAs_sync(fname, data, info, outputFormat) }
}

// Save image once, with auto file name
func Snapshot(q Quantity) {
	fname := fmt.Sprintf(OD+FilenameFormat+"."+SnapshotFormat, q.Name(), autonum[q])
	s, r := q.Slice()
	if r {
		defer cuda.Recycle(s)
	}
	data := s.HostCopy() // must be copy (asyncio)
	SaveQue <- func() { snapshot_sync(fname, data) }
	autonum[q]++
}

// synchronous snapshot
func snapshot_sync(fname string, output *data.Slice) {
	util.FatalErr(draw.RenderFile(fname, output, "auto", "auto", arrowSize))
}

// synchronous save
func saveAs_sync(fname string, s *data.Slice, info data.Meta, format OutputFormat) {
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	util.FatalErr(err)
	out := bufio.NewWriter(f)

	switch format {
	case OVF1_TEXT:
		oommf.WriteOVF1(out, s, info, "text")
	case OVF1_BINARY:
		oommf.WriteOVF1(out, s, info, "binary 4")
	case OVF2_TEXT:
		oommf.WriteOVF2(out, s, info, "text")
	case OVF2_BINARY:
		oommf.WriteOVF2(out, s, info, "binary 4")
	default:
		panic("invalid output format")
	}

	out.Flush()
	f.Close()
}

type OutputFormat int

const (
	OVF1_TEXT OutputFormat = iota + 1
	OVF1_BINARY
	OVF2_TEXT
	OVF2_BINARY
)
