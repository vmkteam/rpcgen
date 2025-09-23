/// Code generated from jsonrpc schema by rpcgen v2.4.5; DO NOT EDIT.

import Foundation

protocol ArithNetworking {
    /// CheckError throws error is isErr true.
    func arithCheckError(isErr: Bool) async -> Result<RpcError>
    /// CheckZenRPCError throws zenrpc error is isErr true.
    func arithCheckZenRPCError(isErr: Bool) async -> Result<RpcError>
    /// Divide divides two numbers.
    func arithDivide(a: Int, b: Int) async -> Result<Quotient, RpcError>
    func arithDoSomething() async -> Result<RpcError>
    func arithDoSomethingWithPoint(p: ModelPoint) async -> Result<ModelPoint, RpcError>
    func arithGetPoints() async -> Result<[model.Point], RpcError>
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
    func arithCheckError(isErr: Bool) async -> Result<RpcError> {
        await request(.arithCheckError(isErr: isErr))
    }

    /// CheckZenRPCError throws zenrpc error is isErr true.
    func arithCheckZenRPCError(isErr: Bool) async -> Result<RpcError> {
        await request(.arithCheckZenRPCError(isErr: isErr))
    }

    /// Divide divides two numbers.
    func arithDivide(a: Int,b: Int) async -> Result<Quotient, RpcError> {
        await request(.arithDivide(a: a, b: b))
    }

    func arithDoSomething() async -> Result<RpcError> {
        await request(.arithDoSomething())
    }

    func arithDoSomethingWithPoint(p: ModelPoint) async -> Result<ModelPoint, RpcError> {
        await request(.arithDoSomethingWithPoint(p: p))
    }

    func arithGetPoints() async -> Result<[model.Point], RpcError> {
        await request(.arithGetPoints())
    }

    /// Multiply multiples two digits and returns result.
    func arithMultiply(a: Int,b: Int) async -> Result<Int, RpcError> {
        await request(.arithMultiply(a: a, b: b))
    }

    /// Pi returns math.Pi.
    func arithPi() async -> Result<Double, RpcError> {
        await request(.arithPi())
    }

    func arithPositive() async -> Result<Bool, RpcError> {
        await request(.arithPositive())
    }

    /// Pow returns x**y, the base-x exponential of y. If Exp is not set then default value is 2.
    func arithPow(base: Double,exp: Double? = nil) async -> Result<Double, RpcError> {
        await request(.arithPow(base: base, exp: exp))
    }

    /// Sum sums two digits and returns error with error code as result and IP from context.
    func arithSum(a: Int,b: Int) async -> Result<Bool, RpcError> {
        await request(.arithSum(a: a, b: b))
    }

    /// SumArray returns sum all items from array
    func arithSumArray(array: [Double]? = nil) async -> Result<Double, RpcError> {
        await request(.arithSumArray(array: array))
    }
}


protocol CatalogueNetworking {
    func catalogueFirst(groups: [Group]) async -> Result<Bool, RpcError>
    func catalogueSecond(campaigns: [Campaign]) async -> Result<Bool, RpcError>
    func catalogueThird() async -> Result<Campaign, RpcError>
}

extension Networking: CatalogueNetworking {
    func catalogueFirst(groups: [Group]) async -> Result<Bool, RpcError> {
        await request(.catalogueFirst(groups: groups))
    }

    func catalogueSecond(campaigns: [Campaign]) async -> Result<Bool, RpcError> {
        await request(.catalogueSecond(campaigns: campaigns))
    }

    func catalogueThird() async -> Result<Campaign, RpcError> {
        await request(.catalogueThird())
    }
}
