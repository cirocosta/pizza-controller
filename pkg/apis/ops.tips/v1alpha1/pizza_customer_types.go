package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Closest",type=string,JSONPath=`.status.closestStoreRef.name`
// +kubebuilder:printcolumn:name="Condition",type=string,JSONPath=`.status.conditions[-1].type`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

type PizzaCustomer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PizzaCustomerSpec   `json:"spec,omitempty"`
	Status PizzaCustomerStatus `json:"status,omitempty"`
}

type PizzaCustomerSpec struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`

	StreetNumber string `json:"streetNumber"`
	StreetName   string `json:"streetName"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`

	CreditCardSecretRef corev1.LocalObjectReference `json:"creditCardSecretRef"`
}

type PizzaCustomerStatus struct {
	ClosestStoreRef corev1.LocalObjectReference `json:"closestStoreRef,omitempty"`
	Conditions      []metav1.Condition          `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true

type PizzaCustomerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PizzaCustomer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PizzaCustomer{}, &PizzaCustomerList{})
}
