package betterhandler

// WriteString writes str into ResponseWriter
func (r response) WriteString(str string) error {
    _, err := r.Write([]byte(str))
    return err
}
