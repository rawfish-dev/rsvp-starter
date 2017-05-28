const SET_OPERATION_SUCCESS = 'SET_OPERATION_SUCCESS'
const SET_OPERATION_FAILURE = 'SET_OPERATION_FAILURE'
const UNSET_OPERATION_RESULT = 'UNSET_OPERATION_RESULT'

const GENERIC_SERVER_ERROR = 'Server error encountered, please try again later.'
const INVALID_SESSION_ERROR = 'Your session has expired. Please login again.'

function setOperationSuccess(message) {
	return {
		type: SET_OPERATION_SUCCESS,
		message
	}
}

function setOperationFailure(message) {
	return {
		type: SET_OPERATION_FAILURE,
		message
	}
}

function unsetOperationResult() {
	return {
		type: UNSET_OPERATION_RESULT
	}
}

function flashOperationResult(message, success) {
  return dispatch => {
    if (success) {
      dispatch(setOperationSuccess(message))
    } else {
      dispatch(setOperationFailure(message))
    }

    // Remove message after some time
    setTimeout(() => {
      dispatch(unsetOperationResult())
    }, 10000)

    return Promise.resolve()
  }
}

function flashGuestOperationFailure(message) {
  return dispatch => {
    dispatch(setOperationFailure(message))

    // Remove message after some time
    setTimeout(() => {
      dispatch(unsetOperationResult())
    }, 10000)

    return Promise.resolve()
  }
}
 
module.exports = {
    SET_OPERATION_SUCCESS,
    SET_OPERATION_FAILURE,
    UNSET_OPERATION_RESULT,
    GENERIC_SERVER_ERROR,
    INVALID_SESSION_ERROR,
    flashOperationResult,
    flashGuestOperationFailure
}
