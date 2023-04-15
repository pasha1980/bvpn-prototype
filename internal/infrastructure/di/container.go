package di

import "github.com/sarulabs/di"

var container di.Container

func Get(name string) any {
	return container.Get(name)
}

func GetContainer() di.Container {
	return container
}

func Set(c di.Container) {
	container = c
}
