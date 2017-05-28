import React, { Component } from 'react';
import { connect } from 'react-redux';
import { reduxForm,reset,change,Field } from 'redux-form';
import { browserHistory } from 'react-router';
import { Row,Col,FormGroup,FormControl,ControlLabel,Button } from 'react-bootstrap';

import {
  toggleInvitationFormVisibility,
  submitInvitationCreate,
  submitInvitationEdit
} from '../../actions/invitation';

import {
  INVITATION_FORM_NEW_MODE,
  INVITATION_FORM_EDIT_MODE,
  INVITATION_MAX_GUESTS,
  INVITATION_STATUS_NOT_SENT,
  INVITATION_STATUS_SENT
} from '../../constants';

import { isEmpty } from '../../validation';

const validationValues = Object.freeze({
  GREETING_MIN_LENGTH: 2,
  GREETING_MAX_LENGTH: 20,
  MINIMUM_GUEST_COUNT: 1,
  MAXIMUM_GUEST_COUNT: 10,
  MAXIMUM_PHONE_NUMBER_LENGTH: 20,
  NOTES_MAXIMUM_LENGTH: 500
});

const validate = values => {
  console.warn("invitation validation", values);

  var errors = {}

  if (isEmpty(values.greeting) || 
      values.greeting.length < validationValues.GREETING_MIN_LENGTH || 
      values.greeting.length > validationValues.GREETING_MAX_LENGTH) {
    errors.greeting = `Please enter a greeting between ${validationValues.GREETING_MIN_LENGTH} to ${validationValues.GREETING_MAX_LENGTH} characters long`;
  }

  let parsedMaximumGuestCount = parseInt(values.maximumGuestCount);
  if (parsedMaximumGuestCount < validationValues.MINIMUM_GUEST_COUNT ||
    parsedMaximumGuestCount > validationValues.MAXIMUM_GUEST_COUNT) {
      errors.maximumGuestCount = `Please choose a maximum guest count between ${validationValues.MINIMUM_GUEST_COUNT} to ${validationValues.MAXIMUM_GUEST_COUNT}`;
  }

  if (values.mobilePhoneNumber && values.mobilePhoneNumber.length > validationValues.MAXIMUM_PHONE_NUMBER_LENGTH) {
      errors.mobilePhoneNumber = `Please enter a mobile phone number less than ${validationValues.MAXIMUM_PHONE_NUMBER_LENGTH} numbers`;
  }

  if (!isEmpty(values.notes) && values.notes.length > validationValues.NOTES_MAXIMUM_LENGTH) {
    errors.notes = `Please enter some notes no longer than ${validationValues.NOTES_MAXIMUM_LENGTH} characters in length`;
  }

  return errors;
}

const categorySelect = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Category:
    </Col>

    <Col lg={8}>
      {(() => {
        if (field.categories && field.categories.length > 0) {
          return <FormControl componentClass="select" {...field.input}>
            {field.categories && field.categories.map((category) => {
              return <option value={category.id}>{category.tag}</option>; 
            })}
          </FormControl>;
        }

        return <p className="margin-top-xs">No categories added</p>
      })()}
    </Col>
  </FormGroup>;

const greetingInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Greeting:
    </Col>

    <Col lg={8}>
      <FormControl
        type="text"
        placeholder="Exact greeting without 'Dear' e.g Mitten and Family" 
        {...field.input}>
      </FormControl>
      {field.meta.touched && field.meta.error && <div className="form-error">{field.meta.error}</div>}
    </Col>
  </FormGroup>;

const maximumGuestCountInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Max No. of Guests:
    </Col>

    <Col lg={8}>
      <FormControl componentClass="select" {...field.input}>
        {[...Array(INVITATION_MAX_GUESTS)].map((_, i) => {
          let optionValue = i+1;
          return <option value={optionValue.toString()}>{optionValue}</option>;
        })}
      </FormControl>
    </Col>
  </FormGroup>;

const notesInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Notes (optional):
    </Col>

    <Col lg={8}>
      <FormControl 
        componentClass="textarea"
        placeholder=""
        rows="5"
        {...field.input}>
      </FormControl>
      {field.meta.touched && field.meta.error && <div className="form-error">{field.meta.error}</div>}
    </Col>
  </FormGroup>;

const mobilePhoneNumberInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Mobile Phone No.:
    </Col>

    <Col lg={8}>
      <FormControl
        type="text" 
        {...field.input}>
      </FormControl>
      {field.meta.touched && field.meta.error && <div className="form-error">{field.meta.error}</div>}
    </Col>
  </FormGroup>;

const statusSelect = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Status:
    </Col>

    <Col lg={8}>
      <FormControl componentClass="select" placeholder="1" {...field.input}>
        <option value={INVITATION_STATUS_NOT_SENT}>Not Sent</option>
        <option value={INVITATION_STATUS_SENT}>Sent</option>
      </FormControl>
    </Col>
  </FormGroup>;

const submit = (values, dispatch, props) => {
  let invitation = {
    id: values.id,
    categoryID: parseInt(values.categoryID),
    greeting: values.greeting,
    maximumGuestCount: parseInt(values.maximumGuestCount),
    notes: values.notes,
    mobilePhoneNumber: values.mobilePhoneNumber,
    status: values.status
  }

  if (props.mode === INVITATION_FORM_NEW_MODE) {
    return dispatch(submitInvitationCreate(invitation))  
  }

  return dispatch(submitInvitationEdit(invitation))
}

const newInvitationMode = {
  header: 'New Invitation',
  submitText: 'Create'
}

const editInvitationMode = {
  header: 'Edit Invitation',
  submitText: 'Update'
}

class InvitationForm extends Component {

  render() {
    const { handleSubmit, submitting } = this.props;
    const formText = this.props.mode === INVITATION_FORM_NEW_MODE ? newInvitationMode : editInvitationMode;

    return <Row>
      <Col lg={6} lgOffset={3}>
        <div className="well">
          <form className="form-horizontal margin-left-xs margin-right-sm" onSubmit={handleSubmit(submit)}>
            <h4>{formText.header} <span>{this.props.invitationID}</span></h4>

            <Field
              name="categoryID"
              component={categorySelect}
              categories={this.props.categories}
            />

            <Field
              name="greeting"
              component={greetingInput}
            />

            <Field
              name="maximumGuestCount"
              component={maximumGuestCountInput}
            />

            <Field
              name="notes"
              component={notesInput}
            />

            <Field
              name="mobilePhoneNumber"
              component={mobilePhoneNumberInput}
            />

            <Field
              name="status"
              component={statusSelect}
            />

            <Row className="margin-top-md">
              <Col className="text-right margin-top-sm" xs={12}>
                <Button bsStyle="default" bsSize="sm" onClick={this.props.onToggleInvitationFormVisibilityClick}>Cancel</Button>
                <Button type="submit" className="margin-left-sm" bsStyle="success" bsSize="sm" disabled={submitting}>{formText.submitText}</Button>
              </Col>
            </Row>
          </form>
        </div>
      </Col>
    </Row>;
  }     
}

const mapStateToProps = (state) => {
  return {};
};

const mapDispatchToProps = (dispatch) => {
  return {
    onToggleInvitationFormVisibilityClick: (mode) => {
      dispatch(toggleInvitationFormVisibility(null, {}))
    }
  };
};

InvitationForm = reduxForm({ 
  form: 'invitationForm',
  validate
})(InvitationForm);

InvitationForm = connect(
  mapStateToProps,
  mapDispatchToProps
)(InvitationForm)

export default InvitationForm;
