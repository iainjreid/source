package nav

type Nav struct {
	Items []*NavItem
}

type NavItem struct {
	Name   string
	Href   string
	Active bool
}

func NewItem(name string, href string, active bool) *NavItem {
	return &NavItem{
		Name:   name,
		Href:   href,
		Active: active,
	}
}
