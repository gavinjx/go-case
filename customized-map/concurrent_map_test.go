package customized_map

import (
    "testing"
    "reflect"
    "fmt"
    "runtime/debug"
    "math/rand"
    "bytes"
    "time"
)

/*************** 功能测试 *************/
//测试Int64型的cmap
func TestInt64Cmap(t *testing.T) {
    newFunc := func() ConcurrentMapIntfs {
        keyType := reflect.TypeOf(int64(2))
        valType := keyType
        return NewConcurrentMap(keyType, valType)
    }

    genFunc := func() interface{} { return rand.Int63n(1000) }

    test(t, newFunc, genFunc, genFunc, reflect.Int64,reflect.Int64)
}
//测试Float64型的cmap
func TestFloat64Cmap(t *testing.T) {
    newFunc := func() ConcurrentMapIntfs {
        keyType := reflect.TypeOf(float64(2))
        valType := keyType
        return NewConcurrentMap(keyType, valType)
    }
    genFunc := func() interface{} { return rand.Float64() }

    test(t, newFunc, genFunc, genFunc, reflect.Float64, reflect.Float64)
}

//测试string型cmap
func TestStringCmap(t *testing.T) {
    newCmap := func() ConcurrentMapIntfs {
        keyType := reflect.TypeOf(string(2))
        valType := keyType
        return NewConcurrentMap(keyType, valType)
    }
    genFunc := func() interface{} { return genRandString() }
    test(t, newCmap, genFunc, genFunc, reflect.String, reflect.String)
}

func genRandString() string {
    var buff bytes.Buffer
    var prev string
    var curr string
    for i := 0; buff.Len() < 3; i++ {
        curr = string(genRandAZAscii())
        if curr == prev {
            continue
        }
        prev = curr
        buff.WriteString(curr)
    }
    return buff.String()
}

func genRandAZAscii() int {
    min := 65 // A
    max := 90 // Z
    rand.Seed(time.Now().UnixNano())
    return min + rand.Intn(max-min)
}

func test(t *testing.T, newConcurrentMap func() ConcurrentMapIntfs, genKey func() interface{}, genVal func() interface{}, keyKind reflect.Kind, valKind reflect.Kind) {
    mapType := fmt.Sprintf("ConcurrentMap<keyType=%s, elemType=%s>", keyKind, valKind)
    defer func() {
        if err := recover(); err != nil {
            debug.PrintStack()
            t.Errorf("Fatal Error: %s: %s\n", mapType, err)
        }
    }()
    t.Logf("Starting Test%s...", mapType)

    // Basic
    cmap := newConcurrentMap()
    expectedLen := 0
    if cmap.Len() != expectedLen {
        t.Errorf("ERROR: The length of %s value %d is not %d!\n", mapType, cmap.Len(), expectedLen)
        t.FailNow()
    }
    expectedLen = 5
    testMap := make(map[interface{}]interface{}, expectedLen)
    var invalidKey interface{}
    for i := 0; i < expectedLen; i++ {
        key := genKey()
        testMap[key] = genVal()
        if invalidKey == nil {
            invalidKey = key
        }
    }
    for key, val := range testMap {
        oldVal, ok := cmap.Put(key, val)
        if !ok {
            t.Errorf("ERROR: Put (%v, %v) to %s value %d is failing!\n", key, val, mapType, cmap)
            t.FailNow()
        }
        if oldVal != nil {
            t.Errorf("ERROR: Already had a (%v, %v) in %s value %d!\n", key, val, mapType, cmap)
            t.FailNow()
        }
        t.Logf("Put (%v, %v) to the %s value %v.", key, val, mapType, cmap)
    }
    if cmap.Len() != expectedLen {
        t.Errorf("ERROR: The length of %s value %d is not %d!\n", mapType, cmap.Len(), expectedLen)
        t.FailNow()
    }
    for key, val := range testMap {
        contains := cmap.Contains(key)
        if !contains {
            t.Errorf("ERROR: The %s value %v do not contains %v!", mapType, cmap, key)
            t.FailNow()
        }
        actualVal := cmap.Get(key)
        if actualVal == nil {
            t.Errorf("ERROR: The %s value %v do not contains %v!", mapType, cmap, key)
            t.FailNow()
        }
        t.Logf("The %s value %v contains key %v.", mapType, cmap, key)
        if actualVal != val {
            t.Errorf("ERROR: The element of %s value %v with key %v do not equals %v!\n", mapType, actualVal, key, val)
            t.FailNow()
        }
        t.Logf("The element of %s value %v to key %v is %v.", mapType, cmap, key, actualVal)
    }
    oldVal := cmap.Remove(invalidKey)
    if oldVal == nil {
        t.Errorf("ERROR: Remove %v from %s value %d is failing!\n", invalidKey, mapType, cmap)
        t.FailNow()
    }
    t.Logf("Removed (%v, %v) from the %s value %v.", invalidKey, oldVal, mapType, cmap)
    delete(testMap, invalidKey)

    // Type
    actualValType := cmap.ValType()
    if actualValType == nil {
        t.Errorf("ERROR: The element type of %s value is nil!\n", mapType)
        t.FailNow()
    }
    actualValKind := actualValType.Kind()
    if actualValKind != valKind {
        t.Errorf("ERROR: The element type of %s value %s is not %s!\n", mapType, actualValKind, valKind)
        t.FailNow()
    }
    t.Logf("The element type of %s value %v is %s.", mapType, cmap, actualValKind)
    actualKeyKind := cmap.KeyType().Kind()
    if actualKeyKind != keyKind {
        t.Errorf("ERROR: The key type of %s value %s is not %s!\n", mapType, actualKeyKind, keyKind)
        t.FailNow()
    }
    t.Logf("The key type of %s value %v is %s.", mapType, cmap, actualKeyKind)

    // Export
    keys := cmap.Keys()
    vals := cmap.Vals()
    pairs := cmap.ToMap()
    for key, elem := range testMap {
        var hasKey bool
        for _, k := range keys {
            if k == key {
                hasKey = true
            }
        }
        if !hasKey {
            t.Errorf("ERROR: The keys of %s value %v do not contains %v!\n", mapType, cmap, key)
            t.FailNow()
        }
        var hasVal bool
        for _, e := range vals {
            if e == elem {
                hasVal = true
            }
        }
        if !hasVal {
            t.Errorf("ERROR: The elems of %s value %v do not contains %v!\n", mapType, cmap, elem)
            t.FailNow()
        }
        var hasPair bool
        for k, e := range pairs {
            if k == key && e == elem {
                hasPair = true
            }
        }
        if !hasPair {
            t.Errorf("ERROR: The elems of %s value %v do not contains (%v, %v)!\n",
                mapType, cmap, key, elem)
            t.FailNow()
        }
    }

    // Clear
    cmap.Clear()
    if cmap.Len() != 0 {
        t.Errorf("ERROR: Clear %s value %d is failing!\n", mapType, cmap)
        t.FailNow()
    }
    t.Logf("The %s value %v has been cleared.", mapType, cmap)
}