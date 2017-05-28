import fetch from 'isomorphic-fetch'

const TOGGLE_CATEGORY_FORM_VISIBILITY = 'TOGGLE_CATEGORY_FORM_VISIBILITY'

const SET_CATEGORIES = 'SET_CATEGORIES'

const SET_CATEGORY_UPDATED = 'SET_CATEGORY_UPDATED'
const EDIT_CATEGORY_SUCCESS_MESSAGE = 'Category was updated successfully.'

const SET_CATEGORY_CREATED = 'SET_CATEGORY_CREATED'
const CREATE_CATEGORY_SUCCESS_MESSAGE = 'Category was created successfully.'

const TOGGLE_CATEGORY_DELETE_CONFIRMATION = 'TOGGLE_CATEGORY_DELETE_CONFIRMATION'
const SET_CATEGORY_DELETED = 'SET_CATEGORY_DELETED'
const DELETE_CATEGORY_SUCCESS_MESSAGE = 'Category was deleted successfully.'

import {
  INVALID_SESSION_ERROR,
  GENERIC_SERVER_ERROR,
  flashOperationResult
} from './general'

import {
	logoutUser
} from './logout'

function toggleCategoryFormVisibility(mode, initialValues) {
  return {
    type: TOGGLE_CATEGORY_FORM_VISIBILITY,
    mode,
    initialValues
  }
}

/* List */

function setCategories(categories) {
	return {
		type: SET_CATEGORIES,
    categories
	}
}

function fetchCategories() {
	let request = {
		method: 'GET',
		headers: { 
			'Content-Type':'application/json',
			'X-Auth-Header': localStorage.getItem('authToken') 
		}
	}

	return dispatch => {
		return fetch('/api/categories', request)
		.then(rawResponse => {
			if (!rawResponse.ok) {
				switch(rawResponse.status) {
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
      // Update categories
      dispatch(setCategories(response))

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("fetch categories error", err)
			}
		})
	}
}

/* Edit */

function setCategoryUpdated(category) {
  return {
    type: SET_CATEGORY_UPDATED,
    category
  }
}

function submitCategoryEdit(category) {
	let request = {
		method: 'PUT',
		headers: { 
			'Content-Type':'application/json',
			'X-Auth-Header': localStorage.getItem('authToken') 
		},
		body: JSON.stringify(category)
	}

	return dispatch => {
		return fetch(`/api/categories/${category.id}`, request)
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
      dispatch(toggleCategoryFormVisibility(null, {}))

			// Dispatch the success action
			dispatch(flashOperationResult(EDIT_CATEGORY_SUCCESS_MESSAGE, true))

      // Update category
      dispatch(setCategoryUpdated(response))

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("update category error", err)
			}
		})
	}
}

/* Create */

function setCategoryCreated(category) {
  return {
    type: SET_CATEGORY_CREATED,
    category
  }
}

function submitCategoryCreate(category) {
	let request = {
		method: 'POST',
		headers: { 
			'Content-Type':'application/json',
			'X-Auth-Header': localStorage.getItem('authToken') 
		},
		body: JSON.stringify(category)
	}

	return dispatch => {
		return fetch('/api/categories', request)
		.then(rawResponse => {
			if (!rawResponse.ok) {
        switch(rawResponse.status) {
          case 400:
            dispatch(flashOperationResult(rawResponse.error))
            break
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
      dispatch(toggleCategoryFormVisibility(null, {}))

      // Update categories
      dispatch(setCategoryCreated(response))

			// Dispatch the success action
			dispatch(flashOperationResult(CREATE_CATEGORY_SUCCESS_MESSAGE, true))

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("update category error", err)
			}
		})
	}
}

/* Delete */

function setCategoryDeleted(categoryID) {
  return {
    type: SET_CATEGORY_DELETED,
    categoryID
  }
}

function toggleCategoryDeleteConfirmation(categoryID) {
  return {
    type: TOGGLE_CATEGORY_DELETE_CONFIRMATION,
    categoryID
  }
}

function submitCategoryDelete(categoryID) {
	let request = {
		method: 'DELETE',
		headers: { 
			'Content-Type':'application/json',
			'X-Auth-Header': localStorage.getItem('authToken') 
		}
	}

	return dispatch => {
		return fetch(`/api/categories/${categoryID}`, request)
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

      dispatch(toggleCategoryDeleteConfirmation(0))

			// Dispatch the success action
			dispatch(flashOperationResult(DELETE_CATEGORY_SUCCESS_MESSAGE, true))

      // Update categories
      dispatch(setCategoryDeleted(categoryID))

			return Promise.resolve()
		}).catch(err => {
			if (err) {
				console.warn("delete category error", err)
			}
		})
	}
}

module.exports = {
	TOGGLE_CATEGORY_FORM_VISIBILITY,
	SET_CATEGORIES,
	SET_CATEGORY_UPDATED,
	SET_CATEGORY_CREATED,
	SET_CATEGORY_DELETED,
	TOGGLE_CATEGORY_DELETE_CONFIRMATION,
	toggleCategoryFormVisibility,
	toggleCategoryDeleteConfirmation,    
	fetchCategories,
	submitCategoryEdit,
	submitCategoryCreate,
	submitCategoryDelete
}
