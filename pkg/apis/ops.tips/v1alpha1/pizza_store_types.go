package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

type PizzaStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PizzaStoreSpec   `json:"spec,omitempty"`
	Status PizzaStoreStatus `json:"status,omitempty"`
}

type PizzaStoreSpec struct {
	ID       string              `json:"id"`
	Phone    string              `json:"phone"`
	Address  string              `json:"address"`
	Products []PizzaStoreProduct `json:"products"`
}

type PizzaStoreProduct struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Size        string `json:"size"`
}

type PizzaStoreStatus struct {
}

// +kubebuilder:object:root=true

type PizzaStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PizzaStore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PizzaStore{}, &PizzaStoreList{})
}
