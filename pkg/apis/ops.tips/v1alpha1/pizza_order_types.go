package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Price",type=string,JSONPath=`.status.price`
// +kubebuilder:printcolumn:name="ID",type=string,JSONPath=`.status.orderID`
// +kubebuilder:printcolumn:name="Condition",type=string,JSONPath=`.status.conditions[-1].type`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

type PizzaOrder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PizzaOrderSpec   `json:"spec,omitempty"`
	Status PizzaOrderStatus `json:"status,omitempty"`
}

type PizzaOrderSpec struct {
	YeahSurePlaceTheOrder bool                        `json:"yeahSurePlaceTheOrder,omitempty"`
	PaymentType           string                      `json:"paymentType"`
	StoreRef              corev1.LocalObjectReference `json:"storeRef"`
	CustomerRef           corev1.LocalObjectReference `json:"customerRef"`
	Products              []PizzaOrderProduct         `json:"products"`
}

type PizzaOrderProduct struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type PizzaOrderStatus struct {
	OrderID    string             `json:"orderID,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	Price      string             `json:"price,omitempty"`
}

// +kubebuilder:object:root=true

type PizzaOrderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PizzaOrder `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PizzaOrder{}, &PizzaOrderList{})
}
