package coll

func Map(origin interface{}, dest interface{}, mapper interface{}) error {
	return MustNew(origin).Map(mapper).All(dest)
}

func Filter(origin interface{}, dest interface{}, filter interface{}) error {
	return MustNew(origin).Filter(filter).All(dest)
}
