/// Code generated from jsonrpc schema by rpcgen v2.8.0; DO NOT EDIT.

import Foundation

protocol ArithNetworking {
    /// CheckError throws error is isErr true.
    /// TEST row 2
    func arithCheckError(isErr: Bool) async -> RpcError?
    /// CheckZenRPCError throws zenrpc error is isErr true.
    /// Second description row
    func arithCheckZenRPCError(isErr: Bool) async -> RpcError?
    /// Divide divides two numbers.
    func arithDivide(a: Int, b: Int) async -> Result<Quotient, RpcError>
    func arithDoSomething() async -> RpcError?
    func arithDoSomethingV2() async -> Result<ExternalData, RpcError>
    func arithDoSomethingWithPoint(p: Point, pp: [Point]) async -> Result<Point, RpcError>
    func arithGetPoints() async -> Result<[Point], RpcError>
    /// Multiply multiples two digits and returns result.
    func arithMultiply(a: Int, b: Int) async -> Result<Int, RpcError>
    /// Pi returns math.Pi.
    func arithPi() async -> Result<Double, RpcError>
    func arithPositive() async -> Result<Bool, RpcError>
    /// Pow returns x**y, the base-x exponential of y. If Exp is not set then default value is 2.
    func arithPow(base: Double, exp: Double?) async -> Result<Double, RpcError>
    /// Sum sums two digits and returns error with error code as result and IP from context.
    func arithSum(a: Int, b: Int) async -> Result<Bool, RpcError>
    /// SumArray returns sum all items from array
    func arithSumArray(array: [Double]?) async -> Result<Double, RpcError>
}

extension Networking: ArithNetworking {
    /// CheckError throws error is isErr true.
    /// TEST row 2
	/// - Returns: Result<RpcError>
    func arithCheckError(isErr: Bool) async -> RpcError? {
        await request(.arithCheckError(isErr: isErr))
    }

    /// CheckZenRPCError throws zenrpc error is isErr true.
    /// Second description row
	/// - Returns: Result<RpcError>
    func arithCheckZenRPCError(isErr: Bool) async -> RpcError? {
        await request(.arithCheckZenRPCError(isErr: isErr))
    }

    /// Divide divides two numbers.
	/// - Parameters:
	///  - a : the a
	///  - b : the b
	/// - Returns: Result<Quotient, RpcError>
    func arithDivide(a: Int, b: Int) async -> Result<Quotient, RpcError> {
        await request(.arithDivide(a: a, b: b))
    }

	/// - Returns: Result<RpcError>
    func arithDoSomething() async -> RpcError? {
        await request(.arithDoSomething())
    }

	/// - Returns: Result<ExternalData, RpcError>
    func arithDoSomethingV2() async -> Result<ExternalData, RpcError> {
        await request(.arithDoSomethingV2())
    }

	/// - Returns: Result<Point, RpcError>
    func arithDoSomethingWithPoint(p: Point, pp: [Point]) async -> Result<Point, RpcError> {
        await request(.arithDoSomethingWithPoint(p: p, pp: pp))
    }

	/// - Returns: Result<[Point], RpcError>
    func arithGetPoints() async -> Result<[Point], RpcError> {
        await request(.arithGetPoints())
    }

    /// Multiply multiples two digits and returns result.
	/// - Returns: Result<Int, RpcError>
    func arithMultiply(a: Int, b: Int) async -> Result<Int, RpcError> {
        await request(.arithMultiply(a: a, b: b))
    }

    /// Pi returns math.Pi.
	/// - Returns: Result<Double, RpcError>
    func arithPi() async -> Result<Double, RpcError> {
        await request(.arithPi())
    }

	/// - Returns: Result<Bool, RpcError>
    func arithPositive() async -> Result<Bool, RpcError> {
        await request(.arithPositive())
    }

    /// Pow returns x**y, the base-x exponential of y. If Exp is not set then default value is 2.
	/// - Parameters:
	///  - exp : exponent could be empty
	/// - Returns: Result<Double, RpcError>
    func arithPow(base: Double, exp: Double? = nil) async -> Result<Double, RpcError> {
        await request(.arithPow(base: base, exp: exp))
    }

    /// Sum sums two digits and returns error with error code as result and IP from context.
	/// - Returns: Result<Bool, RpcError>
    func arithSum(a: Int, b: Int) async -> Result<Bool, RpcError> {
        await request(.arithSum(a: a, b: b))
    }

    /// SumArray returns sum all items from array
	/// - Returns: Result<Double, RpcError>
    func arithSumArray(array: [Double]? = nil) async -> Result<Double, RpcError> {
        await request(.arithSumArray(array: array))
    }
}
