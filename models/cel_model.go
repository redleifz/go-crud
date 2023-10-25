package models

type Cel struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Eventtype   string `json:"eventtype"`
	Eventtime   string `json:"eventtime"`
	Cid_name    string `json:"cid_name"`
	Cid_num     string `json:"cid_num"`
	Cid_ani     string `json:"cid_ani"`
	Cid_rdnis   string `json:"cid_rdnis"`
	Cid_dnid    string `json:"cid_dnid"`
	Exten       string `json:"exten"`
	Context     string `json:"context"`
	Channame    string `json:"channame"`
	Appname     string `json:"appname"`
	Appdata     string `json:"appdata"`
	Amaflags    int    `json:"amaflags"`
	Accountcode string `json:"accountcode"`
	Uniqueid    string `json:"uniqueid"`
	Linkedid    string `json:"linkedid"`
	Peer        string `json:"peer"`
	Userdeftype string `json:"userdeftype"`
	Extra       string `json:"extra"`
}
