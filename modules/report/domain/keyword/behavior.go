package keyword

func (k *Keyword) Word() string  { return k.state.Word }
func (k *Keyword) Kind() string  { return k.state.Kind }
func (k *Keyword) Group() string { return k.state.GroupName }
