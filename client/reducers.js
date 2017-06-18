import { reducer as formReducer } from 'redux-form';

import {
	SET_OPERATION_SUCCESS,
	SET_OPERATION_FAILURE,
  UNSET_OPERATION_RESULT
} from './actions/general';

import { 
  SET_GUEST_RSVP
} from './actions/guest';

import { 
  TOGGLE_RSVP_FORM_VISIBILITY,
	SET_RSVPS,
  SET_RSVP_UPDATED,
  SET_RSVP_CREATED,
  SET_RSVP_DELETED,
  TOGGLE_RSVP_DELETE_CONFIRMATION
} from './actions/rsvp';

import { 
	TOGGLE_CATEGORY_FORM_VISIBILITY,
	SET_CATEGORIES,
  SET_CATEGORY_UPDATED,
  SET_CATEGORY_CREATED,
  SET_CATEGORY_DELETED,
  TOGGLE_CATEGORY_DELETE_CONFIRMATION
} from './actions/category';

import {
  TOGGLE_INVITATION_FORM_VISIBILITY,
  SET_INVITATIONS,
  SET_INVITATION_UPDATED,
  SET_INVITATION_CREATED,
  SET_INVITATION_DELETED,
  TOGGLE_INVITATION_DELETE_CONFIRMATION
} from './actions/invitation';

import { 
	LOGIN_REQUEST,
	LOGIN_SUCCESS,
	LOGIN_FAILURE
} from './actions/login';

import { 
	LOGOUT_SUCCESS
} from './actions/logout';

export function operation(state = null, action) {
	switch(action.type) {
		case SET_OPERATION_SUCCESS:
			return {
				message: action.message,
				success: true
			}
    case SET_OPERATION_FAILURE:
      return {
        message: action.message,
        success: false
      }
    case UNSET_OPERATION_RESULT:
      return null;
    default:
      return state;
	}
}

export function guestRSVP(state = null, action) {
  switch (action.type) {
		case SET_GUEST_RSVP:
      return action.rsvp;
    default:
      return state;
  }
}

export function rsvps(state = [], action) {
	switch (action.type) {
		case SET_RSVPS:
			return action.rsvps;

    case SET_RSVP_UPDATED:
      var copiedRSVPs = state.slice();

      for (var i = 0; i < copiedRSVPs.length; i++) {
        if (copiedRSVPs[i].id === action.rsvp.id) {
          copiedRSVPs[i] = action.rsvp;
        }
      }

      return copiedRSVPs;

    case SET_RSVP_CREATED:
      // Ensure we don't somehow add a rsvp that is already there
      for (var i = 0; i < state.length; i++) {
        if (state[i].id === action.rsvp.id) {
          return state;
        }
      }

      return [action.rsvp, ...state];

    case SET_RSVP_DELETED:
      return state.filter(rsvp => rsvp.id !== action.rsvpID);

		default:
			return state; 
	}
}

export function categories(state = [], action) {
	switch (action.type) {
		case SET_CATEGORIES:
			return action.categories;

    case SET_CATEGORY_UPDATED:
      var copiedCategories = state.slice();

      for (var i = 0; i < copiedCategories.length; i++) {
        if (copiedCategories[i].id === action.category.id) {
          copiedCategories[i] = action.category;
        }
      }

      return copiedCategories;

    case SET_CATEGORY_CREATED:
      // Ensure we don't somehow add a category that is already there
      for (var i = 0; i < state.length; i++) {
        if (state[i].id === action.category.id) {
          return state;
        }
      }

      return [action.category, ...state];

    case SET_CATEGORY_DELETED:
      return state.filter(category => category.id !== action.categoryID);

		default:
			return state; 
	}
}

export function invitations(state = [], action) {
  switch (action.type) {
		case SET_INVITATIONS:
			return action.invitations;

    case SET_INVITATION_UPDATED:
      var copiedInvitations = state.slice();

      for (var i = 0; i < copiedInvitations.length; i++) {
        if (copiedInvitations[i].id === action.invitation.id) {
          copiedInvitations[i] = action.invitation;
        }
      }

      return copiedInvitations;
    
    case SET_INVITATION_CREATED:
      // Ensure we don't somehow add an invitation that is already there
      for (var i = 0; i < state.length; i++) {
        if (state[i].id === action.invitation.id) {
          return state;
        }
      }

      return [action.invitation, ...state];

    case SET_INVITATION_DELETED:
      return state.filter(invitation => invitation.id !== action.invitationID);

    default:
			return state;
  }
}

const defaultFormState = {
  mode: null,
  initialValues: null,
  visible: false
};

export function rsvpForm(state = defaultFormState, action) {
	switch(action.type) {
		case TOGGLE_RSVP_FORM_VISIBILITY:
      if (state.visible === true) {
        // Reset to default if trying to toggle it off
        return defaultFormState;
      }

			return {
        mode: action.mode,
        initialValues: action.initialValues,
        visible: true
      };
		default:
			return state;
	}
}

export function categoryForm(state = defaultFormState, action) {
	switch(action.type) {
		case TOGGLE_CATEGORY_FORM_VISIBILITY:
      if (state.visible === true) {
        // Reset to default if trying to toggle it off
        return defaultFormState;
      }

			return {
        mode: action.mode,
        initialValues: action.initialValues,
        visible: true
      };
		default:
			return state;
	}
}

export function invitationForm(state = defaultFormState, action) {
	switch(action.type) {
		case TOGGLE_INVITATION_FORM_VISIBILITY:
      if (state.visible === true) {
        // Reset to default if trying to toggle it off
        return defaultFormState;
      }

      // Set default values
      var defaultValues = Object.assign({}, action.initialValues, {
        maximumGuestCount: 1
      });

			return {
        mode: action.mode,
        initialValues: defaultValues,
        visible: true
      };
		default:
			return state;
	}
}

const defaultDeleteRSVPConfirmationState = {
  visible: false,
  rsvpID: 0 
};

export function deleteRSVPConfirmation(state = defaultDeleteRSVPConfirmationState, action) {
  switch(action.type) {
    case TOGGLE_RSVP_DELETE_CONFIRMATION:
      if (state.visible === true) {
        // Reset to default if trying to toggle it off
        return defaultDeleteRSVPConfirmationState;
      }

			return {
        visible: true,
        rsvpID: action.rsvpID
      };

    default:
      return state;
  }
}

const defaultDeleteCategoryConfirmationState = {
  visible: false,
  categoryID: 0 
};

export function deleteCategoryConfirmation(state = defaultDeleteCategoryConfirmationState, action) {
  switch(action.type) {
    case TOGGLE_CATEGORY_DELETE_CONFIRMATION:
      if (state.visible === true) {
        // Reset to default if trying to toggle it off
        return defaultDeleteCategoryConfirmationState;
      }

			return {
        visible: true,
        categoryID: action.categoryID
      };

    default:
      return state;
  }
}

const defaultDeleteInvitationConfirmationState = {
  visible: false,
  invitationID: 0 
};

export function deleteInvitationConfirmation(state = defaultDeleteInvitationConfirmationState, action) {
  switch(action.type) {
    case TOGGLE_INVITATION_DELETE_CONFIRMATION:
      if (state.visible === true) {
        // Reset to default if trying to toggle it off
        return defaultDeleteInvitationConfirmationState;
      }

			return {
        visible: true,
        invitationID: action.invitationID
      };

    default:
      return state;
  }
}


// TODO:: Check if the token is expired
export function auth(state = {
    isFetching: false,
    isAuthenticated: localStorage.getItem('authToken') ? true : false
  }, action) {
	  switch (action.type) {
    	case LOGIN_REQUEST:
			return Object.assign({}, state, {
				isFetching: true,
				isAuthenticated: false,
				user: action.credentials
			});
		case LOGIN_SUCCESS:
			return Object.assign({}, state, {
				isFetching: false,
				isAuthenticated: true,
				errorMessage: ''
			});
		case LOGIN_FAILURE:
			return Object.assign({}, state, {
				isFetching: false,
				isAuthenticated: false,
				errorMessage: action.message
			});
		case LOGOUT_SUCCESS:
			return Object.assign({}, state, {
				isFetching: false,
				isAuthenticated: false
			});
		default:
			return state;
  	}
}

export function removeByKey(obj, deleteKey) {
	return Object.keys(obj).reduce((result, key) => {
		if (parseInt(key) !== deleteKey) {
            result[key] = obj[key];
        }
        return result;
    }, {});
}

const gatheredReducers = {operation, guestRSVP, rsvps, categories, invitations, rsvpForm, categoryForm, invitationForm, deleteRSVPConfirmation, deleteCategoryConfirmation, deleteInvitationConfirmation, auth, form: formReducer};

export default gatheredReducers;
