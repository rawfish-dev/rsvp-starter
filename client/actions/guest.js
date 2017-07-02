import fetch from 'isomorphic-fetch'

const SET_GUEST_RSVP = 'SET_GUEST_RSVP'
const SET_GUEST_RSVP_CREATED = 'SET_GUEST_RSVP_CREATED'

import {
  GENERIC_SERVER_ERROR,
  flashGuestOperationFailure
} from './general'

/* Fetch */

function setGuestRSVP(rsvp) {
  return {
    type: SET_GUEST_RSVP,
    rsvp
  }
}

function fetchRSVP(id) {
	let request = {
		method: 'GET',
		headers: { 
			'Content-Type':'application/json' 
		}
	}

	return dispatch => {
		return fetch(`/api/rsvps/${id}`, request)
		.then(rawResponse => {
			if (!rawResponse.ok) {
        dispatch(flashGuestOperationFailure(GENERIC_SERVER_ERROR))

				return Promise.reject()
			}

			return rawResponse.json()
		}).then(response =>  {
      // Update guest rsvp
      dispatch(setGuestRSVP(response))

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("fetch guest rsvp error", err)
			}
		})
	}
}

/* Create */

function submitGuestRSVPCreate(rsvp) {
	let request = {
		method: 'POST',
		headers: { 
			'Content-Type':'application/json' 
		},
		body: JSON.stringify(rsvp)
	}

	return dispatch => {
		return fetch('/api/rsvps', request)
		.then(rawResponse => {
			if (!rawResponse.ok) {
				dispatch(flashGuestOperationFailure(GENERIC_SERVER_ERROR))

				return Promise.reject()
			}

			return rawResponse.json()
		}).then(response =>  {
      // Update rsvps
      dispatch(setGuestRSVP(response))

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("create guest rsvp error", err)
			}
		})
	}
}

module.exports = {
  SET_GUEST_RSVP,
  fetchRSVP,
  submitGuestRSVPCreate
}
