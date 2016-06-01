package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type PAdList struct {
	Id                int     `orm:"column(id);auto"`
	Uid               int     `orm:"column(uid);null"`
	Name              string  `orm:"column(name);size(255)"`
	Web               string  `orm:"column(web)"`
	Country           string  `orm:"column(country);size(255);null"`
	Belong            int     `orm:"column(belong)"`
	Class             string  `orm:"column(class);size(255)"`
	Unclass           string  `orm:"column(unclass);size(255);null"`
	Putin             int     `orm:"column(putin);null"`
	Replace           int     `orm:"column(replace)"`
	Weight            int     `orm:"column(weight);null"`
	Size              int     `orm:"column(size);null"`
	Type              int     `orm:"column(type);null"`
	Code              string  `orm:"column(code);null"`
	FlashType         int     `orm:"column(flash_type);null"`
	UrlFlashAddress   string  `orm:"column(url_flash_address);null"`
	FlashAddress      string  `orm:"column(flash_address);null"`
	IconType          int     `orm:"column(icon_type);null"`
	UrlIcon           string  `orm:"column(url_icon);null"`
	Icon              string  `orm:"column(icon);null"`
	Status            int     `orm:"column(status);null"`
	InjectionLocation int     `orm:"column(injection_location);null"`
	Ecpm              float64 `orm:"column(ecpm);null"`
	ShowNum           int     `orm:"column(show_num);null"`
	MaxShowNum        int     `orm:"column(max_show_num);null"`
	StaticUrl         string  `orm:"column(static_url);null"`
	Webblack          string  `orm:"column(webblack);null"`
	UpdateTime        int     `orm:"column(update_time);null"`
	CreateTime        int     `orm:"column(create_time);null"`
}

func (t *PAdList) TableName() string {
	return "p_ad_list"
}

func init() {
	orm.RegisterModel(new(PAdList))
}

// AddPAdList insert a new PAdList into database and returns
// last inserted Id on success.
func AddPAdList(m *PAdList) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPAdListById retrieves PAdList by Id. Returns error if
// Id doesn't exist
func GetPAdListById(id int) (v *PAdList, err error) {
	o := orm.NewOrm()
	v = &PAdList{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetAllPAdListByWhere(where string) (list []PAdList, err error) {

	o := orm.NewOrm()
	query := "select `id`,`name`,`size`,`web`,`country`,`class`,`show_num`,`max_show_num`,`weight`,`replace`,`ecpm`,`update_time`,`static_url`,`webblack`" +
	" from `p_ad_list` where 1 and `replace`=3 and `status`=1" + where +
	" order by `ecpm` desc, `update_time` desc"
	_, err = o.Raw(query).QueryRows(&list)

	return list, err
}

// GetAllPAdList retrieves all PAdList matches certain condition. Returns empty list if
// no records exist
func GetAllPAdList(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PAdList))
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

	var l []PAdList
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

// UpdatePAdList updates PAdList by Id and returns error if
// the record to be updated doesn't exist
func UpdatePAdListById(m *PAdList) (err error) {
	o := orm.NewOrm()
	v := PAdList{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePAdList deletes PAdList by Id and returns error if
// the record to be deleted doesn't exist
func DeletePAdList(id int) (err error) {
	o := orm.NewOrm()
	v := PAdList{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PAdList{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
