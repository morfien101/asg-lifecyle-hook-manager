package hookmanager

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/morfien101/asg-lifecyle-hook-manager/ec2metadatareader"
)

func awsSession() (*session.Session, error) {
	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		guess, err := ec2metadatareader.Region()
		if err != nil {
			return nil, err
		}
		region = guess
	}
	return session.NewSession(&aws.Config{Region: aws.String(region)})
}
func asgSession(session *session.Session) *autoscaling.AutoScaling {
	return autoscaling.New(session)
}
func newAWSSession() (*autoscaling.AutoScaling, error) {
	basicSession, err := awsSession()
	if err != nil {
		return nil, err
	}
	asgSession := asgSession(basicSession)
	return asgSession, nil
}

// SetContinue will send a continue request to the AutoScaling group
func SetContinue(asgName, hookName, instanceID string) (string, error) {
	clientSession, err := newAWSSession()
	if err != nil {
		return "", err
	}
	action := "CONTINUE"
	return submitHookAction(asgName, action, hookName, instanceID, clientSession)
}

// SetAbandon will send an abandon request to the AutoScaling group
func SetAbandon(asgName, hookName, instanceID string) (string, error) {
	clientSession, err := newAWSSession()
	if err != nil {
		return "", err
	}
	action := "ABANDON"
	return submitHookAction(asgName, action, hookName, instanceID, clientSession)
}

func submitHookAction(asgName, action, hookName, instanceID string, clientSession *autoscaling.AutoScaling) (string, error) {
	input := &autoscaling.CompleteLifecycleActionInput{
		AutoScalingGroupName:  aws.String(asgName),
		LifecycleActionResult: aws.String(action),
		LifecycleHookName:     aws.String(hookName),
		InstanceId:            aws.String(instanceID),
	}

	result, err := clientSession.CompleteLifecycleAction(input)
	if err != nil {
		returnErr := ""
		if awsError, ok := err.(awserr.Error); ok {
			switch awsError.Code() {
			case autoscaling.ErrCodeResourceContentionFault:
				returnErr = fmt.Sprintln(autoscaling.ErrCodeResourceContentionFault, awsError.Error())
			default:
				returnErr = fmt.Sprintln(awsError.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			returnErr = fmt.Sprintln(err.Error())
		}
		return "", fmt.Errorf("%s", returnErr)
	}
	return result.String(), nil
}

// RecordHeartBeat will reset the timeout for the Autoscaling Lifecycle Hook specified
func RecordHeartBeat(asgName, hookName, instanceID string) (string, error) {
	clientSession, err := newAWSSession()
	if err != nil {
		return "", err
	}

	input := &autoscaling.RecordLifecycleActionHeartbeatInput{
		AutoScalingGroupName: aws.String(asgName),
		InstanceId:           aws.String(instanceID),
		LifecycleHookName:    aws.String(hookName),
	}

	result, err := clientSession.RecordLifecycleActionHeartbeat(input)
	if err != nil {
		returnErr := ""
		if awsError, ok := err.(awserr.Error); ok {
			switch awsError.Code() {
			case autoscaling.ErrCodeResourceContentionFault:
				returnErr = fmt.Sprintln(autoscaling.ErrCodeResourceContentionFault, awsError.Error())
			default:
				returnErr = fmt.Sprintln(awsError.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			returnErr = fmt.Sprintln(err.Error())
		}
		return "", fmt.Errorf("%s", returnErr)
	}

	return result.String(), nil
}
