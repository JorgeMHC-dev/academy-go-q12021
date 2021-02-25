package pokemon

//Defines the properties of a Pokemon to be listed

type Pokemon struct {
	ID    int `json:"ID"`
	Name  string `json:"name,omitempty"`
}