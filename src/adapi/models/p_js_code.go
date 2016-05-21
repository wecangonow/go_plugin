package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type PJsCode struct {
	Id         int     `orm:"column(id);auto"`
	Name       string  `orm:"column(name);null"`
	Code       string  `orm:"column(code);null"`
	Weight     int     `orm:"column(weight);null"`
	Class      string  `orm:"column(class);size(255);null"`
	Unclass    string  `orm:"column(unclass);size(255);null"`
	Putin      int     `orm:"column(putin);null"`
	Country    string  `orm:"column(country);size(255);null"`
	AppId      string  `orm:"column(app_id);size(255);null"`
	AppName    string  `orm:"column(app_name);size(45);null"`
	Status     int     `orm:"column(status);null"`
	Ecpm       float64 `orm:"column(ecpm);null"`
	ShowNum    int     `orm:"column(show_num);null"`
	MaxShowNum int     `orm:"column(max_show_num);null"`
	StaticUrl  string  `orm:"column(static_url);size(255);null"`
	Replace    int     `orm:"column(replace);null"`
	Webblack   string  `orm:"column(webblack);null"`
	UpdateTime int     `orm:"column(update_time);null"`
	CreateTime int     `orm:"column(create_time);null"`
}

func (t *PJsCode) TableName() string {
	return "p_js_code"
}

func init() {
	orm.RegisterModel(new(PJsCode))
}

// AddPJsCode insert a new PJsCode into database and returns
// last inserted Id on success.
func AddPJsCode(m *PJsCode) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPJsCodeById retrieves PJsCode by Id. Returns error if
// Id doesn't exist
func GetPJsCodeById(id int) (v *PJsCode, err error) {
	o := orm.NewOrm()
	v = &PJsCode{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}
func GetAllPJsCodeByWhere(where string) (list []PJsCode, err error) {

	o := orm.NewOrm()
	query := "select `id`,`name`,`ecpm`,`country`,`class`,`show_num`,`max_show_num`,`app_id`,`replace`,`status`,`update_time`,`webblack`" +
	         " from `p_js_code` where 1 and `status`=1" + where +
			 " order by `ecpm` desc, `update_time` desc"
	_, err = o.Raw(query).QueryRows(&list)

	return list, err
}

// GetAllPJsCode retrieves all PJsCode matches certain condition. Returns empty list if
// no records exist
func GetAllPJsCode(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PJsCode))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []PJsCode
	qs = qs.OrderBy(sortFields...)
	if _, err := qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdatePJsCode updates PJsCode by Id and returns error if
// the record to be updated doesn't exist
func UpdatePJsCodeById(m *PJsCode) (err error) {
	o := orm.NewOrm()
	v := PJsCode{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePJsCode deletes PJsCode by Id and returns error if
// the record to be deleted doesn't exist
func DeletePJsCode(id int) (err error) {
	o := orm.NewOrm()
	v := PJsCode{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PJsCode{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
