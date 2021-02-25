package pokemon

//Defines the properties of a Pokemon to be listed

type Pokemon struct {
	ID    int `json:"ID,omitempty"`
	Name  string `json:"Name,omitempty"`
}