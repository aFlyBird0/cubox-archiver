package cubox

type Source interface {
	List(items chan *Item) // 需要该方法主动关闭 chan
}
