# CRDs

## PizzaCustomer

The `PizzaCustomer` CRD represents a customer who might want to place an
order at Dominos.

In its spec, one fills the fields that will let Dominos know of your personal
information, address, and a reference to a secret where credit card details
can be found.

```yaml
kind: PizzaCustomer
apiVersion: ops.tips/v1alpha1
spec:
  address: etcetc
  name: bla
```

The reconciler has the responsability of finding stores nearby the customer
so that orders can be placed for it later on.

```yaml
status:
  conditions:
    - type: Ready
      status: "True"
      reason: StoresFound
  closestStoreRef: { name: store-123 }
```

So ultimately, it's a state machine like so:

```mermaid
stateDiagram-v2
  [*] --> Ready
  [*] --> Errored
  Ready --> Errored
  Ready --> Ready
  Errored --> Ready
  Errored --> Errored
```

Where:

- `Ready` implies that stores have been found and orders for that customer can be placed
- `Errored` implies something went wrong

## PizzaStore

A `PizzaStore` object represents a physical Dominos location where one can
order food from.

```yaml
kind: PizzaStore
apiVersion: ops.tips/v1alpha1
spec:
  name: bla
  products:
    - code: 10SCREEN
      name: ""
      description: ""
```

It's _not_ supposed to be created by humans - `PizzaStore` objects are created by the controller.

## PizzaOrder

With a `PizzaOrder` object, you declare the intention to have food from a
Dominos store by filling three fields:

- `spec.storeRef`: reference to a `PizzaStore` object
- `spec.customerRef`: reference to a `PizzaCustomer` object
- `spec.products`: set of products to order from that store

```yaml
kind: PizzaOrder
apiVersion: ops.tips/v1alpha1
spec:
  storeRef: { name: store-123 }
  customerRef: { name: customer-1 }
  serviceType: carryout
  products:
    - code: 10SCREEN
      quantity: 1
  payment:
    creditCardSecretRef:
      name: credit-card
```
