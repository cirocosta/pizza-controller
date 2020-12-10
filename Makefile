build: gen-objects
	cd cmd/controller && go build -v -i

run: build
	./cmd/controller/controller


install:
	kapp deploy --yes -c -a pizza-controller -f ./config/bases/crds.yaml

uninstall:
	kapp delete -a pizza-controller --yes


image:
	docker build -t cirocosta/pizza-controller .
push-image:
	docker push cirocosta/pizza-controller


release:
	ytt -f ./config/bases/ | kbld --images-annotation=false -f- > ./config/release.yaml


gen-manifests:
	controller-gen \
		crd \
		paths=./pkg/apis/ops.tips/v1alpha1 \
		output:stdout > ./config/bases/crds.yaml
	controller-gen \
		rbac:roleName=pizza-controller \
		paths=./pkg/reconciler \
		output:stdout > ./config/bases/role.yaml

gen-objects:
	controller-gen \
		object \
		paths=./pkg/apis/ops.tips/v1alpha1
