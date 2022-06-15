package protobuf

type Request struct {
	Id     string `json:"id"`
	User   User   `json:"user"`
	Device Device `json:"device"`
	Ext    Ext    `json:"ext"`
}

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Device struct {
	IP string `json:"ip"`
	UA string `json:"ua"`
}

type Ext struct {
	Iphash   int32 `json:"iphash"`
	RegionID int32 `json:"region_id"`
}
