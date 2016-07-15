package postgres

import ()

type (
	IPostgresDriver interface {
		Insert(p interface{}) error
		InsertMulti(p interface{}) error
		Delete(p interface{}) error
		DeleteInIds(p interface{}, column string, args ...string) error
		Update(p, query interface{}, cols ...string) error //p对需更新的字段赋值；codition对条件字段赋值
		UpdateById(p interface{}, id string) error
		UpdateCols(p interface{}, query interface{}, cols ...string) error
		UpdateOmit(p interface{}, query string, omit ...string) error
		UpdateOmitById(p interface{}, id string, omit ...string) error
		GetOne(p interface{}) (bool, error)
		GetOneSort(p interface{}, sort ...string) (bool, error)
		GetAll(p, query interface{}) error
		GetPage(p, query interface{}, page, size int) error
		Count(query interface{}) (int64, error)
	}
)

type (
	PostgresDriver struct {
		Name string
	}
)

func NewPostgresDriver(name string) IPostgresDriver {
	return &PostgresDriver{Name: name}
}

func (this *PostgresDriver) Insert(p interface{}) (err error) {
	//defer CatchPanic(&err, "PostgresDriver Insert", this.Name)

	o, err := NewPostgresXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = o.InsertOne(p)
	return err
}

///*
//不能保证同时插入成功
//*/
//func (this *PostgresDriver) InsertSync(p interface{}, q interface{}) (err error) {
//	defer CatchPanic(&err, "PostgresDriver InsertSync", this.Name)
//
//	o := NewXormPostgres()
//	_, err = o.Insert(p, q)
//	return err
//}

/*
批量插入，同步执行
*/
func (this *PostgresDriver) InsertMulti(p interface{}) (err error) {
	//defer CatchPanic(&err, "PostgresDriver InsertMulti", this.Name)
	o, err := NewPostgresXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = o.Insert(p)
	return err
}

func (this *PostgresDriver) Delete(p interface{}) (err error) {
	//defer CatchPanic(&err, "PostgresDriver Delete", this.Name)
	o, err := NewPostgresXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = o.Delete(p)
	return err
}

func (this *PostgresDriver) DeleteInIds(p interface{}, column string, args ...string) (err error) {
	//defer CatchPanic(&err, "PostgresDriver Delete", this.Name)
	o, err := NewPostgresXormInit(this.Name)
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
func (this *PostgresDriver) Update(p interface{}, query interface{}, cols ...string) (err error) {
	o, err := NewPostgresXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	table := o.TableInfo(p)
	if cols == nil || len(cols) == 0 {
		cols = table.ColumnsSeq()
	}
	_, err = o.Cols(cols...).Update(p, query)
	return err
}

/*
更新所有列
*/
func (this *PostgresDriver) UpdateById(p interface{}, id string) (err error) {
	//defer CatchPanic(&err, "PostgresDriver Update", this.Name)
	o, err := NewPostgresXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = o.AllCols().Id(id).Update(p)
	return err
}

/*
更新指定列
*/
func (this *PostgresDriver) UpdateCols(p interface{}, query interface{}, cols ...string) (err error) {
	//defer CatchPanic(&err, "PostgresDriver UpdateCols", this.Name)
	o, err := NewPostgresXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = o.Cols(cols...).Update(p, query)
	return err
}

/*
更新所有列，排除指定列
*/
func (this *PostgresDriver) UpdateOmit(p interface{}, query string, omit ...string) (err error) {
	//defer CatchPanic(&err, "PostgresDriver UpdateOmit", this.Name)
	o, err := NewPostgresXormInit(this.Name)
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
func (this *PostgresDriver) UpdateOmitById(p interface{}, id string, omit ...string) (err error) {
	//defer CatchPanic(&err, "PostgresDriver UpdateOmit", this.Name)
	o, err := NewPostgresXormInit(this.Name)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = o.AllCols().Omit(omit...).Id(id).Update(p)
	return err
}

func (this *PostgresDriver) GetOne(p interface{}) (r bool, err error) {
	//defer CatchPanic(&err, "PostgresDriver GetOne", this.Name)
	o, err := NewPostgresXormInit(this.Name)
	if err != nil {
		return false, err
	}
	defer o.Close()

	r, err = o.Get(p)
	return
}

func (this *PostgresDriver) GetOneSort(p interface{}, sort ...string) (r bool, err error) {
	//defer CatchPanic(&err, "PostgresDriver GetOne", this.Name)

	o, err := NewPostgresXormInit(this.Name)
	if err != nil {
		return false, err
	}
	defer o.Close()

	r, err = o.Desc(sort...).Get(p)
	return
}

func (this *PostgresDriver) GetAll(p interface{}, query interface{}) (err error) {
	//defer CatchPanic(&err, "PostgresDriver GetAll", this.Name)
	o, err := NewPostgresXormInit(this.Name)
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

func (this *PostgresDriver) GetPage(p, query interface{}, skip, limit int) (err error) {
	//defer CatchPanic(&err, "PostgresDriver GetPage", this.Name, skip, limit)
	o, err := NewPostgresXormInit(this.Name)
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

func (this *PostgresDriver) Count(query interface{}) (n int64, err error) {
	//defer CatchPanic(&err, "PostgresDriver GetCount", this.Name)
	o, err := NewPostgresXormInit(this.Name)
	if err != nil {
		return 0, err
	}
	defer o.Close()

	n, err = o.Count(query)

	return
}
