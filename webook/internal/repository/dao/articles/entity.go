package articles

type Article struct {
	Id       int64  `gorm:"primary_key;autoIncrement" bson:"id,omitempty"`
	Title    string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	Content  string `gorm:"type=BLOB" bson:"content,omitempty"`
	AuthorId int64  `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}

// PublishedArticle 衍生类型
type PublishedArticle Article
