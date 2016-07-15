package oracle

import (
	"fmt"
)

var (
	_ = fmt.Printf
)

type (
	IOracleDriver interface {
		Insert(p interface{}) error
		Delete(p interface{}) error
		DeleteIn(p interface{}, column string, args ...string) error
		Update(p, query interface{}, cols ...string) error //p对需更新的字段赋值；codition对条件字段赋值
		UpdateById(p interface{}, id string) error
		UpdateOmit(p interface{}, query string, omit ...string) error
		UpdateOmitById(p interface{}, id string, omit ...string) error
		Get(p interface{}) (bool, error)
		GetSort(p interface{}, ty int, sort ...string) (bool, error)
		Find(p, query interface{}) error
		FindPage(p, query interface{}, page, size int) error
		Count(query interface{}) (int64, error)
	}
)

type (
	OracleDriver struct {
		Name string
	}
)

func NewOracleDriver(name string) IOracleDriver {
	return &OracleDriver{Name: name}
}

func (this *OracleDriver) Insert(p interface{}) (err error) {
	o, err := NewOracleXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = o.Insert(p)
	return err
}

func (this *OracleDriver) Delete(p interface{}) (err error) {
	o, err := NewOracleXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = o.Delete(p)
	return err
}

func (this *OracleDriver) DeleteIn(p interface{}, column string, args ...string) (err error) {
	o, err := NewOracleXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = o.In(column, args).Delete(p)
	return err
}

/*
更新所有列
*/
func (this *OracleDriver) Update(p interface{}, query interface{}, cols ...string) (err error) {
	o, err := NewOracleXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	table := o.TableInfo(p)
	if cols == nil || len(cols) == 0 {
		cols = table.ColumnsSeq()
	} else {
		trs := make([]string, 0)
		for _, rs := range table.ColumnsSeq() {
			for _, s := range cols {
				if rs == o.ColumnMapper.Obj2Table(s) {
					trs = append(trs, s)
					break
				}
			}
		}
		cols = trs
	}
	_, err = o.Cols(cols...).Update(p, query)

	return err
}

/*
更新所有列
*/
func (this *OracleDriver) UpdateById(p interface{}, id string) (err error) {
	o, err := NewOracleXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = o.AllCols().Where(`"ID" = ?`, id).Update(p)
	return err
}

/*
更新所有列，排除指定列
*/
func (this *OracleDriver) UpdateOmit(p interface{}, query string, omit ...string) (err error) {
	o, err := NewOracleXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = o.AllCols().Omit(omit...).Where(query).Update(p)
	return err
}

/*
更新所有列，排除指定列
*/
func (this *OracleDriver) UpdateOmitById(p interface{}, id string, omit ...string) (err error) {
	o, err := NewOracleXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = o.AllCols().Omit(omit...).Where(`"ID" = ?`, id).Update(p)
	return err
}

func (this *OracleDriver) Get(p interface{}) (r bool, err error) {
	o, err := NewOracleXormInit(this.Name)
	if err != nil {
		return false, err
	}
	defer o.Close()

	r, err = o.Get(p)
	return
}

/*
param : ty -1 ASC; 1 : DESC
*/
func (this *OracleDriver) GetSort(p interface{}, ty int, sort ...string) (r bool, err error) {
	o, err := NewOracleXormInit(this.Name)
	if err != nil {
		return false, err
	}
	defer o.Close()

	if ty == -1 {
		r, err = o.Asc(sort...).Get(p)
	} else {
		r, err = o.Desc(sort...).Get(p)
	}
	return
}

func (this *OracleDriver) Find(p interface{}, query interface{}) (err error) {
	o, err := NewOracleXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	if query == nil {
		err = o.Find(p)
	} else {
		err = o.Find(p, query)
	}
	return
}

func (this *OracleDriver) FindPage(p, query interface{}, skip, limit int) (err error) {
	o, err := NewOracleXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	if query == nil {
		err = o.Limit(limit, skip).Find(p)
	} else {
		err = o.Limit(limit, skip).Find(p, query)
	}
	return
}

func (this *OracleDriver) Count(query interface{}) (n int64, err error) {
	o, err := NewOracleXormInit(this.Name)
	if err != nil {
		return 0, err
	}
	defer o.Close()

	n, err = o.Count(query)
	return
}
