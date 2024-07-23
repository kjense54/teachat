package main

// fix incorrect wrapping in viewport by manually resizing strings
func (m model) ChopText(text string, size int) []string {
	if len(text) == 0 {
		return nil
	}
	if len(text) < size {
		return []string{text}
	}
	var chopped []string = make([]string, 0, (len(text)-1)/size+1)
	currentLen := 0
	currentStart := 0
	for i := range text {
		if currentLen == size {
			chopped = append(chopped, text[currentStart:i])
			currentLen = 0
			currentStart = i
		} 
		currentLen++
	}
	// add extra bits at end
	chopped = append(chopped, text[currentStart:])
	return chopped
}
