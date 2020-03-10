local eval = {} 

-- http://dmg.org/pmml/v4-1/TreeModel.html
-- P       Q       AND     OR      XOR
-- True	True	True	True	False
-- True	False	False	True	True
-- True	Unknown	Unknown	True	Unknown
-- False	True	False	True	True
-- False	False	False	False	False
-- False	Unknown	False	Unknown	Unknown
-- Unknown	True	Unknown	True	Unknown
-- Unknown	False	False	Unknown	Unknown
-- Unknown	Unknown	Unknown	Unknown	Unknown




-- Checks if the value is missing
function Unknown(v)
	return v == nil or v == ''
end


-- And performs a logical AND operation on the set of nullable boolean values
-- true     true    true
-- true     false   false
-- true     nil     nil
-- false    true    false
-- false    false   false
-- false    nil     false
-- nil      true    nil
-- nil      false   false
-- nil      nil     nil
function eval.And(arr)
    local result = true
    for i=1, arr.n do
        local v = arr[i]
        if Unknown(v) and result == true then
            result = nil
        elseif Unknown(result) and v == false then
            result = false
        else
            result = result and v
        end
    end
    return result
end

-- Or performs a logical OR operation on a set of nullable boolean values
-- true     true	true
-- true     false	true
-- true     nil	    true
-- false	true	true
-- false	false	false
-- false	nil	    nil
-- nil      true	true
-- nil	    false	nil
-- nil	    nil	    nil
function eval.Or(arr)
    local result = false
    for i=1, arr.n do
        local v = arr[i]
        if Unknown(v) and result == false then
            result = nil
        elseif Unknown(result) and v == false then
            result = nil
        else
            result = result or v
        end
    end
    return result
end

-- Xor performs a logical XOR operation on a set of nullable boolean values
-- true	    true	false
-- true	    false	true
-- true	    nil	    nil
-- false	true	true
-- false	false	false
-- false	nil	    nil
-- nil	    true	nil
-- nil	    false	nil
-- nil	    nil	    nil
function Xor(arr)
    local result = false
    for i=1, arr.n do
        local v = arr[i]
        if Unknown(v) then
            return nil
        end
        result = nand(nand(result, nand(result, v)), nand(v, nand(result, v)))
    end
    return result
end

-- Inverted AND operator
function nand(a, b)
    return not (a and b)
end
 
return eval