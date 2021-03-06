package engine

import (
	"bufio"
	"fmt"
	"github.com/mumax/3/script"
	"github.com/mumax/3/util"
	"os"
)

var Table = *newTable("table") // output handle for tabular data (average magnetization etc.)

func init() {
	DeclFunc("TableAdd", TableAdd, "Add quantity as a column to the data table.")
	DeclFunc("TableAddVar", TableAddVariable, "Add user-defined variable + name + unit to data table.")
	DeclFunc("TableSave", TableSave, "Save the data table right now (appends one line).")
	DeclFunc("TableAutoSave", TableAutoSave, "Auto-save the data table ever period (s).")
	DeclFunc("TablePrint", TablePrint, "Print anyting in the data table")
	Table.Add(&M)
}

type DataTable struct {
	*bufio.Writer
	info
	outputs []TableData
	autosave
}

func newTable(name string) *DataTable {
	t := new(DataTable)
	t.name = name
	return t
}

func TableAdd(col TableData) {
	Table.Add(col)
}

func TableAddVariable(x script.ScalarFunction, name, unit string) {
	Table.AddVariable(x, name, unit)
}

func (t *DataTable) AddVariable(x script.ScalarFunction, name, unit string) {
	TableAdd(&userVar{x, name, unit})
}

type userVar struct {
	value      script.ScalarFunction
	name, unit string
}

func (x *userVar) Name() string       { return x.name }
func (x *userVar) NComp() int         { return 1 }
func (x *userVar) Unit() string       { return x.unit }
func (x *userVar) average() []float64 { return []float64{x.value.Float()} }

func TableSave() {
	Table.Save()
}

func TableAutoSave(period float64) {
	Table.autosave = autosave{period, Time, 0, nil}
}

func (t *DataTable) Add(output TableData) {
	if t.inited() {
		util.Fatal("data table add ", output.Name(), ": need to add quantity before table is output the first time")
	}
	t.outputs = append(t.outputs, output)
}

func (t *DataTable) Save() {
	t.init()
	fmt.Fprint(t, Time)
	for _, o := range t.outputs {
		vec := o.average()
		for _, v := range vec {
			fmt.Fprint(t, "\t", float32(v))
		}
	}
	fmt.Fprintln(t)
	t.Flush()
	t.count++
}

func (t *DataTable) Println(msg ...interface{}) {
	t.init()
	fmt.Fprintln(t, msg...)
}

func TablePrint(msg ...interface{}) {
	Table.Println(msg...)
}

// open writer and write header
func (t *DataTable) init() {
	if !t.inited() {
		f, err := os.OpenFile(OD+t.name+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		util.FatalErr(err)
		t.Writer = bufio.NewWriter(f)

		// write header
		fmt.Fprint(t, "# t (s)")
		for _, o := range t.outputs {
			if o.NComp() == 1 {
				fmt.Fprint(t, "\t", o.Name(), " (", o.Unit(), ")")
			} else {
				for c := 0; c < o.NComp(); c++ {
					fmt.Fprint(t, "\t", o.Name()+string('x'+c), " (", o.Unit(), ")")
				}
			}
		}
		fmt.Fprintln(t)
		t.Flush()
	}
}

func (t *DataTable) inited() bool {
	return t.Writer != nil
}

func (t *DataTable) flush() {
	if t.Writer != nil {
		t.Flush()
	}
}

// can be saved in table
type TableData interface {
	average() []float64
	Name() string
	Unit() string
	NComp() int
}
