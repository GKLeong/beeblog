package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strconv" //转换格式包
	"time"
)

const (
	_DB_NAME      = "beeblog"
	_MYSQL_DRIVER = "mysql"
	_DB_USER      = "root"
	_DB_PWD       = "123456"
)

/*创建模型*/
type BeeCategory struct {
	Id              int64
	Title           string
	ParentId        int64
	Created         time.Time `orm:"column(created);type(timestamp);auto_now_add"`
	View            int64     `orm:"index"`
	TopicTime       time.Time `orm:"column(topic_time);type(timestamp);auto_now_add"`
	TopicCount      int64
	TopiclastUserId int64
}

/*创建模型*/
type BeeTopic struct {
	Id               int64
	Uid              int64
	CateId           int64
	Title            string
	Content          string `orm:"size(5000)"`
	Attachment       string
	Created          time.Time `orm:"index"`
	Updated          time.Time `orm:"index"`
	Views            int64     `orm:"index"`
	Author           string
	ReplyTime        time.Time `orm:"index"`
	ReplyCount       int64
	RepleylastUserId int64
}

/*创建comment模型*/
type BeeComment struct {
	Id      int64
	Tid     int64
	Name    string
	Content string `orm:"size(1000)"`
	Email   string
	Created time.Time `orm:"inedx"`
}

func RegisterDB() {
	/*if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}*/
	/*注册模型*/
	orm.RegisterModel(new(BeeCategory), new(BeeTopic), new(BeeComment))
	/*注册驱动*/
	orm.RegisterDriver(_MYSQL_DRIVER, orm.DRMySQL)
	/*默认数据库必须有个叫做default  数据库名称，数据库连接以及编码设置*/
	orm.RegisterDataBase("default", _MYSQL_DRIVER, "root:123456@/beeblog?charset=utf8")
}

func AddCategory(name string) error {
	o := orm.NewOrm() //获取orm对象
	o.Using("beeblog")
	cate := &BeeCategory{Title: name} //创建category 对象 //目前cate相当于指针
	qs := o.QueryTable("bee_category")
	err := qs.Filter("title", name).One(cate) //one 传递参数是指向某个地址的指针 取址
	if err == nil {
		return err
	}
	//插入操作
	_, err = o.Insert(cate)
	if err != nil {
		return err
	}
	return nil
}

//获取分类  ** return 元素类型为beecategory 的slice
func GetAllCategory() ([]*BeeCategory, error) {
	o := orm.NewOrm()                //获取orm对象
	cates := make([]*BeeCategory, 0) //定义 BeeCategory 的slice
	qs := o.QueryTable("bee_category")
	_, err := qs.All(&cates)
	return cates, err
}

//删除分类
func DelCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64) //id 10：十进制  64:64int
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &BeeCategory{Id: cid}
	_, err = o.Delete(cate)
	if err != nil {
		return err
	}
	return nil
}

//文章操作方法 -----------------------------------------
func AddTopic(title, content, cate_id string) error {
	o := orm.NewOrm()
	cateid, err := strconv.ParseInt(cate_id, 10, 64)
	if err != nil {
		return err
	}
	topic := &BeeTopic{
		Title:     title,
		CateId:    cateid,
		Content:   content,
		Created:   time.Now(),
		Updated:   time.Now(),
		ReplyTime: time.Now(),
	}
	_, err = o.Insert(topic)
	return err
}

//文章读取数据
func GetAllTopic(IsHome bool, cateid string) ([]*BeeTopic, error) {
	o := orm.NewOrm()
	topics := make([]*BeeTopic, 0)
	qs := o.QueryTable("bee_topic")
	var err error
	if IsHome {
		if len(cateid) > 0 && cateid != " " {
			cid, _ := strconv.ParseInt(cateid, 10, 64)
			_, err = qs.Filter("CateId", cid).OrderBy("-id").All(&topics)
		} else {
			_, err = qs.OrderBy("-created").All(&topics)
		}

	} else {
		_, err = qs.All(&topics)
	}

	return topics, err
}

//文章显示数据
func GetTopic(sid string) (*BeeTopic, error) {
	tid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	topic := new(BeeTopic)
	qs := o.QueryTable("bee_topic")
	err = qs.Filter("id", tid).One(topic)
	if err != nil {
		return nil, err
	}
	topic.Views++
	_, err = o.Update(topic)
	if err != nil {
		return nil, err
	}
	return topic, err
}

//编辑文章
func EditTopics(sid string) (*BeeTopic, error) {
	tid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	topic := new(BeeTopic)
	qs := o.QueryTable("bee_topic")
	err = qs.Filter("id", tid).One(topic)
	if err != nil {
		return nil, err
	}
	return topic, err
}

//修改文章
func UpdataTopic(sid, title, content, cate_id string) (*BeeTopic, error) {
	tid, err := strconv.ParseInt(sid, 10, 64)
	cateid, err := strconv.ParseInt(cate_id, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	topic := new(BeeTopic)
	err = o.Read(&BeeTopic{Id: tid})
	if err != nil {
		return nil, err
	} else {
		//		topic := &BeeTopic{
		//			Id:        tid,
		//			Title:     title,
		//			Content:   content,
		//			Created:   time.Now(),
		//			Updated:   time.Now(),
		//			ReplyTime: time.Now()}
		topic.Id = tid
		topic.CateId = cateid
		topic.Title = title
		topic.Content = content
		topic.Created = time.Now()
		topic.Updated = time.Now()
		topic.ReplyTime = time.Now()
		o.Update(topic)
		return topic, err
	}
}

//删除文章
func ArticleDel(sid string) (*BeeTopic, error) {
	tid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	topic := new(BeeTopic)
	o.Delete(&BeeTopic{Id: tid})
	return topic, err
}

/*留言*/
func Replayadd(tid, name, email, content string) (*BeeComment, error) {
	o := orm.NewOrm()
	newtid, err := changeId(tid)
	replay := &BeeComment{
		Tid:     newtid,
		Name:    name,
		Content: content,
		Email:   email,
		Created: time.Now(),
	}
	_, err = o.Insert(replay)
	return nil, err

}

/*获取该文章下的所有留言*/
func GetReplays(tid string) ([]*BeeComment, error) {
	var err error
	tidNum, err := changeId(tid)
	o := orm.NewOrm()
	relays := make([]*BeeComment, 0)
	qs := o.QueryTable("BeeComment")
	if len(tid) > 0 {
		_, err = qs.Filter("tid", tidNum).All(&relays)
		if err != nil {
			return nil, err
		}
	}

	return relays, err
}

/*删除文章中的留言*/
func DeleteRepaly(rid string) (*BeeComment, error) {
	var err error
	nrid, err := changeId(rid)
	o := orm.NewOrm()
	replay := &BeeComment{
		Id: nrid,
	}
	o.Delete(replay)
	return replay, err
}
func changeId(id string) (nid int64, err error) {
	hid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 1, err
	}
	return hid, err
}
