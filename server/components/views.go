package components

import "github.com/a-h/templ"

func IndexView() templ.Component {
	return Assemble(&PageComponents{
		Header: indexHeader(),
		Main:   indexView(),
	})
}

func AdminView(data *AdminViewData) templ.Component {
	return Assemble(&PageComponents{
		Main: adminView(data),
	})
}

func GroupView(data *GroupViewData) templ.Component {
	return Assemble(&PageComponents{
		Header: header(data.Heading),
		Main:   groupView(data),
	})
}

func SeedView(data *SeedViewData) templ.Component {
	return Assemble(&PageComponents{
		Title:  data.Title,
		Header: header(data.Heading),
		Main:   seedView(data),
	})
}
