package peer

import "errors"

type Bit_field []byte

func (b Bit_field) Has_piece(index int) bool{
	if index / 8 >= len(b){
		return false
	}
	
	target_byte := uint8(b[index / 8])
	target_byte >>= 7 - (index % 8)
	target_byte &= 1
	
	if target_byte == 1{
		return true
	}

	return false
}

func (b Bit_field) Have_update(index int) (error){
	if index / 8 >= len(b){
		return errors.New("Invalid index!")
	}

	target_byte := uint8(b[index/8])
	new_value := uint8(1) << (7- (index%8))
	target_byte |= new_value

	b[index/8] = target_byte

	return nil
} 