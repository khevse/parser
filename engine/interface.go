package engine

type ITOCItem interface {
	// Url of the page with content where have urls of target items
	Href() string

	// Initialize urls of target items
	// parameters:
	//  url - string - url of the page with target items
	SetChildren(string)

	// Returns list of target items
	GetChildren() []ITarget
}

type ITarget interface {
	// Url of the page with description
	Href() string

	// Set description of element
	// parameters:
	//  url - string - url of page with description
	SetDescription(string)
}
