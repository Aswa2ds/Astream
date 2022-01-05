package astream

type Interface interface {
	ToInterface() []interface{}
}

func Int8ToInterface(list []int8) []interface{} {
	interfaces := make([]interface{}, 0)
	for _, val := range list {
		interfaces = append(interfaces, val)
	}
	return interfaces
}

func Int16ToInterface(list []int16) []interface{} {
	interfaces := make([]interface{}, 0)
	for _, val := range list {
		interfaces = append(interfaces, val)
	}
	return interfaces
}

func Int32ToInterface(list []int32) []interface{} {
	interfaces := make([]interface{}, 0)
	for _, val := range list {
		interfaces = append(interfaces, val)
	}
	return interfaces
}

func IntToInterface(list []int) []interface{} {
	interfaces := make([]interface{}, 0)
	for _, val := range list {
		interfaces = append(interfaces, val)
	}
	return interfaces
}

func StringToInterface(list []string) []interface{} {
	interfaces := make([]interface{}, 0)
	for _, val := range list {
		interfaces = append(interfaces, val)
	}
	return interfaces
}
