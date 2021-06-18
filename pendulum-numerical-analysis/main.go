package main

import (
	"fmt"
	"math"
	"math/big"

	"gonum.org/v1/plot/plotter"
)

func main() {
	gravity := bigfloatCreate(9.8)
	length := bigfloatCreate(1)

	percentage := bigfloatCreate(2)
	degree := getDegreeLimitForDifferencePercentage(
		gravity,
		length,
		bigfloatCreate(0.000001),
		bigfloatCreate(0.00001),
		percentage,
	)
	fmt.Printf(
		"Maximal degree to be under or equal to %2.5f%% is about %v deg\n",
		percentage,
		degree,
	)
	differencePercentageS := getDifferencePercentageS(
		gravity,
		length,
		bigfloatCreate(1),
		bigfloatCreate(0.0001),
	)
	plotWrite(differencePercentageS, "difference.png")
}

//getDegreeLimitForDifferencePercentage return the maximal degree to not be
//over the given difference percentage between actual and simple pendulum's period
//percentage ∈ [0, inf)
//delta should be small enough that delta's difference percentage is under or equal to percentage
func getDegreeLimitForDifferencePercentage(
	gravity,
	length,
	degreeDelta,
	timeDelta,
	percentage *big.Float,
) *big.Float {
	//at 0 deg both have Inf as period so they have 0% difference which is also the minimum value
	degreeUnderOrEqual := bigfloatCreate(0)
	//the differnce percentage at 90 deg might be the percentage we want
	degreeOver := bigfloatCreate(0).Add(bigfloatCreate(90), degreeDelta)
	degreeCurrent := bigfloatCreate(0).Add(degreeUnderOrEqual, degreeOver)
	degreeCurrent.Quo(degreeCurrent, bigfloatCreate(2))

	//binary search exact value until in error range
	for bigfloatCreate(0).Add(degreeUnderOrEqual, degreeDelta).Cmp(degreeOver) == -1 {
		period := periodGet(gravity, length, degreeCurrent, timeDelta, bigfloatCreate(0))
		differencePercentage := getDifferencePercentage(gravity, length, period)

		//half the search filed
		if differencePercentage.Cmp(percentage) != 1 {
			degreeUnderOrEqual.Copy(degreeCurrent)
		} else {
			degreeOver.Copy(degreeCurrent)
		}

		//new halfing point
		degreeCurrent = degreeCurrent.Add(degreeUnderOrEqual, degreeOver)
		degreeCurrent.Quo(degreeCurrent, bigfloatCreate(2))
	}

	return degreeCurrent
}

//getDifferencePercentageS returns difference percentage between simple and actual pendulum's period for specified degrees:
//X ∈ [degreeDelta, 90]
//Y ∈ [0, Inf)
func getDifferencePercentageS(gravity, length, degreeDelta, timeDelta *big.Float) plotter.XYs {
	//do not set degreeStart to 0 as plot can not draw infinite value
	degreeStart := degreeDelta
	degreeEnd := bigfloatCreate(90)

	var differencePercentageS plotter.XYs

	i := 0
	for {
		//current degree
		//calculate each time to get exact value
		degree := bigfloatCreate(0).Mul(bigfloatCreate(float64(i)), degreeDelta)
		degree.Add(degree, degreeStart)

		//degree may not be equal to degreeEnd either because
		//degreeDelta's value or
		//float precision
		if degree.Cmp(degreeEnd) == 1 {
			break
		}

		period := periodGet(gravity, length, degree, timeDelta, bigfloatCreate(0))
		percentageDifference := getDifferencePercentage(gravity, length, period)

		//save data for plot
		degree64, _ := degree.Float64()
		percentageDifference64, _ := percentageDifference.Float64()
		differencePercentageS = append(differencePercentageS, plotter.XY{
			X: degree64,
			Y: percentageDifference64,
		})

		i++
	}

	return differencePercentageS
}

//getDifferencePercentage returns difference percentage between simple and actual pendulum's period
func getDifferencePercentage(gravity, length, periodActual *big.Float) *big.Float {
	//simple period = 2π√(length/gravity)
	periodSimple := bigfloatCreate(0).Quo(length, gravity)
	periodSimple.Sqrt(periodSimple)
	//Multiply separately to keep pi's high precision
	periodSimple.Mul(periodSimple, bigfloatCreate(math.Pi))
	periodSimple.Mul(periodSimple, bigfloatCreate(2))

	//period percentage difference
	percentageDifference := bigfloatCreate(0).Quo(periodActual, periodSimple)
	percentageDifference.Sub(percentageDifference, bigfloatCreate(1))
	percentageDifference.Mul(percentageDifference, bigfloatCreate(100))

	return percentageDifference
}

//periodGet returns the period based on parameters by numerical analysis using euler's algorithm
func periodGet(gravity, length, degreeStart, timeDelta, angularVelocity *big.Float) *big.Float {
	if degreeStart.Cmp(bigfloatCreate(0)) == 0 {
		return bigfloatCreate(math.Inf(1))
	}

	//constant = gravity / length
	constant := bigfloatCreate(0).Quo(gravity, length)
	constant.Neg(constant)

	//start values
	alphaStart := degreeToRadian(degreeStart)
	alpha := bigfloatCreate(0).Copy(alphaStart)

	i := 0
	for {
		//current time
		//calculate each time to get exact value
		time := bigfloatCreate(0).Mul(bigfloatCreate(float64(i)), timeDelta)

		//period check
		if alpha.Cmp(alphaStart) != -1 && time.Cmp(bigfloatCreate(0)) != 0 {
			return time
		}

		//f'' calc
		alpha64, _ := alpha.Float64()
		angularAcceleration := bigfloatCreate(0).Mul(constant, bigfloatCreate(math.Sin(alpha64)))

		//f' = f' + f''dt
		angularVelocityDelta := bigfloatCreate(0).Mul(angularAcceleration, timeDelta)
		angularVelocity.Add(angularVelocity, angularVelocityDelta)

		//f = f + f'dt
		alphaDelta := bigfloatCreate(0).Mul(angularVelocity, timeDelta)
		alpha.Add(alpha, alphaDelta)

		i++
	}
}
