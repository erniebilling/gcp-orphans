default: all

GO_PACKAGES = $$(go list ./... | grep -v vendor | grep -v pks-api-release)
GO_TESTS = $$(go list ./... | grep -v vendor | grep -v pks-api-release | grep -v smoke)
GO_FILES = $$(find . -name "*.go" | grep -v vendor | grep -v "./.history" | grep -v pks-api-release | uniq)

build:
	go build -o gcp-orphans ./main.go

unit-test:
	@go test ${GO_TESTS}

fmt:
	gofmt -s -l -w $(GO_FILES)

vet:
	@go vet ${GO_PACKAGES}

test: unit-test vet

generate:
	#counterfeiter -o test/fake_kubernetes_client.go k8s.io/client-go/kubernetes.Interface
	# ^ requires having k8s.io/client-go checked out, see https://git.io/vFo28
	#sed -i '' 's/FakeInterface/FakeK8sInterface/g' test/fake_kubernetes_client.go
	go generate ./...

all: fmt test build

