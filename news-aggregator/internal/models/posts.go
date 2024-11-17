package models

type Post struct {
	Id int //`json:"id"`
	//PosrUUID uuid.UUID `json:"post_uuid"`
	Title   string //`json:"title"`
	Content string //`json:"content"`
	PubTime int64  //`json:"publication_date"`
	Link    string //`json:"link"`
}
