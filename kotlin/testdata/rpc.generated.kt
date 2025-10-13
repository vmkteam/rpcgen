/// Code generated from jsonrpc schema by rpcgen v2.4.6; DO NOT EDIT.
package api.model
import java.time.LocalTime
import java.time.ZonedDateTime


data class Point(
    /**
     * coordinate
     */
    val X: Int = 0,
    /**
     * coordinate
     */
    val Y: Int = 0,
    /**
     * version group float - 1
     */
    val baseFloat: Float = 0f,
    /**
     * version id - 2
     */
    val baseId: Long = 0,
    /**
     * version date - 1
     */
    val createdAt: ZonedDateTime = ZonedDateTime.now(),
    val emptyString: String? = null,
    /**
     * version id - 1
     */
    val id: Long = 0,
    /**
     * version group geo coordinate № - 2
     */
    val lat: Float = 0f,
    /**
     * version group geo coordinate № - 3
     */
    val latitude: Float = 0f,
    /**
     * version group geo coordinate № - 2
     */
    val lon: Float = 0f,
    /**
     * version group geo coordinate № - 3
     */
    val longitude: Float = 0f,
    /**
     * version date - 3
     */
    val manualChangedAt: ZonedDateTime = ZonedDateTime.now(),
    val name: String = "",
    /**
     * version group geo coordinate № - 1
     */
    val newLat: Double = 0.0,
    /**
     * version group geo coordinate № - 1
     */
    val newLon: Double = 0.0,
    /**
     * version group float - 2
     */
    val secondFloat: Float = 0f,
    /**
     * version id - 3
     */
    val secondID: Long = 0,
    val secondPoints: List<Point> = emptyList(),
    /**
     * version date - 2
     */
    val updatedAt: ZonedDateTime = ZonedDateTime.now(),
)

data class Quotient(
    /**
     * Quo docs
     */
    val Quo: Int,
    val baseRow: String,
    /**
     * Rem docs
     */
    val rem: Int,
    val rowNil: String?,
)

