package capacity

type ProvisionTRUInput struct {
	Namespace          string
	APSLimit           int32
	MinutesToProvision int32
}

type DeprovisionTRUInput struct {
	Namespace          string
	MinutesToProvision int32
}

type AddTRUInput struct {
	Namespace string
	APSLimit  int32
}

type RemoveTRUInput struct {
	Namespace string
}
