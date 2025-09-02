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
		Header: groupHeader(),
		Main:   groupView(data),
	})
}

func SeedView(data *SeedViewData) templ.Component {
	return Assemble(&PageComponents{
		Title:  data.Title,
		Header: seedHeader(data.Seed.URL),
		Main:   seedView(data),
	})
}

func ErrorView(data *ErrorViewData) templ.Component {
	return Assemble(&PageComponents{
		Title:  data.Title,
		Header: errorHeader(data),
		Main:   errorView(data),
	})
}

func GeneratorView() templ.Component {
	return Assemble(&PageComponents{
		Title:  "Generátor citací",
		Header: header("Generátor citací"),
		Main:   generatorView(),
	})
}
