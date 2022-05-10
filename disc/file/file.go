package file

type HashFunc string
type Tag string

type File struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Path       []Path    `json:"path"`
	Size       uint64    `json:"size"`
	SimpleHash string    `json:"simple_hash"`
	Hash       string    `json:"hash"`
	HashFunc   HashFunc  `json:"hash_func"`
	FileDesc   string    `json:"file_desc"`
	Metas      []Meta    `json:"metas"`
	Explains   []Explain `json:"explains"`
	Tags       []Tag     `json:"tags"`
	Links      []Link    `json:"links"`
	SubFile    []File    `json:"sub_file"`
}

type Path struct {
	Id     string `json:"id"`
	FileId string `json:"file_id"`
	Path   string `json:"path"`
	Disk   Disk   `json:"disk"`
}

type Disk struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	RootPath string `json:"root_path"`
}

type Meta struct {
	Id     string `json:"id"`
	FileId string `json:"file_id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type Explain struct {
	Id     string `json:"id"`
	FileId string `json:"file_id"`
	Text   string `json:"text"`
}

type Link struct {
	Id           string `json:"id"`
	FileIdSource string `json:"file_id_source"`
	FileIdTarget string `json:"file_id_target"`
}
