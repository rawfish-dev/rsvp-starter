import fetch from 'isomorphic-fetch'
import { browserHistory } from 'react-router'

import {
  GENERIC_SERVER_ERROR,
  flashOperationResult
} from './general'

const LOGIN_REQUEST = 'LOGIN_REQUEST'
const LOGIN_SUCCESS = 'LOGIN_SUCCESS'
const LOGIN_FAILURE = 'LOGIN_FAILURE'

const INVALID_CREDENTIALS_ERROR = 'Oops, username and/or password were not correct.'

function requestLogin(credentials) {
  return {
    type: LOGIN_REQUEST,
    isFetching: true,
    isAuthenticated: false,
    credentials
  }
}

function receiveLogin(user) {
  return {
    type: LOGIN_SUCCESS,
    isFetching: false,
    isAuthenticated: true,
    authToken: user.authToken
  }
}

function loginError(message) {
  return {
    type: LOGIN_FAILURE,
    isFetching: false,
    isAuthenticated: false,
    message
  }
}

function loginUser(credentials) {
	let request = {
		method: 'POST',
		headers: { 'Content-Type':'application/json' },
		body: JSON.stringify(credentials)
	}

	return dispatch => {
		// We dispatch requestLogin to kickoff the call to the API
		dispatch(requestLogin(credentials))

		return fetch('/api/sessions', request)
		.then(rawResponse => {
			if (!rawResponse.ok) {
        switch(rawResponse.status) {
          case 401:
            dispatch(loginError(INVALID_CREDENTIALS_ERROR))
            break
          default:
            dispatch(flashOperationResult(GENERIC_SERVER_ERROR, false))
        }

				return Promise.reject()
			}

			return rawResponse.json()
		}).then(response =>  {
			if (!response.authToken) {
        dispatch(flashOperationResult(GENERIC_SERVER_ERROR, false))

				return Promise.reject()
			}

			// If login was successful, set the token in local storage
			localStorage.setItem('authToken', response.authToken)
			localStorage.setItem('username', response.username)

			// Dispatch the success action
			dispatch(receiveLogin(response))

			browserHistory.push('/control_panel')

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("login error", err)
			}
		})
	}
}

module.exports = {
    LOGIN_REQUEST,
    LOGIN_SUCCESS,
    LOGIN_FAILURE,
    INVALID_CREDENTIALS_ERROR,
    loginUser
}
