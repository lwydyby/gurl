package body

func Unmarshal(data []byte) (string, any, error) {

	o, a, err := jsonUnmarshal([]byte(data))
	if err != nil {
		// 出错，再尝试用yaml解析
		o, a, err = yamlUnmarshal([]byte(data))
		if err != nil {
			return "", nil, err
		}
		if o != nil {
			return EncodeYAML, IfElseAny(o != nil, o, a), nil
		}
	}

	return EncodeJSON, IfElseAny(o != nil, o, a), nil
}

func IfElseAny(cond bool, ifVal any, elseVal any) any {
	if cond {
		return ifVal
	}
	return elseVal
}
