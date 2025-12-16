/// Code generated from jsonrpc schema by rpcgen v2.5.x with swift v1.0.0; DO NOT EDIT.

import Foundation

extension RPCAPI: RPCMethod {
    public var rpcMethod: String {
        switch self {
        case .batch(let requests): return requests.compactMap { $0.rpcMethod }.joined(separator: ",")
        case .arithCheckError: return "arith.CheckError"
        case .arithCheckZenRPCError: return "arith.CheckZenRPCError"
        case .arithDivide: return "arith.Divide"
        case .arithDoSomething: return "arith.DoSomething"
        case .arithDoSomethingV2: return "arith.DoSomethingV2"
        case .arithDoSomethingWithPoint: return "arith.DoSomethingWithPoint"
        case .arithGetPoints: return "arith.GetPoints"
        case .arithMultiply: return "arith.Multiply"
        case .arithPi: return "arith.Pi"
        case .arithPositive: return "arith.Positive"
        case .arithPow: return "arith.Pow"
        case .arithSum: return "arith.Sum"
        case .arithSumArray: return "arith.SumArray"
        }
    }
}

extension RPCAPI: RPCParameters {
    public var rpcParameters: [String: Any?]? {
        switch self {
        case .batch:
              return nil
        case let .arithCheckError(isErr, _):
            return ["isErr": isErr]

        case let .arithCheckZenRPCError(isErr, _):
            return ["isErr": isErr]

        case let .arithDivide(a, b, _):
            return ["a": a, "b": b]

        case .arithDoSomething:
            return nil

        case .arithDoSomethingV2:
            return nil

        case let .arithDoSomethingWithPoint(p, pp, _):
            return ["p": p.any, "pp": pp.any]

        case .arithGetPoints:
            return nil

        case let .arithMultiply(a, b, _):
            return ["a": a, "b": b]

        case .arithPi:
            return nil

        case .arithPositive:
            return nil

        case let .arithPow(base, exp, _):
            return ["base": base, "exp": exp]

        case let .arithSum(a, b, _):
            return ["a": a, "b": b]

        case let .arithSumArray(array, _):
            return ["array": array.any]
        }
    }
}

public enum RPCAPI: Codable, Hashable {
    /// Make batch requests.
    case batch(requests: [RPCAPI])

    /// CheckError throws error is isErr true.
    /// TEST row 2
	case arithCheckError(isErr: Bool, requestId: String? = nil)
    /// CheckZenRPCError throws zenrpc error is isErr true.
    /// Second description row
	case arithCheckZenRPCError(isErr: Bool, requestId: String? = nil)
    /// Divide divides two numbers.
    /// - Returns: Quotient?
	case arithDivide(a: Int, b: Int, requestId: String? = nil)
	case arithDoSomething(requestId: String? = nil)
    /// - Returns: ExternalData
	case arithDoSomethingV2(requestId: String? = nil)
    /// - Returns: Point
	case arithDoSomethingWithPoint(p: Point, pp: [Point], requestId: String? = nil)
    /// - Returns: [Point]
	case arithGetPoints(requestId: String? = nil)
    /// Multiply multiples two digits and returns result.
    /// - Returns: Int
	case arithMultiply(a: Int, b: Int, requestId: String? = nil)
    /// Pi returns math.Pi.
    /// - Returns: Double
	case arithPi(requestId: String? = nil)
    /// - Returns: Bool
	case arithPositive(requestId: String? = nil)
    /// Pow returns x**y, the base-x exponential of y. If Exp is not set then default value is 2.
    /// - Returns: Double
	case arithPow(base: Double, exp: Double?, requestId: String? = nil)
    /// Sum sums two digits and returns error with error code as result and IP from context.
    /// - Returns: Bool
	case arithSum(a: Int, b: Int, requestId: String? = nil)
    /// SumArray returns sum all items from array
    /// - Returns: Double
	case arithSumArray(array: [Double]?, requestId: String? = nil)
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


extension RPCAPI {
  public var rpcId: String? {
    switch self {
    case .batch:
      return nil

    case.arithCheckError(_, let requestId),
       .arithCheckZenRPCError(_, let requestId),
       .arithDivide(_, _, let requestId),
       .arithDoSomething(let requestId),
       .arithDoSomethingV2(let requestId),
       .arithDoSomethingWithPoint(_, _, let requestId),
       .arithGetPoints(let requestId),
       .arithMultiply(_, _, let requestId),
       .arithPi(let requestId),
       .arithPositive(let requestId),
       .arithPow(_, _, let requestId),
       .arithSum(_, _, let requestId),
       .arithSumArray(_, let requestId):
          return requestId
    }
  }
}
