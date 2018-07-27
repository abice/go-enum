package tuple

func FirstOfTwo(v1, v2 interface{}) interface{}  { return v1 }
func SecondOfTwo(v1, v2 interface{}) interface{} { return v2 }

func FirstOfThree(v1, v2, v3 interface{}) interface{}  { return v1 }
func SecondOfThree(v1, v2, v3 interface{}) interface{} { return v2 }
func ThirdOfThree(v1, v2, v3 interface{}) interface{}  { return v3 }

func FirstError(v1 error, v2 interface{}) error  { return v1 }
func SecondError(v1 interface{}, v2 error) error { return v2 }

func FirstBool(v1 bool, v2 interface{}) bool  { return v1 }
func SecondBool(v1 interface{}, v2 bool) bool { return v2 }
