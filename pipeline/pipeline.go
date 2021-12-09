package pipeline

type Handler func(value interface{}, next Handler)

func NewPipeline(v interface{}) *Pipline {
	return &Pipline{value: v}
}

type Pipline struct {
	value interface{}
}

func (p *Pipline) Send(v interface{}) *Pipline {
	p.value = v
	return p
}

func (p *Pipline) Through(carry ...Handler) *Pipline {
	p.iter(&p.value, carry...)(p.value, nil)
	return p
}

func (p Pipline) Return() interface{} {
	return p.value
}

func (p *Pipline) iter(value *interface{}, carry ...Handler) Handler {
	l := len(carry)
	if l == 0 {
		return nil
	}

	if l == 1 {
		return func(v interface{}, next Handler) {
			carry[0](v, func(v interface{}, next Handler) {
				p.value = v
			})
		}
	}

	return func(v interface{}, next Handler) {
		p.value = v
		carry[0](v, p.iter(&v, carry[1:]...))
	}
}
