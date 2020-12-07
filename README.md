
<img src="https://procyon-projects.github.io/img/logo.png" width="128">

# Procyon Peas
[![codecov](https://codecov.io/gh/procyon-projects/procyon-peas/branch/master/graph/badge.svg?token=YQYFX375ET)](https://codecov.io/gh/procyon-projects/procyon-peas)
[![Build Status](https://travis-ci.com/procyon-projects/procyon-peas.svg?branch=master)](https://travis-ci.com/procyon-projects/procyon-peas)
[![Gitter](https://badges.gitter.im/procyon-projects/community.svg)](https://gitter.im/procyon-projects/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![PkgGoDev](https://pkg.go.dev/badge/procyon-projects/procyon)](https://pkg.go.dev/github.com/procyon-projects/procyon-peas)

This gives you a basic understanding of Procyon Peas Module. It covers
components provided by the framework, such as Pea Processors and Initializers.

Note that you need to register pea processors and initializers by using the function **core.Register**.

## Pea Definition Registry Processor
It's used to do something after Pea Definition Registry is initialized.
```go
type PeaDefinitionRegistryProcessor interface {
	AfterPeaDefinitionRegistryInitialization(registry PeaDefinitionRegistry)
}
```

## Pea Factory Processor
It's used to do something after Pea Factory is initialized.
```go
type PeaFactoryProcessor interface {
	AfterPeaFactoryInitialization(factory ConfigurablePeaFactory)
}
```

## Pea Processors and Initializers
**BeforePeaInitialization**, **InitializePea** and **AfterPeaInitialization** are invoked respectively. 

### Processor
Pea Processors are used to manipulate the instance while being created. For example, Binding
the configuration properties are done by using Pea Processors.  

You can look into [ConfigurationPropertiesBindingProcessor](https://github.com/procyon-projects/procyon-context/blob/master/processor.go#L44) for more information.
```go
type PeaProcessor interface {
	BeforePeaInitialization(peaName string, pea interface{}) (interface{}, error)
	AfterPeaInitialization(peaName string, pea interface{}) (interface{}, error)
}
```

### Initializer
Pea Initializers are used to initialize Pea instances. You can use to initialize your peas. It is invoked
while the instance are created. 
```go
type PeaInitializer interface {
	InitializePea() error
}
```

## License
Procyon Framework is released under version 2.0 of the Apache License
