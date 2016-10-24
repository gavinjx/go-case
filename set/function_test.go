package set
/*
 * 功能测试
 */
import (
    "testing"
    "runtime/debug"
    "fmt"
    "strings"
    "bytes"
    "time"
    "math/rand"
)

/**************************** HashSet相关功能 **************************/
//必须大写Test开头
func TestHashSetCreation(t *testing.T) {
    defer func() {
        if err := recover(); err != nil {
            debug.PrintStack()
            t.Errorf("Fatal Error: %s\n", err)
        }
    }()
    t.Log("Starting TestHashSetCreation...")
    hs := NewHashSet()
    t.Logf("Create a HashSet value: %v\n", hs)
    if hs == nil {
        t.Errorf("The result of func NewHashSet is nil!\n")
    }
    isSet := IsSet(hs)
    if !isSet {
        t.Errorf("The value of HashSet is not Set!\n")
    } else {
        t.Logf("The HashSet value is aSet.\n")
    }
}

func TestSetOperation(t *testing.T) {
    defer func() {
        if err := recover(); err != nil {
            debug.PrintStack()
            t.Errorf("Fatal Error: %s\n", err)
        }
    }()
    fmt.Println(222)
    t.Logf("Starting TestHashSetOp...")
    hs := NewHashSet()
    if hs.Len() != 0 {
        t.Errorf("ERROR: The length of original HashSet value is not 0!\n")
        t.FailNow()
    }
    randElem := genRandElement()
    expectedElemMap := make(map[interface{}]bool)
    t.Logf("Add %v to the HashSet value %v.\n", randElem, hs)
    hs.Add(randElem)
    expectedElemMap[randElem] = true
    expectedLen := len(expectedElemMap)
    if hs.Len() != expectedLen {
        t.Errorf("ERROR: The length of HashSet value %d is not %d!\n", hs.Len(), expectedLen)
        t.FailNow()
    }
    var result bool
    for i := 0; i < 8; i++ {
        randElem = genRandElement()
        t.Logf("Add %v to the HashSet value %v.\n", randElem, hs)
        result = hs.Add(randElem)
        if expectedElemMap[randElem] && result {
            t.Errorf("ERROR: The element adding (%v => %v) is successful but should be failing!\n",
                randElem, hs)
            t.FailNow()
        }
        if !expectedElemMap[randElem] && !result {
            t.Errorf("ERROR: The element adding (%v => %v) is failing!\n",
                randElem, hs)
            t.FailNow()
        }
        expectedElemMap[randElem] = true
    }
    expectedLen = len(expectedElemMap)
    if hs.Len() != expectedLen {
        t.Errorf("ERROR: The length of HashSet value %d is not %d!\n", hs.Len(), expectedLen)
        t.FailNow()
    }
    for k, _ := range expectedElemMap {
        if !hs.Contains(k) {
            t.Errorf("ERROR: The HashSet value %v do not contains %v!", hs, k)
            t.FailNow()
        }
    }
    number := 2
    for k, _ := range expectedElemMap {
        if number%2 == 0 {
            t.Logf("Remove %v from the HashSet value %v.\n", k, hs)
            hs.Remove(k)
            if hs.Contains(k) {
                t.Errorf("ERROR: The element adding (%v => %v) is failing!\n",
                    randElem, hs)
                t.FailNow()
            }
            delete(expectedElemMap, k)
        }
        number++
    }
    expectedLen = len(expectedElemMap)
    if hs.Len() != expectedLen {
        t.Errorf("ERROR: The length of HashSet value %d is not %d!\n", hs.Len(), expectedLen)
        t.FailNow()
    }
    for _, v := range hs.Elements() {
        if !expectedElemMap[v] {
            t.Errorf("ERROR: The HashSet value %v contains %v!", hs, v)
            t.FailNow()
        }
    }
    hs2 := NewHashSet()
    for k, _ := range expectedElemMap {
        hs2.Add(k)
    }
    if !hs.Same(hs2) {
        t.Errorf("ERROR: HashSet value %v do not same %v!\n", hs, hs2)
        t.FailNow()
    }
    str := hs.String()
    t.Logf("The string of HashSet value %v is '%s'.\n", hs, str)
    for _, v := range hs.Elements() {
        if !strings.Contains(str, fmt.Sprintf("%v", v)) {
            t.Errorf("ERROR: '%s' do not contains '%v'!", str, v)
            t.FailNow()
        }
    }
}

// ----- 随机测试对象生成函数 -----
func genRandSet(newSet func() SetIntfs) (set SetIntfs, elemMap map[interface{}]bool) {
    set = newSet()
    elemMap = make(map[interface{}]bool)
    var enough bool
    for !enough {
        e := genRandElement()
        set.Add(e)
        elemMap[e] = true
        if len(elemMap) >= 3 {
            enough = true
        }
    }
    return
}

func genRandElement() interface{} {
    seed := rand.Int63n(10000)
    switch seed {
    case 0:
        return genRandInt()
    case 1:
        return genRandString()
    case 2:
        return struct {
            num int64
            str string
        }{genRandInt(), genRandString()}
    default:
        const length = 2
        arr := new([length]interface{})
        for i := 0; i < length; i++ {
            if i%2 == 0 {
                arr[i] = genRandInt()
            } else {
                arr[i] = genRandString()
            }
        }
        return *arr
    }
}

func genRandString() string {
    var buff bytes.Buffer
    var prev string
    var curr string
    for i := 0; buff.Len() < 3; i++ {
        curr = string(genRandAZAscii())
        if curr == prev {
            continue
        } else {
            prev = curr
        }
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

func genRandInt() int64 {
    return rand.Int63n(10000)
}
