import React, { Component } from 'react';
import { connect } from 'react-redux';
import { reduxForm,reset,change,Field } from 'redux-form';
import { browserHistory } from 'react-router';
import { Row,Col,FormGroup,FormControl,ControlLabel,Button } from 'react-bootstrap';

import {
  toggleCategoryFormVisibility,
  submitCategoryCreate,
  submitCategoryEdit
} from '../../actions/category';

import {
  CATEGORY_FORM_NEW_MODE,
  CATEGORY_FORM_EDIT_MODE
} from '../../constants';

import { isEmpty } from '../../validation';

const validationValues = Object.freeze({
  TAG_MIN_LENGTH: 1,
  TAG_MAX_LENGTH: 100
})

const validate = values => {
  var errors = {}

  if (isEmpty(values.tag) || 
      values.tag.length < validationValues.TAG_MIN_LENGTH || 
      values.tag.length > validationValues.TAG_MAX_LENGTH) {
    errors.tag = `Please enter a tag between ${validationValues.TAG_MIN_LENGTH} to ${validationValues.TAG_MAX_LENGTH} characters long`;
  }

  return errors;
}

const tagInput = field =>
  <FormGroup>
    <Col componentClass={ControlLabel} lg={4}>
      Tag:
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
  let category = {
    id: values.id,
    tag: values.tag
  }

  if (props.mode === CATEGORY_FORM_NEW_MODE) {
    return dispatch(submitCategoryCreate(category))  
  }

  return dispatch(submitCategoryEdit(category))
}

const newCategoryMode = {
  header: 'New Category',
  submitText: 'Create'
}

const editCategoryMode = {
  header: 'Edit Category',
  submitText: 'Update'
}

class CategoryForm extends Component {

  render() {
    const { handleSubmit, submitting } = this.props;
    const formText = this.props.mode === CATEGORY_FORM_NEW_MODE ? newCategoryMode : editCategoryMode;

    return <Row>
      <Col lg={6} lgOffset={3}>
        <div className="well">
          <form className="form-horizontal margin-left-xs margin-right-sm" onSubmit={handleSubmit(submit)}>
            <h4>{formText.header} <span>{this.props.categoryID}</span></h4>
            
            <Field
              name="tag"
              component={tagInput}
            />

            <Row className="margin-top-md">
              <Col className="text-right margin-top-sm" xs={12}>
                <Button bsStyle="default" bsSize="sm" onClick={this.props.onToggleCategoryFormVisibilityClick}>Cancel</Button>
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
    onToggleCategoryFormVisibilityClick: (mode) => {
      dispatch(toggleCategoryFormVisibility(null, {}))
    }
  };
};

CategoryForm = reduxForm({ 
  form: 'categoryForm',
  validate
})(CategoryForm);

CategoryForm = connect(
  mapStateToProps,
  mapDispatchToProps
)(CategoryForm)

export default CategoryForm;
