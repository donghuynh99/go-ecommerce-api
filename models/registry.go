package models

type Model struct {
	Model interface{}
}

func RegisteredModel() []Model {
	return []Model{
		{Model: Image{}},
		{Model: User{}},
		{Model: Product{}},
		{Model: Address{}},
		{Model: Category{}},
		{Model: OrderItem{}},
		{Model: Order{}},
		{Model: Cart{}},
		{Model: CartItem{}},
	}
}
