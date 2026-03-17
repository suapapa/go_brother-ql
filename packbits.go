package brother_ql

func Packbits(input []byte) []byte {
	var output []byte
	i := 0
	for i < len(input) {
		// Look for a run of identical bytes
		runLen := 1
		for i+runLen < len(input) && runLen < 128 && input[i+runLen] == input[i] {
			runLen++
		}
		if runLen > 1 {
			// Output run header and the byte
			// Header is - (runLen - 1)
			h := byte(257 - runLen) // Equivalent to byte(-(runLen - 1)) in 2's complement
			output = append(output, h)
			output = append(output, input[i])
			i += runLen
		} else {
			// Look for non-run
			nonRunLen := 1
			for i+nonRunLen < len(input) && nonRunLen < 128 {
				// If we encounter a run of 2 or more identical bytes, stop before it
				if i+nonRunLen+1 < len(input) && input[i+nonRunLen] == input[i+nonRunLen+1] {
					break
				}
				nonRunLen++
			}
			// Output non-run header and the bytes
			h := byte(nonRunLen - 1)
			output = append(output, h)
			output = append(output, input[i:i+nonRunLen]...)
			i += nonRunLen
		}
	}
	return output
}
