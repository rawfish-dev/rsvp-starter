import fetch from 'isomorphic-fetch'

const TOGGLE_RSVP_FORM_VISIBILITY = 'TOGGLE_RSVP_FORM_VISIBILITY'

const SET_RSVPS = 'SET_RSVPS'

const SET_RSVP_UPDATED = 'SET_RSVP_UPDATED'
const EDIT_RSVP_SUCCESS_MESSAGE = 'RSVP was updated successfully.'

const SET_RSVP_CREATED = 'SET_RSVP_CREATED'
const CREATE_RSVP_SUCCESS_MESSAGE = 'RSVP was created successfully.'

const TOGGLE_RSVP_DELETE_CONFIRMATION = 'TOGGLE_RSVP_DELETE_CONFIRMATION'
const SET_RSVP_DELETED = 'SET_RSVP_DELETED'
const DELETE_RSVP_SUCCESS_MESSAGE = 'RSVP was deleted successfully.'

import {
  INVALID_SESSION_ERROR,
  GENERIC_SERVER_ERROR,
  flashOperationResult
} from './general'

import {
	logoutUser
} from './logout'

function toggleRSVPFormVisibility(mode, initialValues) {
  return {
    type: TOGGLE_RSVP_FORM_VISIBILITY,
    mode,
    initialValues
  }
}

/* List */

function setRSVPs(rsvps) {
	return {
		type: SET_RSVPS,
    rsvps
	}
}

function fetchRSVPs() {
	let request = {
		method: 'GET',
		headers: { 
			'Content-Type':'application/json',
			'X-Auth-Header': localStorage.getItem('authToken') 
		}
	}

	return dispatch => {
		return fetch('/api/rsvps', request)
		.then(rawResponse => {
			if (!rawResponse.ok) {
				switch(rawResponse.status) {
					case 401:
						dispatch(flashOperationResult(INVALID_SESSION_ERROR, false))
						dispatch(logoutUser())
						break
          default:
            dispatch(flashOperationResult(GENERIC_SERVER_ERROR, false))
				}

				return Promise.reject()
			}

			return rawResponse.json()
		}).then(response =>  {
      // Update rsvps
      dispatch(setRSVPs(response))

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("fetch rsvps error", err)
			}
		})
	}
}

/* Edit */

function setRSVPUpdated(rsvp) {
  return {
    type: SET_RSVP_UPDATED,
    rsvp
  }
}

function submitRSVPEdit(rsvp) {
	let request = {
		method: 'PUT',
		headers: { 
			'Content-Type':'application/json',
			'X-Auth-Header': localStorage.getItem('authToken') 
		},
		body: JSON.stringify(rsvp)
	}

	return dispatch => {
		return fetch(`/api/rsvps/${rsvp.id}`, request)
		.then(rawResponse => {
			if (!rawResponse.ok) {
        switch(rawResponse.status) {
          case 400:
            dispatch(flashOperationResult(rawResponse.error, false))
            break
          case 401:
            dispatch(flashOperationResult(INVALID_SESSION_ERROR, false))
            break
          default:
            dispatch(flashOperationResult(GENERIC_SERVER_ERROR, false))
        }

				return Promise.reject()
			}

			return rawResponse.json()
		}).then(response =>  {
      dispatch(toggleRSVPFormVisibility(null, {}))

			// Dispatch the success action
			dispatch(flashOperationResult(EDIT_RSVP_SUCCESS_MESSAGE, true))

      // Update rsvp
      dispatch(setRSVPUpdated(response))

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("update rsvp error", err)
			}
		})
	}
}

/* Create */

function setRSVPCreated(rsvp) {
  return {
    type: SET_RSVP_CREATED,
    rsvp
  }
}

function submitRSVPCreate(rsvp) {
	let request = {
		method: 'POST',
		headers: { 
			'Content-Type':'application/json',
			'X-Auth-Header': localStorage.getItem('authToken') 
		},
		body: JSON.stringify(rsvp)
	}

	return dispatch => {
		return fetch('/api/rsvps', request)
		.then(rawResponse => {
			if (!rawResponse.ok) {
				switch(rawResponse.status) {
				case 400:
					dispatch(flashOperationResult(rawResponse.error))
					break
				case 401:
					dispatch(flashOperationResult(INVALID_SESSION_ERROR, false))
					break
				default:
					dispatch(flashOperationResult(GENERIC_SERVER_ERROR, false))
				}

				return Promise.reject()
			}

			return rawResponse.json()
		}).then(response =>  {
      dispatch(toggleRSVPFormVisibility(null, {}))

			// Dispatch the success action
			dispatch(flashOperationResult(CREATE_RSVP_SUCCESS_MESSAGE, true))

      // Update rsvps
      dispatch(setRSVPCreated(response))

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("create rsvp error", err)
			}
		})
	}
}

/* Delete */

function setRSVPDeleted(rsvpID) {
  return {
    type: SET_RSVP_DELETED,
    rsvpID
  }
}

function toggleRSVPDeleteConfirmation(rsvpID) {
  return {
    type: TOGGLE_RSVP_DELETE_CONFIRMATION,
    rsvpID
  }
}

function submitRSVPDelete(rsvpID) {
	let request = {
		method: 'DELETE',
		headers: { 
			'Content-Type':'application/json',
			'X-Auth-Header': localStorage.getItem('authToken') 
		}
	}

	return dispatch => {
		return fetch(`/api/rsvps/${rsvpID}`, request)
		.then(rawResponse => {
			if (!rawResponse.ok) {
        switch(rawResponse.status) {
          case 400:
            dispatch(flashOperationResult(rawResponse.error))
            break
          case 401:
            dispatch(flashOperationResult(INVALID_SESSION_ERROR, false))
            break
          default:
            dispatch(flashOperationResult(GENERIC_SERVER_ERROR, false))
        }

				return Promise.reject()
			}

      dispatch(toggleRSVPDeleteConfirmation(0))

			// Dispatch the success action
			dispatch(flashOperationResult(DELETE_RSVP_SUCCESS_MESSAGE, true))

      // Update rsvps
      dispatch(setRSVPDeleted(rsvpID))

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("delete rsvp error", err)
			}
		})
	}
}

module.exports = {
  TOGGLE_RSVP_FORM_VISIBILITY,
  SET_RSVPS,
  SET_RSVP_UPDATED,
  SET_RSVP_CREATED,
  SET_RSVP_DELETED,
	TOGGLE_RSVP_DELETE_CONFIRMATION,
  toggleRSVPFormVisibility,
  toggleRSVPDeleteConfirmation,
  fetchRSVPs,
  submitRSVPEdit,
  submitRSVPCreate,
  submitRSVPDelete
}
