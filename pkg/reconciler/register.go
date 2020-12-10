package reconciler

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"

	v1alpha1 "github.com/cirocosta/pizza-controller/pkg/apis/ops.tips/v1alpha1"
)

func AddToScheme(scheme *runtime.Scheme) error {
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return fmt.Errorf("clientgoscheme addtoscheme: %w", err)
	}

	if err := v1alpha1.AddToScheme(scheme); err != nil {
		return fmt.Errorf("v1alpha1 addtoscheme: %w", err)
	}

	return nil
}

func RegisterReconcilers(mgr manager.Manager) error {
	// if err := RegisterPizzaCustomerReconciler(mgr); err != nil {
	// 	return fmt.Errorf("register pizza customer reconciler: %w")
	// }

	if err := RegisterPizzaOrderReconciler(mgr); err != nil {
		return fmt.Errorf("register pizza order reconciler: %w")
	}

	return nil
}

func RegisterPizzaOrderReconciler(mgr manager.Manager) error {
	c, err := controller.New("pizza-order-reconciler", mgr, controller.Options{
		Reconciler: &PizzaOrderReconciler{
			Log:    mgr.GetLogger().WithName("pizza-order-reconciler"),
			Client: mgr.GetClient(),
		},
	})
	if err != nil {
		return fmt.Errorf("new controller: %w", err)
	}

	if err := c.Watch(
		&source.Kind{Type: &v1alpha1.PizzaOrder{}},
		&handler.EnqueueRequestForObject{},
	); err != nil {
		return fmt.Errorf("watch: %w", err)
	}

	return nil
}

func RegisterPizzaCustomerReconciler(mgr manager.Manager) error {
	c, err := controller.New("pizza-customer-reconciler", mgr, controller.Options{
		Reconciler: &PizzaCustomerReconciler{
			Log:    mgr.GetLogger().WithName("pizza-customer-reconciler"),
			Client: mgr.GetClient(),
		},
	})
	if err != nil {
		return fmt.Errorf("new controller: %w", err)
	}

	if err := c.Watch(
		&source.Kind{Type: &v1alpha1.PizzaCustomer{}},
		&handler.EnqueueRequestForObject{},
	); err != nil {
		return fmt.Errorf("watch: %w", err)
	}

	return nil
}
