import React, { Component } from 'react';
import { connect } from 'react-redux';
import { reduxForm,reset,change,Field,Fields } from 'redux-form';
import { browserHistory } from 'react-router';
import { Row,Col,FormGroup,FormControl,ControlLabel,Button } from 'react-bootstrap';
import ReCAPTCHA from "react-google-recaptcha";

import {
  loginUser
} from '../../actions/login';

import { isEmpty } from '../../validation';

const validate = values => {
  var errors = {}

  if (isEmpty(values.username)) {
    errors.username = `Please enter your username`;
  }

  if (isEmpty(values.password)) {
    errors.password = `Please enter your password`;
  }

  if (isEmpty(values.recaptcha)) {
    errors.recaptcha = `Please click on the checkbox`;
  }

  return errors;
}

const usernameInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Username:
    </Col>

    <Col lg={6}>
      <FormControl
        type="text" 
        {...field.input}>
      </FormControl>
      {field.meta.touched && field.meta.error && <div className="form-error">{field.meta.error}</div>}
    </Col>
  </FormGroup>;

const passwordInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Password:
    </Col>

    <Col lg={6}>
      <FormControl
        type="password" 
        {...field.input}>
      </FormControl>
      {field.meta.touched && field.meta.error && <div className="form-error">{field.meta.error}</div>}
    </Col>
  </FormGroup>;

const recaptchaInput = field =>
  <Col className="margin-top-md margin-bottom-lg" lg={6} lgOffset={4}>
    <ReCAPTCHA
      sitekey="6Ld9PQkUAAAAADM9IjQjKnXKNOOswNJe0NknxJAF"
      onChange={field.onReCAPTCHAChange}
    />
    {field.meta.touched && field.meta.error && <div className="form-error">{field.meta.error}</div>}
  </Col>;

const submit = (values, dispatch) => {
  console.warn("VALUES", values)

  let credentials = {
    username: values.username,
    password: values.password,
    recaptcha: values.recaptcha
  };

  return dispatch(loginUser(credentials));
}

class LoginForm extends Component {

  componentWillMount() {
    if (this.props.auth.isAuthenticated) {
      browserHistory.push('/control_panel');
    }
  }

  render() {
    const { handleSubmit } = this.props;
    const errorMessage = this.props.auth.errorMessage;

    return <div className="margin-top-lg">
      <Row className="padding-top-lg">
        <Col lg={6} lgOffset={3}>
          {errorMessage && <h4 className="text-center text-danger margin-bottom-md">{errorMessage}</h4>}

          <div className="panel">
            <div className="panel-body">
              <form className="form-horizontal margin-left-xs margin-right-sm" onSubmit={handleSubmit(submit)}>
                <Field
                  name="username"
                  component={usernameInput}
                />

                <Field
                  name="password"
                  component={passwordInput}
                />

                <Field
                  name="recaptcha"
                  component={recaptchaInput}
                  onReCAPTCHAChange={this.props.onReCAPTCHAChange}
                />

                <Row className="margin-top-md">
                  <Col className="text-right margin-top-sm" xs={12}>
                    <Button type="submit" className="margin-left-sm" bsStyle="primary" bsSize="small" disabled={this.props.auth.isFetching}>Login</Button>
                  </Col>
                </Row>
              </form>
            </div>
          </div>
        </Col>
      </Row>
    </div>;
  }     
}

const mapStateToProps = (state) => {
  return {
    auth: state.auth
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    onReCAPTCHAChange: (value) => {
      dispatch(change('loginForm', 'recaptcha', value));
    }
  };
};

LoginForm = reduxForm({ 
  form: 'loginForm',
  validate
})(LoginForm);

LoginForm = connect(
  mapStateToProps,
  mapDispatchToProps
)(LoginForm)

export default LoginForm;
