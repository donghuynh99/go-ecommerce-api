package models

type Model struct {
	Model interface{}
}

func RegisteredModel() []Model {
	return []Model{
		{Model: User{}},
		{Model: Product{}},
		{Model: Address{}},
		{Model: Category{}},
		{Model: OrderCustomer{}},
		{Model: OrderItem{}},
		{Model: ProductImage{}},
		{Model: Order{}},
	}
}
