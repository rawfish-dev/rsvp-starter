import fetch from 'isomorphic-fetch'

const SET_RSVP_MODE = 'SET_RSVP_MODE'

const SET_GUEST_RSVP = 'SET_GUEST_RSVP'

const SET_GUEST_RSVP_CREATED = 'SET_GUEST_RSVP_CREATED'

import {
  GENERIC_SERVER_ERROR,
  flashGuestOperationFailure
} from './general'

function setRSVPMode(mode) {
	return {
		type: SET_RSVP_MODE,
		mode
	}
}

/* Fetch */

function setGuestRSVP(rsvp) {
  return {
    type: SET_GUEST_RSVP,
    rsvp
  }
}

const guestRSVPDefaultValues = {
    ableToAttend: true,
    numberAttending: 1,
    specialDietaryRequirements: false
}

function fetchGuestRSVP(id) {
  if (!id) {
    return dispatch => {
      dispatch(setGuestRSVP(guestRSVPDefaultValues))

      return Promise.resolve()
    }
  }

	let request = {
		method: 'GET',
		headers: { 
			'Content-Type':'application/json' 
		}
	}

	return dispatch => {
		return fetch(`/api/p_rsvps/${id}`, request)
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
		return fetch('/api/p_rsvps', request)
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
  SET_RSVP_MODE,
  SET_GUEST_RSVP,
  setRSVPMode,
  fetchGuestRSVP,
  submitGuestRSVPCreate
}
