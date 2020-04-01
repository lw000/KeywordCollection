package table

type TConfig struct {
	Id      int    `db:"id,omitempty"`
	Name    string `db:"name,omitempty"`
	Group   string `db:"group,omitempty"`
	Title   string `db:"title,omitempty"`
	Tip     string `db:"tip,omitempty"`
	Type    string `db:"type,omitempty"`
	Value   string `db:"value,omitempty"`
	Context string `db:"context,omitempty"`
	Rule    string `db:"rule,omitempty"`
	Extend  string `db:"extend,omitempty"`
}
