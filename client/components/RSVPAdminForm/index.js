import React, { Component } from 'react';
import { connect } from 'react-redux';
import { reduxForm,reset,change,Field } from 'redux-form';
import { browserHistory } from 'react-router';
import { Row,Col,FormGroup,FormControl,ControlLabel,Button,Radio } from 'react-bootstrap';

import {
  toggleRSVPFormVisibility,
  submitRSVPCreate,
  submitRSVPEdit
} from '../../actions/rsvp';

import {
  RSVP_FORM_NEW_MODE,
  RSVP_FORM_EDIT_MODE,
  INVITATION_MAX_GUESTS,
  INVITATION_STATUS_NOT_SENT,
  INVITATION_STATUS_SENT
} from '../../constants';

import { isEmpty } from '../../validation';

const validationValues = Object.freeze({
  GREETING_MIN_LENGTH: 2,
  GREETING_MAX_LENGTH: 100,
  MINIMUM_GUEST_COUNT: 1,
  MAXIMUM_GUEST_COUNT: 10,
  MINIMUM_PHONE_NUMBER_LENGTH: 8,
  MAXIMUM_PHONE_NUMBER_LENGTH: 20,
  NOTES_MAXIMUM_LENGTH: 500
});

const validate = values => {
  console.warn("rsvp validation", values);

  var errors = {}

  if (isEmpty(values.fullName) || 
      values.fullName.length < validationValues.GREETING_MIN_LENGTH || 
      values.fullName.length > validationValues.GREETING_MAX_LENGTH) {
    errors.fullName = `Please enter a full name between ${validationValues.GREETING_MIN_LENGTH} to ${validationValues.GREETING_MAX_LENGTH} characters long`;
  }

  let parsedGuestCount = parseInt(values.guestCount);
  if (parsedGuestCount < validationValues.MINIMUM_GUEST_COUNT ||
    parsedGuestCount > validationValues.MAXIMUM_GUEST_COUNT) {
      errors.maximumGuestCount = `Please choose a guest count between ${validationValues.MINIMUM_GUEST_COUNT} to ${validationValues.MAXIMUM_GUEST_COUNT}`;
  }

  if (!isEmpty(values.remarks) && values.remarks.length > validationValues.NOTES_MAXIMUM_LENGTH) {
    errors.remarks = `Please enter some remarks no longer than ${validationValues.NOTES_MAXIMUM_LENGTH} characters in length`;
  }

  if (isEmpty(values.mobilePhoneNumber) ||
    values.mobilePhoneNumber.length < validationValues.MINIMUM_PHONE_NUMBER_LENGTH ||
    values.mobilePhoneNumber.length > validationValues.MAXIMUM_PHONE_NUMBER_LENGTH) {
      errors.mobilePhoneNumber = `Please enter a mobile phone number between ${validationValues.MINIMUM_PHONE_NUMBER_LENGTH} to ${validationValues.MAXIMUM_PHONE_NUMBER_LENGTH} long`;
  }

  return errors;
}

const invitationPrivateIDInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Invitation Private ID:
    </Col>

    <Col lg={8}>
      <FormControl
        type="text"
        placeholder="The complicated text obtained from 'Copy Link'" 
        {...field.input}>
      </FormControl>
    </Col>
  </FormGroup>;

const fullNameInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Full Name:
    </Col>

    <Col lg={8}>
      <FormControl
        type="text"
        placeholder="Exact full name without 'Dear' e.g Mitten and Family" 
        {...field.input}>
      </FormControl>
      {field.meta.touched && field.meta.error && <div className="form-error">{field.meta.error}</div>}
    </Col>
  </FormGroup>;

const attendingInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Attending:
    </Col>

    <Col lg={8}>
      <Radio inline checked={field.input.value} onClick={value => field.input.onChange(true)}>Yes</Radio>
      <Radio inline checked={!field.input.value} onClick={value => field.input.onChange(false)}>No</Radio>
    </Col>
  </FormGroup>;

const specialDietInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Special diet:
    </Col>

    <Col lg={8}>
      <Radio inline checked={field.input.value} onClick={value => field.input.onChange(true)}>Yes</Radio>
      <Radio inline checked={!field.input.value} onClick={value => field.input.onChange(false)}>No</Radio>
    </Col>
  </FormGroup>;

const guestCountInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      No. of Guests:
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

const remarksInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Remarks (optional):
    </Col>

    <Col lg={8}>
      <FormControl 
        componentClass="textarea"
        placeholder="Additional needs e.g. babychairs, wheelchairs etc"
        rows="5"
        {...field.input}>
      </FormControl>
      {field.meta.touched && field.meta.error && <div className="form-error">{field.meta.error}</div>}
    </Col>
  </FormGroup>;

const mobilePhoneNumberInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Mobile phone number:
    </Col>

    <Col lg={6}>
      <FormControl
        type="text" 
        {...field.input}>
      </FormControl>
      {field.meta.touched && field.meta.error && <div className="form-error">{field.meta.error}</div>}
    </Col>
  </FormGroup>;

const submit = (values, dispatch, props) => {
  let rsvp = {
    id: values.id,
    invitationPrivateID: values.invitationPrivateID,
    fullName: values.fullName,
    attending: values.attending,
    specialDiet: values.specialDiet,
    guestCount: parseInt(values.guestCount),
    remarks: values.remarks,
    mobilePhoneNumber: values.mobilePhoneNumber
  }

  if (props.mode === RSVP_FORM_NEW_MODE) {
    return dispatch(submitRSVPCreate(rsvp))  
  }

  return dispatch(submitRSVPEdit(rsvp))
}

const newRSVPMode = {
  header: 'New RSVP',
  submitText: 'Create'
}

const editRSVPMode = {
  header: 'Edit RSVP',
  submitText: 'Update'
}

class RSVPAdminForm extends Component {

  render() {
    const { handleSubmit, submitting } = this.props;
    const formText = this.props.mode === RSVP_FORM_NEW_MODE ? newRSVPMode : editRSVPMode;

    return <Row>
      <Col lg={6} lgOffset={3}>
        <div className="well">
          <form className="form-horizontal margin-left-xs margin-right-sm" onSubmit={handleSubmit(submit)}>
            <h4>{formText.header} <span>{this.props.rsvpID}</span></h4>

            {this.props.mode === RSVP_FORM_NEW_MODE && <Field
              name="invitationPrivateID"
              component={invitationPrivateIDInput} 
              />
            }

            <Field
              name="fullName"
              component={fullNameInput}
            />

            <Field
              name="attending"
              component={attendingInput}
            />

            <Field
              name="specialDiet"
              component={specialDietInput}
            />

            <Field
              name="guestCount"
              component={guestCountInput}  
            />

            <Field
              name="remarks"
              component={remarksInput}
            />

            <Field
              name="mobilePhoneNumber"
              component={mobilePhoneNumberInput}
            />

            <Row className="margin-top-md">
              <Col className="text-right margin-top-sm" xs={12}>
                <Button bsStyle="default" bsSize="sm" onClick={this.props.onToggleRSVPFormVisibilityClick}>Cancel</Button>
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
    onToggleRSVPFormVisibilityClick: (mode) => {
      dispatch(toggleRSVPFormVisibility(null, {}))
    }
  };
};

RSVPAdminForm = reduxForm({ 
  form: 'rsvpAdminForm',
  validate
})(RSVPAdminForm);

RSVPAdminForm = connect(
  mapStateToProps,
  mapDispatchToProps
)(RSVPAdminForm)

export default RSVPAdminForm;
