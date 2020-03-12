local tree = {} 

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
function tree.And(arr)
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
function tree.Or(arr)
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
function tree.Xor(arr)
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

-- The operator surrogate provides a special means to handle logical expressions with missing 
-- values. It is applied to a sequence of predicates. The order of the predicates matters, the 
-- first predicate is the primary, the next predicates are the surrogates. Evaluation order is
-- left-to-right. The cascaded predicates are applied when the primary predicate evaluates to 
-- UNKNOWN.
function tree.Surrogate(arr)
    local result = nil
    for i=1, arr.n do
        result = arr[i]
        if not Unknown(result) then
            return result
        end
    end
    return result
end

-- Checks if the value is present in the set
function tree.IsIn(target, array)
    for i, v in ipairs(array) do
        if v == target then
            return true
        end
    end

    return false
end

-- Checks if the value is missing in the set
function tree.IsNotIn(target, array)
	return not tree.IsIn(target, array)
end

-- NewNode creates a node structure
function tree.NewNode(id, run, ...)
    n = {}
    n.id = id
    n.run = run
    n.children = {}

    -- Create a table with all of the children
    for k, v in pairs({ ... }) do
        n.children[v.id] = v
    end

    -- The evaluation function
    n.eval = function(t, v)
        local done = n.run(t, n, v); if done then
            return
        end
        
        -- Evaluate siblings
        for k, n in pairs(nodes) do
            local done = n.eval(t, v); if done then
                return
            end
        end 
    end
    
    return n
end

-- NewTree creates a new decision tree
function tree.NewTree(strategy, node)
    t = {}
    t.onMiss = strategy
    t.node = node
    t.conf = {}
    t.last = nil

    -- Function which traverses the tree
    t.eval = function(v)
        t.node.eval(t, v)
        return t.last
    end

    -- Function that records an unknown value
    t.miss = function(score, count, confidence)
        t.conf[score] = t.conf[score] or {}
        table.insert(t.conf[score], {count, confidence})
    end
    return t
end

-- If a Node's predicate evaluates to UNKNOWN while traversing the tree, evaluation is stopped
-- and the current winner is returned as the final prediction.
function tree.LastPrediction(t, n, v)
    return true -- stop
end

-- If a Node's predicate value evaluates to UNKNOWN while traversing the tree, abort the scoring
-- process and give no prediction.
function tree.NullPrediction(t, n, v)
    t.last = nil
    return true -- stop
end

-- Comparisons with missing values other than checks for missing values always evaluate to FALSE. 
-- If no rule fires, then use the noTrueChildStrategy to decide on a result. This option requires
-- that missing values be handled after all rules at the Node have been evaluated.
-- Note: In contrast to lastPrediction, evaluation is carried on instead of stopping immediately
-- upon first discovery of a Node who's predicate value cannot be determined due to missing 
-- values.
function tree.None(t, n, v)
    return false -- continue
end

-- If a Node's predicate value evaluates to UNKNOWN while traversing the tree, evaluate the 
-- attribute defaultChild which gives the child to continue traversing with. Requires the 
-- presence of the attribute defaultChild in every non-leaf Node.
function tree.DefaultChild(t, n, v)
    --return t.next(v, {n.children[n.def]})
    return true
end

-- If a Node's predicate value evaluates to UNKNOWN while traversing the tree, the confidences for
-- each class is calculated from scoring it and each of its sibling Nodes in turn (excluding any 
-- siblings whose predicates evaluate to FALSE). The confidences returned for each class from each 
-- sibling Node that was scored are weighted by the proportion of the number of records in that Node,
-- then summed to produce a total confidence for each class. The winner is the class with the highest
-- confidence. Note that weightedConfidence should be applied recursively to deal with situations 
-- where several predicates within the tree evaluate to UNKNOWN during the scoring of a case.
function tree.WeightedConfidence(t, n, v)
    return true
end
    
-- If a Node's predicate value evaluates to UNKNOWN while traversing the tree, we consider evaluation
-- of the Node's predicate being TRUE and follow this Node. In addition, subsequent Nodes to the 
-- initial Node are evaluated as well. This procedure is applied recursively for each Node being 
-- evaluated until a leaf Node is reached. All leaf Nodes being reached by this procedure are 
-- aggregated such that for each value attribute of such a leaf Node's ScoreDistribution element the
-- corresponding recordCount attribute values are accumulated. The value associated with the highest
-- recordCount accumulated through this procedure is predicted. The basic idea of aggregateNodes is 
-- to aggregate all leaf Nodes which may be reached by a record with one or more missing values 
-- considering all possible values. Strategy aggregateNodes calculates a virtual Node and predicts a
-- score according to this virtual Node. Requires the presence of attribute recordCount in all 
-- ScoreDistribution elements.
function tree.AggregateNodes(t, n, v)
    
end
    



-- Checks if the value is missing
function Unknown(v)
	return v == nil or v == ''
end

-- Inverted AND operator
function nand(a, b)
    return not (a and b)
end
 
return tree