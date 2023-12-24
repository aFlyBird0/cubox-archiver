package core

// Source 能够提供 Cubox Item 的数据源
type Source interface {
	List(items chan *Item) // 需要该方法主动关闭 chan
}
