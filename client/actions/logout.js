import fetch from 'isomorphic-fetch'
import { browserHistory } from 'react-router'

const LOGOUT_REQUEST = 'LOGOUT_REQUEST'
const LOGOUT_SUCCESS = 'LOGOUT_SUCCESS'
const LOGOUT_FAILURE = 'LOGOUT_FAILURE'

function requestLogout() {
  return {
    type: LOGOUT_REQUEST,
    isFetching: true,
    isAuthenticated: true
  }
}

function receiveLogout() {
  return {
    type: LOGOUT_SUCCESS,
    isFetching: false,
    isAuthenticated: false
  }
}

function logoutUser() {
	let request = {
		method: 'DELETE',
		headers: {
			'X-Auth-Header': localStorage.getItem('authToken')
		}
	}

	return dispatch => {
		dispatch(requestLogout())

		return fetch('/api/sessions', request)
      .then(rawResponse => {
        if (!rawResponse.ok) {
          // Just warn but allow logout anyway
          console.warn("logout error")
        }

        localStorage.removeItem('authToken')
        localStorage.removeItem('username')

        dispatch(receiveLogout())

        browserHistory.push('/login')

        return Promise.resolve()
      }).catch(err => {
        console.warn("logout error", err)
      })
    }
  }

module.exports = {
    LOGOUT_REQUEST,
    LOGOUT_SUCCESS,
    LOGOUT_FAILURE,
    logoutUser
}
