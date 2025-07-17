package components

import "github.com/a-h/templ"

func IndexView() templ.Component {
	return Assemble(&PageComponents{
		Main: indexView(),
	})
}

func AdminView(data *AdminViewData) templ.Component {
	return Assemble(&PageComponents{
		Main: adminView(data),
	})
}

func GroupView(data *GroupViewData) templ.Component {
	return Assemble(&PageComponents{
		Main: groupView(data),
	})
}

func SeedView(data *SeedViewData, title string) templ.Component {
	return Assemble(&PageComponents{
		Title: title,
		Main:  seedView(data),
	})
}
