package inject

import (
	"reflect"
	"sync"
	"log"
	"strings"
)

type Inject interface {
	ServiceByName(serviceName string) interface{}
	Service(ptr interface{}) interface{}
	Apply(service interface{})
	Inject()
	Map(key string, value interface{})
	MapStruct(key string, value *struct{})
	MapBool(key string, value *bool)
	MapInt(key string, value *int)
	MapInt8(key string, value *int8)
	MapInt16(key string, value *int16)
	MapInt32(key string, value *int32)
	MapInt64(key string, value *int64)

	MapUint(key string, value *uint)
	MapUint8(key string, value *uint8)
	MapUint16(key string, value *uint16)
	MapUint32(key string, value *uint32)
	MapUint64(key string, value *uint64)

	MapFloat32(key string, value *float32)
	MapFloat64(key string, value *float64)

	MapString(key string, value *string)

	MapByte(key string, value *byte)

	MapStructs(key string, value *[]struct{})
	MapBools(key string, values *[]bool)
	MapInts(key string, values *[]int)
	MapInt8s(key string, values *[]int8)
	MapInt16s(key string, values *[]int16)
	MapInt32s(key string, values *[]int32)
	MapInt64s(key string, values *[]int64)

	MapUints(key string, values *[]uint)
	MapUint8s(key string, values *[]uint8)
	MapUint16s(key string, values *[]uint16)
	MapUint32s(key string, values *[]uint32)
	MapUint64s(key string, values *[]uint64)

	MapFloat32s(key string, values *[]float32)
	MapFloat64s(key string, values *[]float64)

	MapStrings(key string, values *[]string)
	MapBytes(key string, values *[]byte)

	Get(key string) interface{}
	GetStruct(key string) *struct{}
	GetBool(key string) *bool
	GetInt(key string) *int
	GetInt8(key string) *int8
	GetInt16(key string) *int16
	GetInt32(key string) *int32
	GetInt64(key string) int64

	GetUint(key string) *uint
	GetUint8(key string) *uint8
	GetUint16(key string) *uint16
	GetUint32(key string) *uint32
	GetUint64(key string) *uint64

	GetFloat32(key string) *float32
	GetFloat64(key string) *float64

	GetString(key string) *string

	GetByte(key string) *byte
	GetStructs(key string) *[]struct{}
	GetBools(key string) *[]bool
	GetInts(key string) *[]int
	GetInt8s(key string) *[]int8
	GetInt16s(key string) *[]int16
	GetInt32s(key string) *[]int32
	GetInt64s(key string) *[]int64

	GetUints(key string) *[]uint
	GetUint8s(key string) *[]uint8
	GetUint16s(key string) *[]uint16
	GetUint32s(key string) *[]uint32
	GetUint64s(key string) *[]uint64

	GetFloat32s(key string) *[]float32
	GetFloat64s(key string) *[]float64

	GetStrings(key string) *[]string
	GetBytes(key string) *[]byte
}

type Injector struct {
	beanMap map[string]interface{}
	Locker  sync.RWMutex
}

func New() *Injector {
	return &Injector{
		beanMap: make(map[string]interface{}),
		Locker:  sync.RWMutex{},
	}
}

func (inject *Injector) ServiceByName(serviceName string) interface{} {
	return reflect.ValueOf(inject.Get(serviceName)).Elem().Interface()
}
func (inject *Injector) ServicePtrByName(serviceName string) interface{} {
	return reflect.ValueOf(inject.Get(serviceName)).Interface()
}
func (inject *Injector) Service(ptr interface{}) interface{} {
	return reflect.ValueOf(inject.Get(reflect.TypeOf(ptr).String())).Elem().Interface()
}
func (inject *Injector) ServiceByPrefixName(prefix string) (services []interface{}) {
	for k, v := range inject.beanMap {
		if strings.HasPrefix(k, prefix) {
			services = append(services, reflect.ValueOf(v).Elem().Interface())
		}
	}
	return services
}
func (inject *Injector) Services() (services []interface{}) {
	for _, v := range inject.beanMap {
		services = append(services, reflect.ValueOf(v).Elem().Interface())
	}
	return services
}
func (inject *Injector) Apply(services ... interface{}) {
	for _, ser := range services {
		inject.replyOnInject(ser)
	}
}
func (inject *Injector) ApplyWithName(name string, service interface{}) {
	if reflect.TypeOf(service).Kind() == reflect.Ptr {
		inject.Map(name, service)
	} else if reflect.TypeOf(service).Kind() == reflect.Func {
		inject.Map(name, service)
	} else if reflect.TypeOf(service).Kind() == reflect.Chan {
		inject.Map(name, service)
	} else if reflect.TypeOf(service).Kind() == reflect.Map {
		inject.Map(name, service)
	} else if reflect.TypeOf(service).Kind() == reflect.Slice {
		inject.Map(name, service)
	} else if reflect.TypeOf(service).Kind() == reflect.Array {
		inject.Map(name, service)
	} else {
		log.Fatal("no support type")
	}
}
func (inject *Injector) RegisterService(services ... interface{}) {
	for _, service := range services {
		inject.Apply(service)
	}
	inject.Inject()
}
func (inject *Injector) RegisterServiceWithName(name string, service interface{}) {
	inject.ApplyWithName(name, service)
	inject.Inject()
}
func (inject *Injector) Inject() {
	for _, value := range inject.beanMap {
		inject.injectWithReply(value)
	}

}

func (inject *Injector) injectWithReply(service interface{}) {
	value := reflect.ValueOf(service)
	vl := value.Elem()
	tp := reflect.TypeOf(service).Elem()
	if vl.Kind() == reflect.Struct {
		for i := 0; i < vl.NumField(); i++ {
			field := tp.Field(i)
			if injectNameTag, ok := field.Tag.Lookup("inject"); ok {
				if injectNameTag == "" {
					injectNameTag = field.Type.String()
				}
				result := inject.Get(injectNameTag)
				if result == nil {
					log.Printf("no found `%s` inject from %s", injectNameTag, tp.Name())
				}
				vl.Field(i).Set(reflect.ValueOf(result))
			}
		}
	}
}

func (inject *Injector) replyOnInject(service interface{}) {
	tp := reflect.TypeOf(service)
	if tp.Kind() == reflect.Ptr {
		vl := reflect.ValueOf(service)
		inject.ApplyWithName(vl.Type().String(), service)
	} else if tp.Kind() == reflect.Struct {
		log.Fatalf("no support struct:%s",tp.Name())
	} else if tp.Kind() == reflect.Interface {
		log.Fatalf("no support interface:%s",tp.Name())
	} else if tp.Kind() == reflect.Chan {
		log.Fatalf("no support chan:%s",tp.Name())
	} else if tp.Kind() == reflect.Slice {
		log.Fatalf("no support slice:%s",tp.Name())
	} else if tp.Kind() == reflect.Array {
		log.Fatalf("no support array:%s",tp.Name())
	} else {
		log.Fatalf("no support type:%s",tp.Name())
	}
}

func (inject *Injector) Map(key string, value interface{}) {
	inject.Locker.Lock()
	defer inject.Locker.Unlock()
	inject.beanMap[key] = value
}
func (inject *Injector) MapStruct(key string, value *struct{}) {
	inject.Map(key, value)
}
func (inject *Injector) MapBool(key string, value *bool) {
	inject.Map(key, value)
}
func (inject *Injector) MapInt(key string, value *int) {
	inject.Map(key, value)
}
func (inject *Injector) MapInt8(key string, value *int8) {
	inject.Map(key, value)
}
func (inject *Injector) MapInt16(key string, value *int16) {
	inject.Map(key, value)
}
func (inject *Injector) MapInt32(key string, value *int32) {
	inject.Map(key, value)
}
func (inject *Injector) MapInt64(key string, value *int64) {
	inject.Map(key, value)
}

func (inject *Injector) MapUint(key string, value *uint) {
	inject.Map(key, value)
}
func (inject *Injector) MapUint8(key string, value *uint8) {
	inject.Map(key, value)
}
func (inject *Injector) MapUint16(key string, value *uint16) {
	inject.Map(key, value)
}
func (inject *Injector) MapUint32(key string, value *uint32) {
	inject.Map(key, value)
}
func (inject *Injector) MapUint64(key string, value *uint64) {
	inject.Map(key, value)
}

func (inject *Injector) MapFloat32(key string, value *float32) {
	inject.Map(key, value)
}
func (inject *Injector) MapFloat64(key string, value *float64) {
	inject.Map(key, value)
}

func (inject *Injector) MapString(key string, value *string) {
	inject.Map(key, value)
}

func (inject *Injector) MapByte(key string, value *byte) {
	inject.Map(key, value)
}

func (inject *Injector) MapStructs(key string, values *[]struct{}) {
	inject.Map(key, values)
}

func (inject *Injector) MapBools(key string, values *[]bool) {
	inject.Map(key, values)
}

func (inject *Injector) MapInts(key string, values *[]int) {
	inject.Map(key, values)
}
func (inject *Injector) MapInt8s(key string, values *[]int8) {
	inject.Map(key, values)
}
func (inject *Injector) MapInt16s(key string, values *[]int16) {
	inject.Map(key, values)
}
func (inject *Injector) MapInt32s(key string, values *[]int32) {
	inject.Map(key, values)
}
func (inject *Injector) MapInt64s(key string, values *[]int64) {
	inject.Map(key, values)
}

func (inject *Injector) MapUints(key string, values *[]uint) {
	inject.Map(key, values)
}
func (inject *Injector) MapUint8s(key string, values *[]uint8) {
	inject.Map(key, values)
}
func (inject *Injector) MapUint16s(key string, values *[]uint16) {
	inject.Map(key, values)
}
func (inject *Injector) MapUint32s(key string, values *[]uint32) {
	inject.Map(key, values)
}
func (inject *Injector) MapUint64s(key string, values *[]uint64) {
	inject.Map(key, values)
}

func (inject *Injector) MapFloat32s(key string, values *[]float32) {
	inject.Map(key, values)
}
func (inject *Injector) MapFloat64s(key string, values *[]float64) {
	inject.Map(key, values)
}

func (inject *Injector) MapStrings(key string, values *[]string) {
	inject.Map(key, values)
}
func (inject *Injector) MapBytes(key string, values *[]byte) {
	inject.Map(key, values)
}

func (inject *Injector) Get(key string) interface{} {
	inject.Locker.RLock()
	defer inject.Locker.RUnlock()
	value := inject.beanMap[key]
	return value
}

func (inject *Injector) GetStruct(key string) *struct{} {
	return inject.Get(key).(*struct{})
}

func (inject *Injector) GetBool(key string) *bool {
	return inject.Get(key).(*bool)
}

func (inject *Injector) GetInt(key string) *int {
	return inject.Get(key).(*int)
}
func (inject *Injector) GetInt8(key string) *int8 {
	return inject.Get(key).(*int8)
}
func (inject *Injector) GetInt16(key string) *int16 {
	return inject.Get(key).(*int16)
}
func (inject *Injector) GetInt32(key string) *int32 {
	return inject.Get(key).(*int32)
}
func (inject *Injector) GetInt64(key string) *int64 {
	return inject.Get(key).(*int64)
}

func (inject *Injector) GetUint(key string) *uint {
	return inject.Get(key).(*uint)
}
func (inject *Injector) GetUint8(key string) *uint8 {
	return inject.Get(key).(*uint8)
}
func (inject *Injector) GetUint16(key string) *uint16 {
	return inject.Get(key).(*uint16)
}
func (inject *Injector) GetUint32(key string) *uint32 {
	return inject.Get(key).(*uint32)
}
func (inject *Injector) GetUint64(key string) *uint64 {
	return inject.Get(key).(*uint64)
}

func (inject *Injector) GetFloat32(key string) *float32 {
	return inject.Get(key).(*float32)
}
func (inject *Injector) GetFloat64(key string) *float64 {
	return inject.Get(key).(*float64)
}

func (inject *Injector) GetString(key string) *string {
	return inject.Get(key).(*string)
}

func (inject *Injector) GetByte(key string) *byte {
	return inject.Get(key).(*byte)
}

func (inject *Injector) GetStructs(key string) *[]struct{} {
	return inject.Get(key).(*[]struct{})
}

func (inject *Injector) GetInts(key string) *[]int {
	return inject.Get(key).(*[]int)
}
func (inject *Injector) GetInt8s(key string) *[]int8 {
	return inject.Get(key).(*[]int8)
}
func (inject *Injector) GetInt16s(key string) *[]int16 {
	return inject.Get(key).(*[]int16)
}
func (inject *Injector) GetInt32s(key string) *[]int32 {
	return inject.Get(key).(*[]int32)
}
func (inject *Injector) GetInt64s(key string) *[]int64 {
	return inject.Get(key).(*[]int64)
}

func (inject *Injector) GetUints(key string) *[]uint {
	return inject.Get(key).(*[]uint)
}
func (inject *Injector) GetUint8s(key string) *[]uint8 {
	return inject.Get(key).(*[]uint8)
}
func (inject *Injector) GetUint16s(key string) *[]uint16 {
	return inject.Get(key).(*[]uint16)
}
func (inject *Injector) GetUint32s(key string) *[]uint32 {
	return inject.Get(key).(*[]uint32)
}
func (inject *Injector) GetUint64s(key string) *[]uint64 {
	return inject.Get(key).(*[]uint64)
}

func (inject *Injector) GetFloat32s(key string) *[]float32 {
	return inject.Get(key).(*[]float32)
}
func (inject *Injector) GetFloat64s(key string) *[]float64 {
	return inject.Get(key).(*[]float64)
}

func (inject *Injector) GetStrings(key string) *[]string {
	return inject.Get(key).(*[]string)
}
func (inject *Injector) GetBytes(key string) *[]byte {
	return inject.Get(key).(*[]byte)
}
func (inject *Injector) GetBools(key string) *[]bool {
	return inject.Get(key).(*[]bool)
}
