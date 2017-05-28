import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Row,Col,Table,Button,Modal } from 'react-bootstrap';

import {
  toggleCategoryFormVisibility,
  toggleCategoryDeleteConfirmation,
  submitCategoryDelete
} from '../../actions/category';

import {
  CATEGORY_FORM_NEW_MODE,
  CATEGORY_FORM_EDIT_MODE
} from '../../constants';

import CategoryForm from '../CategoryForm';

class Categories extends Component {
    render() {
        return <div className="panel-body">
          <Row className="margin-bottom-sm">
            <Col lg={8}>
              <h4>Categories</h4>
            </Col>

            <Col lg={4}>
              <div className="text-right">
                <Button bsStyle="primary" bsSize="small" className="margin-top-sm" onClick={this.props.onAddCategoryClick}><i className="fa fa-plus fa-fw fa-lg"></i> Add Category</Button>
              </div>
            </Col>
          </Row>

          {this.props.deleteCategoryConfirmation.visible && <div className="static-modal">
            <Modal.Dialog>
              <Modal.Header>
                <Modal.Title>Delete Category {this.props.deleteCategoryConfirmation.categoryID}</Modal.Title>
              </Modal.Header>

              <Modal.Body>
                Are you sure you want to delete this category?
              </Modal.Body>

              <Modal.Footer>
                <Button bsSize="sm" onClick={() => {this.props.onToggleDeleteCategory(0)}}>Close</Button>
                <Button bsSize="sm" bsStyle="danger" onClick={() => {this.props.onDeleteCategory(this.props.deleteCategoryConfirmation.categoryID)}}>Delete</Button>
              </Modal.Footer>
            </Modal.Dialog>
          </div>}

          {(() => {
            if (this.props.categoryForm !== null && this.props.categoryForm.initialValues !== null) {
              return <CategoryForm 
                mode={this.props.categoryForm.mode}
                categoryID={this.props.categoryForm.initialValues.id}
                initialValues={this.props.categoryForm.initialValues} />;
            }
          })()}

          <Table striped={true} responsive>
            <thead>
              <tr>
                <th>ID</th>
                <th>Tag</th>
                <th>Invitation Count</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {this.props.categories && this.props.categories.map((category) => {
                  return <tr key={category.id}>
                    <td>{category.id}</td>
                    <td>{category.tag}</td>
                    <td>{category.total}</td>
                    <td>
                      <Button bsStyle="danger" bsSize="xs" className="margin-right-sm" disabled={category.total > 0} onClick={() => {this.props.onToggleDeleteCategory(category.id)}}>Delete</Button>
                      <Button bsStyle="default" bsSize="xs" onClick={() => {this.props.onEditCategoryClick(category)}}>Edit</Button>
                    </td>
                  </tr>;
              })}
            </tbody>
          </Table>
        </div>;
    }
}

const mapStateToProps = (state) => {
  return {
    categoryForm: state.categoryForm,
    deleteCategoryConfirmation: state.deleteCategoryConfirmation
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    onAddCategoryClick: () => {
      dispatch(toggleCategoryFormVisibility(CATEGORY_FORM_NEW_MODE, {}))
    },
    onEditCategoryClick: (category) => {
      dispatch(toggleCategoryFormVisibility(CATEGORY_FORM_EDIT_MODE, category))
    },
    onToggleDeleteCategory: (categoryID) => {
      dispatch(toggleCategoryDeleteConfirmation(categoryID))
    },
    onDeleteCategory: (categoryID) => {
      dispatch(submitCategoryDelete(categoryID))
    }
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(Categories);
