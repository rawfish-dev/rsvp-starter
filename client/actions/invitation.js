import fetch from 'isomorphic-fetch'

const TOGGLE_INVITATION_FORM_VISIBILITY = 'TOGGLE_INVITATION_FORM_VISIBILITY'

const SET_INVITATIONS = 'SET_INVITATIONS'

const SET_INVITATION_UPDATED = 'SET_INVITATION_UPDATED'
const EDIT_INVITATION_SUCCESS_MESSAGE = 'Invitation was updated successfully.'

const SET_INVITATION_CREATED = 'SET_INVITATION_CREATED'
const CREATE_INVITATION_SUCCESS_MESSAGE = 'Invitation was created successfully.'

const TOGGLE_INVITATION_DELETE_CONFIRMATION = 'TOGGLE_INVITATION_DELETE_CONFIRMATION'
const SET_INVITATION_DELETED = 'SET_INVITATION_DELETED'
const DELETE_INVITATION_SUCCESS_MESSAGE = 'Invitation was deleted successfully.'

import {
  INVALID_SESSION_ERROR,
  GENERIC_SERVER_ERROR,
  flashOperationResult
} from './general'

import {
	logoutUser
} from './logout'

function toggleInvitationFormVisibility(mode, initialValues) {
  return {
    type: TOGGLE_INVITATION_FORM_VISIBILITY,
    mode,
    initialValues
  }
}

/* List */

function setInvitations(invitations) {
	return {
		type: SET_INVITATIONS,
    invitations
	}
}

function fetchInvitations() {
	let request = {
		method: 'GET',
		headers: { 
			'Content-Type':'application/json',
			'X-Auth-Header': localStorage.getItem('authToken') 
		}
	}

	return dispatch => {
		return fetch('/api/invitations', request)
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
      // Update invitations
      dispatch(setInvitations(response))

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("fetch invitations error", err)
			}
		})
	}
}

/* Edit */

function setInvitationUpdated(invitation) {
  return {
    type: SET_INVITATION_UPDATED,
    invitation
  }
}

function submitInvitationEdit(invitation) {
	let request = {
		method: 'PUT',
		headers: { 
			'Content-Type':'application/json',
			'X-Auth-Header': localStorage.getItem('authToken') 
		},
		body: JSON.stringify(invitation)
	}

	return dispatch => {
		return fetch(`/api/invitations/${invitation.id}`, request)
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
      dispatch(toggleInvitationFormVisibility(null, {}))

			// Dispatch the success action
			dispatch(flashOperationResult(EDIT_INVITATION_SUCCESS_MESSAGE, true))

      // Update invitation
      dispatch(setInvitationUpdated(response))

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("update invitation error", err)
			}
		})
	}
}

/* Create */

function setInvitationCreated(invitation) {
  return {
    type: SET_INVITATION_CREATED,
    invitation
  }
}

function submitInvitationCreate(invitation) {
	let request = {
		method: 'POST',
		headers: { 
			'Content-Type':'application/json',
			'X-Auth-Header': localStorage.getItem('authToken') 
		},
		body: JSON.stringify(invitation)
	}

	return dispatch => {
		return fetch('/api/invitations', request)
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
      dispatch(toggleInvitationFormVisibility(null, {}))

			// Dispatch the success action
			dispatch(flashOperationResult(CREATE_INVITATION_SUCCESS_MESSAGE, true))

      // Update invitations
      dispatch(setInvitationCreated(response))

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("create invitation error", err)
			}
		})
	}
}

/* Delete */

function setInvitationDeleted(invitationID) {
  return {
    type: SET_INVITATION_DELETED,
    invitationID
  }
}

function toggleInvitationDeleteConfirmation(invitationID) {
  return {
    type: TOGGLE_INVITATION_DELETE_CONFIRMATION,
    invitationID
  }
}

function submitInvitationDelete(invitationID) {
	let request = {
		method: 'DELETE',
		headers: { 
			'Content-Type':'application/json',
			'X-Auth-Header': localStorage.getItem('authToken') 
		}
	}

	return dispatch => {
		return fetch(`/api/invitations/${invitationID}`, request)
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

      dispatch(toggleInvitationDeleteConfirmation(0))

			// Dispatch the success action
			dispatch(flashOperationResult(DELETE_INVITATION_SUCCESS_MESSAGE, true))

      // Update invitations
      dispatch(setInvitationDeleted(invitationID))

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("delete invitation error", err)
			}
		})
	}
}

module.exports = {
  TOGGLE_INVITATION_FORM_VISIBILITY,
  SET_INVITATIONS,
  SET_INVITATION_UPDATED,
  SET_INVITATION_CREATED,
	SET_INVITATION_DELETED,
	TOGGLE_INVITATION_DELETE_CONFIRMATION,
  toggleInvitationFormVisibility,
	toggleInvitationDeleteConfirmation,
  fetchInvitations,
  submitInvitationEdit,
  submitInvitationCreate,
	submitInvitationDelete
}
