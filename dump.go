package trash

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"time"
)

type Dump struct {
	Errs   map[string]Err
	Format string
	Writer io.Writer
}

func NewDump(w io.Writer, format string) *Dump {
	return &Dump{make(map[string]Err), format, w}
}

func (d *Dump) Insert(t string, err string) {
	d.Errs[time.Now().String()] = d.NewErr(t, err)
}

func (d *Dump) Remove(err Err) {
	delete(d.Errs, err.Error())
}

func (d *Dump) Get(err string) Err {
	return d.Errs[err]
}

func (d *Dump) Catch(err string, message string) {
	key := time.Now().String()
	errorr := d.NewErr(err, message)
	d.Errs[key] = errorr
	d.Errs[key].Send(d.Writer)
	d.Errs[key].Log()
}

func (d *Dump) Dump() error {
	data, err := json.MarshalIndent(d.Errs, "", "  ")
	if err != nil {
		return err
	}
	_, err = d.Writer.Write(data)
	if err != nil {
		return err
	}
	d.Errs = make(map[string]Err)
	return nil
}

func (d *Dump) NewErr(err string, message string) Err {
	checksum := base64.StdEncoding.EncodeToString([]byte(err + time.Now().String()))
	switch d.Format {
	case "json":
		d.Insert(err, message)
		return JsonErr{errData: errData{checksum, err, message, 0}}
	case "xml":
		d.Insert(err, message)
		return XmlErr{errData: errData{checksum, err, message, 0}}
	default:
		return nil
	}
}
