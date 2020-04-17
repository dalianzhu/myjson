package myjson

import (
    sysjson "encoding/json"
    "errors"
    "fmt"
    "gopkg.in/go-playground/validator.v9"
    "log"
    "reflect"
    "strings"
)

/* Validate 引入 validator
{
	"name": "gt=5;name must greater than 5",
	"info": [{
			"name": "gt=5",
			"years": "gt=2;years must greater than 2"
		},
		"required"
	]
}
*/
func Validate(origin *MyJson, rules *MyJson) error {
    validate := validator.New()
    debugf("Validate origin:%v", origin)
    debugf("Validate rules:%v", rules)
    var err error

    if rules.IsMap() {
        err = validateMap(validate, origin.data, rules.data)
    }
    if rules.IsSlice() {
        err = validateSlice(validate, origin.data, rules.data)
    }

    return err
}

// origin rule永远为同一种类型，否则是rule出错
func validateMap(validate *validator.Validate, origin interface{}, rulesInfos interface{}) error {
    //debugf("validateMap origin:%v", origin)
    //debugf("validateMap rules:%v", rulesInfos)
    switch v := rulesInfos.(type) {
    case map[string]interface{}:
        ov, ok := origin.(map[string]interface{})
        if !ok {
            log.Printf("origin:%v is not map\n", origin)
            return fmt.Errorf("origin is not map: %v", origin)
        }

        for key := range v {
            if key == "_required" {
                continue
            }
            debugf("map key:%v", key)
            err := validateValue(validate, key, ov[key], v[key])
            if err != nil {
                return err
            }
        }
    default:
        return fmt.Errorf("rulesInfos:%v is not maps", origin)
    }
    return nil
}

func validateSlice(validate *validator.Validate, origin interface{}, ruleInfos interface{}) error {
    switch v := ruleInfos.(type) {
    case []interface{}:

        ov, ok := origin.([]interface{})
        if !ok {
            return fmt.Errorf("validateSlice origin is not slice: %v", origin)
        }

        for i := range ov {
            ret := validateValue(validate, "", ov[i], v[0])
            if ret != nil {
                return ret
            }
        }
    default:
        return fmt.Errorf("ruleInfos:%v is not maps", origin)
    }
    return nil
}

type ruleCache struct {
    rule     string
    errMsg   string
    required bool
}
type isRequired struct {
    isRequired bool
    errMsg     string
}

// cutSemicolon, 按分号切割字符串，注意\;会被忽略
func cutSemicolon(str string) []string {
    var ret []string
    tp := ""
    runeArray := []rune(str)
    for i := 0; i < len(runeArray); i++ {
        c := string(runeArray[i])
        if i+1 < len(runeArray) && string(runeArray[i+1]) == ";" && c == ";" {
            // 处理 ;;
            tp += ";"
            i += 1
            continue
        }
        if c == ";" {
            ret = append(ret, tp)
            tp = ""
            continue
        }
        tp += c
    }
    if tp != "" {
        ret = append(ret, tp)
    }
    return ret
}

func checkSliceRequired(ruleSlice []interface{}) (bool, string) {
    if len(ruleSlice) <= 1 {
        return false, ""
    }
    v, ok := ruleSlice[1].(*isRequired)
    if ok {
        debugf("checkSliceRequired:%v, %v", v.isRequired, v.errMsg)
        return v.isRequired, v.errMsg
    }
    return false, ""
}

func validateValue(validate *validator.Validate,
    key string, origin interface{}, info interface{}) error {
    debugf("validateValue origin:%v, info:%v", origin, info)

    switch v := info.(type) {
    case []interface{}:
        if len(v) == 0 {
            return fmt.Errorf("info:%v rules is empty", info)
        }
        isRequired, requiredMsg := checkSliceRequired(v)
        if origin == nil {
            if isRequired {
                if requiredMsg == "" {
                    return fmt.Errorf("key:%v value cannot be empty", key)
                } else {
                    return errors.New(requiredMsg)
                }
            }
            return nil
        }
        return validateSlice(validate, origin, v)
    case map[string]interface{}:
        requiredMsg, isRequired := v["_required"]
        if origin == nil {
            if isRequired {
                requiredMsgStr := ToStr(requiredMsg)
                if requiredMsgStr != "" {
                    return errors.New(requiredMsgStr)
                }
                return fmt.Errorf("key:%v cannot be empty", key)
            }
            return nil
        }
        return validateMap(validate, origin, v)
    }

    cache := info.(*ruleCache)
    debugf("infoDataArray,origin:%v,originType:%v rule:%v, msg:%v", origin, reflect.TypeOf(origin), cache.rule, cache.errMsg)
    _, ok := origin.(sysjson.Number)
    if ok {
        originNum, _ := ToInt(origin)
        origin = originNum
    }
    if cache.required == false && origin == nil {
        return nil
    }
    if err := validate.Var(origin, cache.rule); err != nil {
        debugf("validateValue default return false,%v", err)
        if cache.errMsg != "" {
            return errors.New(cache.errMsg)
        }
        return fmt.Errorf("%v check error!! rules:%v", origin, cache.rule)
    }
    return nil
}

func NewValidateRules(rulesData string) *MyJson {
    rules := NewJson(rulesData)

    if rules.IsMap() {
        initRulesMap(rules.data.(map[string]interface{}))
    } else if rules.IsSlice() {
        initRulesSlice(rules.data.([]interface{}))
    }
    return rules
}

func initRulesMap(mapValue map[string]interface{}) {
    for key := range mapValue {
        value := mapValue[key]
        initRulesValue(value, key, -1, mapValue)
    }
}

func initRulesSlice(sliceValue []interface{}) {
    for i, value := range sliceValue {
        initRulesValue(value, "", i, sliceValue)
    }
}

func initRulesValue(i interface{}, key string, index int, parent interface{}) {
    switch v := i.(type) {
    case map[string]interface{}:
        initRulesMap(v)
    case []interface{}:
        initRulesSlice(v)
    // 说明info是一个原始的规则信息
    case string:
        debugf("info is string, key:%v", key)
        infoDataArray := cutSemicolon(v)
        cache := new(ruleCache)
        cache.required = false
        if len(infoDataArray) >= 2 {
            cache.rule = infoDataArray[0]
            cache.errMsg = infoDataArray[1]

        } else if len(infoDataArray) > 0 {
            cache.rule = infoDataArray[0]
        }
        if strings.Contains(cache.rule, "required") {
            cache.required = true
        }

        switch pv := parent.(type) {
        case []interface{}:
            if index == 0 {
                pv[0] = cache
            } else if index == 1 {
                // 处理数组的required
                required := new(isRequired)
                debugf("slice is required msg:%v", v)
                tagArrays := cutSemicolon(v)
                if len(tagArrays) >= 1 && tagArrays[0] == "_required" {
                    required.isRequired = true
                }
                if len(tagArrays) >= 2 {
                    required.errMsg = tagArrays[1]
                }
                pv[1] = required
            }
        case map[string]interface{}:
            pv[key] = cache
        }
    }
}
