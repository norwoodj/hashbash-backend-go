package dao

import "fmt"

type RainbowTableExistsError struct {
	Name string
}

func (err RainbowTableExistsError) Error() string {
	return fmt.Sprintf("Rainbow table with name %s already exists!", err.Name)
}

type RainbowTableNotExistsError struct {
	ID int16
}

func (err RainbowTableNotExistsError) Error() string {
	return fmt.Sprintf("No rainbow table with id %d already exists!", err.ID)
}

type InvalidHashError struct {
	Hash             string
	HashFunctionName string
}

func (err InvalidHashError) Error() string {
	return fmt.Sprintf("Invalid %s hash format provided: %s", err.HashFunctionName, err.Hash)
}

func IsInvalidHashError(err error) bool {
	if err != nil {
		switch err.(type) {
		case InvalidHashError:
			return true
		}
	}

	return false
}

func IsRainbowTableExistsError(err error) bool {
	if err != nil {
		switch err.(type) {
		case RainbowTableExistsError:
			return true
		}
	}

	return false
}

func IsRainbowTableNotExistsError(err error) bool {
	if err != nil {
		switch err.(type) {
		case RainbowTableNotExistsError:
			return true
		}
	}

	return false
}
