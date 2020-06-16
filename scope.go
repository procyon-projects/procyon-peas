package peas

/* If you don't implement PeaScope interface to your struct,
 * it is considered as singleton.
 */
const (
	PeaSingletonScope string = "singleton" /* it is only supported now by default */
)

//type PeaScope interface {
//	GetScopeName() string
//}
