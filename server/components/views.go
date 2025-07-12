package components

import "github.com/a-h/templ"

func IndexView() templ.Component {
	return wrapMain(indexView())
}

func AdminView(data *AdminViewData) templ.Component {
	return wrapMain(adminView(data))
}

func GroupView(data *GroupViewData) templ.Component {
	return wrapMain(groupView(data))
}
