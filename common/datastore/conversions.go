package datastore

func toInt(raw_data interface{}, definition DataStoreDefinition) (interface{}, error) {
	data, ok := raw_data.(float64)

	if !ok {
		return nil, ErrTypeMismatch
	}

	return int64(data), nil
}

func toFloat(raw_data interface{}, definition DataStoreDefinition) (interface{}, error) {
	data, ok := raw_data.(float64)

	if !ok {
		return nil, ErrTypeMismatch
	}

	return data, nil
}

func toString(raw_data interface{}, definition DataStoreDefinition) (interface{}, error) {
	data, ok := raw_data.(string)

	if !ok {
		return nil, ErrTypeMismatch
	}

	return data, nil
}

func toBoolean(raw_data interface{}, definition DataStoreDefinition) (interface{}, error) {
	data, ok := raw_data.(bool)

	if !ok {
		return nil, ErrTypeMismatch
	}

	return data, nil
}

func toObject(raw_data interface{}, definition DataStoreDefinition) (interface{}, error) {
	unfiltered_data, ok := raw_data.(map[string]interface{})

	if !ok {
		return nil, ErrTypeMismatch
	}

	data := make(map[string]interface{})

	for _, field := range definition.Fields {
		unfiltered_fields, exists := unfiltered_data[field.Name]

		if exists {
			var err error
			data[field.Name], err = buildDataFromDefinition(unfiltered_fields, field)

			if err != nil {
				return nil, err
			}
		} else {
			data[field.Name] = getDefaultValue(field.Type)
		}
	}

	return data, nil
}

func toIntArray(raw_data interface{}, definition DataStoreDefinition) (interface{}, error) {
	data, ok := raw_data.([]interface{})

	if !ok {
		return nil, ErrTypeMismatch
	}

	int_data := make([]int64, len(data))

	for i := 0; i < len(data); i++ {
		element, ok := data[i].(float64)

		if !ok {
			return nil, ErrTypeMismatch
		}

		int_data[i] = int64(element)
	}

	return int_data, nil
}

func toFloatArray(raw_data interface{}, definition DataStoreDefinition) (interface{}, error) {
	data, ok := raw_data.([]interface{})

	if !ok {
		return nil, ErrTypeMismatch
	}

	float_data := make([]float64, len(data))

	for i := 0; i < len(data); i++ {
		element, ok := data[i].(float64)

		if !ok {
			return nil, ErrTypeMismatch
		}

		float_data[i] = element
	}

	return float_data, nil
}

func toStringArray(raw_data interface{}, definition DataStoreDefinition) (interface{}, error) {
	data, ok := raw_data.([]interface{})

	if !ok {
		return nil, ErrTypeMismatch
	}

	string_data := make([]string, len(data))

	for i := 0; i < len(data); i++ {
		element, ok := data[i].(string)

		if !ok {
			return nil, ErrTypeMismatch
		}

		string_data[i] = element
	}

	return string_data, nil
}

func toBooleanArray(raw_data interface{}, definition DataStoreDefinition) (interface{}, error) {
	data, ok := raw_data.([]interface{})

	if !ok {
		return nil, ErrTypeMismatch
	}

	bool_data := make([]bool, len(data))

	for i := 0; i < len(data); i++ {
		element, ok := data[i].(bool)

		if !ok {
			return nil, ErrTypeMismatch
		}

		bool_data[i] = element
	}

	return bool_data, nil
}

func toObjectArray(raw_data interface{}, definition DataStoreDefinition) (interface{}, error) {
	unfiltered_data, ok := raw_data.([]interface{})

	if !ok {
		return nil, ErrTypeMismatch
	}

	data := make([]map[string]interface{}, len(unfiltered_data))

	for i := 0; i < len(unfiltered_data); i++ {
		element := make(map[string]interface{})

		for _, field := range definition.Fields {
			unfiltered_data_element, ok := unfiltered_data[i].(map[string]interface{})

			if !ok {
				return nil, ErrTypeMismatch
			}

			unfiltered_fields, exists := unfiltered_data_element[field.Name]

			if exists {
				var err error
				element[field.Name], err = buildDataFromDefinition(unfiltered_fields, field)

				if err != nil {
					return nil, err
				}
			} else {
				element[field.Name] = getDefaultValue(field.Type)
			}
		}

		data[i] = element
	}

	return data, nil
}

func getDefaultValue(t string) interface{} {
	value, exists := defaultValues[t]

	if !exists {
		return nil
	}

	return value
}
