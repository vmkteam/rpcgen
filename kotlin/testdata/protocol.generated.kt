/// Code generated from jsonrpc schema by rpcgen v2.4.6; DO NOT EDIT.
package api

import com.google.gson.reflect.TypeToken
import api.model.*
import java.time.ZonedDateTime

interface Api : Transport {

    /**
     * CheckError throws error is isErr true.
     * TEST row 2
     *
     * Коды ошибок:
     *
     *    "500": "test error",
     *
     *
     */
    fun arithCheckError(
        isErr: Boolean,
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<Nothing>>() {},
        "arith.CheckError",
        "isErr" to isErr,
    )

    /**
     * CheckZenRPCError throws zenrpc error is isErr true.
     * Вторая строка описания
     *
     * Коды ошибок:
     *
     *    "500": "test error",
     *
     *
     */
    fun arithCheckZenRPCError(
        isErr: Boolean,
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<Nothing>>() {},
        "arith.CheckZenRPCError",
        "isErr" to isErr,
    )

    /**
     * Divide divides two numbers.
     *
     * Коды ошибок:
     *
     *    "401": "we do not serve 1",
     *
     *
     * @param a the a
     * @param b the b
     */
    fun arithDivide(
        a: Int,
        b: Int,
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<Quotient>>() {},
        "arith.Divide",
        "a" to a,
        "b" to b,
    )

    fun arithDoSomething(
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<Nothing>>() {},
        "arith.DoSomething",
    )

    /**
     * @return Point test description in return
     */
    fun arithDoSomethingWithPoint(
        p: Point,
        pp: List<Point>,
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<Point>>() {},
        "arith.DoSomethingWithPoint",
        "p" to p,
        "pp" to pp,
    )

    fun arithGetPoints(
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<List<Point>>>() {},
        "arith.GetPoints",
    )

    /**
     * Multiply multiples two digits and returns result.
     */
    fun arithMultiply(
        a: Int,
        b: Int,
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<Int>>() {},
        "arith.Multiply",
        "a" to a,
        "b" to b,
    )

    /**
     * Pi returns math.Pi.
     */
    fun arithPi(
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<Float>>() {},
        "arith.Pi",
    )

    fun arithPositive(
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<Boolean>>() {},
        "arith.Positive",
    )

    /**
     * Pow returns x**y, the base-x exponential of y. If Exp is not set then default value is 2.
     * @param exp exponent could be empty
     */
    fun arithPow(
        base: Float,
        exp: Float?,
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<Float>>() {},
        "arith.Pow",
        "base" to base,
        "exp" to exp,
    )

    /**
     * Sum sums two digits and returns error with error code as result and IP from context.
     */
    fun arithSum(
        a: Int,
        b: Int,
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<Boolean>>() {},
        "arith.Sum",
        "a" to a,
        "b" to b,
    )

    /**
     * SumArray returns sum all items from array
     */
    fun arithSumArray(
        array: List<Float>?,
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<Float>>() {},
        "arith.SumArray",
        "array" to array,
    )
}
