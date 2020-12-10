# pizza-controller

making kubernetes do what it was _always_ meant to do: order pizza.

- https://gum.co/kubernetes-crds to get up to speed with custom resources and controllers
- [docs](./docs) for CRD reference


## usage

First, create a secret with the credit card information (_yeah, this is fine,
trust me_) to be used during payment:

```yaml
kind: Secret
apiVersion: v1
metadata:
  name: credit-card
spec:
  number: 123343132314232
  expiration: 12-02-2020
  cc: 123
```

then, create a `PizzaCustomer`, the representation of _you_, the customer:


```yaml
kind: PizzaCustomer
apiVersion: ops.tips/v1alpha1
metadata:
  name: you
spec:
  firstName: barack
  lastName: obama
  email: obama@gov.gov
  phone: ""
  streetNumber:
  streetName:
  city:
  state:
  zip:
``` 

With the `PizzaCustomer` object created, we can see what's the closest store available
for it:


```console
$ kubectl get pizzacustomer
NAME              CLOSEST
you               store-123
```

Looking at the PizzaStore object, we can check out its menu:

```console
$ kubectl get pizzastore store-123 -o yaml
kind: PizzaStore
metadata:
  name: store-123
spec:
  address: |
    51 Niagara St
    Toronto, ON M5V1C3
  id: "10391"
  phone: 416-364-3939
  products:
    - description: Unique Lymon (lemon-lime) flavor, clear, clean and crisp with no caffeine.
      id: 2LSPRITE
      name: Sprite
      size: 2 Litre
```

Knowing what's available for us, we can place the order:

```yaml
kind: PizzaOrder
apiVersion: ops.tips/v1
metadata:
  name: ma-pizza
spec:
  yeahSurePlaceThisOrder: true  # otherwise, it'll just calculate the price
  storeRef: {name: store-123}
  customerRef: {name: you}
  payment:
    creditCardSecretRef: {name: cc}
  items:
    - ticker: 10SCREEN
      quantity: 1
```

To keep track of what's going on with your pizza, check out the order's status:

```console
$ kubectl get pizzaorder ma-pizza
NAME    PRICE      ID                     CONDITION     AGE
order   9.030000   Wlz6HcE6BPlfQNlxDAXa   OrderPlaced   68m
```


## what's next?

are you _really_ into ordering pizza using `kubectl`?

like, really? are you sure?

here's what's missing:

- well, any .. tests :horse:
- order tracking (it's `xml`-based - SOAP stuff)
- being more flexible with non-canadian folks (it's currently hardcoded for
  Canada, but could easily be changed)


## Installation

1. apply the manifest

```
git clone https://github.com/cirocosta/pizza-controller
cd pizza-controller

# using `kapp` - recommended :)
#
kapp deploy -a pizza-controller -f ./config/release.yaml


# OR .. plain old kubectl
#
kubectl apply -f ./config/release.yaml
```

## thanks

- https://github.com/ggrammar/pizzapi
- https://github.com/harrybrwn/apizza


## LICENSE

MIT
