# Coll

Coll is a collection library for Go.

    cc := coll.MustNew(testMapData)
    collectionWithoutEmpty := cc.Filter(func(value string) bool {
        return value != ""
    }).Filter(func(value string, key string) bool {
        return key != ""
    })
    collectionWithoutEmpty.Each(func(value, key string) {
        if value == "" || key == "" {
            t.Errorf("test failed: %s=>%s", key, value)
        }
    })