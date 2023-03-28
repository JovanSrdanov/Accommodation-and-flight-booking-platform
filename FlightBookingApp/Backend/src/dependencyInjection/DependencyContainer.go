package dependencyInjection

import "fmt"

type DependencyContainer struct {
	container map[string]interface{}
}

func NewDependencyContainer() *DependencyContainer {
	return &DependencyContainer{
		container: make(map[string]interface{}),
	}
}

func (depCont *DependencyContainer) AddRepository(entityName string, dependency interface{}) {
	completeName := entityName + "-repository"
	depCont.container[completeName] = dependency
}

func (depCont *DependencyContainer) GetRepository(entityName string) interface{} {
	completeName := entityName + "-repository"
	return depCont.container[completeName]
}

func (depCont *DependencyContainer) AddService(entityName string, dependency interface{}) {
	completeName := entityName + "-service"
	depCont.container[completeName] = dependency
}

func (depCont *DependencyContainer) GetService(entityName string) interface{} {
	completeName := entityName + "-service"
	return depCont.container[completeName]
}

func (depCont *DependencyContainer) AddController(entityName string, dependency interface{}) {
	completeName := entityName + "-controller"
	depCont.container[completeName] = dependency
}
func (depCont *DependencyContainer) GetController(entityName string) interface{} {
	completeName := entityName + "-controller"
	return depCont.container[completeName]
}

func (depCont *DependencyContainer) AddEntityDependencyBundle(entityName string, repository, service, controller interface{}) {
	depCont.AddRepository(entityName, repository)
	depCont.AddService(entityName, service)
	depCont.AddController(entityName, controller)
}

func (depCont *DependencyContainer) PrintAllDependencies() {
	keys := make([]string, len(depCont.container))
	i := 0
	for k := range depCont.container {
		keys[i] = k
		i++
	}
	fmt.Printf("Dependency container content: ")
	fmt.Println(keys)
}
