package testdata

import (
	"context"
	"errors"
	"math"

	"github.com/vmkteam/zenrpc/v2"
)

//go:generate zenrpc

type ArithService struct{ zenrpc.Service }

type Point struct {
	X, Y            int       // coordinate
	Z               int       `json:"-"`
	CartId          string    `json:"cartOd"`          // version string id - 1
	ID              int       `json:"id"`              // version id - 1
	BaseID          int64     `json:"baseId"`          // version id - 2
	SecondID        int32     `json:"secondID"`        // version id - 3
	SecondIDs       []int     `json:"secondIDs"`       // version multi id
	CartIDs         []string  `json:"cartIDs"`         // version multi string id
	CreatedAt       string    `json:"createdAt"`       // version date - 1
	UpdatedAt       string    `json:"updatedAt"`       // version date - 2
	ManualChangedAt string    `json:"manualChangedAt"` // version date - 3
	StartAt         string    `json:"startAt"`         // version date - 4
	NewLat          float64   `json:"newLat"`          // version group geo coordinate № - 1
	NewLon          float64   `json:"newLon"`          // version group geo coordinate № - 1
	Lat             float64   `json:"lat"`             // version group geo coordinate № - 2
	Lon             float64   `json:"lon"`             // version group geo coordinate № - 2
	Latitude        float64   `json:"latitude"`        // version group geo coordinate № - 3
	Longitude       float64   `json:"longitude"`       // version group geo coordinate № - 3
	BaseFloat       float64   `json:"baseFloat"`       // version group float - 1
	SecondFloat     float32   `json:"secondFloat"`     // version group float - 2
	EmptyString     *string   `json:"emptyString"`
	Name            string    `json:"name"`
	SecondPoints    []Point   `json:"secondPoints"`
	NextQuotient    Quotient  `json:"nextQuotient"`
	SecondQuotient  *Quotient `json:"secondQuotient"`
	Class           string    `json:"class"`
}

type SecondPoint struct {
}

func (as ArithService) GetByID(ctx context.Context, cartId string, categoryId int, baseID int, id int) (*Point, error) {
	return &Point{}, nil
}

func (as ArithService) GetSecondIds(ctx context.Context) []int {
	return nil
}

func (as ArithService) GetByLatLong(ctx context.Context, categoryId int, baseID int, lat, lon float64) (*Point, error) {
	return &Point{}, nil
}

func (as ArithService) GetByTime(ctx context.Context, createdAt string, updateAt string, startAt string) (*Point, error) {
	return &Point{}, nil
}

// Sum sums two digits and returns error with error code as result and IP from context.
func (as ArithService) Sum(ctx context.Context, a, b int) (bool, *zenrpc.Error) {
	r, _ := zenrpc.RequestFromContext(ctx)

	return true, zenrpc.NewStringError(a+b, r.Host)
}

func (as ArithService) Positive() (bool, *zenrpc.Error) {
	return true, nil
}

func (ArithService) DoSomething() {
	// some optimistic operations
}
func (ArithService) DoSomethingV2() ExternalData {
	// some optimistic operations
	return ExternalData{}
}

func (ArithService) GetPoints() []Point {
	return []Point{}
}

//zenrpc:return Point test description in return
func (ArithService) DoSomethingWithPoint(p Point, pp []Point) Point {
	// some optimistic operations
	return p
}

// Multiply multiples two digits and returns result.
func (as ArithService) Multiply(a, b int) int {
	return a * b
}

// CheckError throws error is isErr true.
// TEST row 2
//
//zenrpc:500 test error
func (ArithService) CheckError(isErr bool) error {
	if isErr {
		return errors.New("test")
	}

	return nil
}

// CheckZenRPCError throws zenrpc error is isErr true.
// Second description row
//
//zenrpc:500 test error
func (ArithService) CheckZenRPCError(isErr bool) *zenrpc.Error {
	if isErr {
		return zenrpc.NewStringError(500, "test")
	}

	return nil
}

// Quotient docs
type Quotient struct {
	// Quo docs
	Quo int

	// Rem docs
	Rem     int             `json:"rem"`
	BaseRow string          `json:"baseRow"`
	RowNil  *string         `json:"rowNil"`
	Data    CycleInitStruct `json:"data"`
}

type CycleInitStruct struct {
	IsCycleInit bool `json:"isCycleInit"`
}
type ExternalData struct {
	Name string `json:"name"`
}

// Divide divides two numbers.
//
//zenrpc:a			the a
//zenrpc:b 			the b
//zenrpc:quo		result is Quotient, should be named var
//zenrpc:401 		we do not serve 1
func (as *ArithService) Divide(a, b int) (quo *Quotient, err error) {
	if b == 0 {
		return nil, errors.New("divide by zero")
	} else if b == 1 {
		return nil, zenrpc.NewError(401, errors.New("we do not serve 1"))
	}

	return &Quotient{
		Quo: a / b,
		Rem: a % b,
	}, nil
}

// Pow returns x**y, the base-x exponential of y. If Exp is not set then default value is 2.
//
//zenrpc:exp=2 	exponent could be empty
func (as *ArithService) Pow(base float64, exp *float64) float64 {
	return math.Pow(base, *exp)
}

// Pi returns math.Pi.
func (ArithService) Pi() float64 {
	return math.Pi
}

// SumArray returns sum all items from array
//
//zenrpc:array=[]float64{1,2,4}
func (as *ArithService) SumArray(array *[]float64) float64 {
	var sum float64

	for _, i := range *array {
		sum += i
	}
	return sum
}
