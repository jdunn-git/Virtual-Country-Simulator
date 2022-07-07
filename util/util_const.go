package util

const (
	TransformType string = "Transform"
	TransferType         = "Transfer"

	FailureCost float64 = -0.05

	// Gamma is a rate of change for the Discounted reward at each step.
	//  This is just an initial value
	Gamma float64 = 0.5

	// These values are used for the LogisticRegression calculation
	X0 float64 = 0.0
	K          = 1.0
	L          = 1.0
)
