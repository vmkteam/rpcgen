/// Code generated from jsonrpc schema by rpcgen v2.7.0; DO NOT EDIT.
package api

import com.google.gson.reflect.TypeToken
import java.time.ZonedDateTime
import java.time.LocalTime
import api.model.*
import api.debug.model.*
import api.TransportOption
import api.Transport

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
     * @return 
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
     * Second description row
     *
     * Коды ошибок:
     *
     *    "500": "test error",
     *
     *
     * @return 
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
     * @param a
     * @param b
     * @return 
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

    fun arithDoSomethingV2(
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<ExternalData>>() {},
        "arith.DoSomethingV2",
    )

    /**
     * @return 
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

    /**
     * @return 
     */
    fun arithGetByID(
        cartId: String,
        categoryId: Long,
        baseID: Long,
        id: Long,
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<Point>>() {},
        "arith.GetByID",
        "cartId" to cartId,
        "categoryId" to categoryId,
        "baseID" to baseID,
        "id" to id,
    )

    /**
     * @return 
     */
    fun arithGetByLatLong(
        categoryId: Long,
        baseID: Long,
        lat: Double,
        lon: Double,
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<Point>>() {},
        "arith.GetByLatLong",
        "categoryId" to categoryId,
        "baseID" to baseID,
        "lat" to lat,
        "lon" to lon,
    )

    /**
     * @return 
     */
    fun arithGetByTime(
        createdAt: ZonedDateTime,
        updateAt: ZonedDateTime,
        startAt: ZonedDateTime,
        vararg transportOptions: TransportOption,
    ) = request(
        transportOptions,
        object : TypeToken<ApiResponse<Point>>() {},
        "arith.GetByTime",
        "createdAt" to createdAt,
        "updateAt" to updateAt,
        "startAt" to startAt,
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
     * @return 
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
     * @return 
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
     * @param exp
     * @return 
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
     * @return 
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
     * @return 
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
