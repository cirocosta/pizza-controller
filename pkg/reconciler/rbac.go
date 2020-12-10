package reconciler

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups=ops.tips,resources=pizzacustomers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ops.tips,resources=pizzacustomers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ops.tips,resources=pizzaorders,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ops.tips,resources=pizzaorders/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ops.tips,resources=pizzastores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ops.tips,resources=pizzastores/status,verbs=get;update;patch
