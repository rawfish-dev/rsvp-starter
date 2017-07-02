import React, { Component } from 'react';
import { connect } from 'react-redux';
import { reduxForm,reset,change,Field,Fields } from 'redux-form';
import { Row,Col,FormGroup,FormControl,ControlLabel,Radio,Button,Alert } from 'react-bootstrap';
import ReCAPTCHA from "react-google-recaptcha";

const { textarea } = require('./styles.css');

import {
  submitGuestRSVPCreate
} from '../../actions/guest';

import { isEmpty,isIncluded } from '../../validation';

const validationValues = Object.freeze({
  FULL_NAME_MIN_LENGTH: 2,
  FULL_NAME_MAX_LENGTH: 100,
  MAXIMUM_GUEST_COUNT_MINIMUM: 1,
  MAXIMUM_GUEST_COUNT_MAXIMUM: 10,
  MINIMUM_PHONE_NUMBER_LENGTH: 8,
  MAXIMUM_PHONE_NUMBER_LENGTH: 20,
  REMARKS_MAXIMUM_LENGTH: 500
});

const validate = values => {
  var errors = {}

  if (isEmpty(values.fullName) || 
      values.fullName.length < validationValues.FULL_NAME_MIN_LENGTH || 
      values.fullName.length > validationValues.FULL_NAME_MAX_LENGTH) {
    errors.fullName = `Please enter a name between ${validationValues.FULL_NAME_MIN_LENGTH} to ${validationValues.FULL_NAME_MAX_LENGTH} characters long`;
  }

  if (!isEmpty(values.remarks) && values.remarks.length > validationValues.REMARKS_MAXIMUM_LENGTH) {
    errors.remarks = `Please enter some remarks no longer than ${validationValues.REMARKS_MAXIMUM_LENGTH} characters in length`;
  }

  if (isEmpty(values.mobilePhoneNumber) ||
    values.mobilePhoneNumber.length < validationValues.MINIMUM_PHONE_NUMBER_LENGTH ||
    values.mobilePhoneNumber.length > validationValues.MAXIMUM_PHONE_NUMBER_LENGTH) {
      errors.mobilePhoneNumber = `Please enter a mobile phone number between ${validationValues.MINIMUM_PHONE_NUMBER_LENGTH} to ${validationValues.MAXIMUM_PHONE_NUMBER_LENGTH} in length`;
  }

  if (isEmpty(values.reCAPTCHA)) {
    errors.reCAPTCHA = `Please click on the checkbox`;
  }

  return errors;
}

const fullNameInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Full Name:
    </Col>

    <Col lg={6}>
      <FormControl type="text" {...field.input} />

      {field.meta.touched && field.meta.error && <div className="form-error">{field.meta.error}</div>}
    </Col>
  </FormGroup>;

const attendingInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Will you be able to attend:
    </Col>

    <Col lg={8}>
      <Radio inline checked={field.input.value} onClick={value => field.input.onChange(true)}>Yes</Radio>
      <Radio inline checked={!field.input.value} onClick={value => field.input.onChange(false)}>No</Radio>
    </Col>
  </FormGroup>;

const guestCountInput = field =>
  <FormGroup controlId="formControlsSelect">
    <Col componentClass={ControlLabel} lg={4}>
      How many people will be attending:
    </Col>

    <Col lg={2}>
      <FormControl componentClass="select" className="margin-top-sm" {...field.input}>
        {[...Array(this.props.fields.guestCount.initialValue)].map((_, i) => {
          let optionValue = i+1;
          return <option value={optionValue.toString()}>{optionValue}</option>;
        })}
      </FormControl>
    </Col>
  </FormGroup>;

const specialDietInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Special dietary requirements:
    </Col>

    <Col lg={8}>
      <Radio inline checked={field.input.value} onClick={value => field.input.onChange(true)}>Yes</Radio>
      <Radio inline checked={!field.input.value} onClick={value => field.input.onChange(false)}>No</Radio>

      {field.input.value === true && <div className="well margin-top-md margin-bottom-impt-xs">
        <p>Thank you for informing us. We will contact you shortly after you complete the form to understand your needs<i className="fa fa-smile-o fa-fw fa-lg"></i></p>
      </div>}
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
        placeholder="Please use this space to tell us if you have any additional needs e.g. babychairs, wheelchairs etc"
        rows="5"
        style={textarea}
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
      <small className="text-muted">* For receiving reminders and updates.</small>
      {field.meta.touched && field.meta.error && <div className="form-error">{field.meta.error}</div>}
    </Col>
  </FormGroup>;

const recaptchaInput = field =>
  <Col className="margin-top-md margin-bottom-lg" lg={8} lgOffset={4}>
    <ReCAPTCHA
      sitekey="6Ld9PQkUAAAAADM9IjQjKnXKNOOswNJe0NknxJAF"
      onChange={field.onReCAPTCHAChange}
    />
    {field.meta.touched && field.meta.error && <div className="form-error">{field.meta.error}</div>}
  </Col>;

const submit = (values, dispatch) => {
  let rsvp = {
    invitationPrivateID: values.invitationPrivateID,
    fullName: values.fullName,
    attending: values.attending,
    guestCount: parseInt(values.guestCount),
    reCAPTCHA: values.reCAPTCHA,
    remarks: values.remarks,
    specialDiet: values.specialDiet,
    mobilePhoneNumber: values.mobilePhoneNumber
  }

  return dispatch(submitGuestRSVPCreate(rsvp))
}

class RSVPForm extends Component {

  render() {
    const { handleSubmit, submitting } = this.props;
    const guestGreeting = this.props.fields.fullName.initialValue;

    return <div>
      <div className="panel">
        <div className="panel-body">
          <div className="margin-bottom-xs margin-right-sm text-muted text-right">
            <small>Not you? Click <a href="/">here</a> to RSVP</small>
          </div>

          <p className="landing-sub-title margin-top-lg margin-left-md">Dear {guestGreeting} -</p>

          <div className="text-center text-muted margin-top-lg margin-bottom-lg">
            <p>Kindly complete the form to RSVP :)</p>
          </div>

          <form className="form-horizontal margin-left-xs margin-right-sm" onSubmit={handleSubmit(submit)}>
            <Field
              name="fullName"
              component={fullNameInput}
            />

            <Field
              name="attending"
              component={attendingInput}
            />

            <Field
              name="guestCount"
              component={guestCountInput}
            />

            <Field
              name="specialDiet"
              component={specialDietInput}
            />

            <Field
              name="remarks"
              component={remarksInput}
            />

            <Field
              name="mobilePhoneNumber"
              component={mobilePhoneNumberInput}
            />

            <Field
              name="recaptcha"
              component={recaptchaInput}
              onReCAPTCHAChange={this.props.onReCAPTCHAChange}
            />

            <div className="margin-top-md margin-bottom-md">
              {this.props.operation && <Row>
                <Col lg={6} lgOffset={3}>
                  <Alert bsStyle="danger" className="text-center">
                    <p>{this.props.operation.message}</p>
                  </Alert>
                </Col>
              </Row>}
            </div>

            <Row className="margin-top-md">
              <Col xs={6}>
                <small className="text-muted">Facing difficulties? Contact ## Contact Name ## @ ## Contact Phone Number ##</small>
              </Col>

              <Col className="text-right margin-top-sm" xs={6}>
                <Button bsStyle="default" bsSize="small" onClick={this.props.onResetClick}>Clear</Button>
                <Button type="submit" className="margin-left-sm" bsStyle="primary" bsSize="small" disabled={submitting}>Submit</Button>
              </Col>
            </Row>
          </form>
        </div>
      </div>
    </div>;
  }
}

const mapStateToProps = (state) => {
  return {
    operation: state.operation
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    onResetClick: () => {
      dispatch(reset('rsvpForm'));
    },
    onReCAPTCHAChange: (value) => {
      dispatch(change('rsvpForm', 'reCAPTCHA', value));
    }
  };
};

RSVPForm = reduxForm({ 
  form: 'rsvpForm',
  validate
})(RSVPForm);

RSVPForm = connect(
  mapStateToProps,
  mapDispatchToProps
)(RSVPForm)

export default RSVPForm;
