# Golang Contributions

## Common
 - [Google style guide](https://github.com/google/styleguide/tree/gh-pages/go)
 - code must be formatted with ```go fmt```
 - use **PascalCase** for functions that will be available outside the package, **camelCase** for inner functions
 - use ```any``` instead of ```interface{}```
    ```go
    // Bad
    func foo(a interface{})

    //Good
    func foo(a any)
    ```
 - don't ignore errors
    ```go
        // Bad
        res, _ := DoSmth()

        //Good
        res, err := DoSmth()
        if err != nil {
            //
        }
    ```
 - use ```WaitGroup.Go(func)``` instead of 
    ```go   
    WaitGroup.Add(1)
    go func() {
        defer WaitGroup.Done()

    }()
    ```

## Structure
Source - https://habr.com/ru/articles/911018/
```
├── cmd/ 
│   ├── api/   
│   └── worker/
├── internal/ 
│   ├── app/ 
│   ├── domain/ 
│   │   ├── models/
│   │   ├── rules/
│   │   ├── events/
│   │   └── ports/
│   │       ├── repository
│   │       └── service
│   ├── infrastructure/ 
│   │   ├── adapters/
│   │   │   ├── cache      
│   │   │   ├── logger    
│   │   │   └── router    
│   │   ├── clients/
│   │   ├── persistence/
│   │   └── services/
│   └── interfaces/ 
│       └── http/
│           ├── dto/
│           ├── handlers/
│           └── server/
│               └── middleware/
├── config/ 
└── pkg/    
    └── testutils/
```
 
## Commits
 - use ```go mod tidy``` to delete unused imports
 - use ```go test -v ./...``` to run all tests
 
