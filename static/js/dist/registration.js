(function () {
  'use strict';

  function _defineProperty(obj, key, value) {
    if (key in obj) {
      Object.defineProperty(obj, key, {
        value: value,
        enumerable: true,
        configurable: true,
        writable: true
      });
    } else {
      obj[key] = value;
    }

    return obj;
  }

  function _objectSpread(target) {
    for (var i = 1; i < arguments.length; i++) {
      var source = arguments[i] != null ? arguments[i] : {};
      var ownKeys = Object.keys(source);

      if (typeof Object.getOwnPropertySymbols === 'function') {
        ownKeys = ownKeys.concat(Object.getOwnPropertySymbols(source).filter(function (sym) {
          return Object.getOwnPropertyDescriptor(source, sym).enumerable;
        }));
      }

      ownKeys.forEach(function (key) {
        _defineProperty(target, key, source[key]);
      });
    }

    return target;
  }

  var createObserver = function createObserver() {
    var subscribers = [];

    var subscribe = function subscribe(callback) {
      var index = subscribers.length;
      subscribers.push(callback);
      var isUnsubscribed = false;

      var unsubscribe = function unsubscribe() {
        if (isUnsubscribed) {
          return;
        }

        subscribers.splice(index, 1);
        isUnsubscribed = true;
      };

      return unsubscribe;
    };

    var emit = function emit() {
      for (var _len = arguments.length, values = new Array(_len), _key = 0; _key < _len; _key++) {
        values[_key] = arguments[_key];
      }

      subscribers.forEach(function (c) {
        return c.apply(void 0, values);
      });
    };

    return {
      subscribe: subscribe,
      emit: emit
    };
  };
  var createState = function createState(initialState) {
    var state = initialState;

    var _createObserver = createObserver(),
        subscribe = _createObserver.subscribe,
        emit = _createObserver.emit;

    var update = function update(newState) {
      var prevState = state;
      state = _objectSpread({}, state, newState);
      emit(state, prevState);
    };

    var getState = function getState() {
      return state;
    };

    return {
      subscribe: subscribe,
      update: update,
      getState: getState
    };
  };
  var createValidator = function createValidator(rules) {
    var validate = function validate(value) {
      var failedRule = rules.find(function (r) {
        return !r.check(value);
      });

      if (!failedRule) {
        return;
      }

      return failedRule.message;
    };

    return {
      validate: validate
    };
  };
  var createRepeatPattern = function createRepeatPattern(pattern, times) {
    var regExp = new RegExp(pattern, 'g');
    return function (value) {
      return (value.match(regExp) || []).length >= times;
    };
  };

  var emailRegex = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

  var letter = '[a-zA-Z]';
  var number = /\d/;
  var upperCase = '[A-Z]';
  var checkIfEnoughLetters = createRepeatPattern(letter, 7);
  var checkIfEnoughNumbers = createRepeatPattern(number, 1);
  var checkIfEnoughUpperCases = createRepeatPattern(upperCase, 1);
  var passwordRules = [{
    check: checkIfEnoughNumbers,
    message: '*at least 1 number'
  }, {
    check: checkIfEnoughUpperCases,
    message: '*at least 1 upper case'
  }, {
    check: checkIfEnoughLetters,
    message: '*at least 7 letters'
  }];
  var emailRules = [{
    check: function check(v) {
      return emailRegex.test(v);
    },
    message: '*email is invalid'
  }];
  var strongPasswordValidator = createValidator(passwordRules);
  var emailValidator = createValidator(emailRules);

  var createErrorView = function createErrorView(errorId) {
    if (!errorId || typeof errorId !== 'string') {
      throw new Error('Error Id is invalid');
    }

    var errorElem = document.getElementById(errorId);

    if (!errorElem) {
      throw new Error("There is no element with Id: ".concat(errorId));
    }

    var hiddenClass = 'hidden';

    var showError = function showError() {
      errorElem.classList.remove(hiddenClass);
    };

    var hideError = function hideError() {
      errorElem.classList.add(hiddenClass);
    };

    var setError = function setError(error) {
      errorElem.innerHTML = error;
    };

    var getErrorMessage = function getErrorMessage() {
      return errorElem.innerHTML;
    };

    return {
      showError: showError,
      hideError: hideError,
      setError: setError,
      getErrorMessage: getErrorMessage
    };
  };
  var createInputView = function createInputView(inputId) {
    if (!inputId || typeof inputId !== 'string') {
      throw new Error('Input Id is invalid');
    }

    var input = document.getElementById(inputId);

    if (!input) {
      throw new Error("There is no element with Id: ".concat(inputId));
    }

    var subscribeOnInputChange = function subscribeOnInputChange(cb) {
      input.addEventListener('input', cb);
      return function () {
        return input.removeEventListener('input', cb);
      };
    };

    return {
      subscribeOnInputChange: subscribeOnInputChange,

      get value() {
        return input.value;
      }

    };
  };
  var createFormView = function createFormView(formId) {
    if (!formId || typeof formId !== 'string') {
      throw new Error('Form Id is invalid');
    }

    var form = document.getElementById(formId);

    if (!form) {
      throw new Error("There is no element with Id ".concat(formId));
    }

    var subscribeOnSubmit = function subscribeOnSubmit(cb) {
      form.addEventListener('submit', cb);
      return function () {
        return form.removeEventListener('submit', cb);
      };
    };

    return {
      subscribeOnSubmit: subscribeOnSubmit,
      submit: function submit() {
        return form.submit();
      }
    };
  };

  var formView = createFormView('form');
  var emailView = createInputView('email');
  var passwordView = createInputView('password');
  var emailErrorView = createErrorView('email-error');
  var passwordErrorView = createErrorView('password-error');
  var errorView = createErrorView('error');

  if (errorView.getErrorMessage()) {
    setTimeout(function () {
      return errorView.hideError();
    }, 5000);
  }

  var initialState = {
    isSilentMode: true,
    passwordError: '',
    emailError: '',
    isPasswordValid: false,
    isEmailValid: false
  };

  var _createState = createState(initialState),
      getState = _createState.getState,
      subscribe = _createState.subscribe,
      update = _createState.update;

  formView.subscribeOnSubmit(function (e) {
    e.preventDefault();

    if (getState().isSilentMode) {
      validateFields();
    }

    update({
      isSilentMode: false
    });

    var _getState = getState(),
        isPasswordValid = _getState.isPasswordValid,
        isEmailValid = _getState.isEmailValid;

    if (isPasswordValid && isEmailValid) {
      form.submit();
    }
  });
  subscribe(function (state, prevState) {
    if (state.isSilentMode) {
      return;
    }

    var emailError = state.emailError,
        passwordError = state.passwordError;

    if (emailError) {
      emailErrorView.setError(emailError);
      emailErrorView.showError();
    } else if (prevState.emailError) {
      emailErrorView.hideError();
    }

    if (passwordError) {
      passwordErrorView.setError(passwordError);
      passwordErrorView.showError();
    } else if (prevState.passwordError) {
      passwordErrorView.hideError();
    }
  });

  var validateEmail = function validateEmail(value) {
    var newError = emailValidator.validate(value);

    var _getState2 = getState(),
        emailError = _getState2.emailError;

    if (emailError === newError) {
      return;
    }

    var isEmailValid = !newError;
    update({
      emailError: newError,
      isEmailValid: isEmailValid
    });
  };

  var validatePassword = function validatePassword(value) {
    var newError = strongPasswordValidator.validate(value);

    var _getState3 = getState(),
        passwordError = _getState3.passwordError;

    if (passwordError === newError) {
      return;
    }

    var isPasswordValid = !newError;
    update({
      passwordError: newError,
      isPasswordValid: isPasswordValid
    });
  };

  emailView.subscribeOnInputChange(function (e) {
    return validateEmail(e.target.value);
  });
  passwordView.subscribeOnInputChange(function (e) {
    return validatePassword(e.target.value);
  });

  var validateFields = function validateFields() {
    validatePassword(passwordView.value);
    validateEmail(emailView.value);
  };

}());
