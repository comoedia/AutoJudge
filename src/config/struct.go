package config

type FileRoot struct {
	Schema      string `json:"$schema"`
	Problem     `json:"problem"`
	Submit      `json:"submit"`
	Endpoint    `json:"endpoint"`
	Credentials `json:"credentials"`
}

type Problem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Limit       `json:"limit"`
	Try         int `json:"try"`
}

type Limit struct {
	Time   string `json:"time"`
	Memory string `json:"memory"`
}

type Submit struct {
	Lang   `json:"lang"`
	Before `json:"before"`
	File   string `json:"file"`
}

type Lang struct {
	Index int    `json:"index"`
	Str   string `json:"str"`
}

type Before struct {
	Test string `json:"test"`
	Run  string `json:"run"`
}

type Endpoint struct {
	Host      string `json:"host"`
	Resources `json:"resources"`
}

type Resources struct {
	Problem     string `json:"problem"`
	Submissions string `json:"submissions"`
	Submit      string `json:"submit"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
