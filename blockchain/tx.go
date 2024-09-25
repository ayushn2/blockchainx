package blockchain

// Since blockchains are public we don't want to save any sensitive information on the blocks

type TxOutput struct{
	Value int
	PubKey string
}


// Input just references output
type TxInput struct{
	ID []byte
	Out int
	Sig string
}

func (in *TxInput) CanUnlock(data string) bool{
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool{
	return out.PubKey == data
}