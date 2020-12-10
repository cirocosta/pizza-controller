package reconciler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1alpha1 "github.com/cirocosta/pizza-controller/pkg/apis/ops.tips/v1alpha1"
	"github.com/cirocosta/pizza-controller/pkg/dominos"
	"github.com/go-logr/logr"
)

type PizzaOrderReconciler struct {
	Log    logr.Logger
	Client client.Client
}

func (r *PizzaOrderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, err error) {
	log := r.Log.WithValues("name", req.NamespacedName)

	log.Info("start")
	defer func() {
		if err != nil {
			log.Error(err, "finished")
		} else {
			log.Info("finished")
		}
	}()

	order, err := r.GetPizzaOrder(ctx, req.Name, req.Namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			return
		}

		err = fmt.Errorf("get pizza order: %w", err)
		return
	}

	err = r.ReconcilePizzaOrder(ctx, order)
	if err != nil {
		err = fmt.Errorf("reconcile pizza order: %w", err)
		return
	}

	return ctrl.Result{
		RequeueAfter: 3 * time.Minute,
	}, nil
}

func (r *PizzaOrderReconciler) ReconcilePizzaOrder(
	ctx context.Context,
	order *v1alpha1.PizzaOrder,
) error {
	if r.IsOrderAlreadyPlaced(order) {
		return nil
	}

	client, err := dominos.NewClient(dominos.CanadaURL, true)
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}

	dominosOrder, err := r.AssembleDominosOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("assemble dominos order: %w", err)
	}

	if !r.IsOrderAlreadyPriced(order) {
		price, err := client.PriceOrder(ctx, *dominosOrder)
		if err != nil {
			return fmt.Errorf("price order: %w", err)
		}

		order.Status.Price = price
		order.Status.Conditions = append(order.Status.Conditions, metav1.Condition{
			Type:               "OrderPriced",
			Status:             metav1.ConditionTrue,
			Reason:             "OrderPriced",
			LastTransitionTime: metav1.Now(),
		})
		if err := r.Client.Status().Update(ctx, order); err != nil {
			return fmt.Errorf("price status update: %w", err)
		}

		return nil
	}

	price, err := strconv.ParseFloat(order.Status.Price, 64)
	if err != nil {
		return fmt.Errorf("parse float '%s': %w",
			order.Status.Price, err,
		)
	}

	dominosOrder.Amount = price

	if order.Spec.YeahSurePlaceTheOrder {
		orderID, err := client.PlaceOrder(ctx, *dominosOrder)
		if err != nil {
			return fmt.Errorf("place order: %w", err)
		}

		order.Status.OrderID = orderID
		order.Status.Conditions = append(order.Status.Conditions, metav1.Condition{
			Type:               "OrderPlaced",
			Status:             metav1.ConditionTrue,
			Reason:             "OrderPlaced",
			LastTransitionTime: metav1.Now(),
		})
		if err := r.Client.Status().Update(ctx, order); err != nil {
			return fmt.Errorf("price status update: %w", err)
		}
	}

	return nil
}

func (r *PizzaOrderReconciler) AssembleDominosOrder(
	ctx context.Context,
	order *v1alpha1.PizzaOrder,
) (*dominos.Order, error) {
	customer, err := r.GetPizzaCustomer(ctx,
		order.Spec.CustomerRef.Name, order.Namespace,
	)
	if err != nil {
		return nil, fmt.Errorf("get pizza customer '%s': %w",
			order.Spec.CustomerRef.Name, err,
		)
	}

	store, err := r.GetPizzaStore(ctx,
		order.Spec.StoreRef.Name, order.Namespace,
	)
	if err != nil {
		return nil, fmt.Errorf("get pizza store '%s': %w",
			order.Spec.StoreRef.Name, err,
		)
	}

	cc := &dominos.CreditCard{}
	if order.Spec.YeahSurePlaceTheOrder {
		cc, err = r.GetCreditCardInfo(ctx,
			customer.Spec.CreditCardSecretRef.Name,
			order.Namespace,
		)
		if err != nil {
			return nil, fmt.Errorf("get credit card info: %w", err)
		}
	}

	products := []dominos.Product{}
	for _, product := range order.Spec.Products {
		products = append(products, dominos.Product{
			ID: product.ID,
		})
	}

	return &dominos.Order{
		StoreID: store.Spec.ID,
		PersonalInformation: dominos.PersonalInformation{
			FirstName: customer.Spec.FirstName,
			LastName:  customer.Spec.LastName,
			Email:     customer.Spec.Email,
			Phone:     customer.Spec.Phone,
		},
		CreditCard: *cc,
		Address: dominos.Address{
			StreetName:   customer.Spec.StreetName,
			StreetNumber: customer.Spec.StreetNumber,
			City:         customer.Spec.City,
			State:        customer.Spec.State,
			Zip:          customer.Spec.Zip,
		},
		Products: products,
		Service:  dominos.ServiceCarryout,
	}, nil
}

func (r *PizzaOrderReconciler) GetCreditCardInfo(
	ctx context.Context,
	name, namespace string,
) (*dominos.CreditCard, error) {
	obj := &corev1.Secret{}
	if err := r.Client.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, obj); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	number, found := obj.Data["number"]
	if !found {
		return nil, fmt.Errorf("'number' not found in cc info")
	}

	expiration, found := obj.Data["expiration"]
	if !found {
		return nil, fmt.Errorf("'expiration' not found in cc info")
	}

	securityCode, found := obj.Data["securityCode"]
	if !found {
		return nil, fmt.Errorf("'securityCode' not found in cc info")
	}

	var cardType dominos.CreditCardType
	cardTypeStr, found := obj.Data["cardType"]
	if !found {
		return nil, fmt.Errorf("'cardType' not found in cc info")
	}
	switch strings.ToLower(string(cardTypeStr)) {
	case "mastercard":
		cardType = dominos.CreditCardTypeMastercard
	case "visa":
		cardType = dominos.CreditCardTypeVisa
	default:
		return nil, fmt.Errorf("unknown card type '%s'", cardType)
	}

	zip, found := obj.Data["zip"]
	if !found {
		return nil, fmt.Errorf("'zip' not found in cc info")
	}

	return &dominos.CreditCard{
		Type:         cardType,
		Expiration:   string(expiration),
		Number:       string(number),
		PostalCode:   string(zip),
		SecurityCode: string(securityCode),
	}, nil
}

func (r *PizzaOrderReconciler) GetPizzaStore(
	ctx context.Context,
	name, namespace string,
) (*v1alpha1.PizzaStore, error) {
	obj := &v1alpha1.PizzaStore{}
	if err := r.Client.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, obj); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return obj, nil
}

func (r *PizzaOrderReconciler) GetPizzaCustomer(
	ctx context.Context,
	name, namespace string,
) (*v1alpha1.PizzaCustomer, error) {
	obj := &v1alpha1.PizzaCustomer{}
	if err := r.Client.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, obj); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return obj, nil
}

func (r *PizzaOrderReconciler) GetPizzaOrder(
	ctx context.Context,
	name, namespace string,
) (*v1alpha1.PizzaOrder, error) {
	obj := &v1alpha1.PizzaOrder{}
	if err := r.Client.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, obj); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return obj, nil
}

func (r *PizzaOrderReconciler) IsOrderAlreadyPriced(order *v1alpha1.PizzaOrder) bool {
	for _, cond := range order.Status.Conditions {
		if cond.Type == "OrderPriced" {
			return true
		}
	}

	return false
}

func (r *PizzaOrderReconciler) IsOrderAlreadyPlaced(order *v1alpha1.PizzaOrder) bool {
	for _, cond := range order.Status.Conditions {
		if cond.Type == "OrderPlaced" {
			return true
		}
	}

	return false
}
