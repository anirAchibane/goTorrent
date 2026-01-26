package peer

type Bit_field []byte

func (b Bit_field) Has_piece(index int) bool{
	if index / 8 >= len(b){
		return false
	}
	
	target_byte := uint8(b[index / 8])
	target_byte >>= 8 - (index % 8) - 1
	target_byte &= 1
	
	if target_byte == 1{
		return true
	}
	return false
}