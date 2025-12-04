/// Code generated from jsonrpc schema by rpcgen v2.7.0; DO NOT EDIT.

import Foundation

extension RPCAPI: RPCMethod {
    public var rpcMethod: String {
        switch self {
        case .batch(let requests): return requests.compactMap { $0.rpcMethod }.joined(separator: ",")
        case .catalogueCheckError: return "catalogue.CheckError"
        case .catalogueCheckZenRPCError: return "catalogue.CheckZenRPCError"
        case .catalogueDivide: return "catalogue.Divide"
        case .catalogueDoSomething: return "catalogue.DoSomething"
        case .catalogueDoSomethingV2: return "catalogue.DoSomethingV2"
        case .catalogueDoSomethingWithPoint: return "catalogue.DoSomethingWithPoint"
        case .catalogueGetPoints: return "catalogue.GetPoints"
        case .catalogueMultiply: return "catalogue.Multiply"
        case .cataloguePi: return "catalogue.Pi"
        case .cataloguePositive: return "catalogue.Positive"
        case .cataloguePow: return "catalogue.Pow"
        case .catalogueSum: return "catalogue.Sum"
        case .catalogueSumArray: return "catalogue.SumArray"
        }
    }
}

extension RPCAPI: RPCParameters {
    public var rpcParameters: [String: Any?]? {
        switch self {
        case .batch:
              return nil
        case let .catalogueCheckError(isErr):
            return ["isErr": isErr]

        case let .catalogueCheckZenRPCError(isErr):
            return ["isErr": isErr]

        case let .catalogueDivide(a, b):
            return ["a": a, "b": b]

        case .catalogueDoSomething:
            return nil

        case .catalogueDoSomethingV2:
            return nil

        case let .catalogueDoSomethingWithPoint(p, pp):
            return ["p": p.any, "pp": pp.any]

        case .catalogueGetPoints:
            return nil

        case let .catalogueMultiply(a, b):
            return ["a": a, "b": b]

        case .cataloguePi:
            return nil

        case .cataloguePositive:
            return nil

        case let .cataloguePow(base, exp):
            return ["base": base, "exp": exp]

        case let .catalogueSum(a, b):
            return ["a": a, "b": b]

        case let .catalogueSumArray(array):
            return ["array": array.any]
        }
    }
}

public enum RPCAPI: Codable, Hashable {
    /// Make batch requests.
    case batch(requests: [RPCAPI])

    /// CheckError throws error is isErr true.
    /// TEST row 2
case catalogueCheckError(isErr: Bool, requestId: String? = nil)
    /// CheckZenRPCError throws zenrpc error is isErr true.
    /// Second description row
case catalogueCheckZenRPCError(isErr: Bool, requestId: String? = nil)
    /// Divide divides two numbers.
    /// - Returns: Quotient?
case catalogueDivide(a: Int, b: Int, requestId: String? = nil)
case catalogueDoSomething
    /// - Returns: ExternalData
case catalogueDoSomethingV2
    /// - Returns: Point
case catalogueDoSomethingWithPoint(p: Point, pp: [Point], requestId: String? = nil)
    /// - Returns: [Point]
case catalogueGetPoints
    /// Multiply multiples two digits and returns result.
    /// - Returns: Int
case catalogueMultiply(a: Int, b: Int, requestId: String? = nil)
    /// Pi returns math.Pi.
    /// - Returns: Double
case cataloguePi
    /// - Returns: Bool
case cataloguePositive
    /// Pow returns x**y, the base-x exponential of y. If Exp is not set then default value is 2.
    /// - Returns: Double
case cataloguePow(base: Double, exp: Double?, requestId: String? = nil)
    /// Sum sums two digits and returns error with error code as result and IP from context.
    /// - Returns: Bool
case catalogueSum(a: Int, b: Int, requestId: String? = nil)
    /// SumArray returns sum all items from array
    /// - Returns: Double
case catalogueSumArray(array: [Double]?, requestId: String? = nil)
}


public struct CycleInitStruct: Codable, Hashable {
    @DecodableDefault.False
    var isCycleInit: Bool
    init(isCycleInit: Bool) {
        self.isCycleInit = isCycleInit
    }
}

public struct ExternalData: Codable, Hashable {
    @DecodableDefault.EmptyString
    var name: String
    init(name: String) {
        self.name = name
    }
}

public struct Point: Codable, Hashable {
    /// coordinate
    @DecodableDefault.IntegerZero
    var X: Int
    /// coordinate
    @DecodableDefault.IntegerZero
    var Y: Int
    /// version group float - 1
    @DecodableDefault.DoubleZero
    var baseFloat: Double
    /// version id - 2
    @DecodableDefault.IntegerZero
    var baseId: Int
    @DecodableDefault.EmptyString
    var `class`: String
    /// version date - 1
    @DecodableDefault.EmptyString
    var createdAt: String
    var emptyString: String?
    /// version id - 1
    @DecodableDefault.IntegerZero
    var id: Int
    /// version group geo coordinate № - 2
    @DecodableDefault.DoubleZero
    var lat: Double
    /// version group geo coordinate № - 3
    @DecodableDefault.DoubleZero
    var latitude: Double
    /// version group geo coordinate № - 2
    @DecodableDefault.DoubleZero
    var lon: Double
    /// version group geo coordinate № - 3
    @DecodableDefault.DoubleZero
    var longitude: Double
    /// version date - 3
    @DecodableDefault.EmptyString
    var manualChangedAt: String
    @DecodableDefault.EmptyString
    var name: String
    /// version group geo coordinate № - 1
    @DecodableDefault.DoubleZero
    var newLat: Double
    /// version group geo coordinate № - 1
    @DecodableDefault.DoubleZero
    var newLon: Double
    var nextQuotient: Quotient
    /// version group float - 2
    @DecodableDefault.DoubleZero
    var secondFloat: Double
    /// version id - 3
    @DecodableDefault.IntegerZero
    var secondID: Int
    @DecodableDefault.EmptyList
    var secondPoints: [Point]
    var secondQuotient: Quotient?
    /// version date - 2
    @DecodableDefault.EmptyString
    var updatedAt: String
    init(X: Int, Y: Int, baseFloat: Double, baseId: Int, class: String, createdAt: String, emptyString: String? = nil, id: Int, lat: Double, latitude: Double, lon: Double, longitude: Double, manualChangedAt: String, name: String, newLat: Double, newLon: Double, nextQuotient: Quotient, secondFloat: Double, secondID: Int, secondPoints: [Point], secondQuotient: Quotient? = nil, updatedAt: String) {
        self.X = X
        self.Y = Y
        self.baseFloat = baseFloat
        self.baseId = baseId
        self.class = `class`
        self.createdAt = createdAt
        self.emptyString = emptyString
        self.id = id
        self.lat = lat
        self.latitude = latitude
        self.lon = lon
        self.longitude = longitude
        self.manualChangedAt = manualChangedAt
        self.name = name
        self.newLat = newLat
        self.newLon = newLon
        self.nextQuotient = nextQuotient
        self.secondFloat = secondFloat
        self.secondID = secondID
        self.secondPoints = secondPoints
        self.secondQuotient = secondQuotient
        self.updatedAt = updatedAt
    }
}

public struct Quotient: Codable, Hashable {
    /// Quo docs
    @DecodableDefault.IntegerZero
    var Quo: Int
    @DecodableDefault.EmptyString
    var baseRow: String
    var data: CycleInitStruct
    /// Rem docs
    @DecodableDefault.IntegerZero
    var rem: Int
    var rowNil: String?
    init(Quo: Int, baseRow: String, data: CycleInitStruct, rem: Int, rowNil: String? = nil) {
        self.Quo = Quo
        self.baseRow = baseRow
        self.data = data
        self.rem = rem
        self.rowNil = rowNil
    }
}

