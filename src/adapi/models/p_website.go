package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type PWebsite struct {
	Id         int    `orm:"column(id);auto"`
	ClassId    int    `orm:"column(class_id)"`
	Web        string `orm:"column(web);size(255)"`
	UpdateTime int    `orm:"column(update_time);null"`
	CreateTime int    `orm:"column(create_time);null"`
}

func (t *PWebsite) TableName() string {
	return "p_website"
}

func init() {
	orm.RegisterModel(new(PWebsite))
}

// AddPWebsite insert a new PWebsite into database and returns
// last inserted Id on success.
func AddPWebsite(m *PWebsite) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPWebsiteById retrieves PWebsite by Id. Returns error if
// Id doesn't exist
func GetPWebsiteById(id int) (v *PWebsite, err error) {
	o := orm.NewOrm()
	v = &PWebsite{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetMyAllPWebsite() (list orm.ParamsList, err error) {

	o := orm.NewOrm()
	_, err = o.Raw("select web from p_website").ValuesFlat(&list)

	return list, err
}

func GetPWebsitesByLikeWeb(web string) (v []*PWebsite, err error) {

	o := orm.NewOrm()
	_, err = o.Raw("select * from p_website where web like '%?%'", web).QueryRows(&v)

	return v, err
}
func GetPWebsiteByWeb(web string) (v *PWebsite, err error) {

	o := orm.NewOrm()
	err = o.Raw("select * from p_website where web = ?", web).QueryRow(&v)

	return v, err
}
// GetAllPWebsite retrieves all PWebsite matches certain condition. Returns empty list if
// no records exist
func GetAllPWebsite(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PWebsite))
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

	var l []PWebsite
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

// UpdatePWebsite updates PWebsite by Id and returns error if
// the record to be updated doesn't exist
func UpdatePWebsiteById(m *PWebsite) (err error) {
	o := orm.NewOrm()
	v := PWebsite{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePWebsite deletes PWebsite by Id and returns error if
// the record to be deleted doesn't exist
func DeletePWebsite(id int) (err error) {
	o := orm.NewOrm()
	v := PWebsite{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PWebsite{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
