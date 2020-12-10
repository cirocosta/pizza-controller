package reconciler

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	v1alpha1 "github.com/cirocosta/pizza-controller/pkg/apis/ops.tips/v1alpha1"
	"github.com/cirocosta/pizza-controller/pkg/dominos"
	"github.com/go-logr/logr"
)

type PizzaCustomerReconciler struct {
	Log    logr.Logger
	Client client.Client
}

func (r *PizzaCustomerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, err error) {
	log := r.Log.WithValues("name", req.NamespacedName)

	log.Info("start")
	defer func() {
		if err != nil {
			log.Error(err, "finished")
		} else {
			log.Info("finished")
		}
	}()

	customer, err := r.GetPizzaCustomer(ctx, req.Name, req.Namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			return
		}

		err = fmt.Errorf("get pizza customer: %w", err)
		return
	}

	err = r.ReconcilePizzaCustomer(ctx, customer)
	if err != nil {
		err = fmt.Errorf("reconcile pizza customer: %w", err)
		return
	}

	return ctrl.Result{
		RequeueAfter: 3 * time.Minute,
	}, nil
}

func (r *PizzaCustomerReconciler) ReconcilePizzaCustomer(
	ctx context.Context,
	customer *v1alpha1.PizzaCustomer,
) error {
	client, err := dominos.NewClient(dominos.CanadaURL, false)
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}

	stores, err := client.StoresNearby(ctx, dominos.Address{
		StreetNumber: customer.Spec.StreetNumber,
		StreetName:   customer.Spec.StreetName,
		City:         customer.Spec.City,
		State:        customer.Spec.State,
		Zip:          customer.Spec.Zip,
	}, dominos.ServiceDelivery)
	if err != nil {
		return fmt.Errorf("stores nearby: %w", err)
	}

	if len(stores) >= 3 {
		stores = stores[:3]
	}

	refs := []*corev1.LocalObjectReference{}
	for _, store := range stores {
		products, err := client.StoreMenu(ctx, store.ID)
		if err != nil {
			return fmt.Errorf("store menu '%s': %w", store.ID, err)
		}

		pizzaStore := r.AssemblePizzaStore(customer, store, products)
		pizzaStoreRef, err := r.FindOrCreate(ctx, pizzaStore)
		if err != nil {
			return fmt.Errorf("find or create: %w", err)
		}

		refs = append(refs, pizzaStoreRef)
	}

	customer.Status.ClosestStoreRef = *refs[0]
	customer.Status.Conditions = []metav1.Condition{
		{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			Reason:             "StoresFound",
			LastTransitionTime: metav1.Now(),
		},
	}

	if err := r.Client.Status().Update(ctx, customer); err != nil {
		return fmt.Errorf("status update: %w", err)
	}

	return nil
}

func (r *PizzaCustomerReconciler) FindOrCreate(
	ctx context.Context,
	obj controllerutil.Object,
) (*corev1.LocalObjectReference, error) {
	if err := r.Client.Create(ctx, obj); err != nil {
		if !errors.IsAlreadyExists(err) {
			return nil, fmt.Errorf("create: %w", err)
		}

		if err := r.Client.Get(ctx, client.ObjectKey{
			Name:      obj.GetName(),
			Namespace: obj.GetNamespace(),
		}, obj); err != nil {
			return nil, fmt.Errorf("get: %w", err)
		}
	}

	return &corev1.LocalObjectReference{
		Name: obj.GetName(),
	}, nil
}

func (r *PizzaCustomerReconciler) AssemblePizzaStore(
	customer *v1alpha1.PizzaCustomer,
	store *dominos.Store,
	products []*dominos.Product,
) *v1alpha1.PizzaStore {

	specProducts := []v1alpha1.PizzaStoreProduct{}
	for _, product := range products {
		specProducts = append(specProducts, v1alpha1.PizzaStoreProduct{
			Name:        product.Name,
			ID:          product.ID,
			Description: product.Description,
			Size:        product.Size,
		})
	}

	return &v1alpha1.PizzaStore{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "store-" + strings.ToLower(store.ID),
			Namespace: customer.Namespace,
		},
		Spec: v1alpha1.PizzaStoreSpec{
			Address:  store.Address,
			ID:       store.ID,
			Phone:    store.Phone,
			Products: specProducts,
		},
	}
}

func (r *PizzaCustomerReconciler) GetPizzaCustomer(
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
