package routes

type author struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	LastName  string `json:"lastName"`
	CreatedAt string `json:"createdAt"`
}

type queryBooks struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	PublishedAt string `json:"publishedAt"`
	Name        string `json:"name"`
	LastName    string `json:"lastName"`
}

type book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	AuthorID    int    `json:"authorID"`
	PublishedAt string `json:"publishedAt"`
}
