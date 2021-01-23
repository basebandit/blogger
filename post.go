package blogger

type frontMatter struct {
	Author string   `yaml:"author"`
	Title  string   `yaml:"title"`
	Date   string   `yaml:"date"`
	Tags   []string `yaml:"tags"`
}
