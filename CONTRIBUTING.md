Contributions are welcome! Working on code that generates code can be challenging (in a good way!). Since most changes are to generated code, I typically start by modifying one of the generated test mocks [here](generator/testmoqs) (ignore the `DO NOT EDIT!` warning in this case). Then fixup any broken test in the same package, add some tests, rinse and repeat. When that one mock is functioning as expected, next update the code-generation code in the parent package so that it generates identical code. Iterating on the code-generation code typically involve running a command line to generate an example mock:
```shell
go run main.go --debug --destination xxx.go github.com/myshkin5/moqueries/generator/testmoqs Usual
```

The previous command generates a scratch mock in a `xxx.go` file (don't check this file in!) which can then be compared to your hand-altered mock (assuming you modified the [`moq_usual_test.go`](generator/testmoqs/moq_usual_test.go)). When the scratch mock matches your hand-altered mock, follow the instructions below to regenerate all the mocks and run all the tests before submitting a PR.

# CI/CD
Moqueries uses CircleCI for CI/CD builds. If you don't see any jobs run when creating or updating a PR, [CircleCI says](https://circleci.com/docs/2.0/oss/#build-pull-requests-from-forked-repositories) that you are probably following your fork which will then trigger jobs under your own account. Please unfollow your fork.

You can iterate more quickly by locally running some of the [same commands CircleCI is running](.circleci/config.yml) (all the following commands assume you are running them from the root moqueries working directory):
1. Build the binary:
    ```shell
    go build -o $GOPATH/bin/moqueries github.com/myshkin5/moqueries
    ```
2. Generate all the mocks:
    ```shell
    go generate ./...
    ```
3. Run all the tests:
    ```shell
    go test ./...
    ```
4. Run all the linters (may require installing [`golangci-lint`](https://golangci-lint.run/usage/install/#local-installation)):
    ```shell
    golangci-lint run
    ```

# Have fun!
